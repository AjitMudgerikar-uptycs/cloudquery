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

const SecuritycenterSubscriptionPricing string = "azure_securitycenter_subscription_pricing"

// SecuritycenterSubscriptionPricingColumns returns the list of columns in the table
func SecuritycenterSubscriptionPricingColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("free_trial_remaining_time"),
		table.TextColumn("pricing_tier"),
		table.TextColumn("type"),
	}
}

//SecuritycenterSubscriptionPricingGenerate returns the rows in the table for all configured accounts
func SecuritycenterSubscriptionPricingGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": SecuritycenterSubscriptionPricing,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountSecuritycenterSubscriptionPricing(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": SecuritycenterSubscriptionPricing,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountSecuritycenterSubscriptionPricing(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountSecuritycenterSubscriptionPricing(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap[SecuritycenterSubscriptionPricing]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": SecuritycenterSubscriptionPricing,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	setSecuritycenterSubscriptionPricingtoTable(session, "", &resultMap, tableConfig)

	return resultMap, nil
}
func setSecuritycenterSubscriptionPricingtoTable(session *azure.AzureSession, rg string, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {

	resources, err := getSecuritycenterSubscriptionPricingData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      SecuritycenterSubscriptionPricing,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get pricing list from api")
	}

	for _, contact := range *resources.Value {
		structs.DefaultTagName = "json"
		resMap := structs.Map(contact)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     SecuritycenterSubscriptionPricing,
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
func getSecuritycenterSubscriptionPricingData(session *azure.AzureSession, asclocation string) (result azuresecurity.PricingList, err error) {

	svcClient := azuresecurity.NewPricingsClient(session.SubscriptionId, asclocation)
	svcClient.Authorizer = session.Authorizer
	return svcClient.List(context.Background())

}
