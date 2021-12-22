package keyvault

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/fatih/structs"
)

const keyvaultKey string = "azure_keyvault_key"

// KeyvaultKeyColumns returns the list of columns in the table
func KeyvaultKeyColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("attributes"),
		// table.BigIntColumn("attributes_created"),
		// table.TextColumn("attributes_enabled"),
		// table.BigIntColumn("attributes_exp"),
		// table.BigIntColumn("attributes_nbf"),
		// table.TextColumn("attributes_recovery_level"),
		// table.BigIntColumn("attributes_updated"),
		table.TextColumn("curve_name"),
		table.TextColumn("key_ops"),
		table.IntegerColumn("key_size"),
		table.TextColumn("key_uri"),
		table.TextColumn("key_uri_with_version"),
		table.TextColumn("kty"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// KeyvaultKeysGenerate returns the rows in the table for all configured accounts
func KeyvaultKeysGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": keyvaultKey,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountKeyvaultKeys(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": keyvaultKey,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountKeyvaultKeys(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountKeyvaultKeys(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	var wg sync.WaitGroup
	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}
	groups, err := azure.GetGroups(session)

	if err != nil {
		return resultMap, err
	}

	wg.Add(len(groups))

	tableConfig, ok := utilities.TableConfigurationMap[keyvaultKey]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": keyvaultKey,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setKeyvaultKeyToTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setKeyvaultKeyToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	resources, err := getKeyvaultVaultData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      keyvaultKey,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get keyvault vault list from api")
	}

	for _, vault := range *resources.Response().Value {
		structs.DefaultTagName = "json"
		setKeyvaultKeyToTableHelper(session, rg, wg, resultMap, tableConfig, *vault.Name)
	}
}
func setKeyvaultKeyToTableHelper(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, vaultName string) {

	KeysList := make([]keyvault.Key, 0)
	resourceItr, err := getKeyvaultKeyHelperData(session, rg, vaultName)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     keyvaultKey,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get list from api")

	}
	resource := resourceItr.Response().Value
	KeysList = append(KeysList, *resource...)
	for _, KeyList := range KeysList {

		structs.DefaultTagName = "json"
		resMap := structs.Map(KeyList)
		byteArr, err := json.Marshal(resMap)

		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     keyvaultKey,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to marshal response")
			continue
		}

		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := azure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
func getKeyvaultKeyHelperData(session *azure.AzureSession, rg string, vaultName string) (result keyvault.KeyListResultPage, err error) {

	svcClient := keyvault.NewKeysClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.List(context.Background(), rg, vaultName)

}
