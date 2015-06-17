package identity

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/projectatomic/appinfra-next/pkg/user/api"
)

// Registry is an interface implemented by things that know how to store Identity objects.
type Registry interface {
	// ListIdentities obtains a list of Identities having labels which match selector.
	ListIdentities(ctx kapi.Context, selector labels.Selector) (*api.IdentityList, error)
	// GetIdentity returns a specific Identity
	GetIdentity(ctx kapi.Context, name string) (*api.Identity, error)
	// CreateIdentity creates a Identity
	CreateIdentity(ctx kapi.Context, Identity *api.Identity) (*api.Identity, error)
	// UpdateIdentity updates an existing Identity
	UpdateIdentity(ctx kapi.Context, Identity *api.Identity) (*api.Identity, error)
}

func identityName(provider, identity string) string {
	// TODO: normalize?
	return provider + ":" + identity
}

// Storage is an interface for a standard REST Storage backend
// TODO: move me somewhere common
type Storage interface {
	rest.Lister
	rest.Getter

	Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error)
	Update(ctx kapi.Context, obj runtime.Object) (runtime.Object, bool, error)
}

// storage puts strong typing around storage calls
type storage struct {
	Storage
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s Storage) Registry {
	return &storage{s}
}

func (s *storage) ListIdentities(ctx kapi.Context, label labels.Selector) (*api.IdentityList, error) {
	obj, err := s.List(ctx, label, fields.Everything())
	if err != nil {
		return nil, err
	}
	return obj.(*api.IdentityList), nil
}

func (s *storage) GetIdentity(ctx kapi.Context, name string) (*api.Identity, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Identity), nil
}

func (s *storage) CreateIdentity(ctx kapi.Context, Identity *api.Identity) (*api.Identity, error) {
	obj, err := s.Create(ctx, Identity)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Identity), nil
}

func (s *storage) UpdateIdentity(ctx kapi.Context, Identity *api.Identity) (*api.Identity, error) {
	obj, _, err := s.Update(ctx, Identity)
	if err != nil {
		return nil, err
	}
	return obj.(*api.Identity), nil
}
