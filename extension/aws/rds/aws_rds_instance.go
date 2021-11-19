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

// ListCertificatesColumns returns the list of columns in the table
func ListInstanceColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.IntegerColumn("allocated_storage"),
		table.TextColumn("associated_roles"),
		// table.TextColumn("associated_roles_feature_name"),
		// table.TextColumn("associated_roles_role_arn"),
		// table.TextColumn("associated_roles_status"),
		table.TextColumn("auto_minor_version_upgrade"),
		table.TextColumn("availability_zone"),
		table.TextColumn("aws_backup_recovery_point_arn"),
		table.IntegerColumn("backup_retention_period"),
		table.TextColumn("ca_certificate_identifier"),
		table.TextColumn("character_set_name"),
		table.TextColumn("copy_tags_to_snapshot"),
		table.TextColumn("customer_owned_ip_enabled"),
		table.TextColumn("db_cluster_identifier"),
		table.TextColumn("db_instance_arn"),
		table.TextColumn("db_instance_automated_backups_replications"),
		// table.TextColumn("db_instance_automated_backups_replications_db_instance_automated_backups_arn"),
		table.TextColumn("db_instance_class"),
		table.TextColumn("db_instance_identifier"),
		table.TextColumn("db_instance_status"),
		table.TextColumn("db_name"),
		table.TextColumn("db_parameter_groups"),
		// table.TextColumn("db_parameter_groups_db_parameter_group_name"),
		// table.TextColumn("db_parameter_groups_parameter_apply_status"),
		table.TextColumn("db_security_groups"),
		// table.TextColumn("db_security_groups_db_security_group_name"),
		// table.TextColumn("db_security_groups_status"),
		table.TextColumn("db_subnet_group"),
		// table.TextColumn("db_subnet_group_db_subnet_group_arn"),
		// table.TextColumn("db_subnet_group_db_subnet_group_description"),
		// table.TextColumn("db_subnet_group_db_subnet_group_name"),
		// table.TextColumn("db_subnet_group_subnet_group_status"),
		// table.TextColumn("db_subnet_group_subnets"),
		// table.TextColumn("db_subnet_group_subnets_subnet_availability_zone"),
		// table.TextColumn("db_subnet_group_subnets_subnet_availability_zone_name"),
		// table.TextColumn("db_subnet_group_subnets_subnet_identifier"),
		// table.TextColumn("db_subnet_group_subnets_subnet_outpost"),
		// table.TextColumn("db_subnet_group_subnets_subnet_outpost_arn"),
		// table.TextColumn("db_subnet_group_subnets_subnet_status"),
		// table.TextColumn("db_subnet_group_vpc_id"),
		table.IntegerColumn("db_instance_port"),
		table.TextColumn("dbi_resource_id"),
		table.TextColumn("deletion_protection"),
		table.TextColumn("domain_memberships"),
		// table.TextColumn("domain_memberships_domain"),
		// table.TextColumn("domain_memberships_fqdn"),
		// table.TextColumn("domain_memberships_iam_role_name"),
		// table.TextColumn("domain_memberships_status"),
		table.TextColumn("enabled_cloudwatch_logs_exports"),
		table.TextColumn("endpoint"),
		// table.TextColumn("endpoint_address"),
		// table.TextColumn("endpoint_hosted_zone_id"),
		// table.IntegerColumn("endpoint_port"),
		table.TextColumn("engine"),
		table.TextColumn("engine_version"),
		table.TextColumn("enhanced_monitoring_resource_arn"),
		table.TextColumn("iam_database_authentication_enabled"),
		table.TextColumn("instance_create_time"),
		// table.BigIntColumn("instance_create_time_ext"),
		// table.TextColumn("instance_create_time_loc"),
		// table.BigIntColumn("instance_create_time_loc_cache_end"),
		// table.BigIntColumn("instance_create_time_loc_cache_start"),
		// table.TextColumn("instance_create_time_loc_cache_zone"),
		// table.TextColumn("instance_create_time_loc_cache_zone_is_dst"),
		// table.TextColumn("instance_create_time_loc_cache_zone_name"),
		// table.IntegerColumn("instance_create_time_loc_cache_zone_offset"),
		// table.TextColumn("instance_create_time_loc_extend"),
		// table.TextColumn("instance_create_time_loc_name"),
		// table.TextColumn("instance_create_time_loc_tx"),
		// table.IntegerColumn("instance_create_time_loc_tx_index"),
		// table.TextColumn("instance_create_time_loc_tx_isstd"),
		// table.TextColumn("instance_create_time_loc_tx_isutc"),
		// table.BigIntColumn("instance_create_time_loc_tx_when"),
		// table.TextColumn("instance_create_time_loc_zone"),
		// table.TextColumn("instance_create_time_loc_zone_is_dst"),
		// table.TextColumn("instance_create_time_loc_zone_name"),
		// table.IntegerColumn("instance_create_time_loc_zone_offset"),
		// table.BigIntColumn("instance_create_time_wall"),
		table.IntegerColumn("iops"),
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
		table.TextColumn("license_model"),
		table.TextColumn("listener_endpoint"),
		// table.TextColumn("listener_endpoint_address"),
		// table.TextColumn("listener_endpoint_hosted_zone_id"),
		// table.IntegerColumn("listener_endpoint_port"),
		table.TextColumn("master_username"),
		table.IntegerColumn("max_allocated_storage"),
		table.IntegerColumn("monitoring_interval"),
		table.TextColumn("monitoring_role_arn"),
		table.TextColumn("multi_az"),
		table.TextColumn("nchar_character_set_name"),
		table.TextColumn("option_group_memberships"),
		// table.TextColumn("option_group_memberships_option_group_name"),
		// table.TextColumn("option_group_memberships_status"),
		table.TextColumn("pending_modified_values"),
		// table.IntegerColumn("pending_modified_values_allocated_storage"),
		// table.IntegerColumn("pending_modified_values_backup_retention_period"),
		// table.TextColumn("pending_modified_values_ca_certificate_identifier"),
		// table.TextColumn("pending_modified_values_db_instance_class"),
		// table.TextColumn("pending_modified_values_db_instance_identifier"),
		// table.TextColumn("pending_modified_values_db_subnet_group_name"),
		// table.TextColumn("pending_modified_values_engine_version"),
		// table.TextColumn("pending_modified_values_iam_database_authentication_enabled"),
		// table.IntegerColumn("pending_modified_values_iops"),
		// table.TextColumn("pending_modified_values_license_model"),
		// table.TextColumn("pending_modified_values_master_user_password"),
		// table.TextColumn("pending_modified_values_multi_az"),
		// table.TextColumn("pending_modified_values_pending_cloudwatch_logs_exports"),
		// table.TextColumn("pending_modified_values_pending_cloudwatch_logs_exports_log_types_to_disable"),
		// table.TextColumn("pending_modified_values_pending_cloudwatch_logs_exports_log_types_to_enable"),
		// table.IntegerColumn("pending_modified_values_port"),
		// table.TextColumn("pending_modified_values_processor_features"),
		// table.TextColumn("pending_modified_values_processor_features_name"),
		// table.TextColumn("pending_modified_values_processor_features_value"),
		// table.TextColumn("pending_modified_values_storage_type"),
		table.TextColumn("performance_insights_enabled"),
		table.TextColumn("performance_insights_kms_key_id"),
		table.IntegerColumn("performance_insights_retention_period"),
		table.TextColumn("preferred_backup_window"),
		table.TextColumn("preferred_maintenance_window"),
		table.TextColumn("processor_features"),
		// table.TextColumn("processor_features_name"),
		// table.TextColumn("processor_features_value"),
		table.IntegerColumn("promotion_tier"),
		table.TextColumn("publicly_accessible"),
		table.TextColumn("read_replica_db_cluster_identifiers"),
		table.TextColumn("read_replica_db_instance_identifiers"),
		table.TextColumn("read_replica_source_db_instance_identifier"),
		table.TextColumn("replica_mode"),
		table.TextColumn("secondary_availability_zone"),
		table.TextColumn("status_infos"),
		// table.TextColumn("status_infos_message"),
		// table.TextColumn("status_infos_normal"),
		// table.TextColumn("status_infos_status"),
		// table.TextColumn("status_infos_status_type"),
		table.TextColumn("storage_encrypted"),
		table.TextColumn("storage_type"),
		table.TextColumn("tag_list"),
		// table.TextColumn("tag_list_key"),
		// table.TextColumn("tag_list_value"),
		table.TextColumn("tde_credential_arn"),
		table.TextColumn("timezone"),
		table.TextColumn("vpc_security_groups"),
		// table.TextColumn("vpc_security_groups_status"),
		// table.TextColumn("vpc_security_groups_vpc_security_group_id"),
		table.TextColumn("values"),
	}
}

func DescribeDBInstances(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAws.Accounts) == 0 && extaws.ShouldProcessAccount("aws_rds_instance", utilities.AwsAccountID) {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_rds_instance",
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountDBInstances(osqCtx, queryContext, nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAws.Accounts {
			if !extaws.ShouldProcessAccount("aws_rds_instance", account.ID) {
				continue
			}
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_instance",
				"account":   account.ID,
			}).Info("processing account")
			results, err := processAccountDBInstances(osqCtx, queryContext, &account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processRegionDescribeInstance(osqCtx context.Context, queryContext table.QueryContext, tableConfig *utilities.TableConfig, account *utilities.ExtensionConfigurationAwsAccount, region types.Region) ([]map[string]string, error) {
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
		"tableName": "aws_rds_instance",
		"account":   accountId,
		"region":    *region.RegionName,
	}).Debug("processing region")

	svc := rds.NewFromConfig(*sess)
	params := &rds.DescribeDBInstancesInput{}

	paginator := rds.NewDescribeDBInstancesPaginator(svc, params)

	for {
		page, err := paginator.NextPage(osqCtx)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_instance",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeDBInstances",
				"errString": err.Error(),
			}).Error("failed to process region")
			return resultMap, err
		}
		byteArr, err := json.Marshal(page)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": "aws_rds_instance",
				"account":   accountId,
				"region":    *region.RegionName,
				"task":      "DescribeDBInstances",
				"errString": err.Error(),
			}).Error("failed to marshal response")
			return nil, err
		}
		table := utilities.NewTable(byteArr, tableConfig)
		for _, row := range table.Rows {
			if !extaws.ShouldProcessRow(osqCtx, queryContext, "aws_rds_instance", accountId, *region.RegionName, row) {
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

func processAccountDBInstances(osqCtx context.Context, queryContext table.QueryContext, account *utilities.ExtensionConfigurationAwsAccount) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	awsSession, err := extaws.GetAwsConfig(account, "us-east-1")
	if err != nil {
		return resultMap, err
	}
	regions, err := extaws.FetchRegions(osqCtx, awsSession)
	if err != nil {
		return resultMap, err
	}
	tableConfig, ok := utilities.TableConfigurationMap["aws_rds_instance"]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": "aws_rds_instance",
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}
	for _, region := range regions {
		accountId := utilities.AwsAccountID
		if account != nil {
			accountId = account.ID
		}
		if !extaws.ShouldProcessRegion("aws_rds_instance", accountId, *region.RegionName) {
			continue
		}
		result, err := processRegionDescribeInstance(osqCtx, queryContext, tableConfig, account, region)
		if err != nil {
			continue
		}
		resultMap = append(resultMap, result...)
	}
	return resultMap, nil
}
