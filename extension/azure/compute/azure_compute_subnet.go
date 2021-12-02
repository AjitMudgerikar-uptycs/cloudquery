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

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-03-01/network"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/azure"

	//extazure "github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

var tableName2 string = "azure_compute_subnet"

// VirtualSubnetColumns returns the list of columns in the table
func VirtualSubnetColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("address_prefix"),
		table.TextColumn("address_prefixes"),
		table.TextColumn("application_gateway_ip_configurations"),
		// table.TextColumn("application_gateway_ip_configurations_etag"),
		// table.TextColumn("application_gateway_ip_configurations_id"),
		// table.TextColumn("application_gateway_ip_configurations_name"),
		// table.TextColumn("application_gateway_ip_configurations_type"),
		table.TextColumn("delegations"),
		// table.TextColumn("delegations_etag"),
		// table.TextColumn("delegations_id"),
		// table.TextColumn("delegations_name"),
		// table.TextColumn("delegations_type"),
		table.TextColumn("ip_allocations"),
		// table.TextColumn("ip_allocations_id"),
		table.TextColumn("ip_configuration_profiles"),
		// table.TextColumn("ip_configuration_profiles_etag"),
		// table.TextColumn("ip_configuration_profiles_id"),
		// table.TextColumn("ip_configuration_profiles_name"),
		// table.TextColumn("ip_configuration_profiles_type"),
		// table.TextColumn("ip_configurations"),
		// table.TextColumn("ip_configurations_etag"),
		// table.TextColumn("ip_configurations_id"),
		// table.TextColumn("ip_configurations_name"),
		table.TextColumn("nat_gateway"),
		// table.TextColumn("nat_gateway_id"),
		table.TextColumn("network_security_group"),
		// table.TextColumn("network_security_group_etag"),
		// table.TextColumn("network_security_group_id"),
		// table.TextColumn("network_security_group_location"),
		// table.TextColumn("network_security_group_name"),
		// table.TextColumn("network_security_group_tags"),
		// table.TextColumn("network_security_group_type"),
		table.TextColumn("private_endpoint_network_policies"),
		table.TextColumn("private_endpoints"),
		// table.TextColumn("private_endpoints_etag"),
		// table.TextColumn("private_endpoints_extended_location"),
		// table.TextColumn("private_endpoints_extended_location_name"),
		// table.TextColumn("private_endpoints_extended_location_type"),
		// table.TextColumn("private_endpoints_id"),
		// table.TextColumn("private_endpoints_location"),
		// table.TextColumn("private_endpoints_name"),
		// table.TextColumn("private_endpoints_tags"),
		// table.TextColumn("private_endpoints_type"),
		table.TextColumn("private_link_service_network_policies"),
		// table.TextColumn("provisioning_state"),
		table.TextColumn("purpose"),
		table.TextColumn("resource_navigation_links"),
		// table.TextColumn("resource_navigation_links_etag"),
		// table.TextColumn("resource_navigation_links_id"),
		// table.TextColumn("resource_navigation_links_name"),
		// table.TextColumn("resource_navigation_links_type"),
		table.TextColumn("route_table"),
		// table.TextColumn("route_table_etag"),
		// table.TextColumn("route_table_id"),
		// table.TextColumn("route_table_location"),
		// table.TextColumn("route_table_name"),
		// table.TextColumn("route_table_tags"),
		// table.TextColumn("route_table_type"),
		table.TextColumn("service_association_links"),
		// table.TextColumn("service_association_links_etag"),
		// table.TextColumn("service_association_links_id"),
		// table.TextColumn("service_association_links_name"),
		// table.TextColumn("service_association_links_type"),
		table.TextColumn("service_endpoint_policies"),
		// table.TextColumn("service_endpoint_policies_etag"),
		// table.TextColumn("service_endpoint_policies_id"),
		// table.TextColumn("service_endpoint_policies_kind"),
		// table.TextColumn("service_endpoint_policies_location"),
		// table.TextColumn("service_endpoint_policies_name"),
		// table.TextColumn("service_endpoint_policies_tags"),
		// table.TextColumn("service_endpoint_policies_type"),
		table.TextColumn("service_endpoints"),
		// table.TextColumn("service_endpoints_locations"),
		table.TextColumn("service_endpoints_provisioning_state"),
		// table.TextColumn("service_endpoints_service"),
		table.TextColumn("type"),
	}
}

// VirtualSubnetsGenerate returns the rows in the table for all configured accounts
func VirtualSubnetsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": tableName2,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountVirtualSubnets(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": tableName2,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountVirtualSubnets(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountVirtualSubnets(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[tableName]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": tableName2,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getVirtualNetworksForSubnet(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}
func getVirtualNetworksForSubnet(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := network.NewVirtualNetworksClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     tableName2,
				"resourceGroup": rg,
				"errString":     err.Error(),
			}).Error("failed to get resource list")
			continue
		}

		resource := resourceItr.Value()

		getVirtualSubnets(session, rg, wg, resultMap, tableConfig, *resource.Name)

	}
}

func getVirtualSubnets(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig, networkName string) {

	svcClient := network.NewSubnetsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListComplete(context.Background(), rg, networkName); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     tableName2,
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
				"tableName":     tableName2,
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
