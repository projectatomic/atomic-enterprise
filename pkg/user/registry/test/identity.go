package test

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kerrs "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	"github.com/projectatomic/appinfra-next/pkg/user/api"
)

type Action struct {
	Name   string
	Object interface{}
}

type IdentityRegistry struct {
	GetErr map[string]error
	Get    map[string]*api.Identity

	CreateErr error
	Create    *api.Identity

	UpdateErr error
	Update    *api.Identity

	ListErr error
	List    *api.IdentityList

	Actions *[]Action
}

func NewIdentityRegistry() *IdentityRegistry {
	return &IdentityRegistry{
		GetErr:  map[string]error{},
		Get:     map[string]*api.Identity{},
		Actions: &[]Action{},
	}
}

func (r *IdentityRegistry) GetIdentity(ctx kapi.Context, name string) (*api.Identity, error) {
	*r.Actions = append(*r.Actions, Action{"GetIdentity", name})
	if identity, ok := r.Get[name]; ok {
		return identity, nil
	}
	if err, ok := r.GetErr[name]; ok {
		return nil, err
	}
	return nil, kerrs.NewNotFound("Identity", name)
}

func (r *IdentityRegistry) CreateIdentity(ctx kapi.Context, u *api.Identity) (*api.Identity, error) {
	*r.Actions = append(*r.Actions, Action{"CreateIdentity", u})
	if r.Create == nil && r.CreateErr == nil {
		return u, nil
	}
	return r.Create, r.CreateErr
}

func (r *IdentityRegistry) UpdateIdentity(ctx kapi.Context, u *api.Identity) (*api.Identity, error) {
	*r.Actions = append(*r.Actions, Action{"UpdateIdentity", u})
	if r.Update == nil && r.UpdateErr == nil {
		return u, nil
	}
	return r.Update, r.UpdateErr
}

func (r *IdentityRegistry) ListIdentities(ctx kapi.Context, labels labels.Selector) (*api.IdentityList, error) {
	*r.Actions = append(*r.Actions, Action{"ListIdentities", labels})
	if r.List == nil && r.ListErr == nil {
		return &api.IdentityList{}, nil
	}
	return r.List, r.ListErr
}
