package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/fatih/structs"
)

const azureDnsRecordSet string = "azure_dns_record_set"

// DnsRecordSetColumns returns the list of columns in the table
func DnsRecordSetColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.TextColumn("aaaa_records"),
		// table.TextColumn("aaaa_records_ipv6_address"),
		table.TextColumn("a_records"),
		// table.TextColumn("a_records_ipv4_address"),
		table.TextColumn("cname_record"),
		// table.TextColumn("cname_record_cname"),
		table.TextColumn("mx_records"),
		// table.TextColumn("mx_records_exchange"),
		// table.IntegerColumn("mx_records_preference"),
		table.TextColumn("ns_records"),
		// table.TextColumn("ns_records_nsdname"),
		table.TextColumn("ptr_records"),
		// table.TextColumn("ptr_records_ptrdname"),
		table.TextColumn("soa_record"),
		// table.TextColumn("soa_record_email"),
		// table.BigIntColumn("soa_record_expire_time"),
		// table.TextColumn("soa_record_host"),
		// table.BigIntColumn("soa_record_minimum_ttl"),
		// table.BigIntColumn("soa_record_refresh_time"),
		// table.BigIntColumn("soa_record_retry_time"),
		// table.BigIntColumn("soa_record_serial_number"),
		table.TextColumn("srv_records"),
		// table.IntegerColumn("srv_records_port"),
		// table.IntegerColumn("srv_records_priority"),
		// table.TextColumn("srv_records_target"),
		// table.IntegerColumn("srv_records_weight"),
		table.BigIntColumn("ttl"),
		table.TextColumn("txt_records"),
		// table.TextColumn("txt_records_value"),
		table.TextColumn("caa_records"),
		// table.IntegerColumn("caa_records_flags"),
		// table.TextColumn("caa_records_tag"),
		// table.TextColumn("caa_records_value"),
		table.TextColumn("fqdn"),
		table.TextColumn("metadata"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("target_resource"),
		// table.TextColumn("target_resource_id"),
		table.TextColumn("type"),
	}
}

// DnsRecordSetGenerate returns the rows in the table for all configured accounts
func DnsRecordSetGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureDnsRecordSet,
			"account":   "default",
		}).Info("processing account")
		results, err := processDnsRecordSet(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureDnsRecordSet,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processDnsRecordSet(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processDnsRecordSet(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureDnsRecordSet]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureDnsRecordSet,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go collectDnsZonetoTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func collectDnsZonetoTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	for resourcesItr, err := getDnsZoneData(session, rg); resourcesItr.NotDone(); resourcesItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":      azureDnsRecordSet,
				"rescourceGroup": rg,
				"errString":      err.Error(),
			}).Error("failed to get DNS zones")
		}

		resource := resourcesItr.Value()

		setDnsRecordSettoTable(session, rg, *resource.Name, resultMap, tableConfig)
	}
}

func setDnsRecordSettoTable(session *azure.AzureSession, rg string, zone string, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {

	for resourcesItr, err := getDnsRecordSetData(session, rg, zone); resourcesItr.NotDone(); resourcesItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":      azureDnsRecordSet,
				"rescourceGroup": rg,
				"errString":      err.Error(),
			}).Error("failed to get DNS zones")
		}

		resource := resourcesItr.Value()

		structs.DefaultTagName = "json"
		resMap := structs.Map(resource)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     azureDnsRecordSet,
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
func getDnsRecordSetData(session *azure.AzureSession, rg string, zone string) (result dns.RecordSetListResultIterator, err error) {
	svcClient := dns.NewRecordSetsClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListAllByDNSZoneComplete(context.Background(), rg, zone, nil, "")
}
