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

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/fatih/structs"
)

const storageFileService string = "azure_storage_file_service"

// StorageFileServiceColumns returns the list of columns in the table
func StorageFileServiceColumns() []table.ColumnDefinition {
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
		// table.IntegerColumn("cors_cors_rules_max_age_in_seconds"),
		table.TextColumn("protocol_settings"),
		// table.TextColumn("protocol_settings_smb"),
		// table.TextColumn("protocol_settings_smb_authentication_methods"),
		// table.TextColumn("protocol_settings_smb_channel_encryption"),
		// table.TextColumn("protocol_settings_smb_kerberos_ticket_encryption"),
		// table.TextColumn("protocol_settings_smb_multichannel"),
		// table.TextColumn("protocol_settings_smb_multichannel_enabled"),
		// table.TextColumn("protocol_settings_smb_versions"),
		table.TextColumn("share_delete_retention_policy"),
		// table.IntegerColumn("share_delete_retention_policy_days"),
		// table.TextColumn("share_delete_retention_policy_enabled"),
		table.TextColumn("sku"),
		table.TextColumn("sku_name"),
		table.TextColumn("sku_tier"),
		table.TextColumn("type"),
	}
}

// StorageFileServicesGenerate returns the rows in the table for all configured accounts
func StorageFileServicesGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageFileService,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountStorageFileServices(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": storageFileService,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountStorageFileServices(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountStorageFileServices(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[storageFileService]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageFileService,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getAccountsForStorageFileServices(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getAccountsForStorageFileServices(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := storage.NewAccountsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListByResourceGroupComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageFileService,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		// b, _ := json.MarshalIndent(resource, "", " ")
		// fmt.Println(string(b))
		fmt.Println(*resource.Name)
		getStorageFileServicesHelper(session, rg, wg, resultMap, tableConfig, *resource.Name)
	}
}
func getStorageFileServicesHelper(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string) {
	//defer wg.Done()

	svcClient := storage.NewFileServicesClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	resourceItr, err := svcClient.List(context.Background(), rg, accountName) //resourceItr.NotDone(); err = resourceItr.Next() {
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     storageFileService,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get resource list")
		//continue
	}

	//resource := resourceItr.Value()

	structs.DefaultTagName = "json"
	resMap := structs.Map(resourceItr)
	byteArr, err := json.Marshal(resMap)
	//b, _ := json.MarshalIndent(byteArr, "", " ")
	//fmt.Println(string(byteArr))
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     storageFileService,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to marshal response")
		//continue
	}

	table := utilities.NewTable(byteArr, tableConfig)

	for _, row := range table.Rows {
		result := azure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
		*resultMap = append(*resultMap, result)
	}
}
