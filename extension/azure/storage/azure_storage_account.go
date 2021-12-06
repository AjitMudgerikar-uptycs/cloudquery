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

const storageAccount string = "azure_storage_account"

// StorageAccountsColumns returns the list of columns in the table
func StorageAccountColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("extended_location"),
		table.TextColumn("extended_location_name"),
		table.TextColumn("extended_location_type"),
		table.TextColumn("id"),
		table.TextColumn("identity"),
		table.TextColumn("identity_principal_id"),
		table.TextColumn("identity_tenant_id"),
		table.TextColumn("identity_type"),
		table.TextColumn("identity_user_assigned_identities"),
		table.TextColumn("kind"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		table.TextColumn("properties"),
		table.TextColumn("access_tier"),
		table.TextColumn("allow_blob_public_access"),
		table.TextColumn("allow_cross_tenant_replication"),
		table.TextColumn("allow_shared_key_access"),
		table.TextColumn("azure_files_identity_based_authentication"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties_azure_storage_sid"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties_domain_guid"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties_domain_name"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties_domain_sid"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties_forest_name"),
		// table.TextColumn("azure_files_identity_based_authentication_active_directory_properties_net_bios_domain_name"),
		// table.TextColumn("azure_files_identity_based_authentication_default_share_permission"),
		// table.TextColumn("azure_files_identity_based_authentication_directory_service_options"),
		table.TextColumn("blob_restore_status"),
		// table.TextColumn("blob_restore_status_failure_reason"),
		// table.TextColumn("blob_restore_status_parameters"),
		// table.TextColumn("blob_restore_status_parameters_blob_ranges"),
		// table.TextColumn("blob_restore_status_parameters_blob_ranges_end_range"),
		// table.TextColumn("blob_restore_status_parameters_blob_ranges_start_range"),
		// table.TextColumn("blob_restore_status_parameters_time_to_restore"),
		// table.TextColumn("blob_restore_status_restore_id"),
		// table.TextColumn("blob_restore_status_status"),
		table.TextColumn("creation_time"),
		table.TextColumn("custom_domain"),
		// table.TextColumn("custom_domain_name"),
		// table.TextColumn("custom_domain_use_sub_domain_name"),
		table.TextColumn("encryption"),
		// table.TextColumn("encryption_identity"),
		// table.TextColumn("encryption_identity_user_assigned_identity"),
		// table.TextColumn("encryption_key_source"),
		// table.TextColumn("encryption_keyvaultproperties"),
		// table.TextColumn("encryption_keyvaultproperties_current_versioned_key_identifier"),
		// table.TextColumn("encryption_keyvaultproperties_keyname"),
		// table.TextColumn("encryption_keyvaultproperties_keyvaulturi"),
		// table.TextColumn("encryption_keyvaultproperties_keyversion"),
		// table.TextColumn("encryption_keyvaultproperties_last_key_rotation_timestamp"),
		// table.TextColumn("encryption_require_infrastructure_encryption"),
		// table.TextColumn("encryption_services"),
		// table.TextColumn("encryption_services_blob"),
		// table.TextColumn("encryption_services_blob_enabled"),
		// table.TextColumn("encryption_services_blob_key_type"),
		// table.TextColumn("encryption_services_blob_last_enabled_time"),
		// table.TextColumn("encryption_services_file"),
		// table.TextColumn("encryption_services_file_enabled"),
		// table.TextColumn("encryption_services_file_key_type"),
		// table.TextColumn("encryption_services_file_last_enabled_time"),
		// table.TextColumn("encryption_services_queue"),
		// table.TextColumn("encryption_services_queue_enabled"),
		// table.TextColumn("encryption_services_queue_key_type"),
		// table.TextColumn("encryption_services_queue_last_enabled_time"),
		// table.TextColumn("encryption_services_table"),
		// table.TextColumn("encryption_services_table_enabled"),
		// table.TextColumn("encryption_services_table_key_type"),
		// table.TextColumn("encryption_services_table_last_enabled_time"),
		table.TextColumn("failover_in_progress"),
		table.TextColumn("geo_replication_stats"),
		// table.TextColumn("geo_replication_stats_can_failover"),
		// table.TextColumn("geo_replication_stats_last_sync_time"),
		// table.TextColumn("geo_replication_stats_status"),
		table.TextColumn("is_hns_enabled"),
		table.TextColumn("is_nfs_v3_enabled"),
		table.TextColumn("key_creation_time"),
		// table.TextColumn("key_creation_time_key1"),
		// table.TextColumn("key_creation_time_key2"),
		table.TextColumn("key_policy"),
		// table.IntegerColumn("key_policy_key_expiration_period_in_days"),
		table.TextColumn("large_file_shares_state"),
		table.TextColumn("last_geo_failover_time"),
		table.TextColumn("minimum_tls_version"),
		table.TextColumn("network_acls"),
		// table.TextColumn("network_acls_bypass"),
		// table.TextColumn("network_acls_default_action"),
		// table.TextColumn("network_acls_ip_rules"),
		// table.TextColumn("network_acls_ip_rules_action"),
		// table.TextColumn("network_acls_ip_rules_value"),
		// table.TextColumn("network_acls_resource_access_rules"),
		// table.TextColumn("network_acls_resource_access_rules_resource_id"),
		// table.TextColumn("network_acls_resource_access_rules_tenant_id"),
		// table.TextColumn("network_acls_virtual_network_rules"),
		// table.TextColumn("network_acls_virtual_network_rules_action"),
		// table.TextColumn("network_acls_virtual_network_rules_id"),
		// table.TextColumn("network_acls_virtual_network_rules_state"),
		table.TextColumn("primary_endpoints"),
		// table.TextColumn("primary_endpoints_blob"),
		// table.TextColumn("primary_endpoints_dfs"),
		// table.TextColumn("primary_endpoints_file"),
		// table.TextColumn("primary_endpoints_internet_endpoints"),
		// table.TextColumn("primary_endpoints_internet_endpoints_blob"),
		// table.TextColumn("primary_endpoints_internet_endpoints_dfs"),
		// table.TextColumn("primary_endpoints_internet_endpoints_file"),
		// table.TextColumn("primary_endpoints_internet_endpoints_web"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints_blob"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints_dfs"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints_file"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints_queue"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints_table"),
		// table.TextColumn("primary_endpoints_microsoft_endpoints_web"),
		// table.TextColumn("primary_endpoints_queue"),
		// table.TextColumn("primary_endpoints_table"),
		// table.TextColumn("primary_endpoints_web"),
		table.TextColumn("primary_location"),
		table.TextColumn("private_endpoint_connections"),
		// table.TextColumn("private_endpoint_connections_id"),
		// table.TextColumn("private_endpoint_connections_name"),
		// table.TextColumn("private_endpoint_connections_type"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("routing_preference"),
		// table.TextColumn("routing_preference_publish_internet_endpoints"),
		// table.TextColumn("routing_preference_publish_microsoft_endpoints"),
		// table.TextColumn("routing_preference_routing_choice"),
		table.TextColumn("sas_policy"),
		// table.TextColumn("sas_policy_expiration_action"),
		// table.TextColumn("sas_policy_sas_expiration_period"),
		table.TextColumn("secondary_endpoints"),
		// table.TextColumn("secondary_endpoints_blob"),
		// table.TextColumn("secondary_endpoints_dfs"),
		// table.TextColumn("secondary_endpoints_file"),
		// table.TextColumn("secondary_endpoints_internet_endpoints"),
		// table.TextColumn("secondary_endpoints_internet_endpoints_blob"),
		// table.TextColumn("secondary_endpoints_internet_endpoints_dfs"),
		// table.TextColumn("secondary_endpoints_internet_endpoints_file"),
		// table.TextColumn("secondary_endpoints_internet_endpoints_web"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints_blob"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints_dfs"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints_file"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints_queue"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints_table"),
		// table.TextColumn("secondary_endpoints_microsoft_endpoints_web"),
		// table.TextColumn("secondary_endpoints_queue"),
		// table.TextColumn("secondary_endpoints_table"),
		// table.TextColumn("secondary_endpoints_web"),
		table.TextColumn("secondary_location"),
		table.TextColumn("status_of_primary"),
		table.TextColumn("status_of_secondary"),
		table.TextColumn("supports_https_traffic_only"),
		table.TextColumn("sku"),
		table.TextColumn("sku_name"),
		table.TextColumn("sku_tier"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// StorageAccountsGenerate returns the rows in the table for all configured accounts
func StorageAccountsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageAccount,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountStorageAccounts(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": storageAccount,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountStorageAccounts(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountStorageAccounts(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[storageAccount]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": storageAccount,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setStorageAccountDataToTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setStorageAccountDataToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	for resourceItr, err := getStorageAccountData(session, rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     storageAccount,
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

func getStorageAccountData(session *azure.AzureSession, rg string) (result storage.AccountListResultIterator, err error) {
	svcClient := storage.NewAccountsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListByResourceGroupComplete(context.Background(), rg)
}
