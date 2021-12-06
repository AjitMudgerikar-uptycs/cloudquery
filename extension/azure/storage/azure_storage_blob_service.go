package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

const storageBlobService string = "azure_storage_blob_service"

// StorageBlobServiceColumns returns the list of columns in the table
func StorageBlobServiceColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("automatic_snapshot_policy_enabled"),
		table.TextColumn("change_feed"),
		// table.TextColumn("change_feed_enabled"),
		// table.IntegerColumn("change_feed_retention_in_days"),
		table.TextColumn("container_delete_retention_policy"),
		// table.IntegerColumn("container_delete_retention_policy_days"),
		// table.TextColumn("container_delete_retention_policy_enabled"),
		table.TextColumn("cors"),
		// table.TextColumn("cors_cors_rules"),
		// table.TextColumn("cors_cors_rules_allowed_headers"),
		// table.TextColumn("cors_cors_rules_allowed_methods"),
		// table.TextColumn("cors_cors_rules_allowed_origins"),
		// table.TextColumn("cors_cors_rules_exposed_headers"),
		// table.IntegerColumn("cors_cors_rules_max_age_in_seconds"),
		table.TextColumn("default_service_version"),
		table.TextColumn("delete_retention_policy"),
		// table.IntegerColumn("delete_retention_policy_days"),
		// table.TextColumn("delete_retention_policy_enabled"),
		table.TextColumn("is_versioning_enabled"),
		table.TextColumn("last_access_time_tracking_policy"),
		// table.TextColumn("last_access_time_tracking_policy_blob_type"),
		// table.TextColumn("last_access_time_tracking_policy_enable"),
		// table.TextColumn("last_access_time_tracking_policy_name"),
		// table.IntegerColumn("last_access_time_tracking_policy_tracking_granularity_in_days"),
		table.TextColumn("restore_policy"),
		// table.IntegerColumn("restore_policy_days"),
		// table.TextColumn("restore_policy_enabled"),
		// table.TextColumn("restore_policy_last_enabled_time"),
		// table.TextColumn("restore_policy_min_restore_time"),
		table.TextColumn("sku"),
		table.TextColumn("sku_name"),
		table.TextColumn("sku_tier"),
		table.TextColumn("type"),
	}
}

// StorageBlobServicesGenerate returns the rows in the table for all configured accounts
func StorageBlobServicesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageBlobService,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountStorageBlobServices(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": storageBlobService,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountStorageBlobServices(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountStorageBlobServices(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[storageBlobService]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageBlobService,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getAccountsForStorageBlobServices(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getAccountsForStorageBlobServices(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()
	svcClient := storage.NewAccountsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	var blobServices []storage.BlobServiceProperties
	for resourceItr, err := svcClient.ListByResourceGroupComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobService,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		blobService, err := getStorageBlobServicesHelper(session, rg, wg, resultMap, tableConfig, *resource.Name)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobService,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}
		blobServices = append(blobServices, *blobService...)
	}
	for _, blobService := range blobServices {
		structs.DefaultTagName = "json"
		resMap := structs.Map(blobService)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobService,
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
func getStorageBlobServicesHelper(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string) (*[]storage.BlobServiceProperties, error) {

	svcClient := storage.NewBlobServicesClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	resourceItr, err := svcClient.List(context.Background(), rg, accountName)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     storageBlobService,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get resource list")
		return nil, err
	}
	resource := resourceItr.Value
	return resource, nil
}
