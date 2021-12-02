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

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/Uptycs/basequery-go/plugin/table"
	extazure "github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

var azureComputeSecurityGroup = "azure_compute_security_group"

// SecurityGroupsColumns returns the list of columns in the table
func SecurityGroupsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("default_security_rules"),
		// table.TextColumn("default_security_rules_etag"),
		// table.TextColumn("default_security_rules_id"),
		// table.TextColumn("default_security_rules_name"),
		// table.TextColumn("default_security_rules_type"),
		table.TextColumn("flow_logs"),
		// table.TextColumn("flow_logs_etag"),
		// table.TextColumn("flow_logs_id"),
		// table.TextColumn("flow_logs_location"),
		// table.TextColumn("flow_logs_name"),
		// table.TextColumn("flow_logs_tags"),
		// table.TextColumn("flow_logs_type"),
		table.TextColumn("network_interfaces"),
		// table.TextColumn("network_interfaces_etag"),
		// table.TextColumn("network_interfaces_extended_location"),
		// table.TextColumn("network_interfaces_extended_location_name"),
		// table.TextColumn("network_interfaces_extended_location_type"),
		// table.TextColumn("network_interfaces_id"),
		// table.TextColumn("network_interfaces_location"),
		// table.TextColumn("network_interfaces_name"),
		// table.TextColumn("network_interfaces_tags"),
		// table.TextColumn("network_interfaces_type"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("resource_guid"),
		table.TextColumn("security_rules"),
		// table.TextColumn("security_rules_etag"),
		// table.TextColumn("security_rules_id"),
		// table.TextColumn("security_rules_name"),
		// table.TextColumn("security_rules_type"),
		table.TextColumn("subnets"),
		// table.TextColumn("subnets_etag"),
		// table.TextColumn("subnets_id"),
		// table.TextColumn("subnets_name"),
		// table.TextColumn("subnets_type"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// SecurityGroupsGenerate returns the rows in the table for all configured accounts
func SecurityGroupsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureComputeSecurityGroup,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountSecurityGroups(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureComputeSecurityGroup,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountSecurityGroups(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountSecurityGroups(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureComputeSecurityGroup]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureComputeSecurityGroup,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go getSecurityGroups(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func getSecurityGroups(session *extazure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	svcClient := network.NewSecurityGroupsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer

	for resourceItr, err := svcClient.ListComplete(context.Background(), rg); resourceItr.NotDone(); err = resourceItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureComputeSecurityGroup,
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
				"tableName":     azureComputeSecurityGroup,
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
