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

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-03-01/network"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/azure"

	//extazure "github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

var azureNetworkLoadBalancer string = "azure_network_load_balancer"

// NetworkLoadBalancerColumns returns the list of columns in the table
func NetworkLoadBalancerColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("extended_location"),
		table.TextColumn("extended_location_name"),
		table.TextColumn("extended_location_type"),
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("backend_address_pools"),
		// table.TextColumn("backend_address_pools_etag"),
		// table.TextColumn("backend_address_pools_id"),
		// table.TextColumn("backend_address_pools_name"),
		// table.TextColumn("backend_address_pools_type"),
		table.TextColumn("frontend_ip_configurations"),
		// table.TextColumn("frontend_ip_configurations_etag"),
		// table.TextColumn("frontend_ip_configurations_id"),
		// table.TextColumn("frontend_ip_configurations_name"),
		// table.TextColumn("frontend_ip_configurations_type"),
		// table.TextColumn("frontend_ip_configurations_zones"),
		table.TextColumn("inbound_nat_pools"),
		// table.TextColumn("inbound_nat_pools_etag"),
		// table.TextColumn("inbound_nat_pools_id"),
		// table.TextColumn("inbound_nat_pools_name"),
		// table.TextColumn("inbound_nat_pools_type"),
		table.TextColumn("inbound_nat_rules"),
		// table.TextColumn("inbound_nat_rules_etag"),
		// table.TextColumn("inbound_nat_rules_id"),
		// table.TextColumn("inbound_nat_rules_name"),
		// table.TextColumn("inbound_nat_rules_type"),
		table.TextColumn("load_balancing_rules"),
		// table.TextColumn("load_balancing_rules_etag"),
		// table.TextColumn("load_balancing_rules_id"),
		// table.TextColumn("load_balancing_rules_name"),
		// table.TextColumn("load_balancing_rules_type"),
		table.TextColumn("outbound_rules"),
		// table.TextColumn("outbound_rules_etag"),
		// table.TextColumn("outbound_rules_id"),
		// table.TextColumn("outbound_rules_name"),
		// table.TextColumn("outbound_rules_type"),
		table.TextColumn("probes"),
		// table.TextColumn("probes_etag"),
		// table.TextColumn("probes_id"),
		// table.TextColumn("probes_name"),
		// table.TextColumn("probes_type"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("resource_guid"),
		table.TextColumn("sku"),
		table.TextColumn("sku_name"),
		table.TextColumn("sku_tier"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// NetworkLoadBalancersGenerate returns the rows in the table for all configured accounts
func NetworkLoadBalancersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureNetworkLoadBalancer,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountNetworkLoadBalancers(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureNetworkLoadBalancer,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountNetworkLoadBalancers(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountNetworkLoadBalancers(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureNetworkLoadBalancer]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureNetworkLoadBalancer,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getNetworkLoadBalancers(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getNetworkLoadBalancers(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := network.NewLoadBalancersClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureNetworkLoadBalancer,
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
				"tableName":     azureNetworkLoadBalancer,
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
