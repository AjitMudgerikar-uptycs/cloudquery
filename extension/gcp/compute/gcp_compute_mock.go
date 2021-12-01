/**
 * Copyright (c) 2020-present, The cloudquery authors
 *
 * This source code is licensed as defined by the LICENSE file found in the
 * root directory of this source tree.
 *
 * SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)
 */

package compute

import (
	"context"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// GcpComputeMock routes gcp_compute_* table's API invocation to a mock implementation to enable unit tests
type GcpComputeMock struct {
	svc compute.Service

	disksSvc         compute.DisksService
	disksAggList     compute.DisksAggregatedListCall
	instancesSvc     compute.InstancesService
	instancesAggList compute.InstancesAggregatedListCall
	networksSvc      compute.NetworksService
	networksList     compute.NetworksListCall

	instancesPage compute.InstanceAggregatedList
	disksPage     compute.DiskAggregatedList
	networksPage  compute.NetworkList

	itemsKey string
}

// NewGcpComputeMock returns new instance of GcpComputeMock
func NewGcpComputeMock() *GcpComputeMock {
	var mock = GcpComputeMock{}

	mock.itemsKey = "test"

	diskItems := make(map[string]compute.DisksScopedList)
	diskItems[mock.itemsKey] = compute.DisksScopedList{Disks: make([]*compute.Disk, 0)}
	mock.disksPage = compute.DiskAggregatedList{Items: diskItems}

	instanceItems := make(map[string]compute.InstancesScopedList)
	instanceItems[mock.itemsKey] = compute.InstancesScopedList{Instances: make([]*compute.Instance, 0)}
	mock.instancesPage = compute.InstanceAggregatedList{Items: instanceItems}

	networkItems := make([]*compute.Network, 0)
	mock.networksPage = compute.NetworkList{Items: networkItems}

	return &mock
}

// NewService returns compute.Service or error
func (gcp *GcpComputeMock) NewService(ctx context.Context, opts ...option.ClientOption) (*compute.Service, error) {
	return &gcp.svc, nil
}

// NewDisksService returns *compute.DisksService
func (gcp *GcpComputeMock) NewDisksService(svc *compute.Service) *compute.DisksService {
	return &gcp.disksSvc
}

// DisksAggregatedList returns *compute.DisksAggregatedListCall
func (gcp *GcpComputeMock) DisksAggregatedList(apiSvc *compute.DisksService, projectID string) *compute.DisksAggregatedListCall {
	return &gcp.disksAggList
}

// DisksPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) DisksPages(ctx context.Context, listCall *compute.DisksAggregatedListCall, cb callbackDisksPages) error {
	cb(&gcp.disksPage)
	return nil
}

func (gcp *GcpComputeMock) addDisks(inList []*compute.Disk) {
	disks := gcp.disksPage.Items[gcp.itemsKey]
	disks.Disks = append(disks.Disks, inList...)
	gcp.disksPage.Items[gcp.itemsKey] = disks
}

func (gcp *GcpComputeMock) clearDisks() {
	disks := gcp.disksPage.Items[gcp.itemsKey]
	disks.Disks = make([]*compute.Disk, 0)
	gcp.disksPage.Items[gcp.itemsKey] = disks
}

// NewInstancesService returns *compute.InstancesService
func (gcp *GcpComputeMock) NewInstancesService(svc *compute.Service) *compute.InstancesService {
	return &gcp.instancesSvc
}

// InstancesAggregatedList returns *compute.InstancesAggregatedListCall
func (gcp *GcpComputeMock) InstancesAggregatedList(apiSvc *compute.InstancesService, projectID string) *compute.InstancesAggregatedListCall {
	return &gcp.instancesAggList
}

// InstancesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) InstancesPages(ctx context.Context, listCall *compute.InstancesAggregatedListCall, cb callbackInstancesPages) error {
	cb(&gcp.instancesPage)
	return nil
}

func (gcp *GcpComputeMock) addInstances(inList []*compute.Instance) {
	instances := gcp.instancesPage.Items[gcp.itemsKey]
	instances.Instances = append(instances.Instances, inList...)
	gcp.instancesPage.Items[gcp.itemsKey] = instances
}

func (gcp *GcpComputeMock) clearInstances() {
	instances := gcp.instancesPage.Items[gcp.itemsKey]
	instances.Instances = make([]*compute.Instance, 0)
	gcp.instancesPage.Items[gcp.itemsKey] = instances
}

// NewNetworksService returns *compute.NetworksService
func (gcp *GcpComputeMock) NewNetworksService(svc *compute.Service) *compute.NetworksService {
	return &gcp.networksSvc
}

// NetworksList returns *compute.NetworksListCall
func (gcp *GcpComputeMock) NetworksList(apiSvc *compute.NetworksService, projectID string) *compute.NetworksListCall {
	return &gcp.networksList
}

// NetworksPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) NetworksPages(ctx context.Context, listCall *compute.NetworksListCall, cb callbackNetworksPages) error {
	cb(&gcp.networksPage)
	return nil
}

func (gcp *GcpComputeMock) addNetworks(inList []*compute.Network) {
	gcp.networksPage.Items = inList
}

func (gcp *GcpComputeMock) clearNetworks() {
	gcp.networksPage.Items = make([]*compute.Network, 0)
}

// NewImagesService returns *compute.ImagesService
func (gcp *GcpComputeMock) NewImagesService(svc *compute.Service) *compute.ImagesService {
	return nil
}

// ImagesList returns *compute.ImagesListCall
func (gcp *GcpComputeMock) ImagesList(apiSvc *compute.ImagesService, projectID string) *compute.ImagesListCall {
	return nil
}

// ImagesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) ImagesPages(ctx context.Context, listCall *compute.ImagesListCall, cb callbackImagesPages) error {
	return nil
}

// NewInterconnectsService returns *compute.InterconnectsService
func (gcp *GcpComputeMock) NewInterconnectsService(svc *compute.Service) *compute.InterconnectsService {
	return nil
}

// InterconnectsList returns *compute.InterconnectsListCall
func (gcp *GcpComputeMock) InterconnectsList(apiSvc *compute.InterconnectsService, projectID string) *compute.InterconnectsListCall {
	return nil
}

// InterconnectsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) InterconnectsPages(ctx context.Context, listCall *compute.InterconnectsListCall, cb callbackInterconnectsPages) error {
	return nil
}

// NewRoutesService returns *compute.RoutesService
func (gcp *GcpComputeMock) NewRoutesService(svc *compute.Service) *compute.RoutesService {
	return nil
}

// RoutesList returns *compute.RoutesListCall
func (gcp *GcpComputeMock) RoutesList(apiSvc *compute.RoutesService, projectID string) *compute.RoutesListCall {
	return nil
}

// RoutesPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) RoutesPages(ctx context.Context, listCall *compute.RoutesListCall, cb callbackRoutesPages) error {
	return nil
}

// NewReservationsService returns *compute.ReservationsService
func (gcp *GcpComputeMock) NewReservationsService(svc *compute.Service) *compute.ReservationsService {
	return nil
}

// ReservationsAggregatedList returns *compute.ReservationsAggregatedListCall
func (gcp *GcpComputeMock) ReservationsAggregatedList(apiSvc *compute.ReservationsService, projectID string) *compute.ReservationsAggregatedListCall {
	return nil
}

// ReservationsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) ReservationsPages(ctx context.Context, listCall *compute.ReservationsAggregatedListCall, cb callbackReservationsPages) error {
	return nil
}

// NewRoutersService returns *compute.RoutersService
func (gcp *GcpComputeMock) NewRoutersService(svc *compute.Service) *compute.RoutersService {
	return nil
}

// RoutersAggregatedList returns *compute.RoutersAggregatedListCall
func (gcp *GcpComputeMock) RoutersAggregatedList(apiSvc *compute.RoutersService, projectID string) *compute.RoutersAggregatedListCall {
	return nil
}

// RoutersPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) RoutersPages(ctx context.Context, listCall *compute.RoutersAggregatedListCall, cb callbackRoutersPages) error {
	return nil
}

// NewVpnTunnelsService returns *compute.VpnTunnelsService
func (gcp *GcpComputeMock) NewVpnTunnelsService(svc *compute.Service) *compute.VpnTunnelsService {
	return nil
}

// VpnTunnelsAggregatedList returns *compute.VpnTunnelsAggregatedListCall
func (gcp *GcpComputeMock) VpnTunnelsAggregatedList(apiSvc *compute.VpnTunnelsService, projectID string) *compute.VpnTunnelsAggregatedListCall {
	return nil
}

// VpnTunnelsPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) VpnTunnelsPages(ctx context.Context, listCall *compute.VpnTunnelsAggregatedListCall, cb callbackVpnTunnelsPages) error {
	return nil
}

// NewVpnGatewaysService returns *compute.VpnGatewaysService
func (gcp *GcpComputeMock) NewVpnGatewaysService(svc *compute.Service) *compute.VpnGatewaysService {
	return nil
}

// VpnGatewaysAggregatedList returns *compute.VpnGatewaysAggregatedListCall
func (gcp *GcpComputeMock) VpnGatewaysAggregatedList(apiSvc *compute.VpnGatewaysService, projectID string) *compute.VpnGatewaysAggregatedListCall {
	return nil
}

// VpnGatewaysPages invokes cb for each page of results.
// Returns error on failure
func (gcp *GcpComputeMock) VpnGatewaysPages(ctx context.Context, listCall *compute.VpnGatewaysAggregatedListCall, cb callbackVpnGatewaysPages) error {
	return nil
}
