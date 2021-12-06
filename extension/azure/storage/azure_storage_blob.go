package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/fatih/structs"
	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	azureazblob "github.com/Azure/azure-storage-blob-go/azblob"
)

const storageBlob string = "azure_storage_blob"

// StorageBlobColumns returns the list of columns in the table
func StorageBlobColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("access_tier"),
		table.TextColumn("access_tier_change_time"),
		table.BigIntColumn("access_tier_change_time_ext"),
		table.TextColumn("access_tier_change_time_loc"),
		// table.BigIntColumn("access_tier_change_time_loc_cache_end"),
		// table.BigIntColumn("access_tier_change_time_loc_cache_start"),
		// table.TextColumn("access_tier_change_time_loc_cache_zone"),
		// table.TextColumn("access_tier_change_time_loc_cache_zone_is_dst"),
		// table.TextColumn("access_tier_change_time_loc_cache_zone_name"),
		// table.IntegerColumn("access_tier_change_time_loc_cache_zone_offset"),
		// table.TextColumn("access_tier_change_time_loc_extend"),
		// table.TextColumn("access_tier_change_time_loc_name"),
		// table.TextColumn("access_tier_change_time_loc_tx"),
		// table.IntegerColumn("access_tier_change_time_loc_tx_index"),
		// table.TextColumn("access_tier_change_time_loc_tx_isstd"),
		// table.TextColumn("access_tier_change_time_loc_tx_isutc"),
		// table.BigIntColumn("access_tier_change_time_loc_tx_when"),
		// table.TextColumn("access_tier_change_time_loc_zone"),
		// table.TextColumn("access_tier_change_time_loc_zone_is_dst"),
		// table.TextColumn("access_tier_change_time_loc_zone_name"),
		// table.IntegerColumn("access_tier_change_time_loc_zone_offset"),
		table.BigIntColumn("access_tier_change_time_wall"),
		table.TextColumn("access_tier_inferred"),
		table.TextColumn("archive_status"),
		table.BigIntColumn("blob_sequence_number"),
		table.TextColumn("blob_type"),
		table.TextColumn("cache_control"),
		table.TextColumn("container_id"),
		table.TextColumn("container_name"),
		table.TextColumn("container_type"),
		table.TextColumn("content_disposition"),
		table.TextColumn("content_encoding"),
		table.TextColumn("content_language"),
		table.BigIntColumn("content_length"),
		table.TextColumn("content_md5"),
		table.TextColumn("content_type"),
		table.TextColumn("copy_id"),
		table.TextColumn("copy_status"),
		table.TextColumn("creation_time"),
		table.BigIntColumn("creation_time_ext"),
		table.TextColumn("creation_time_loc"),
		// table.BigIntColumn("creation_time_loc_cache_end"),
		// table.BigIntColumn("creation_time_loc_cache_start"),
		// table.TextColumn("creation_time_loc_cache_zone"),
		// table.TextColumn("creation_time_loc_cache_zone_is_dst"),
		// table.TextColumn("creation_time_loc_cache_zone_name"),
		// table.IntegerColumn("creation_time_loc_cache_zone_offset"),
		// table.TextColumn("creation_time_loc_extend"),
		// table.TextColumn("creation_time_loc_name"),
		// table.TextColumn("creation_time_loc_tx"),
		// table.IntegerColumn("creation_time_loc_tx_index"),
		// table.TextColumn("creation_time_loc_tx_isstd"),
		// table.TextColumn("creation_time_loc_tx_isutc"),
		// table.BigIntColumn("creation_time_loc_tx_when"),
		// table.TextColumn("creation_time_loc_zone"),
		// table.TextColumn("creation_time_loc_zone_is_dst"),
		// table.TextColumn("creation_time_loc_zone_name"),
		// table.IntegerColumn("creation_time_loc_zone_offset"),
		table.BigIntColumn("creation_time_wall"),
		table.TextColumn("customer_provided_key_sha256"),
		table.TextColumn("deleted"),
		table.TextColumn("deleted_time"),
		table.BigIntColumn("deleted_time_ext"),
		table.TextColumn("deleted_time_loc"),
		// table.BigIntColumn("deleted_time_loc_cache_end"),
		// table.BigIntColumn("deleted_time_loc_cache_start"),
		// table.TextColumn("deleted_time_loc_cache_zone"),
		// table.TextColumn("deleted_time_loc_cache_zone_is_dst"),
		// table.TextColumn("deleted_time_loc_cache_zone_name"),
		// table.IntegerColumn("deleted_time_loc_cache_zone_offset"),
		// table.TextColumn("deleted_time_loc_extend"),
		// table.TextColumn("deleted_time_loc_name"),
		// table.TextColumn("deleted_time_loc_tx"),
		// table.IntegerColumn("deleted_time_loc_tx_index"),
		// table.TextColumn("deleted_time_loc_tx_isstd"),
		// table.TextColumn("deleted_time_loc_tx_isutc"),
		// table.BigIntColumn("deleted_time_loc_tx_when"),
		// table.TextColumn("deleted_time_loc_zone"),
		// table.TextColumn("deleted_time_loc_zone_is_dst"),
		// table.TextColumn("deleted_time_loc_zone_name"),
		// table.IntegerColumn("deleted_time_loc_zone_offset"),
		table.BigIntColumn("deleted_time_wall"),
		table.TextColumn("encryption_scope"),
		table.TextColumn("etag"),
		table.TextColumn("expires_on"),
		table.BigIntColumn("expires_on_ext"),
		table.TextColumn("expires_on_loc"),
		// table.BigIntColumn("expires_on_loc_cache_end"),
		// table.BigIntColumn("expires_on_loc_cache_start"),
		// table.TextColumn("expires_on_loc_cache_zone"),
		// table.TextColumn("expires_on_loc_cache_zone_is_dst"),
		// table.TextColumn("expires_on_loc_cache_zone_name"),
		// table.IntegerColumn("expires_on_loc_cache_zone_offset"),
		// table.TextColumn("expires_on_loc_extend"),
		// table.TextColumn("expires_on_loc_name"),
		// table.TextColumn("expires_on_loc_tx"),
		// table.IntegerColumn("expires_on_loc_tx_index"),
		// table.TextColumn("expires_on_loc_tx_isstd"),
		// table.TextColumn("expires_on_loc_tx_isutc"),
		// table.BigIntColumn("expires_on_loc_tx_when"),
		// table.TextColumn("expires_on_loc_zone"),
		// table.TextColumn("expires_on_loc_zone_is_dst"),
		// table.TextColumn("expires_on_loc_zone_name"),
		// table.IntegerColumn("expires_on_loc_zone_offset"),
		table.BigIntColumn("expires_on_wall"),
		table.TextColumn("is_sealed"),
		table.TextColumn("is_snapshot"),
		table.TextColumn("last_modified"),
		table.BigIntColumn("last_modified_ext"),
		table.TextColumn("last_modified_loc"),
		// table.BigIntColumn("last_modified_loc_cache_end"),
		// table.BigIntColumn("last_modified_loc_cache_start"),
		// table.TextColumn("last_modified_loc_cache_zone"),
		// table.TextColumn("last_modified_loc_cache_zone_is_dst"),
		// table.TextColumn("last_modified_loc_cache_zone_name"),
		// table.IntegerColumn("last_modified_loc_cache_zone_offset"),
		// table.TextColumn("last_modified_loc_extend"),
		// table.TextColumn("last_modified_loc_name"),
		// table.TextColumn("last_modified_loc_tx"),
		// table.IntegerColumn("last_modified_loc_tx_index"),
		// table.TextColumn("last_modified_loc_tx_isstd"),
		// table.TextColumn("last_modified_loc_tx_isutc"),
		// table.BigIntColumn("last_modified_loc_tx_when"),
		// table.TextColumn("last_modified_loc_zone"),
		// table.TextColumn("last_modified_loc_zone_is_dst"),
		// table.TextColumn("last_modified_loc_zone_name"),
		// table.IntegerColumn("last_modified_loc_zone_offset"),
		table.BigIntColumn("last_modified_wall"),
		table.TextColumn("lease_duration"),
		table.TextColumn("lease_state"),
		table.TextColumn("lease_status"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		table.TextColumn("rehydrate_priority"),
		table.IntegerColumn("remaining_retention_days"),
		table.TextColumn("server_encrypted"),
		table.TextColumn("storage_account_id"),
		table.TextColumn("storage_account_name"),
		table.TextColumn("tags"),
		table.TextColumn("tags_key"),
		table.TextColumn("tags_value"),
		table.TextColumn("tags_xml_name"),
		// table.TextColumn("tags_xml_name_local"),
		// table.TextColumn("tags_xml_name_space"),
	}
}

// StorageBlobGenerate returns the rows in the table for all configured accounts
func StorageBlobGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageBlob,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountStorageBlob(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": storageBlob,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountStorageBlob(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountStorageBlob(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[storageBlob]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageBlob,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go addStorageAccountsForBlob(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func addStorageAccountsForBlob(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()
	for resourceItr, err := getStorageAccountData(session, rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobContainer,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		addStorageAccountKeysForBlob(session, rg, wg, resultMap, tableConfig, *resource.Name)
	}
}

func addStorageAccountKeysForBlob(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string) {

	svcClient := storage.NewAccountsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	accountClient, err := svcClient.ListKeys(context.Background(), rg, accountName, storage.ListKeyExpandKerb)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     storageBlobContainer,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get resource list")
	}

	addStorageBlobContainerForBlob(session, rg, wg, resultMap, tableConfig, accountName, *((*accountClient.Keys)[0].Value))
}

func addStorageBlobContainerForBlob(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string, accountKey string) {

	for resourceItr, err := getStorageBlobContainerData(session, rg, accountName); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobContainer,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		getStorageBlob(session, rg, wg, resultMap, tableConfig, accountName, accountKey, *resource.Name)
	}
}

func getStorageBlob(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string, accountKey string, containerName string) {
	credential, err := azureazblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     storageBlobContainer,
			"resourceGroup": rg,
			"errString":     err.Error(),
			"accountName": accountName,
		}).Error("failed to get credentials")
		return
	}

	p := azureazblob.NewPipeline(credential, azureazblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))

	serviceURL := azureazblob.NewServiceURL(*u, p)
	containerURL := serviceURL.NewContainerURL(containerName)

	for marker := (azureazblob.Marker{}); marker.NotDone(); {

		listBlob, err := containerURL.ListBlobsFlatSegment(context.Background(), marker, azureazblob.ListBlobsSegmentOptions{})
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobContainer,
				"resourceGroup": rg,
				"errString":     err.Error(),
				"accountName": accountName,
			}).Error("failed to get blob")
			return
		}

		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {

			structs.DefaultTagName = "json"
			resMap := structs.Map(blobInfo)
			byteArr, err := json.Marshal(resMap)
			if err != nil {
				utilities.GetLogger().WithFields(log.Fields{
					"tableName":     storageAccount,
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
}

