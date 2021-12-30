/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package graphrbac

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	//"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/graphrbac/graphrbac"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/extension/azure"

	//extazure "github.com/Uptycs/cloudquery/extension/azure"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

var azureGraphrbacUser string = "azure_graphrbac_user"

// GraphrbacUserColumns returns the list of columns in the table
func GraphrbacUserColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_enabled"),
		table.TextColumn("deletion_timestamp"),
		table.TextColumn("display_name"),
		table.TextColumn("given_name"),
		table.TextColumn("immutable_id"),
		table.TextColumn("mail"),
		table.TextColumn("mail_nickname"),
		table.TextColumn("object_id"),
		table.TextColumn("object_type"),
		table.TextColumn("sign_in_names"),
		table.TextColumn("sign_in_names_type"),
		table.TextColumn("sign_in_names_value"),
		table.TextColumn("surname"),
		table.TextColumn("usage_location"),
		table.TextColumn("user_principal_name"),
		table.TextColumn("user_type"),
	}
}

// GraphrbacUsersGenerate returns the rows in the table for all configured accounts
func GraphrbacUsersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureGraphrbacUser,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountGraphrbacUsers(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureGraphrbacUser,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountGraphrbacUsers(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}
func processAccountGraphrbacUsers(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)

	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}

	tableConfig, ok := utilities.TableConfigurationMap[azureGraphrbacUser]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureGraphrbacUser,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	setGraphrbacUserstoTable(account.TenantID, session, &resultMap, tableConfig)

	return resultMap, nil
}

func setGraphrbacUserstoTable(tenantId string, session *azure.AzureSession, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {

	for resourcesItr, err := getGraphrbacUsersData(session, tenantId); resourcesItr.NotDone(); resourcesItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureGraphrbacUser,
				"TenantId":  tenantId,
				"errString": err.Error(),
			}).Error("failed to get DNS zones")
		}

		resource := resourcesItr.Value()

		structs.DefaultTagName = "json"
		resMap := structs.Map(resource)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureGraphrbacUser,
				"TenantId":  tenantId,
				"errString": err.Error(),
			}).Error("failed to marshal response")
			continue
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			result := azure.RowToMap(row, session.SubscriptionId, "", "", tableConfig)
			*resultMap = append(*resultMap, result)
		}
	}
}
func getGraphrbacUsersData(session *azure.AzureSession, tenantId string) (result graphrbac.UserListResultIterator, err error) {
	svcClient := graphrbac.NewUsersClient(tenantId)
	svcClient.Authorizer = session.GraphAuthorizer
	return svcClient.ListComplete(context.Background(), "", "")
}
