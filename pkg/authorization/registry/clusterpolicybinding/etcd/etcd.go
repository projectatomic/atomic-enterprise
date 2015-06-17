package etcd

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	etcdgeneric "github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic/etcd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/tools"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	"github.com/projectatomic/appinfra-next/pkg/authorization/registry/clusterpolicybinding"
)

const ClusterPolicyBindingPath = "/authorization/cluster/policybindings"

type REST struct {
	*etcdgeneric.Etcd
}

// NewStorage returns a RESTStorage object that will work against nodes.
func NewStorage(h tools.EtcdHelper) *REST {
	store := &etcdgeneric.Etcd{
		NewFunc:      func() runtime.Object { return &authorizationapi.ClusterPolicyBinding{} },
		NewListFunc:  func() runtime.Object { return &authorizationapi.ClusterPolicyBindingList{} },
		EndpointName: "clusterpolicybinding",
		KeyRootFunc: func(ctx kapi.Context) string {
			return ClusterPolicyBindingPath
		},
		KeyFunc: func(ctx kapi.Context, id string) (string, error) {
			return etcdgeneric.NoNamespaceKeyFunc(ctx, ClusterPolicyBindingPath, id)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*authorizationapi.ClusterPolicyBinding).Name, nil
		},
		PredicateFunc: func(label labels.Selector, field fields.Selector) generic.Matcher {
			return clusterpolicybinding.Matcher(label, field)
		},

		CreateStrategy: clusterpolicybinding.Strategy,
		UpdateStrategy: clusterpolicybinding.Strategy,

		Helper: h,
	}

	return &REST{store}
}
