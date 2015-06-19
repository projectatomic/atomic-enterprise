package proxy

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	authorizationapi "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	clusterpolicyregistry "github.com/projectatomic/appinfra-next/pkg/authorization/registry/clusterpolicy"
	roleregistry "github.com/projectatomic/appinfra-next/pkg/authorization/registry/role"
	rolestorage "github.com/projectatomic/appinfra-next/pkg/authorization/registry/role/policybased"
)

type ClusterRoleStorage struct {
	roleStorage rolestorage.VirtualStorage
}

func NewClusterRoleStorage(clusterPolicyRegistry clusterpolicyregistry.Registry) *ClusterRoleStorage {
	return &ClusterRoleStorage{rolestorage.VirtualStorage{clusterpolicyregistry.NewSimulatedRegistry(clusterPolicyRegistry), roleregistry.ClusterStrategy, roleregistry.ClusterStrategy}}
}

func (s *ClusterRoleStorage) New() runtime.Object {
	return &authorizationapi.ClusterRole{}
}
func (s *ClusterRoleStorage) NewList() runtime.Object {
	return &authorizationapi.ClusterRoleList{}
}

func (s *ClusterRoleStorage) List(ctx kapi.Context, label labels.Selector, field fields.Selector) (runtime.Object, error) {
	ret, err := s.roleStorage.List(ctx, label, field)
	if ret == nil {
		return nil, err
	}
	return authorizationapi.ToClusterRoleList(ret.(*authorizationapi.RoleList)), err
}

func (s *ClusterRoleStorage) Get(ctx kapi.Context, name string) (runtime.Object, error) {
	ret, err := s.roleStorage.Get(ctx, name)
	if ret == nil {
		return nil, err
	}

	return authorizationapi.ToClusterRole(ret.(*authorizationapi.Role)), err
}
func (s *ClusterRoleStorage) Delete(ctx kapi.Context, name string, options *kapi.DeleteOptions) (runtime.Object, error) {
	ret, err := s.roleStorage.Delete(ctx, name, options)
	if ret == nil {
		return nil, err
	}

	return ret.(*kapi.Status), err
}

func (s *ClusterRoleStorage) Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error) {
	clusterObj := obj.(*authorizationapi.ClusterRole)
	convertedObj := authorizationapi.ToRole(clusterObj)

	ret, err := s.roleStorage.Create(ctx, convertedObj)
	if ret == nil {
		return nil, err
	}

	return authorizationapi.ToClusterRole(ret.(*authorizationapi.Role)), err
}

func (s *ClusterRoleStorage) Update(ctx kapi.Context, obj runtime.Object) (runtime.Object, bool, error) {
	clusterObj := obj.(*authorizationapi.ClusterRole)
	convertedObj := authorizationapi.ToRole(clusterObj)

	ret, created, err := s.roleStorage.Update(ctx, convertedObj)
	if ret == nil {
		return nil, created, err
	}

	return authorizationapi.ToClusterRole(ret.(*authorizationapi.Role)), created, err
}
