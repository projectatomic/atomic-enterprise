package test

import (
	"github.com/projectatomic/appinfra-next/pkg/user/api"
)

type UserIdentityMappingRegistry struct {
	Err                        error
	Created                    bool
	UserIdentityMapping        *api.UserIdentityMapping
	CreatedUserIdentityMapping *api.UserIdentityMapping
}

func (r *UserIdentityMappingRegistry) GetUserIdentityMapping(name string) (*api.UserIdentityMapping, error) {
	return r.UserIdentityMapping, r.Err
}

func (r *UserIdentityMappingRegistry) CreateOrUpdateUserIdentityMapping(mapping *api.UserIdentityMapping) (*api.UserIdentityMapping, bool, error) {
	r.CreatedUserIdentityMapping = mapping
	return r.CreatedUserIdentityMapping, r.Created, r.Err
}
