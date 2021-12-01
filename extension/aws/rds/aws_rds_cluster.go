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

func ListClustersColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("activity_stream_kinesis_stream_name"),
		table.TextColumn("activity_stream_kms_key_id"),
		table.TextColumn("activity_stream_mode"),
		table.TextColumn("activity_stream_status"),
		table.BigIntColumn("allocated_storage"),
		table.TextColumn("associated_roles"),
		// table.TextColumn("associated_roles_feature_name"),
		// table.TextColumn("associated_roles_role_arn"),
		// table.TextColumn("associated_roles_status"),
		table.TextColumn("availability_zones"),
		table.BigIntColumn("backtrack_consumed_change_records"),
		table.BigIntColumn("backtrack_window"),
		table.BigIntColumn("backup_retention_period"),
		table.BigIntColumn("capacity"),
		table.TextColumn("character_set_name"),
		table.TextColumn("clone_group_id"),
		table.TextColumn("cluster_create_time"),
		// table.BigIntColumn("cluster_create_time_ext"),
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
		// table.BigIntColumn("cluster_create_time_wall"),
		table.TextColumn("copy_tags_to_snapshot"),
		table.TextColumn("cross_account_clone"),
		table.TextColumn("custom_endpoints"),
		table.TextColumn("db_cluster_arn"),
		table.TextColumn("db_cluster_identifier"),
		table.TextColumn("db_cluster_members"),
		// table.TextColumn("db_cluster_members_db_cluster_parameter_group_status"),
		// table.TextColumn("db_cluster_members_db_instance_identifier"),
		// table.TextColumn("db_cluster_members_is_cluster_writer"),
		// table.BigIntColumn("db_cluster_members_promotion_tier"),
		table.TextColumn("db_cluster_option_group_memberships"),
		// table.TextColumn("db_cluster_option_group_memberships_db_cluster_option_group_name"),
		// table.TextColumn("db_cluster_option_group_memberships_status"),
		table.TextColumn("db_cluster_parameter_group"),
		table.TextColumn("db_subnet_group"),
		table.TextColumn("database_name"),
		table.TextColumn("db_cluster_resource_id"),
		table.TextColumn("deletion_protection"),
		table.TextColumn("domain_memberships"),
		// table.TextColumn("domain_memberships_domain"),
		// table.TextColumn("domain_memberships_fqdn"),
		// table.TextColumn("domain_memberships_iam_role_name"),
		// table.TextColumn("domain_memberships_status"),
		table.TextColumn("earliest_backtrack_time"),
		// table.BigIntColumn("earliest_backtrack_time_ext"),
		// table.TextColumn("earliest_backtrack_time_loc"),
		// table.BigIntColumn("earliest_backtrack_time_loc_cache_end"),
		// table.BigIntColumn("earliest_backtrack_time_loc_cache_start"),
		// table.TextColumn("earliest_backtrack_time_loc_cache_zone"),
		// table.TextColumn("earliest_backtrack_time_loc_cache_zone_is_dst"),
		// table.TextColumn("earliest_backtrack_time_loc_cache_zone_name"),
		// table.IntegerColumn("earliest_backtrack_time_loc_cache_zone_offset"),
		// table.TextColumn("earliest_backtrack_time_loc_extend"),
		// table.TextColumn("earliest_backtrack_time_loc_name"),
		// table.TextColumn("earliest_backtrack_time_loc_tx"),
		// table.IntegerColumn("earliest_backtrack_time_loc_tx_index"),
		// table.TextColumn("earliest_backtrack_time_loc_tx_isstd"),
		// table.TextColumn("earliest_backtrack_time_loc_tx_isutc"),
		// table.BigIntColumn("earliest_backtrack_time_loc_tx_when"),
		// table.TextColumn("earliest_backtrack_time_loc_zone"),
		// table.TextColumn("earliest_backtrack_time_loc_zone_is_dst"),
		// table.TextColumn("earliest_backtrack_time_loc_zone_name"),
		// table.IntegerColumn("earliest_backtrack_time_loc_zone_offset"),
		// table.BigIntColumn("earliest_backtrack_time_wall"),
		table.TextColumn("earliest_restorable_time"),
		// table.BigIntColumn("earliest_restorable_time_ext"),
		// table.TextColumn("earliest_restorable_time_loc"),
		// table.BigIntColumn("earliest_restorable_time_loc_cache_end"),
		// table.BigIntColumn("earliest_restorable_time_loc_cache_start"),
		// table.TextColumn("earliest_restorable_time_loc_cache_zone"),
		// table.TextColumn("earliest_restorable_time_loc_cache_zone_is_dst"),
		// table.TextColumn("earliest_restorable_time_loc_cache_zone_name"),
		// table.IntegerColumn("earliest_restorable_time_loc_cache_zone_offset"),
		// table.TextColumn("earliest_restorable_time_loc_extend"),
		// table.TextColumn("earliest_restorable_time_loc_name"),
		// table.TextColumn("earliest_restorable_time_loc_tx"),
		// table.IntegerColumn("earliest_restorable_time_loc_tx_index"),
		// table.TextColumn("earliest_restorable_time_loc_tx_isstd"),
		// table.TextColumn("earliest_restorable_time_loc_tx_isutc"),
		// table.BigIntColumn("earliest_restorable_time_loc_tx_when"),
		// table.TextColumn("earliest_restorable_time_loc_zone"),
		// table.TextColumn("earliest_restorable_time_loc_zone_is_dst"),
		// table.TextColumn("earliest_restorable_time_loc_zone_name"),
		// table.IntegerColumn("earliest_restorable_time_loc_zone_offset"),
		// table.BigIntColumn("earliest_restorable_time_wall"),
		table.TextColumn("enabled_cloudwatch_logs_exports"),
		table.TextColumn("endpoint"),
		table.TextColumn("engine"),
		table.TextColumn("engine_mode"),
		table.TextColumn("engine_version"),
		table.TextColumn("global_write_forwarding_requested"),
		table.TextColumn("global_write_forwarding_status"),
		table.TextColumn("hosted_zone_id"),
		table.TextColumn("http_endpoint_enabled"),
		table.TextColumn("iam_database_authentication_enabled"),
		table.TextColumn("kms_key_id"),
		table.TextColumn("latest_restorable_time"),
		// table.BigIntColumn("latest_restorable_time_ext"),
		// table.TextColumn("latest_restorable_time_loc"),
		// table.BigIntColumn("latest_restorable_time_loc_cache_end"),
		// table.BigIntColumn("latest_restorable_time_loc_cache_start"),
		// table.TextColumn("latest_restorable_time_loc_cache_zone"),
		// table.TextColumn("latest_restorable_time_loc_cache_zone_is_dst"),
		// table.TextColumn("latest_restorable_time_loc_cache_zone_name"),
		// table.IntegerColumn("latest_restorable_time_loc_cache_zone_offset"),
		// table.TextColumn("latest_restorable_time_loc_extend"),
		// table.TextColumn("latest_restorable_time_loc_name"),
		// table.TextColumn("latest_restorable_time_loc_tx"),
		// table.IntegerColumn("latest_restorable_time_loc_tx_index"),
		// table.TextColumn("latest_restorable_time_loc_tx_isstd"),
		// table.TextColumn("latest_restorable_time_loc_tx_isutc"),
		// table.BigIntColumn("latest_restorable_time_loc_tx_when"),
		// table.TextColumn("latest_restorable_time_loc_zone"),
		// table.TextColumn("latest_restorable_time_loc_zone_is_dst"),
		// table.TextColumn("latest_restorable_time_loc_zone_name"),
		// table.IntegerColumn("latest_restorable_time_loc_zone_offset"),
		// table.BigIntColumn("latest_restorable_time_wall"),
		table.TextColumn("master_username"),
		table.TextColumn("multi_az"),
		table.TextColumn("pending_modified_values"),
		// table.TextColumn("pending_modified_values_db_cluster_identifier"),
		// table.TextColumn("pending_modified_values_engine_version"),
		// table.TextColumn("pending_modified_values_iam_database_authentication_enabled"),
		// table.TextColumn("pending_modified_values_master_user_password"),
		// table.TextColumn("pending_modified_values_pending_cloudwatch_logs_exports"),
		// table.TextColumn("pending_modified_values_pending_cloudwatch_logs_exports_log_types_to_disable"),
		// table.TextColumn("pending_modified_values_pending_cloudwatch_logs_exports_log_types_to_enable"),
		table.TextColumn("percent_progress"),
		table.BigIntColumn("port"),
		table.TextColumn("preferred_backup_window"),
		table.TextColumn("preferred_maintenance_window"),
		table.TextColumn("read_replica_identifiers"),
		table.TextColumn("reader_endpoint"),
		table.TextColumn("replication_source_identifier"),
		table.TextColumn("scaling_configuration_info"),
		// table.TextColumn("scaling_configuration_info_auto_pause"),
		// table.BigIntColumn("scaling_configuration_info_max_capacity"),
		// table.BigIntColumn("scaling_configuration_info_min_capacity"),
		// table.BigIntColumn("scaling_configuration_info_seconds_until_auto_pause"),
		// table.TextColumn("scaling_configuration_info_timeout_action"),
		table.TextColumn("status"),
		table.TextColumn("storage_encrypted"),
		table.TextColumn("tag_list"),
		// table.TextColumn("tag_list_key"),
		// table.TextColumn("tag_list_value"),
		table.TextColumn("vpc_security_groups"),
		// table.TextColumn("vpc_security_groups_status"),
		// table.TextColumn("vpc_security_groups_vpc_security_group_id"),
	}
}

func DescribeClustersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_rds_cluster", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_rds_cluster",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDescribeClusters(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_rds_cluster", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_cluster",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDescribeClusters(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeClusters(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_rds_cluster",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := rds.NewFromConfig(*sess)
	params := &rds.DescribeDBClustersInput{}

	paginator := rds.NewDescribeDBClustersPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_cluster",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeDBClusters",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_cluster",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeDBClusters",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_rds_cluster", accountId, *region.RegionName, row) {
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

func processAccountDescribeClusters(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_rds_cluster"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_rds_cluster",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_rds_cluster", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeClusters(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
