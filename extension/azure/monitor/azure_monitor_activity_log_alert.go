package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	azuremonitor "github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"
	"github.com/fatih/structs"
)

const monitorActivityLogAlert string = "azure_monitor_activity_log_alert"

// monitorActivityLogAlertColumns returns the list of columns in the table
func MonitorActivityLogAlertColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("id"),
		table.TextColumn("identity"),
		table.TextColumn("kind"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("actions"),
		// table.TextColumn("actions_action_groups"),
		// table.TextColumn("actions_action_groups_arm_role_receivers"),
		// table.TextColumn("actions_action_groups_arm_role_receivers_name"),
		// table.TextColumn("actions_action_groups_arm_role_receivers_role_id"),
		// table.TextColumn("actions_action_groups_arm_role_receivers_use_common_alert_schema"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_automation_account_id"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_is_global_runbook"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_name"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_runbook_name"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_service_uri"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_use_common_alert_schema"),
		// table.TextColumn("actions_action_groups_automation_runbook_receivers_webhook_resource_id"),
		// table.TextColumn("actions_action_groups_azure_app_push_receivers"),
		// table.TextColumn("actions_action_groups_azure_app_push_receivers_email_address"),
		// table.TextColumn("actions_action_groups_azure_app_push_receivers_name"),
		// table.TextColumn("actions_action_groups_azure_function_receivers"),
		// table.TextColumn("actions_action_groups_azure_function_receivers_function_app_resource_id"),
		// table.TextColumn("actions_action_groups_azure_function_receivers_function_name"),
		// table.TextColumn("actions_action_groups_azure_function_receivers_http_trigger_url"),
		// table.TextColumn("actions_action_groups_azure_function_receivers_name"),
		// table.TextColumn("actions_action_groups_azure_function_receivers_use_common_alert_schema"),
		// table.TextColumn("actions_action_groups_email_receivers"),
		// table.TextColumn("actions_action_groups_email_receivers_email_address"),
		// table.TextColumn("actions_action_groups_email_receivers_name"),
		// table.TextColumn("actions_action_groups_email_receivers_status"),
		// table.TextColumn("actions_action_groups_email_receivers_use_common_alert_schema"),
		// table.TextColumn("actions_action_groups_enabled"),
		// table.TextColumn("actions_action_groups_group_short_name"),
		// table.TextColumn("actions_action_groups_itsm_receivers"),
		// table.TextColumn("actions_action_groups_itsm_receivers_connection_id"),
		// table.TextColumn("actions_action_groups_itsm_receivers_name"),
		// table.TextColumn("actions_action_groups_itsm_receivers_region"),
		// table.TextColumn("actions_action_groups_itsm_receivers_ticket_configuration"),
		// table.TextColumn("actions_action_groups_itsm_receivers_workspace_id"),
		// table.TextColumn("actions_action_groups_logic_app_receivers"),
		// table.TextColumn("actions_action_groups_logic_app_receivers_callback_url"),
		// table.TextColumn("actions_action_groups_logic_app_receivers_name"),
		// table.TextColumn("actions_action_groups_logic_app_receivers_resource_id"),
		// table.TextColumn("actions_action_groups_logic_app_receivers_use_common_alert_schema"),
		// table.TextColumn("actions_action_groups_sms_receivers"),
		// table.TextColumn("actions_action_groups_sms_receivers_country_code"),
		// table.TextColumn("actions_action_groups_sms_receivers_name"),
		// table.TextColumn("actions_action_groups_sms_receivers_phone_number"),
		// table.TextColumn("actions_action_groups_sms_receivers_status"),
		// table.TextColumn("actions_action_groups_voice_receivers"),
		// table.TextColumn("actions_action_groups_voice_receivers_country_code"),
		// table.TextColumn("actions_action_groups_voice_receivers_name"),
		// table.TextColumn("actions_action_groups_voice_receivers_phone_number"),
		// table.TextColumn("actions_action_groups_webhook_receivers"),
		// table.TextColumn("actions_action_groups_webhook_receivers_identifier_uri"),
		// table.TextColumn("actions_action_groups_webhook_receivers_name"),
		// table.TextColumn("actions_action_groups_webhook_receivers_object_id"),
		// table.TextColumn("actions_action_groups_webhook_receivers_service_uri"),
		// table.TextColumn("actions_action_groups_webhook_receivers_tenant_id"),
		// table.TextColumn("actions_action_groups_webhook_receivers_use_aad_auth"),
		// table.TextColumn("actions_action_groups_webhook_receivers_use_common_alert_schema"),
		table.TextColumn("condition"),
		// table.TextColumn("condition_all_of"),
		// table.TextColumn("condition_all_of_any_of"),
		// table.TextColumn("condition_all_of_any_of_contains_any"),
		// table.TextColumn("condition_all_of_any_of_equals"),
		// table.TextColumn("condition_all_of_any_of_field"),
		// table.TextColumn("condition_all_of_contains_any"),
		// table.TextColumn("condition_all_of_equals"),
		// table.TextColumn("condition_all_of_field"),
		table.TextColumn("description"),
		table.TextColumn("enabled"),
		table.TextColumn("scopes"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// monitorActivityLogAlertsGenerate returns the rows in the table for all configured accounts
func MonitorActivityLogAlertsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": monitorActivityLogAlert,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountMonitorActivityLogAlerts(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": monitorActivityLogAlert,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountMonitorActivityLogAlerts(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountMonitorActivityLogAlerts(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[monitorActivityLogAlert]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": monitorActivityLogAlert,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setMonitorActivityLogAlertsToTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setMonitorActivityLogAlertsToTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	resources, err := getMonitorActivityLogAlertData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      monitorActivityLogAlert,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get monitor activityLogAlert list from api")
	}

	for _, activityLogAlert := range *resources.Response().Value {
		structs.DefaultTagName = "json"
		resMap := structs.Map(activityLogAlert)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     monitorActivityLogAlert,
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
func getMonitorActivityLogAlertData(session *azure.AzureSession, rg string) (result azuremonitor.AlertRuleListPage, err error) {

	svcClient := azuremonitor.NewActivityLogAlertsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListByResourceGroup(context.Background(), rg)

}
