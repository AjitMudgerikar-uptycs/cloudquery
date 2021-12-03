/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package compute

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	extazure "github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/fatih/structs"
)

var azureComputeDisk = "azure_compute_disk"

// DiskColumns returns the list of columns in the table
func DiskColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("extended_location"),
		//table.TextColumn("extended_location_name"),
		//table.TextColumn("extended_location_type"),
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("managed_by"),
		//table.TextColumn("managed_by_extended"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("bursting_enabled"),
		table.TextColumn("creation_data"),
		// table.TextColumn("creation_data_create_option"),
		// table.TextColumn("creation_data_gallery_image_reference"),
		// table.TextColumn("creation_data_gallery_image_reference_id"),
		// table.TextColumn("creation_data_gallery_image_reference_lun"),
		// table.TextColumn("creation_data_image_reference"),
		// table.TextColumn("creation_data_image_reference_id"),
		// table.TextColumn("creation_data_image_reference_lun"),
		// table.TextColumn("creation_data_logical_sector_size"),
		// table.TextColumn("creation_data_source_resource_id"),
		// table.TextColumn("creation_data_source_unique_id"),
		// table.TextColumn("creation_data_source_uri"),
		// table.TextColumn("creation_data_storage_account_id"),
		// table.BigIntColumn("creation_data_upload_size_bytes"),
		table.TextColumn("disk_access_id"),
		table.BigIntColumn("disk_iops_read_only"),
		table.BigIntColumn("disk_iops_read_write"),
		table.BigIntColumn("disk_m_bps_read_only"),
		table.BigIntColumn("disk_m_bps_read_write"),
		table.BigIntColumn("disk_size_bytes"),
		table.TextColumn("disk_size_gb"),
		table.TextColumn("disk_state"),
		table.TextColumn("encryption"),
		//table.TextColumn("encryption_settings_collection"),
		// table.TextColumn("encryption_settings_collection_enabled"),
		// table.TextColumn("encryption_settings_collection_encryption_settings"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_version"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_disk_encryption_key"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_disk_encryption_key_secret_url"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_disk_encryption_key_source_vault"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_disk_encryption_key_source_vault_id"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_key_encryption_key"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_key_encryption_key_key_url"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_key_encryption_key_source_vault"),
		// table.TextColumn("encryption_settings_collection_encryption_settings_key_encryption_key_source_vault_id"),
		table.TextColumn("encryption_disk_encryption_set_id"),
		table.TextColumn("encryption_type"),
		table.TextColumn("hyper_v_generation"),
		table.TextColumn("max_shares"),
		table.TextColumn("network_access_policy"),
		table.TextColumn("os_type"),
		table.TextColumn("property_updates_in_progress"),
		// table.TextColumn("property_updates_in_progress_target_tier"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("purchase_plan"),
		// table.TextColumn("purchase_plan_name"),
		// table.TextColumn("purchase_plan_product"),
		// table.TextColumn("purchase_plan_promotion_code"),
		// table.TextColumn("purchase_plan_publisher"),
		table.TextColumn("security_profile"),
		// table.TextColumn("security_profile_security_type"),
		table.TextColumn("share_info"),
		// table.TextColumn("share_info_vm_uri"),
		table.TextColumn("supports_hibernation"),
		table.TextColumn("tier"),
		table.TextColumn("time_created"),
		table.TextColumn("unique_id"),
		table.TextColumn("sku"),
		// table.TextColumn("sku_name"),
		// table.TextColumn("sku_tier"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
		table.TextColumn("zones"),
	}
}

// DiskGenerate returns the rows in the table for all configured accounts
func DiskGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureComputeDisk,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDisk(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureComputeDisk,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountDisk(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountDisk(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureComputeDisk]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureComputeDisk,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getDisk(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getDisk(session *extazure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := compute.NewDisksClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListByResourceGroupComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureComputeDisk,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()
		structs.DefaultTagName = "json"
		resMap := structs.Map(resource)
		utilities.GetLogger().Error(resMap)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureComputeDisk,
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
