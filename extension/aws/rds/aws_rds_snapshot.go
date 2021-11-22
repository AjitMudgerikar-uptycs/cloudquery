/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package rds

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/utilities"

	"github.com/Uptycs/basequery-go/plugin/table"
	extaws "github.com/Uptycs/cloudquery/extension/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// ListSnapshotsColumns returns the list of columns in the table
func ListSnapshotsColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("DBClusterIdentifier"),
		table.TextColumn("DBClusterSnapshotIdentifier"),
		table.TextColumn("PercentProgress"),
		table.TextColumn("StorageEncryption"),
		table.TextColumn("Engine"),
		table.TextColumn("SnapshotCreateTime"),
		table.TextColumn("VpcId"),
		table.TextColumn("ClusterCreateTime"),
		table.TextColumn("SnapshotType"),
		table.TextColumn("Status"),
		table.IntegerColumn("allocated_storage"),
		table.TextColumn("availability_zones"),
		//table.BigIntColumn("cluster_create_time_ext"),
		// table.TextColumn("cluster_create_time_loc"),
		// table.BigIntColumn("cluster_create_time_loc_cache_end"),
		// table.BigIntColumn("cluster_create_time_loc_cache_start"),
		// table.TextColumn("cluster_create_time_loc_cache_zone"),
		// table.TextColumn("cluster_create_time_loc_cache_zone_is_dst"),
		// table.TextColumn("cluster_create_time_loc_cache_zone_name"),
		// table.IntegerColumn("cluster_create_time_loc_cache_zone_offset"),
		// table.TextColumn("cluster_create_time_loc_extend"),
		// table.TextColumn("cluster_create_time_loc_name"),
		// table.TextColumn("cluster_create_time_loc_tx"),
		// table.IntegerColumn("cluster_create_time_loc_tx_index"),
		// table.TextColumn("cluster_create_time_loc_tx_isstd"),
		// table.TextColumn("cluster_create_time_loc_tx_isutc"),
		// table.BigIntColumn("cluster_create_time_loc_tx_when"),
		// table.TextColumn("cluster_create_time_loc_zone"),
		// table.TextColumn("cluster_create_time_loc_zone_is_dst"),
		// table.TextColumn("cluster_create_time_loc_zone_name"),
		// table.IntegerColumn("cluster_create_time_loc_zone_offset"),
		table.BigIntColumn("cluster_create_time_wall"),
		table.TextColumn("db_cluster_snapshot_arn"),
		table.TextColumn("engine_mode"),
		table.TextColumn("engine_version"),
		table.TextColumn("iam_database_authentication_enabled"),
		table.TextColumn("kms_key_id"),
		table.TextColumn("license_model"),
		table.TextColumn("master_username"),
		table.IntegerColumn("port"),
		// table.BigIntColumn("snapshot_create_time_ext"),
		// table.TextColumn("snapshot_create_time_loc"),
		// table.BigIntColumn("snapshot_create_time_loc_cache_end"),
		// table.BigIntColumn("snapshot_create_time_loc_cache_start"),
		// table.TextColumn("snapshot_create_time_loc_cache_zone"),
		// table.TextColumn("snapshot_create_time_loc_cache_zone_is_dst"),
		// table.TextColumn("snapshot_create_time_loc_cache_zone_name"),
		// table.IntegerColumn("snapshot_create_time_loc_cache_zone_offset"),
		// table.TextColumn("snapshot_create_time_loc_extend"),
		// table.TextColumn("snapshot_create_time_loc_name"),
		// table.TextColumn("snapshot_create_time_loc_tx"),
		// table.IntegerColumn("snapshot_create_time_loc_tx_index"),
		// table.TextColumn("snapshot_create_time_loc_tx_isstd"),
		// table.TextColumn("snapshot_create_time_loc_tx_isutc"),
		// table.BigIntColumn("snapshot_create_time_loc_tx_when"),
		// table.TextColumn("snapshot_create_time_loc_zone"),
		// table.TextColumn("snapshot_create_time_loc_zone_is_dst"),
		// table.TextColumn("snapshot_create_time_loc_zone_name"),
		// table.IntegerColumn("snapshot_create_time_loc_zone_offset"),
		// table.BigIntColumn("snapshot_create_time_wall"),
		table.TextColumn("source_db_cluster_snapshot_arn"),
		table.TextColumn("tag_list"),
		// table.TextColumn("tag_list_key"),
		// table.TextColumn("tag_list_value"),
		table.TextColumn("values"),
	}
}

// DescribeSnapshotsGenerate returns the rows in the table for all configured accounts
func DescribeSnapshotsGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_rds_snapshot", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_rds_snapshot",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeSnapshots(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_rds_snapshot", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_snapshot",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeSnapshots(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeSnapshots(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	sess, err := extaws.GetAwsConfig(account, *region.RegionName)
	if err != nil {
		return resultMap, err
	}

	accountId := utilities.AwsAccountID
	if account != nil {
		accountId = account.ID
	}

	utilities.GetLogger().WithFields(log.Fields{
		"tableName": "aws_rds_snapshot",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := rds.NewFromConfig(*sess)
	params := &rds.DescribeDBClusterSnapshotsInput{}

	paginator := rds.NewDescribeDBClusterSnapshotsPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_snapshot",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeDBClusterSnapshots",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_snapshot",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeDBClusterSnapshots",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_rds_snapshot", accountId, *region.RegionName, row) {
				continue
			}
			result := extaws.RowToMap(row, accountId, *region.RegionName, tableConfig)
			resultMap = append(resultMap, result)
		}
		if !paginator.HasMorePages() {
			break
		}
	}
	return resultMap, nil
}

func processAccountDescribeSnapshots(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_rds_snapshot"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_rds_snapshot",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_rds_snapshot", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeSnapshots(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
