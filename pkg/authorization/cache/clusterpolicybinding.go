package cache

import (
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	errors "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/cache"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	"github.com/projectatomic/appinfra-next/pkg/authorization/client"
	clusterbindingregistry "github.com/projectatomic/appinfra-next/pkg/authorization/registry/clusterpolicybinding"
)

type readOnlyClusterPolicyBindingCache struct {
	registry  clusterbindingregistry.WatchingRegistry
	indexer   cache.Indexer
	reflector cache.Reflector

	keyFunc cache.KeyFunc
}

func NewReadOnlyClusterPolicyBindingCache(registry clusterbindingregistry.WatchingRegistry) readOnlyClusterPolicyBindingCache {
	ctx := kapi.WithNamespace(kapi.NewContext(), kapi.NamespaceAll)

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"namespace": cache.MetaNamespaceIndexFunc})

	reflector := cache.NewReflector(
		&cache.ListWatch{
			ListFunc: func() (runtime.Object, error) {
				return registry.ListClusterPolicyBindings(ctx, labels.Everything(), fields.Everything())
			},
			WatchFunc: func(resourceVersion string) (watch.Interface, error) {
				return registry.WatchClusterPolicyBindings(ctx, labels.Everything(), fields.Everything(), resourceVersion)
			},
		},
		&authorizationapi.ClusterPolicyBinding{},
		indexer,
		2*time.Minute,
	)

	return readOnlyClusterPolicyBindingCache{
		registry:  registry,
		indexer:   indexer,
		reflector: *reflector,

		keyFunc: cache.MetaNamespaceKeyFunc,
	}
}

// Run begins watching and synchronizing the cache
func (c *readOnlyClusterPolicyBindingCache) Run() {
	c.reflector.Run()
}

// RunUntil starts a watch and handles watch events. Will restart the watch if it is closed.
// RunUntil starts a goroutine and returns immediately. It will exit when stopCh is closed.
func (c *readOnlyClusterPolicyBindingCache) RunUntil(stopChannel <-chan struct{}) {
	c.reflector.RunUntil(stopChannel)
}

// LastSyncResourceVersion exposes the LastSyncResourceVersion of the internal reflector
func (c *readOnlyClusterPolicyBindingCache) LastSyncResourceVersion() string {
	return c.reflector.LastSyncResourceVersion()
}

func (c *readOnlyClusterPolicyBindingCache) List(label labels.Selector, field fields.Selector) (*authorizationapi.ClusterPolicyBindingList, error) {
	clusterPolicyBindingList := &authorizationapi.ClusterPolicyBindingList{}
	returnedList := c.indexer.List()
	for i := range returnedList {
		clusterPolicyBinding, castOK := returnedList[i].(*authorizationapi.ClusterPolicyBinding)
		if !castOK {
			return clusterPolicyBindingList, errors.NewInvalid("ClusterPolicyBinding", "clusterPolicyBinding", []error{})
		}
		if label.Matches(labels.Set(clusterPolicyBinding.Labels)) && field.Matches(ClusterPolicyBindingToSelectableFields(clusterPolicyBinding)) {
			clusterPolicyBindingList.Items = append(clusterPolicyBindingList.Items, *clusterPolicyBinding)
		}
	}
	return clusterPolicyBindingList, nil
}

func (c *readOnlyClusterPolicyBindingCache) Get(name string) (*authorizationapi.ClusterPolicyBinding, error) {
	keyObj := &authorizationapi.ClusterPolicyBinding{ObjectMeta: kapi.ObjectMeta{Name: name}}
	key, _ := c.keyFunc(keyObj)

	item, exists, getErr := c.indexer.GetByKey(key)
	if getErr != nil {
		return &authorizationapi.ClusterPolicyBinding{}, getErr
	}
	if !exists {
		existsErr := errors.NewNotFound("ClusterPolicyBinding", name)
		return &authorizationapi.ClusterPolicyBinding{}, existsErr
	}
	clusterPolicyBinding, castOK := item.(*authorizationapi.ClusterPolicyBinding)
	if !castOK {
		castErr := errors.NewInvalid("ClusterPolicyBinding", name, []error{})
		return &authorizationapi.ClusterPolicyBinding{}, castErr
	}
	return clusterPolicyBinding, nil
}

func newReadOnlyClusterPolicyBindings(cache readOnlyAuthorizationCache) client.ReadOnlyClusterPolicyBindingInterface {
	return &cache.readOnlyClusterPolicyBindingCache
}
