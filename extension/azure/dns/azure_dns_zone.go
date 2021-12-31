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

const azureDnsZone string = "azure_dns_zone"

// DnsZoneColumns returns the list of columns in the table
func DnsZoneColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("location"),
		table.TextColumn("name"),
		// table.TextColumn("properties"),
		table.BigIntColumn("max_number_of_record_sets"),
		table.BigIntColumn("max_number_of_records_per_record_set"),
		table.TextColumn("name_servers"),
		table.BigIntColumn("number_of_record_sets"),
		table.TextColumn("registration_virtual_networks"),
		// table.TextColumn("registration_virtual_networks_id"),
		table.TextColumn("resolution_virtual_networks"),
		// table.TextColumn("resolution_virtual_networks_id"),
		table.TextColumn("zone_type"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// DnsZoneGenerate returns the rows in the table for all configured accounts
func DnsZoneGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureDnsZone,
			"account":   "default",
		}).Info("processing account")
		results, err := processDnsZone(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureDnsZone,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processDnsZone(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processDnsZone(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureDnsZone]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureDnsZone,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setDnsZonetoTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setDnsZonetoTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	for resourcesItr, err := getDnsZoneData(session, rg); resourcesItr.NotDone(); resourcesItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":      azureDnsZone,
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
				"tableName":     azureDnsZone,
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

func getDnsZoneData(session *azure.AzureSession, rg string) (result dns.ZoneListResultIterator, err error) {
	svcClient := dns.NewZonesClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListByResourceGroupComplete(context.Background(), rg, nil)
}
