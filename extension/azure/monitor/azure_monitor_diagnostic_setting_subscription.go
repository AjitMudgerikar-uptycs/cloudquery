package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	extazure "github.com/Uptycs/cloudquery/extension/azure"
	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	azuremonitor "github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"

	azurestorage "github.com/Uptycs/cloudquery/extension/azure/storage"
)

const azureMonitorDiagnosticSettingsSubscription = "azure_monitor_diagnostic_settings_resource"

// DiagnosticSettingsSubscriptionColumns returns the list of columns in the table
func DiagnosticSettingsSubscriptionColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("event_hub_authorization_rule_id"),
		table.TextColumn("event_hub_name"),
		table.TextColumn("log_analytics_destination_type"),
		table.TextColumn("logs"),
		table.TextColumn("logs_category"),
		table.TextColumn("logs_enabled"),
		table.TextColumn("logs_retention_policy"),
		// table.IntegerColumn("logs_retention_policy_days"),
		// table.TextColumn("logs_retention_policy_enabled"),
		table.TextColumn("metrics"),
		table.TextColumn("metrics_category"),
		table.TextColumn("metrics_enabled"),
		table.TextColumn("metrics_retention_policy"),
		// table.IntegerColumn("metrics_retention_policy_days"),
		// table.TextColumn("metrics_retention_policy_enabled"),
		table.TextColumn("metrics_time_grain"),
		table.TextColumn("service_bus_rule_id"),
		table.TextColumn("storage_account_id"),
		table.TextColumn("workspace_id"),
	}
}

// DiagnosticSettingsSubscriptionGenerate returns the rows in the table for all configured accounts
func DiagnosticSettingsSubscriptionGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureMonitorDiagnosticSettingsSubscription,
			"account":   "default",
		}).Info("processing account")
		results, err := processDignosticSettingsSubscription(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureMonitorDiagnosticSettingsSubscription,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processDignosticSettingsSubscription(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processDignosticSettingsSubscription(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	var wg sync.WaitGroup
	session, err := extazure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}
	groups, err := extazure.GetGroups(session)

	if err != nil {
		return resultMap, err
	}

	wg.Add(len(groups))

	tableConfig, ok := utilities.TableConfigurationMap[azureMonitorDiagnosticSettingsSubscription]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureMonitorDiagnosticSettingsSubscription,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getStorageAccountIdForSubscription(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}
func getStorageAccountIdForSubscription(session *extazure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	diagnosticSettings := make([]azuremonitor.DiagnosticSettingsResource, 0)

	for resourceItr, err := azurestorage.GetStorageAccounts(session, rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureMonitorDiagnosticSettingsSubscription,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}
		resource := resourceItr.Value()

		getDiagnosticSettingSubscription(session, rg, *resource.ID, &diagnosticSettings)
	}

	addDignosticSettingsSubscription(session, rg, resultMap, tableConfig, diagnosticSettings)
}

func getDiagnosticSettingSubscription(session *extazure.AzureSession, rg string, resourceURI string, diagnosticSettings *[]azuremonitor.DiagnosticSettingsResource) {

	svcClient := azuremonitor.NewDiagnosticSettingsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	resourceItr, err := svcClient.List(context.Background(), resourceURI)

	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     azureMonitorDiagnosticSettingsSubscription,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get resource list")
	}

	resource := resourceItr.Value

	*diagnosticSettings = append(*diagnosticSettings, *resource...)
}

func addDignosticSettingsSubscription(session *extazure.AzureSession, rg string, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, diagnosticSettings []azuremonitor.DiagnosticSettingsResource) {
	for _, diagnosticSetting := range diagnosticSettings {
		structs.DefaultTagName = "json"
		resMap := structs.Map(diagnosticSetting)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureMonitorDiagnosticSettingsSubscription,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to marshal response")
			continue
		}

		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := extazure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
