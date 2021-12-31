package graphrbac

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/fatih/structs"
)

const azureGraphrbacServicePrincipal string = "azure_graphrbac_service_principal"

// GraphrbacServicePrincipalColumns returns the list of columns in the table
func GraphrbacServicePrincipalColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("account_enabled"),
		table.TextColumn("alternative_names"),
		table.TextColumn("app_display_name"),
		table.TextColumn("app_id"),
		table.TextColumn("app_owner_tenant_id"),
		table.TextColumn("app_role_assignment_required"),
		table.TextColumn("app_roles"),
		table.TextColumn("app_roles_allowed_member_types"),
		table.TextColumn("app_roles_description"),
		table.TextColumn("app_roles_display_name"),
		table.TextColumn("app_roles_id"),
		table.TextColumn("app_roles_is_enabled"),
		table.TextColumn("app_roles_value"),
		table.TextColumn("deletion_timestamp"),
		table.TextColumn("display_name"),
		table.TextColumn("error_url"),
		table.TextColumn("homepage"),
		table.TextColumn("key_credentials"),
		table.TextColumn("key_credentials_custom_key_identifier"),
		table.TextColumn("key_credentials_end_date"),
		table.TextColumn("key_credentials_key_id"),
		table.TextColumn("key_credentials_start_date"),
		table.TextColumn("key_credentials_type"),
		table.TextColumn("key_credentials_usage"),
		table.TextColumn("key_credentials_value"),
		table.TextColumn("logout_url"),
		table.TextColumn("oauth2_permissions"),
		table.TextColumn("oauth2_permissions_admin_consent_description"),
		table.TextColumn("oauth2_permissions_admin_consent_display_name"),
		table.TextColumn("oauth2_permissions_id"),
		table.TextColumn("oauth2_permissions_is_enabled"),
		table.TextColumn("oauth2_permissions_type"),
		table.TextColumn("oauth2_permissions_user_consent_description"),
		table.TextColumn("oauth2_permissions_user_consent_display_name"),
		table.TextColumn("oauth2_permissions_value"),
		table.TextColumn("object_id"),
		table.TextColumn("object_type"),
		table.TextColumn("password_credentials"),
		table.TextColumn("password_credentials_custom_key_identifier"),
		table.TextColumn("password_credentials_end_date"),
		table.TextColumn("password_credentials_key_id"),
		table.TextColumn("password_credentials_start_date"),
		table.TextColumn("password_credentials_value"),
		table.TextColumn("preferred_token_signing_key_thumbprint"),
		table.TextColumn("publisher_name"),
		table.TextColumn("reply_urls"),
		table.TextColumn("saml_metadata_url"),
		table.TextColumn("service_principal_names"),
		table.TextColumn("service_principal_type"),
		table.TextColumn("tags"),
	}
}

// GraphrbacServicePrincipalGenerate returns the rows in the table for all configured accounts
func GraphrbacServicePrincipalGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureGraphrbacServicePrincipal,
			"account":   "default",
		}).Info("processing account")
		results, err := processGraphrbacServicePrincipal(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureGraphrbacServicePrincipal,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processGraphrbacServicePrincipal(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processGraphrbacServicePrincipal(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)

	session, err := azure.GetAuthSession(account)
	if err != nil {
		return resultMap, err
	}

	tableConfig, ok := utilities.TableConfigurationMap[azureGraphrbacServicePrincipal]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureGraphrbacServicePrincipal,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	setGraphrbacServicePrincipaltoTable(account.TenantID, session, &resultMap, tableConfig)

	return resultMap, nil
}

func setGraphrbacServicePrincipaltoTable(tenantId string, session *azure.AzureSession, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {

	for resourcesItr, err := getGraphrbacServicePrincipalData(session, tenantId); resourcesItr.NotDone(); resourcesItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureGraphrbacServicePrincipal,
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
				"tableName": azureGraphrbacServicePrincipal,
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
func getGraphrbacServicePrincipalData(session *azure.AzureSession, tenantId string) (result graphrbac.ServicePrincipalListResultIterator, err error) {
	svcClient := graphrbac.NewServicePrincipalsClient(tenantId)
	svcClient.Authorizer = session.GraphAuthorizer
	return svcClient.ListComplete(context.Background(), "")
}
