package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2020-12-01/redis"
	"github.com/fatih/structs"
)

const azureRedisCache string = "azure_redis_cache"

// RedisCacheColumns returns the list of columns in the table
func RedisCacheColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("etag"),
		table.TextColumn("id"),
		table.TextColumn("name"),
		table.TextColumn("type"),
		// table.TextColumn("properties"),
		table.TextColumn("access_keys"),
		// table.TextColumn("access_keys_primary_key"),
		// table.TextColumn("access_keys_secondary_key"),
		table.TextColumn("enable_non_ssl_port"),
		table.TextColumn("host_name"),
		table.TextColumn("instances"),
		// table.TextColumn("instances_is_master"),
		// table.TextColumn("instances_is_primary"),
		// table.IntegerColumn("instances_non_ssl_port"),
		// table.IntegerColumn("instances_shard_id"),
		// table.IntegerColumn("instances_ssl_port"),
		// table.TextColumn("instances_zone"),
		table.TextColumn("linked_servers"),
		// table.TextColumn("linked_servers_id"),
		table.TextColumn("minimum_tls_version"),
		table.IntegerColumn("port"),
		table.TextColumn("private_endpoint_connections"),
		// table.TextColumn("private_endpoint_connections_id"),
		// table.TextColumn("private_endpoint_connections_name"),
		// table.TextColumn("private_endpoint_connections_type"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("public_network_access"),
		table.TextColumn("redis_configuration"),
		table.TextColumn("redis_version"),
		table.IntegerColumn("replicas_per_master"),
		table.IntegerColumn("replicas_per_primary"),
		table.IntegerColumn("shard_count"),
		table.TextColumn("sku"),
		// table.IntegerColumn("sku_capacity"),
		// table.TextColumn("sku_family"),
		// table.TextColumn("sku_name"),
		table.IntegerColumn("ssl_port"),
		table.TextColumn("static_ip"),
		table.TextColumn("subnet_id"),
		table.TextColumn("tenant_settings"),
	}
}

// RedisCacheGenerate returns the rows in the table for all configured accounts
func RedisCacheGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureRedisCache,
			"account":   "default",
		}).Info("processing account")
		results, err := processRedisCache(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": azureRedisCache,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processRedisCache(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRedisCache(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[azureRedisCache]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": azureRedisCache,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setRedisCachetoTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setRedisCachetoTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	for resourcesItr, err := getRedisCacheData(session, rg); resourcesItr.NotDone(); resourcesItr.Next() {
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":      azureRedisCache,
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
				"tableName":     azureRedisCache,
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
func getRedisCacheData(session *azure.AzureSession, rg string) (result redis.ListResultIterator, err error) {
	svcClient := redis.NewClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListByResourceGroupComplete(context.Background(), rg)
}
