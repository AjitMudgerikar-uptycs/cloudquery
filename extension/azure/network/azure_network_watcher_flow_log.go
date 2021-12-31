/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package network

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/azure"

	//extazure "github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

var azureNetworkWatcherFlowLog string = "azure_network_watcher_flow_log"

// AzureNetworkWatcherFlowLogColumns returns the list of columns in the table
func AzureNetworkWatcherFlowLogColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("enabled"),
		table.TextColumn("flow_analytics_configuration"),
		// table.TextColumn("flow_analytics_configuration_network_watcher_flow_analytics_configuration"),
		// table.TextColumn("flow_analytics_configuration_network_watcher_flow_analytics_configuration_enabled"),
		// table.IntegerColumn("flow_analytics_configuration_network_watcher_flow_analytics_configuration_traffic_analytics_interval"),
		// table.TextColumn("flow_analytics_configuration_network_watcher_flow_analytics_configuration_workspace_id"),
		// table.TextColumn("flow_analytics_configuration_network_watcher_flow_analytics_configuration_workspace_region"),
		// table.TextColumn("flow_analytics_configuration_network_watcher_flow_analytics_configuration_workspace_resource_id"),
		table.TextColumn("format"),
		// table.TextColumn("format_type"),
		// table.IntegerColumn("format_version"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("retention_policy"),
		// table.IntegerColumn("retention_policy_days"),
		// table.TextColumn("retention_policy_enabled"),
		table.TextColumn("storage_id"),
		table.TextColumn("target_resource_guid"),
		table.TextColumn("target_resource_id"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// AzureNetworkWatcherFlowLogsGenerate returns the rows in the table for all configured accounts
func AzureNetworkWatcherFlowLogsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureNetworkWatcherFlowLog,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountAzureNetworkWatcherFlowLogs(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureNetworkWatcherFlowLog,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountAzureNetworkWatcherFlowLogs(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountAzureNetworkWatcherFlowLogs(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureNetworkWatcherFlowLog]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureNetworkWatcherFlowLog,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getWatcherNameForFlowLogs(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getWatcherNameForFlowLogs(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	resources, err := GetWatcherName(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      azureNetworkWatcherFlowLog,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get list from api")
	}

	for _, watcher := range *resources.Value {
		setFlowLogToTableHelper(session, rg, wg, resultMap, tableConfig, *watcher.Name)
	}
}
func setFlowLogToTableHelper(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, watcherName string) {

	for resourceItr, err := getWatcherFlowLogHelperData(session, rg, watcherName); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureNetworkWatcherFlowLog,
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
				"tableName":     azureNetworkWatcherFlowLog,
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
func getWatcherFlowLogHelperData(session *azure.AzureSession, rg string, watcherName string) (result network.FlowLogListResultIterator, err error) {

	svcClient := network.NewFlowLogsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListComplete(context.Background(), rg, watcherName)

}

func GetWatcherName(session *azure.AzureSession, rg string) (result network.WatcherListResult, err error) {
	svcClient := network.NewWatchersClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.List(context.Background(), rg)
}
