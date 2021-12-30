package securitycenter

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	azuresecurity "github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/fatih/structs"
)

const SecuritycenterSetting string = "azure_securitycenter_setting"

// SecuritycenterSettingColumns returns the list of columns in the table
func SecuritycenterSettingColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("kind"),
		table.TextColumn("name"),
		table.TextColumn("type"),
	}
}

//SecuritycenterSettingGenerate returns the rows in the table for all configured accounts
func SecuritycenterSettingGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": SecuritycenterSetting,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountSecuritycenterSetting(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": SecuritycenterSetting,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountSecuritycenterSetting(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountSecuritycenterSetting(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap[SecuritycenterSetting]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": SecuritycenterSetting,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	setSecuritycenterSettingtoTable(session, "", &resultMap, tableConfig)

	return resultMap, nil
}
func setSecuritycenterSettingtoTable(session *azure.AzureSession, rg string, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {

	resources, err := getSecuritycenterSettingData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      SecuritycenterSetting,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get setting list from api")
	}

	for _, contact := range resources.Values() {
		structs.DefaultTagName = "json"
		resMap := structs.Map(contact)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     SecuritycenterSetting,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to marshal response")
			continue
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := azure.RowToMap(row, session.SubscriptionId, "", "", tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
func getSecuritycenterSettingData(session *azure.AzureSession, asclocation string) (result azuresecurity.SettingsListPage, err error) {

	svcClient := azuresecurity.NewSettingsClient(session.SubscriptionId, asclocation)
	svcClient.Authorizer = session.Authorizer
	return svcClient.List(context.Background())

}
