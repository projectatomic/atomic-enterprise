package allocation

import (
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"

	osclient "github.com/projectatomic/appinfra-next/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/route"
)

// RouteAllocationControllerFactory creates a RouteAllocationController
// that allocates router shards to specific routes.
type RouteAllocationControllerFactory struct {
	// OSClient is is an OpenShift client.
	OSClient osclient.Interface

	// KubeClient is a Kubernetes client.
	KubeClient kclient.Interface
}

// Create a RouteAllocationController instance.
func (factory *RouteAllocationControllerFactory) Create(plugin route.AllocationPlugin) *RouteAllocationController {
	return &RouteAllocationController{Plugin: plugin}
}
