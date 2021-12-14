package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
)

const storageQueueService string = "azure_storage_queue_service"

// StorageQueueServicesColumns returns the list of columns in the table
func StorageQueueServicesColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("cors"),
		// table.TextColumn("cors_cors_rules"),
		// table.TextColumn("cors_cors_rules_allowed_headers"),
		// table.TextColumn("cors_cors_rules_allowed_methods"),
		// table.TextColumn("cors_cors_rules_allowed_origins"),
		// table.TextColumn("cors_cors_rules_exposed_headers"),
		// table.TextColumn("cors_cors_rules_max_age_in_seconds"),
		table.TextColumn("type"),
	}
}

// StorageQueueServicesGenerate returns the rows in the table for all configured accounts
func StorageQueueServicesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageQueueService,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountStorageQueueServices(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": storageQueueService,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountStorageQueueServices(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountStorageQueueServices(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[storageQueueService]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageQueueService,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getStorageAccountsForStorageQueueServices(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}
func getStorageAccountsForStorageQueueServices(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	for resourceItr, err := getStorageAccountData(session, rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageQueueService,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		setStorageQueueServicesToTable(session, rg, wg, resultMap, tableConfig, *resource.Name)
	}
}
func setStorageQueueServicesToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string) {

	resource, err := getStorageQueueServicesData(session, rg, accountName)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     storageQueueService,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get Queueservice list")
	}

	for _, Queueservice := range *resource.Value {

		structs.DefaultTagName = "json"
		resMap := structs.Map(Queueservice)
		byteArr, err := json.Marshal(resMap)

		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageQueueService,
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

func getStorageQueueServicesData(session *azure.AzureSession, rg string, accountName string) (result storage.ListQueueServices, err error) {

	svcClient := storage.NewQueueServicesClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.List(context.Background(), rg, accountName)

}
