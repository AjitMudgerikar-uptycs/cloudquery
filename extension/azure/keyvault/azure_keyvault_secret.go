package keyvault

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/fatih/structs"
)

const keyvaultSecret string = "azure_keyvault_secret"

// KeyvaultSecretColumns returns the list of columns in the table
func KeyvaultSecretColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("attributes"),
		table.TextColumn("attributes_created"),
		// table.BigIntColumn("attributes_created_ext"),
		// table.TextColumn("attributes_created_loc"),
		// table.BigIntColumn("attributes_created_loc_cache_end"),
		// table.BigIntColumn("attributes_created_loc_cache_start"),
		// table.TextColumn("attributes_created_loc_cache_zone"),
		// table.TextColumn("attributes_created_loc_cache_zone_is_dst"),
		// table.TextColumn("attributes_created_loc_cache_zone_name"),
		// table.IntegerColumn("attributes_created_loc_cache_zone_offset"),
		// table.TextColumn("attributes_created_loc_extend"),
		// table.TextColumn("attributes_created_loc_name"),
		// table.TextColumn("attributes_created_loc_tx"),
		// table.IntegerColumn("attributes_created_loc_tx_index"),
		// table.TextColumn("attributes_created_loc_tx_isstd"),
		// table.TextColumn("attributes_created_loc_tx_isutc"),
		// table.BigIntColumn("attributes_created_loc_tx_when"),
		// table.TextColumn("attributes_created_loc_zone"),
		// table.TextColumn("attributes_created_loc_zone_is_dst"),
		// table.TextColumn("attributes_created_loc_zone_name"),
		// table.IntegerColumn("attributes_created_loc_zone_offset"),
		// table.BigIntColumn("attributes_created_wall"),
		table.TextColumn("attributes_enabled"),
		table.TextColumn("attributes_exp"),
		// table.BigIntColumn("attributes_exp_ext"),
		// table.TextColumn("attributes_exp_loc"),
		// table.BigIntColumn("attributes_exp_loc_cache_end"),
		// table.BigIntColumn("attributes_exp_loc_cache_start"),
		// table.TextColumn("attributes_exp_loc_cache_zone"),
		// table.TextColumn("attributes_exp_loc_cache_zone_is_dst"),
		// table.TextColumn("attributes_exp_loc_cache_zone_name"),
		// table.IntegerColumn("attributes_exp_loc_cache_zone_offset"),
		// table.TextColumn("attributes_exp_loc_extend"),
		// table.TextColumn("attributes_exp_loc_name"),
		// table.TextColumn("attributes_exp_loc_tx"),
		// table.IntegerColumn("attributes_exp_loc_tx_index"),
		// table.TextColumn("attributes_exp_loc_tx_isstd"),
		// table.TextColumn("attributes_exp_loc_tx_isutc"),
		// table.BigIntColumn("attributes_exp_loc_tx_when"),
		// table.TextColumn("attributes_exp_loc_zone"),
		// table.TextColumn("attributes_exp_loc_zone_is_dst"),
		// table.TextColumn("attributes_exp_loc_zone_name"),
		// table.IntegerColumn("attributes_exp_loc_zone_offset"),
		// table.BigIntColumn("attributes_exp_wall"),
		table.TextColumn("attributes_nbf"),
		// table.BigIntColumn("attributes_nbf_ext"),
		// table.TextColumn("attributes_nbf_loc"),
		// table.BigIntColumn("attributes_nbf_loc_cache_end"),
		// table.BigIntColumn("attributes_nbf_loc_cache_start"),
		// table.TextColumn("attributes_nbf_loc_cache_zone"),
		// table.TextColumn("attributes_nbf_loc_cache_zone_is_dst"),
		// table.TextColumn("attributes_nbf_loc_cache_zone_name"),
		// table.IntegerColumn("attributes_nbf_loc_cache_zone_offset"),
		// table.TextColumn("attributes_nbf_loc_extend"),
		// table.TextColumn("attributes_nbf_loc_name"),
		// table.TextColumn("attributes_nbf_loc_tx"),
		// table.IntegerColumn("attributes_nbf_loc_tx_index"),
		// table.TextColumn("attributes_nbf_loc_tx_isstd"),
		// table.TextColumn("attributes_nbf_loc_tx_isutc"),
		// table.BigIntColumn("attributes_nbf_loc_tx_when"),
		// table.TextColumn("attributes_nbf_loc_zone"),
		// table.TextColumn("attributes_nbf_loc_zone_is_dst"),
		// table.TextColumn("attributes_nbf_loc_zone_name"),
		// table.IntegerColumn("attributes_nbf_loc_zone_offset"),
		// table.BigIntColumn("attributes_nbf_wall"),
		table.IntegerColumn("attributes_recoverable_days"),
		table.TextColumn("attributes_recovery_level"),
		table.TextColumn("attributes_updated"),
		// table.BigIntColumn("attributes_updated_ext"),
		// table.TextColumn("attributes_updated_loc"),
		// table.BigIntColumn("attributes_updated_loc_cache_end"),
		// table.BigIntColumn("attributes_updated_loc_cache_start"),
		// table.TextColumn("attributes_updated_loc_cache_zone"),
		// table.TextColumn("attributes_updated_loc_cache_zone_is_dst"),
		// table.TextColumn("attributes_updated_loc_cache_zone_name"),
		// table.IntegerColumn("attributes_updated_loc_cache_zone_offset"),
		// table.TextColumn("attributes_updated_loc_extend"),
		// table.TextColumn("attributes_updated_loc_name"),
		// table.TextColumn("attributes_updated_loc_tx"),
		// table.IntegerColumn("attributes_updated_loc_tx_index"),
		// table.TextColumn("attributes_updated_loc_tx_isstd"),
		// table.TextColumn("attributes_updated_loc_tx_isutc"),
		// table.BigIntColumn("attributes_updated_loc_tx_when"),
		// table.TextColumn("attributes_updated_loc_zone"),
		// table.TextColumn("attributes_updated_loc_zone_is_dst"),
		// table.TextColumn("attributes_updated_loc_zone_name"),
		// table.IntegerColumn("attributes_updated_loc_zone_offset"),
		// table.BigIntColumn("attributes_updated_wall"),
		table.TextColumn("content_type"),
		table.TextColumn("id"),
		table.TextColumn("kid"),
		table.TextColumn("managed"),
		table.TextColumn("tags"),
		table.TextColumn("value"),
	}
}

// KeyvaultSecretsGenerate returns the rows in the table for all configured accounts
func KeyvaultSecretsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": keyvaultSecret,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountKeyvaultSecrets(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": keyvaultSecret,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountKeyvaultSecrets(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountKeyvaultSecrets(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[keyvaultSecret]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": keyvaultSecret,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setKeyvaultSecretToTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setKeyvaultSecretToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	resources, err := getKeyvaultVaultData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      keyvaultSecret,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get keyvault vault list from api")
	}

	for _, vault := range *resources.Response().Value {
		structs.DefaultTagName = "json"
		setKeyvaultSecretToTableHelper(session, rg, wg, resultMap, tableConfig, *vault.Name)
	}
}
func setKeyvaultSecretToTableHelper(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, vaultName string) {

	SecretsList := make([]keyvault.SecretListResult, 0)
	vaultBaseURL := "https://" + vaultName + ".vault.azure.net"
	resourceItr, err := getKeyvaultSecretHelperData(session, rg, vaultBaseURL)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     keyvaultSecret,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get list from api")

	}
	resource := resourceItr.Response()
	SecretsList = append(SecretsList, resource)
	for _, SecretList := range SecretsList {

		structs.DefaultTagName = "json"
		resMap := structs.Map(SecretList)
		byteArr, err := json.Marshal(resMap)

		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     keyvaultSecret,
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
func getKeyvaultSecretHelperData(session *azure.AzureSession, rg string, vaultBaseURL string) (result keyvault.SecretListResultPage, err error) {

	var top int32 = 25
	svcClient := keyvault.New()
	svcClient.Authorizer = session.VaultAuthorizer
	return svcClient.GetSecrets(context.Background(), vaultBaseURL, &top)
}
