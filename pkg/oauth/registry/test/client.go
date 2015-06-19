package test

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"

	"github.com/projectatomic/appinfra-next/pkg/oauth/api"
)

type ClientRegistry struct {
	Err               error
	Clients           *api.OAuthClientList
	Client            *api.OAuthClient
	DeletedClientName string
}

func (r *ClientRegistry) ListClients(ctx kapi.Context, labels labels.Selector) (*api.OAuthClientList, error) {
	return r.Clients, r.Err
}

func (r *ClientRegistry) GetClient(ctx kapi.Context, name string) (*api.OAuthClient, error) {
	return r.Client, r.Err
}

func (r *ClientRegistry) CreateClient(ctx kapi.Context, client *api.OAuthClient) (*api.OAuthClient, error) {
	return r.Client, r.Err
}

func (r *ClientRegistry) UpdateClient(ctx kapi.Context, client *api.OAuthClient) (*api.OAuthClient, error) {
	return r.Client, r.Err
}

func (r *ClientRegistry) DeleteClient(ctx kapi.Context, name string) error {
	r.DeletedClientName = name
	return r.Err
}
