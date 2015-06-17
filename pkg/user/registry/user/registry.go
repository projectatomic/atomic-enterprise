package user

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	"github.com/projectatomic/appinfra-next/pkg/user/api"
)

// Registry is an interface implemented by things that know how to store User objects.
type Registry interface {
	// ListUsers obtains a list of users having labels which match selector.
	ListUsers(ctx kapi.Context, selector labels.Selector) (*api.UserList, error)
	// GetUser returns a specific user
	GetUser(ctx kapi.Context, name string) (*api.User, error)
	// CreateUser creates a user
	CreateUser(ctx kapi.Context, user *api.User) (*api.User, error)
	// UpdateUser updates an existing user
	UpdateUser(ctx kapi.Context, user *api.User) (*api.User, error)
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

func (s *storage) ListUsers(ctx kapi.Context, label labels.Selector) (*api.UserList, error) {
	obj, err := s.List(ctx, label, fields.Everything())
	if err != nil {
		return nil, err
	}
	return obj.(*api.UserList), nil
}

func (s *storage) GetUser(ctx kapi.Context, name string) (*api.User, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.User), nil
}

func (s *storage) CreateUser(ctx kapi.Context, user *api.User) (*api.User, error) {
	obj, err := s.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return obj.(*api.User), nil
}

func (s *storage) UpdateUser(ctx kapi.Context, user *api.User) (*api.User, error) {
	obj, _, err := s.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return obj.(*api.User), nil
}
