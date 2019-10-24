// Copyright 2015 flannel authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// +build !windows

package vxlan

import (
	"encoding/json"
	"sync"

	log "github.com/golang/glog"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/context"

	"syscall"

	"github.com/zdnscloud/images/flannel-directrouting/backend"
	"github.com/zdnscloud/images/flannel-directrouting/pkg/ip"
	"github.com/zdnscloud/images/flannel-directrouting/subnet"
)

type network struct {
	backend.SimpleNetwork
	dev       *vxlanDevice
	subnetMgr subnet.Manager
	routes    []netlink.Route
}

const (
	encapOverhead     = 50
	routeCheckRetries = 10
)

func newNetwork(subnetMgr subnet.Manager, extIface *backend.ExternalInterface, dev *vxlanDevice, _ ip.IP4Net, lease *subnet.Lease) (*network, error) {
	nw := &network{
		SimpleNetwork: backend.SimpleNetwork{
			SubnetLease: lease,
			ExtIface:    extIface,
		},
		subnetMgr: subnetMgr,
		dev:       dev,
		routes:    make([]netlink.Route, 0),
	}

	return nw, nil
}

func (nw *network) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	log.V(0).Info("watching for new subnet leases")
	events := make(chan []subnet.Event)
	wg.Add(1)
	go func() {
		subnet.WatchLeases(ctx, nw.subnetMgr, nw.SubnetLease, events)
		log.V(1).Info("WatchLeases exited")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		nw.routeCheck(ctx)
		wg.Done()
	}()
	defer wg.Wait()

	for {
		select {
		case evtBatch := <-events:
			nw.handleSubnetEvents(evtBatch)

		case <-ctx.Done():
			return
		}
	}
}

func (nw *network) MTU() int {
	return nw.ExtIface.Iface.MTU - encapOverhead
}

type vxlanLeaseAttrs struct {
	VtepMAC hardwareAddr
}

func (nw *network) handleSubnetEvents(batch []subnet.Event) {
	for _, event := range batch {
		sn := event.Lease.Subnet
		attrs := event.Lease.Attrs
		if attrs.BackendType != "vxlan" {
			log.Warningf("ignoring non-vxlan subnet(%s): type=%v", sn, attrs.BackendType)
			continue
		}

		var vxlanAttrs vxlanLeaseAttrs
		if err := json.Unmarshal(attrs.BackendData, &vxlanAttrs); err != nil {
			log.Error("error decoding subnet lease JSON: ", err)
			continue
		}

		// This route is used when traffic should be vxlan encapsulated
		vxlanRoute := netlink.Route{
			LinkIndex: nw.dev.link.Attrs().Index,
			Scope:     netlink.SCOPE_UNIVERSE,
			Dst:       sn.ToIPNet(),
			Gw:        sn.IP.ToIP(),
		}
		vxlanRoute.SetFlag(syscall.RTNH_F_ONLINK)

		// directRouting is where the remote host is on the same subnet so vxlan isn't required.
		directRoute := netlink.Route{
			Dst: sn.ToIPNet(),
			Gw:  attrs.PublicIP.ToIP(),
		}
		var directRoutingOK = false
		if nw.dev.directRouting {
			if dr, err := ip.DirectRouting(attrs.PublicIP.ToIP()); err != nil {
				log.Error(err)
			} else {
				directRoutingOK = dr
			}
		}

		switch event.Type {
		case subnet.EventAdded:
			if directRoutingOK {
				log.V(2).Infof("Adding direct route to routes: %s PublicIP: %s", sn, attrs.PublicIP)
				nw.addToRouteList(directRoute)
			}
		case subnet.EventRemoved:
			if directRoutingOK {
				log.V(2).Infof("Removing direct route to routes: %s PublicIP: %s", sn, attrs.PublicIP)
				nw.removeFromRouteList(directRoute)
			}
		default:
			log.Error("internal error: unknown event type: ", int(event.Type))
		}
	}
}
