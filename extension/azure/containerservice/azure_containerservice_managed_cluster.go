package containerservice

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Uptycs/cloudquery/extension/azure"

	"github.com/Uptycs/basequery-go/plugin/table"
	"github.com/Uptycs/cloudquery/utilities"

	azurecontainerservice "github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-10-01/containerservice"
	"github.com/fatih/structs"
)

const managedCluster string = "azure_containerservice_managed_cluster"

type managed_Cluster = struct {
	Sku              *azurecontainerservice.ManagedClusterSKU
	ExtendedLocation *azurecontainerservice.ExtendedLocation
	Identity         *azurecontainerservice.ManagedClusterIdentity
	ID               *string
	Name             *string
	Type             *string
	Location         *string
	Tags             map[string]*string
	// Properties
	ProvisioningState       *string
	PowerState              *azurecontainerservice.PowerState
	MaxAgentPools           *int32
	KubernetesVersion       *string
	DNSPrefix               *string
	FqdnSubdomain           *string
	Fqdn                    *string
	PrivateFQDN             *string
	AzurePortalFQDN         *string
	AgentPoolProfiles       *[]azurecontainerservice.ManagedClusterAgentPoolProfile
	LinuxProfile            *azurecontainerservice.LinuxProfile
	WindowsProfile          *azurecontainerservice.ManagedClusterWindowsProfile
	ServicePrincipalProfile *azurecontainerservice.ManagedClusterServicePrincipalProfile
	AddonProfiles           map[string]*azurecontainerservice.ManagedClusterAddonProfile
	PodIdentityProfile      *azurecontainerservice.ManagedClusterPodIdentityProfile
	NodeResourceGroup       *string
	EnableRBAC              *bool
	EnablePodSecurityPolicy *bool
	NetworkProfile          *azurecontainerservice.NetworkProfile
	AadProfile              *azurecontainerservice.ManagedClusterAADProfile
	AutoUpgradeProfile      *azurecontainerservice.ManagedClusterAutoUpgradeProfile
	AutoScalerProfile       *azurecontainerservice.ManagedClusterPropertiesAutoScalerProfile
	APIServerAccessProfile  *azurecontainerservice.ManagedClusterAPIServerAccessProfile
	DiskEncryptionSetID     *string
	IdentityProfile         map[string]*azurecontainerservice.UserAssignedIdentity
	PrivateLinkResources    *[]azurecontainerservice.PrivateLinkResource
	DisableLocalAccounts    *bool
	HTTPProxyConfig         *azurecontainerservice.ManagedClusterHTTPProxyConfig
	SecurityProfile         *azurecontainerservice.ManagedClusterSecurityProfile
	PublicNetworkAccess     *azurecontainerservice.PublicNetworkAccess
}

func getManagedClusterinfo(clusters azurecontainerservice.ManagedCluster) managed_Cluster {
	managedClusterinfo := managed_Cluster{
		Sku:                     clusters.Sku,
		ExtendedLocation:        clusters.ExtendedLocation,
		Identity:                clusters.Identity,
		ID:                      clusters.ID,
		Name:                    clusters.Name,
		Type:                    clusters.Type,
		Location:                clusters.Location,
		Tags:                    clusters.Tags,
		ProvisioningState:       clusters.ProvisioningState,
		PowerState:              clusters.PowerState,
		MaxAgentPools:           clusters.MaxAgentPools,
		KubernetesVersion:       clusters.KubernetesVersion,
		DNSPrefix:               clusters.DNSPrefix,
		FqdnSubdomain:           clusters.FqdnSubdomain,
		Fqdn:                    clusters.Fqdn,
		PrivateFQDN:             clusters.PrivateFQDN,
		AzurePortalFQDN:         clusters.AzurePortalFQDN,
		AgentPoolProfiles:       clusters.AgentPoolProfiles,
		LinuxProfile:            clusters.LinuxProfile,
		WindowsProfile:          clusters.WindowsProfile,
		ServicePrincipalProfile: clusters.ServicePrincipalProfile,
		AddonProfiles:           clusters.AddonProfiles,
		PodIdentityProfile:      clusters.PodIdentityProfile,
		NodeResourceGroup:       clusters.NodeResourceGroup,
		EnableRBAC:              clusters.EnableRBAC,
		EnablePodSecurityPolicy: clusters.EnablePodSecurityPolicy,
		NetworkProfile:          clusters.NetworkProfile,
		AadProfile:              clusters.AadProfile,
		AutoUpgradeProfile:      clusters.AutoUpgradeProfile,
		AutoScalerProfile:       clusters.AutoScalerProfile,
		APIServerAccessProfile:  clusters.APIServerAccessProfile,
		DiskEncryptionSetID:     clusters.DiskEncryptionSetID,
		IdentityProfile:         clusters.IdentityProfile,
		PrivateLinkResources:    clusters.PrivateLinkResources,
		DisableLocalAccounts:    clusters.DisableLocalAccounts,
		HTTPProxyConfig:         clusters.HTTPProxyConfig,
		SecurityProfile:         clusters.SecurityProfile,
		PublicNetworkAccess:     &clusters.PublicNetworkAccess,
	}
	if clusters.PrivateLinkResources != nil {
		managedClusterinfo.PrivateLinkResources = clusters.PrivateLinkResources
	}
	return managedClusterinfo
}

// ContainerserviceManagedClustersColumns returns the list of columns in the table
func ContainerserviceManagedClustersColumns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("api_server_access_profile"),
		// table.TextColumn("api_server_access_profile_authorized_ip_ranges"),
		// table.TextColumn("api_server_access_profile_disable_run_command"),
		// table.TextColumn("api_server_access_profile_enable_private_cluster"),
		// table.TextColumn("api_server_access_profile_enable_private_cluster_public_fqdn"),
		// table.TextColumn("api_server_access_profile_private_dns_zone"),
		table.TextColumn("aad_profile"),
		// table.TextColumn("aad_profile_admin_group_object_i_ds"),
		// table.TextColumn("aad_profile_client_app_id"),
		// table.TextColumn("aad_profile_enable_azure_rbac"),
		// table.TextColumn("aad_profile_managed"),
		// table.TextColumn("aad_profile_server_app_id"),
		// table.TextColumn("aad_profile_server_app_secret"),
		// table.TextColumn("aad_profile_tenant_id"),
		table.TextColumn("auto_scaler_profile"),
		// table.TextColumn("auto_scaler_profile_balance-similar-node-groups"),
		// table.TextColumn("auto_scaler_profile_expander"),
		// table.TextColumn("auto_scaler_profile_max-empty-bulk-delete"),
		// table.TextColumn("auto_scaler_profile_max-graceful-termination-sec"),
		// table.TextColumn("auto_scaler_profile_max-node-provision-time"),
		// table.TextColumn("auto_scaler_profile_max-total-unready-percentage"),
		// table.TextColumn("auto_scaler_profile_new-pod-scale-up-delay"),
		// table.TextColumn("auto_scaler_profile_ok-total-unready-count"),
		// table.TextColumn("auto_scaler_profile_scale-down-delay-after-add"),
		// table.TextColumn("auto_scaler_profile_scale-down-delay-after-delete"),
		// table.TextColumn("auto_scaler_profile_scale-down-delay-after-failure"),
		// table.TextColumn("auto_scaler_profile_scale-down-unneeded-time"),
		// table.TextColumn("auto_scaler_profile_scale-down-unready-time"),
		// table.TextColumn("auto_scaler_profile_scale-down-utilization-threshold"),
		// table.TextColumn("auto_scaler_profile_scan-interval"),
		// table.TextColumn("auto_scaler_profile_skip-nodes-with-local-storage"),
		// table.TextColumn("auto_scaler_profile_skip-nodes-with-system-pods"),
		table.TextColumn("auto_upgrade_profile"),
		// table.TextColumn("auto_upgrade_profile_upgrade_channel"),
		table.TextColumn("azure_portal_fqdn"),
		table.TextColumn("dns_prefix"),
		// table.TextColumn("disable_local_accounts"),
		table.TextColumn("disk_encryption_set_id"),
		table.TextColumn("enable_pod_security_policy"),
		table.TextColumn("enable_rbac"),
		table.TextColumn("extended_location"),
		// table.TextColumn("extended_location_name"),
		// table.TextColumn("extended_location_type"),
		table.TextColumn("fqdn"),
		table.TextColumn("fqdn_subdomain"),
		table.TextColumn("http_proxy_config"),
		// table.TextColumn("http_proxy_config_http_proxy"),
		// table.TextColumn("http_proxy_config_https_proxy"),
		// table.TextColumn("http_proxy_config_no_proxy"),
		// table.TextColumn("http_proxy_config_trusted_ca"),
		table.TextColumn("id"),
		table.TextColumn("identity"),
		// table.TextColumn("identity_profile"),
		// table.TextColumn("identity_principal_id"),
		// table.TextColumn("identity_tenant_id"),
		// table.TextColumn("identity_type"),
		// table.TextColumn("identity_user_assigned_identities"),
		table.TextColumn("kubernetes_version"),
		table.TextColumn("agent_pool_profiles"),
		table.TextColumn("windows_profile"),
		table.TextColumn("service_principal_profile"),
		table.TextColumn("addon_profiles"),
		table.TextColumn("pod_identity_profile"),
		table.TextColumn("linux_profile"),
		// table.TextColumn("linux_profile_admin_username"),
		// table.TextColumn("linux_profile_ssh"),
		// table.TextColumn("linux_profile_ssh_public_keys"),
		// table.TextColumn("linux_profile_ssh_public_keys_key_data"),
		table.TextColumn("location"),
		table.IntegerColumn("max_agent_pools"),
		table.TextColumn("name"),
		table.TextColumn("network_profile"),
		// table.TextColumn("network_profile_dns_service_ip"),
		// table.TextColumn("network_profile_docker_bridge_cidr"),
		// table.TextColumn("network_profile_ip_families"),
		// table.TextColumn("network_profile_load_balancer_profile"),
		// table.IntegerColumn("network_profile_load_balancer_profile_allocated_outbound_ports"),
		// table.TextColumn("network_profile_load_balancer_profile_effective_outbound_i_ps"),
		// table.TextColumn("network_profile_load_balancer_profile_effective_outbound_i_ps_id"),
		// table.TextColumn("network_profile_load_balancer_profile_enable_multiple_standard_load_balancers"),
		// table.IntegerColumn("network_profile_load_balancer_profile_idle_timeout_in_minutes"),
		// table.TextColumn("network_profile_load_balancer_profile_managed_outbound_i_ps"),
		// table.IntegerColumn("network_profile_load_balancer_profile_managed_outbound_i_ps_count"),
		// table.IntegerColumn("network_profile_load_balancer_profile_managed_outbound_i_ps_count_i_pv6"),
		// table.TextColumn("network_profile_load_balancer_profile_outbound_ip_prefixes"),
		// table.TextColumn("network_profile_load_balancer_profile_outbound_ip_prefixes_public_ip_prefixes"),
		// table.TextColumn("network_profile_load_balancer_profile_outbound_ip_prefixes_public_ip_prefixes_id"),
		// table.TextColumn("network_profile_load_balancer_profile_outbound_i_ps"),
		// table.TextColumn("network_profile_load_balancer_profile_outbound_i_ps_public_i_ps"),
		// table.TextColumn("network_profile_load_balancer_profile_outbound_i_ps_public_i_ps_id"),
		// table.TextColumn("network_profile_load_balancer_sku"),
		// table.TextColumn("network_profile_nat_gateway_profile"),
		// table.TextColumn("network_profile_nat_gateway_profile_effective_outbound_i_ps"),
		// table.TextColumn("network_profile_nat_gateway_profile_effective_outbound_i_ps_id"),
		// table.IntegerColumn("network_profile_nat_gateway_profile_idle_timeout_in_minutes"),
		// table.TextColumn("network_profile_nat_gateway_profile_managed_outbound_ip_profile"),
		// table.IntegerColumn("network_profile_nat_gateway_profile_managed_outbound_ip_profile_count"),
		// table.TextColumn("network_profile_network_mode"),
		// table.TextColumn("network_profile_network_plugin"),
		// table.TextColumn("network_profile_network_policy"),
		// table.TextColumn("network_profile_outbound_type"),
		// table.TextColumn("network_profile_pod_cidr"),
		// table.TextColumn("network_profile_pod_cidrs"),
		// table.TextColumn("network_profile_service_cidr"),
		// table.TextColumn("network_profile_service_cidrs"),
		table.TextColumn("node_resource_group"),
		table.TextColumn("power_state"),
		table.TextColumn("power_state_code"),
		table.TextColumn("private_fqdn"),
		table.TextColumn("private_link_resources"),
		// table.TextColumn("private_link_resources_group_id"),
		// table.TextColumn("private_link_resources_id"),
		// table.TextColumn("private_link_resources_name"),
		// table.TextColumn("private_link_resources_private_link_service_id"),
		// table.TextColumn("private_link_resources_required_members"),
		// table.TextColumn("private_link_resources_type"),
		table.TextColumn("provisioning_state"),
		table.TextColumn("public_network_access"),
		table.TextColumn("security_profile"),
		// table.TextColumn("security_profile_azure_defender"),
		// table.TextColumn("security_profile_azure_defender_enabled"),
		// table.TextColumn("security_profile_azure_defender_log_analytics_workspace_resource_id"),
		table.TextColumn("sku"),
		// table.TextColumn("sku_name"),
		// table.TextColumn("sku_tier"),
		table.TextColumn("tags"),
		table.TextColumn("type"),
	}
}

// ContainerserviceManagedClustersGenerate returns the rows in the table for all configured accounts
func ContainerserviceManagedClustersGenerate(osqCtx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	resultMap := make([]map[string]string, 0)
	if len(utilities.ExtConfiguration.ExtConfAzure.Accounts) == 0 {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": managedCluster,
			"account":   "default",
		}).Info("processing account")
		results, err := processAccountContainerserviceManagedClusters(nil)
		if err != nil {
			return resultMap, err
		}
		resultMap = append(resultMap, results...)
	} else {
		for _, account := range utilities.ExtConfiguration.ExtConfAzure.Accounts {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName": managedCluster,
				"account":   account.SubscriptionID,
			}).Info("processing account")
			results, err := processAccountContainerserviceManagedClusters(&account)
			if err != nil {
				continue
			}
			resultMap = append(resultMap, results...)
		}
	}

	return resultMap, nil
}

func processAccountContainerserviceManagedClusters(account *utilities.ExtensionConfigurationAzureAccount) ([]map[string]string, error) {
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

	tableConfig, ok := utilities.TableConfigurationMap[managedCluster]
	if !ok {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName": managedCluster,
		}).Error("failed to get table configuration")
		return resultMap, fmt.Errorf("table configuration not found")
	}

	for _, group := range groups {
		go setContainerserviceManagedClusterstoTable(session, group, &wg, &resultMap, tableConfig)
	}
	wg.Wait()
	return resultMap, nil
}

func setContainerserviceManagedClusterstoTable(session *azure.AzureSession, rg string, wg *sync.WaitGroup, resultMap *[]map[string]string, tableConfig *utilities.TableConfig) {
	defer wg.Done()

	resources, err := getContainerserviceManagedClustersData(session, rg)
	if err != nil {
		utilities.GetLogger().WithFields(log.Fields{
			"tableName":      managedCluster,
			"rescourceGroup": rg,
			"errString":      err.Error(),
		}).Error("failed to get Managed Cluster list from api")
	}

	for _, ManagedCluster := range resources.Values() {
		resource := getManagedClusterinfo(ManagedCluster)
		structs.DefaultTagName = "json"
		resMap := structs.Map(resource)
		byteArr, err := json.Marshal(resMap)
		if err != nil {
			utilities.GetLogger().WithFields(log.Fields{
				"tableName":     managedCluster,
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
func getContainerserviceManagedClustersData(session *azure.AzureSession, rg string) (result azurecontainerservice.ManagedClusterListResultPage, err error) {

	svcClient := azurecontainerservice.NewManagedClustersClient(session.SubscriptionId)
	svcClient.Authorizer = session.Authorizer
	return svcClient.ListByResourceGroup(context.Background(), rg)

}
