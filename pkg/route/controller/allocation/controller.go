package allocation

import (
	"github.com/golang/glog"

	"github.com/projectatomic/appinfra-next/pkg/route"
	routeapi "github.com/projectatomic/appinfra-next/pkg/route/api"
)

// RouteAllocationController abstracts the details of how routes are
// allocated to router shards.
type RouteAllocationController struct {
	Plugin route.AllocationPlugin
}

// AllocateRouterShard allocates a router shard for the given route.
func (c *RouteAllocationController) AllocateRouterShard(route *routeapi.Route) (*routeapi.RouterShard, error) {

	glog.V(4).Infof("Allocating router shard for Route: %s [alias=%s]",
		route.ServiceName, route.Host)

	shard, err := c.Plugin.Allocate(route)

	if err != nil {
		glog.Errorf("unable to allocate router shard: %v", err)
		return shard, err
	}

	glog.V(4).Infof("Route %s allocated to shard %s [suffix=%s]",
		route.ServiceName, shard.ShardName, shard.DNSSuffix)

	return shard, err
}

// GenerateHostname generates a host name for the given route and router shard combination.
func (c *RouteAllocationController) GenerateHostname(route *routeapi.Route, shard *routeapi.RouterShard) string {
	glog.V(4).Infof("Generating host name for Route: %s",
		route.ServiceName)

	s := c.Plugin.GenerateHostname(route, shard)

	glog.V(4).Infof("Route: %s, generated host name/alias=%s",
		route.ServiceName, s)

	return s
}
