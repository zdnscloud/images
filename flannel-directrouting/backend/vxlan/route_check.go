package vxlan

import (
	"bytes"
	log "github.com/golang/glog"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/context"
	"net"
	"time"
)

func (nw *network) routeCheck(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(routeCheckRetries * time.Second):
			nw.checkSubnetExistInRoutes()
		}
	}
}

func (nw *network) checkSubnetExistInRoutes() {
	routeList, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err == nil {
		for _, route := range nw.routes {
			exist := false
			for _, r := range routeList {
				if r.Dst == nil {
					continue
				}
				if routeEqual(r, route) {
					exist = true
					break
				}
			}

			if !exist {
				if err := netlink.RouteAdd(&route); err != nil {
					if nerr, ok := err.(net.Error); !ok {
						log.Errorf("Error recovering route to %v: %v, %v", route.Dst, route.Gw, nerr)
					}
					continue
				} else {
					log.Infof("Route recovered %v : %v", route.Dst, route.Gw)
				}
			}
		}
	} else {
		log.Errorf("Error fetching route list. Will automatically retry: %v", err)
	}
}
func routeEqual(x, y netlink.Route) bool {
	//if x.Dst.IP.Equal(y.Dst.IP) && x.Gw.Equal(y.Gw) && bytes.Equal(x.Dst.Mask, y.Dst.Mask) && x.LinkIndex == y.LinkIndex {
	if x.Dst.IP.Equal(y.Dst.IP) && x.Gw.Equal(y.Gw) && bytes.Equal(x.Dst.Mask, y.Dst.Mask) {
		return true
	}
	return false
}

func (nw *network) addToRouteList(route netlink.Route) {
	for _, r := range nw.routes {
		if routeEqual(r, route) {
			return
		}
	}
	nw.routes = append(nw.routes, route)
}

func (nw *network) removeFromRouteList(route netlink.Route) {
	for index, r := range nw.routes {
		if routeEqual(r, route) {
			nw.routes = append(nw.routes[:index], nw.routes[index+1:]...)
			return
		}
	}
}
