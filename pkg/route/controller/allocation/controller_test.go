package allocation

import (
	"fmt"
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	routeapi "github.com/projectatomic/appinfra-next/pkg/route/api"
)

type TestAllocationPlugin struct {
	Name string
}

func (p *TestAllocationPlugin) Allocate(route *routeapi.Route) (*routeapi.RouterShard, error) {

	return &routeapi.RouterShard{ShardName: "test", DNSSuffix: "openshift.test"}, nil
}

func (p *TestAllocationPlugin) GenerateHostname(route *routeapi.Route, shard *routeapi.RouterShard) string {
	if len(route.ServiceName) > 0 && len(route.Namespace) > 0 {
		return fmt.Sprintf("%s-%s.%s", route.ServiceName, route.Namespace, shard.DNSSuffix)
	}

	return "test-test-test.openshift.test"
}

func TestRouteAllocationController(t *testing.T) {
	tests := []struct {
		name  string
		route *routeapi.Route
	}{
		{
			name: "No Name",
			route: &routeapi.Route{
				ObjectMeta: kapi.ObjectMeta{
					Namespace: "namespace",
				},
				ServiceName: "service",
			},
		},
		{
			name: "No namespace",
			route: &routeapi.Route{
				ObjectMeta: kapi.ObjectMeta{
					Name: "name",
				},
				ServiceName: "nonamespace",
			},
		},
		{
			name: "No service name",
			route: &routeapi.Route{
				ObjectMeta: kapi.ObjectMeta{
					Name:      "name",
					Namespace: "foo",
				},
			},
		},
		{
			name: "Valid route",
			route: &routeapi.Route{
				ObjectMeta: kapi.ObjectMeta{
					Name:      "name",
					Namespace: "foo",
				},
				Host:        "www.example.org",
				ServiceName: "serviceName",
			},
		},
	}

	plugin := &TestAllocationPlugin{Name: "test allocation plugin"}
	fac := &RouteAllocationControllerFactory{nil, nil}
	allocator := fac.Create(plugin)
	for _, tc := range tests {
		shard, err := allocator.AllocateRouterShard(tc.route)
		if err != nil {
			t.Errorf("Test case %s got an error %s", tc.name, err)
			continue
		}
		name := allocator.GenerateHostname(tc.route, shard)
		if len(name) <= 0 {
			t.Errorf("Test case %s got %d length name", tc.name, len(name))
		}
	}
}
