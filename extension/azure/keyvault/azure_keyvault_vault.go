package keyvault

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/fatih/structs"
)

const keyvaultVault string = "azure_keyvault_vault"

// KeyvaultVaultColumns returns the list of columns in the table
func KeyvaultVaultColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("properties_access_policies"),
		// table.TextColumn("properties_access_policies_application_id"),
		// table.TextColumn("properties_access_policies_object_id"),
		// table.TextColumn("properties_access_policies_permissions"),
		// table.TextColumn("properties_access_policies_permissions_certificates"),
		// table.TextColumn("properties_access_policies_permissions_keys"),
		// table.TextColumn("properties_access_policies_permissions_secrets"),
		// table.TextColumn("properties_access_policies_permissions_storage"),
		// table.TextColumn("properties_access_policies_tenant_id"),
		table.TextColumn("properties_create_mode"),
		table.TextColumn("properties_enable_purge_protection"),
		table.TextColumn("properties_enable_rbac_authorization"),
		table.TextColumn("properties_enable_soft_delete"),
		table.TextColumn("properties_enabled_for_deployment"),
		table.TextColumn("properties_enabled_for_disk_encryption"),
		table.TextColumn("properties_enabled_for_template_deployment"),
		table.TextColumn("properties_network_acls"),
		// table.TextColumn("properties_network_acls_bypass"),
		// table.TextColumn("properties_network_acls_default_action"),
		// table.TextColumn("properties_network_acls_ip_rules"),
		// table.TextColumn("properties_network_acls_ip_rules_value"),
		// table.TextColumn("properties_network_acls_virtual_network_rules"),
		// table.TextColumn("properties_network_acls_virtual_network_rules_id"),
		table.TextColumn("properties_private_endpoint_connections"),
		table.TextColumn("properties_sku"),
		// table.TextColumn("properties_sku_family"),
		// table.TextColumn("properties_sku_name"),
		table.IntegerColumn("properties_soft_delete_retention_in_days"),
		table.TextColumn("properties_tenant_id"),
		table.TextColumn("properties_vault_uri"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// KeyvaultVaultsGenerate returns the rows in the table for all configured accounts
func KeyvaultVaultsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": keyvaultVault,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountKeyvaultVaults(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": keyvaultVault,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountKeyvaultVaults(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountKeyvaultVaults(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[keyvaultVault]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": keyvaultVault,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setKeyvaultVaultToTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setKeyvaultVaultToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	resources, err := getKeyvaultVaultData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      keyvaultVault,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get keyvault vault list from api")
	}

	for _, vault := range *resources.Response().Value {
		structs.DefaultTagName = "json"
		resMap := structs.Map(vault)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     keyvaultVault,
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
func getKeyvaultVaultData(session *azure.AzureSession, rg string) (result keyvault.VaultListResultPage, err error) {

	var top int32 = 1
	svcClient := keyvault.NewVaultsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListByResourceGroup(context.Background(), rg, &top)

}
