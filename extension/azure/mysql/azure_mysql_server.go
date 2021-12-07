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

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
)

var azureMysqlServer = "azure_mysql_server"

// MysqlServerColumns returns the list of columns in the table
func MysqlServerColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("identity"),
		table.TextColumn("identity_principal_id"),
		table.TextColumn("identity_tenant_id"),
		table.TextColumn("identity_type"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("administrator_login"),
		table.TextColumn("byok_enforcement"),
		table.TextColumn("earliest_restore_date"),
		table.TextColumn("fully_qualified_domain_name"),
		table.TextColumn("infrastructure_encryption"),
		table.TextColumn("master_server_id"),
		table.TextColumn("minimal_tls_version"),
		table.TextColumn("private_endpoint_connections"),
		// table.TextColumn("private_endpoint_connections_id"),
		// table.TextColumn("private_endpoint_connections_properties"),
		// table.TextColumn("private_endpoint_connections_properties_private_endpoint"),
		// table.TextColumn("private_endpoint_connections_properties_private_endpoint_id"),
		// table.TextColumn("private_endpoint_connections_properties_private_link_service_connection_state"),
		// table.TextColumn("private_endpoint_connections_properties_private_link_service_connection_state_actions_required"),
		// table.TextColumn("private_endpoint_connections_properties_private_link_service_connection_state_description"),
		// table.TextColumn("private_endpoint_connections_properties_private_link_service_connection_state_status"),
		// table.TextColumn("private_endpoint_connections_properties_provisioning_state"),
		table.TextColumn("public_network_access"),
		table.IntegerColumn("replica_capacity"),
		table.TextColumn("replication_role"),
		table.TextColumn("ssl_enforcement"),
		table.TextColumn("storage_profile"),
		// table.IntegerColumn("storage_profile_backup_retention_days"),
		// table.TextColumn("storage_profile_geo_redundant_backup"),
		// table.TextColumn("storage_profile_storage_autogrow"),
		// table.IntegerColumn("storage_profile_storage_mb"),
		table.TextColumn("user_visible_state"),
		table.TextColumn("version"),
		table.TextColumn("sku"),
		table.IntegerColumn("sku_capacity"),
		table.TextColumn("sku_family"),
		table.TextColumn("sku_name"),
		table.TextColumn("sku_size"),
		table.TextColumn("sku_tier"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// MysqlServerGenerate returns the rows in the table for all configured accounts
func MysqlServerGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureMysqlServer,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountMysqlServer(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureMysqlServer,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountMysqlServer(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountMysqlServer(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureMysqlServer]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureMysqlServer,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, group := range groups {
		go getMysqlServer(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getMysqlServer(session *extazure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := mysql.NewServersClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	resourceItr, err := svcClient.List(context.Background())
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     azureMysqlServer,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to get resource list")
	}
	resource := resourceItr.Value
	utilities.GetLogger().Error(resource)
	byteArr, err := json.Marshal(resource)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":     azureMysqlServer,
			"resourceGroup": rg,
			"errString":     err.Error(),
		}).Error("failed to marshal response")
	}
	table := utilities.NewTable(byteArr, tableConfig)
	for _, row := range table.Rows {
		result := extazure.RowToMap(row, session.SubscriptionId, "", rg, tableConfig)
		*resultMap = append(*resultMap, result)
	}
}
