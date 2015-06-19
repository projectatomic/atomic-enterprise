package role

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
)

// Registry is an interface for things that know how to store Roles.
type Registry interface {
	// ListRoles obtains list of policyRoles that match a selector.
	ListRoles(ctx kapi.Context, label labels.Selector, field fields.Selector) (*authorizationapi.RoleList, error)
	// GetRole retrieves a specific policyRole.
	GetRole(ctx kapi.Context, id string) (*authorizationapi.Role, error)
	// CreateRole creates a new policyRole.
	CreateRole(ctx kapi.Context, policyRole *authorizationapi.Role) (*authorizationapi.Role, error)
	// UpdateRole updates a policyRole.
	UpdateRole(ctx kapi.Context, policyRole *authorizationapi.Role) (*authorizationapi.Role, bool, error)
	// DeleteRole deletes a policyRole.
	DeleteRole(ctx kapi.Context, id string) error
}

// Storage is an interface for a standard REST Storage backend
type Storage interface {
	rest.Getter
	rest.Lister
	rest.CreaterUpdater
	rest.GracefulDeleter
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

func (s *storage) ListRoles(ctx kapi.Context, label labels.Selector, field fields.Selector) (*authorizationapi.RoleList, error) {
	obj, err := s.List(ctx, label, field)
	if err != nil {
		return nil, err
	}

	return obj.(*authorizationapi.RoleList), nil
}

func (s *storage) CreateRole(ctx kapi.Context, node *authorizationapi.Role) (*authorizationapi.Role, error) {
	obj, err := s.Create(ctx, node)
	return obj.(*authorizationapi.Role), err
}

func (s *storage) UpdateRole(ctx kapi.Context, node *authorizationapi.Role) (*authorizationapi.Role, bool, error) {
	obj, created, err := s.Update(ctx, node)
	return obj.(*authorizationapi.Role), created, err
}

func (s *storage) GetRole(ctx kapi.Context, name string) (*authorizationapi.Role, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*authorizationapi.Role), nil
}

func (s *storage) DeleteRole(ctx kapi.Context, name string) error {
	_, err := s.Delete(ctx, name, nil)
	return err
}
