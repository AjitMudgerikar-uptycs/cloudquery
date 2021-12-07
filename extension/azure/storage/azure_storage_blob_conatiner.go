/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

const storageBlobContainer string = "azure_storage_blob_container"

// StorageBlobContainerColumns returns the list of columns in the table
func StorageBlobContainerColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("properties"),
		table.TextColumn("default_encryption_scope"),
		table.TextColumn("deleted"),
		table.TextColumn("deleted_time"),
		table.TextColumn("deny_encryption_scope_override"),
		table.TextColumn("has_immutability_policy"),
		table.TextColumn("has_legal_hold"),
		table.TextColumn("immutability_policy"),
		// table.TextColumn("immutability_policy_etag"),
		// table.TextColumn("immutability_policy_update_history"),
		// table.IntegerColumn("immutability_policy_update_history_immutability_period_since_creation_in_days"),
		// table.TextColumn("immutability_policy_update_history_object_identifier"),
		// table.TextColumn("immutability_policy_update_history_tenant_id"),
		// table.TextColumn("immutability_policy_update_history_timestamp"),
		// table.TextColumn("immutability_policy_update_history_update"),
		// table.TextColumn("immutability_policy_update_history_upn"),
		table.TextColumn("immutable_storage_with_versioning"),
		// table.TextColumn("immutable_storage_with_versioning_enabled"),
		// table.TextColumn("immutable_storage_with_versioning_migration_state"),
		// table.TextColumn("immutable_storage_with_versioning_time_stamp"),
		table.TextColumn("last_modified_time"),
		table.TextColumn("lease_duration"),
		table.TextColumn("lease_state"),
		table.TextColumn("lease_status"),
		table.TextColumn("legal_hold"),
		// table.TextColumn("legal_hold_has_legal_hold"),
		// table.TextColumn("legal_hold_tags"),
		// table.TextColumn("legal_hold_tags_object_identifier"),
		// table.TextColumn("legal_hold_tags_tag"),
		// table.TextColumn("legal_hold_tags_tenant_id"),
		// table.TextColumn("legal_hold_tags_timestamp"),
		// table.TextColumn("legal_hold_tags_upn"),
		table.TextColumn("metadata"),
		table.TextColumn("public_access"),
		table.IntegerColumn("remaining_retention_days"),
		table.TextColumn("version"),
		table.TextColumn("type"),
	}
}

// StorageBlobContainerGenerate returns the rows in the table for all configured accounts
func StorageBlobContainerGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageBlobContainer,
			"account":   "default",
		}).Info("processing account")
		results, err := processStorageBlobContainer(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": storageBlobContainer,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processStorageBlobContainer(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processStorageBlobContainer(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[storageBlobContainer]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageBlobContainer,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getStorageAccountsForBlobContainer(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getStorageAccountsForBlobContainer(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
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
		setStorageBlobContainerToTable(session, rg, wg, resultMap, tableConfig, *resource.Name)
	}
}

func setStorageBlobContainerToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, accountName string) {

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
		structs.DefaultTagName = "json"
		resMap := structs.Map(resource)
		byteArr, err := json.Marshal(resMap)

		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageBlobContainer,
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

func getStorageBlobContainerData(session *azure.AzureSession, rg string, accountName string) (result storage.ListContainerItemsIterator, err error) {
	svcClient := storage.NewBlobContainersClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListComplete(context.Background(), rg, accountName, "", "", storage.ListContainersIncludeDeleted)
}