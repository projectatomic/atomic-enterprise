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
	bindingregistry "github.com/projectatomic/appinfra-next/pkg/authorization/registry/policybinding"
)

type readOnlyPolicyBindingCache struct {
	registry  bindingregistry.WatchingRegistry
	indexer   cache.Indexer
	reflector cache.Reflector

	keyFunc cache.KeyFunc
}

func NewReadOnlyPolicyBindingCache(registry bindingregistry.WatchingRegistry) readOnlyPolicyBindingCache {
	ctx := kapi.WithNamespace(kapi.NewContext(), kapi.NamespaceAll)

	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"namespace": cache.MetaNamespaceIndexFunc})

	reflector := cache.NewReflector(
		&cache.ListWatch{
			ListFunc: func() (runtime.Object, error) {
				return registry.ListPolicyBindings(ctx, labels.Everything(), fields.Everything())
			},
			WatchFunc: func(resourceVersion string) (watch.Interface, error) {
				return registry.WatchPolicyBindings(ctx, labels.Everything(), fields.Everything(), resourceVersion)
			},
		},
		&authorizationapi.PolicyBinding{},
		indexer,
		2*time.Minute,
	)

	return readOnlyPolicyBindingCache{
		registry:  registry,
		indexer:   indexer,
		reflector: *reflector,

		keyFunc: cache.MetaNamespaceKeyFunc,
	}
}

// Run begins watching and synchronizing the cache
func (c *readOnlyPolicyBindingCache) Run() {
	c.reflector.Run()
}

// RunUntil starts a watch and handles watch events. Will restart the watch if it is closed.
// RunUntil starts a goroutine and returns immediately. It will exit when stopCh is closed.
func (c *readOnlyPolicyBindingCache) RunUntil(stopChannel <-chan struct{}) {
	c.reflector.RunUntil(stopChannel)
}

// LastSyncResourceVersion exposes the LastSyncResourceVersion of the internal reflector
func (c *readOnlyPolicyBindingCache) LastSyncResourceVersion() string {
	return c.reflector.LastSyncResourceVersion()
}

func (c *readOnlyPolicyBindingCache) List(label labels.Selector, field fields.Selector, namespace string) (*authorizationapi.PolicyBindingList, error) {
	var returnedList []interface{}
	if namespace == kapi.NamespaceAll {
		returnedList = c.indexer.List()
	} else {
		items, err := c.indexer.Index("namespace", &authorizationapi.PolicyBinding{ObjectMeta: kapi.ObjectMeta{Namespace: namespace}})
		returnedList = items
		if err != nil {
			return &authorizationapi.PolicyBindingList{}, errors.NewInvalid("PolicyBindingList", "policyBindingList", []error{err})
		}
	}
	policyBindingList := &authorizationapi.PolicyBindingList{}
	for i := range returnedList {
		policyBinding, castOK := returnedList[i].(*authorizationapi.PolicyBinding)
		if !castOK {
			return policyBindingList, errors.NewInvalid("PolicyBindingList", "policyBindingList", []error{})
		}
		if label.Matches(labels.Set(policyBinding.Labels)) && field.Matches(PolicyBindingToSelectableFields(policyBinding)) {
			policyBindingList.Items = append(policyBindingList.Items, *policyBinding)
		}
	}
	return policyBindingList, nil
}

func (c *readOnlyPolicyBindingCache) Get(name, namespace string) (*authorizationapi.PolicyBinding, error) {
	keyObj := &authorizationapi.PolicyBinding{ObjectMeta: kapi.ObjectMeta{Namespace: namespace, Name: name}}
	key, _ := c.keyFunc(keyObj)

	item, exists, getErr := c.indexer.GetByKey(key)
	if getErr != nil {
		return &authorizationapi.PolicyBinding{}, getErr
	}
	if !exists {
		existsErr := errors.NewNotFound("PolicyBinding", name)
		return &authorizationapi.PolicyBinding{}, existsErr
	}
	policyBinding, castOK := item.(*authorizationapi.PolicyBinding)
	if !castOK {
		castErr := errors.NewInvalid("PolicyBinding", name, []error{})
		return &authorizationapi.PolicyBinding{}, castErr
	}
	return policyBinding, nil
}

// readOnlyPolicyBindings wraps readOnlyPolicyBindingCache in order to expose only List() and Get()
type readOnlyPolicyBindings struct {
	readOnlyPolicyBindingCache *readOnlyPolicyBindingCache
	namespace                  string
}

func newReadOnlyPolicyBindings(cache readOnlyAuthorizationCache, namespace string) client.ReadOnlyPolicyBindingInterface {
	return &readOnlyPolicyBindings{
		readOnlyPolicyBindingCache: &cache.readOnlyPolicyBindingCache,
		namespace:                  namespace,
	}
}

func (p *readOnlyPolicyBindings) List(label labels.Selector, field fields.Selector) (*authorizationapi.PolicyBindingList, error) {
	return p.readOnlyPolicyBindingCache.List(label, field, p.namespace)
}

func (p *readOnlyPolicyBindings) Get(name string) (*authorizationapi.PolicyBinding, error) {
	return p.readOnlyPolicyBindingCache.Get(name, p.namespace)
}
