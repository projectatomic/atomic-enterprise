package oauthpassword

import (
	"fmt"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	"github.com/RangelReale/osincli"
	"github.com/golang/glog"
	authapi "github.com/projectatomic/appinfra-next/pkg/auth/api"
	"github.com/projectatomic/appinfra-next/pkg/auth/authenticator"
)

type Authenticator struct {
	providerName string
	mapper       authapi.UserIdentityMapper
	client       *osincli.Client
}

func New(providerName string, client *osincli.Client, identityMapper authapi.UserIdentityMapper) authenticator.Password {
	return &Authenticator{providerName, identityMapper, client}
}

func (a *Authenticator) AuthenticatePassword(username, password string) (user.Info, bool, error) {
	areq := a.client.NewAccessRequest(osincli.PASSWORD, nil)
	areq.CustomParameters["username"] = username
	areq.CustomParameters["password"] = password
	token, err := areq.GetToken()
	if err != nil {
		if oerr, ok := err.(*osincli.Error); ok {
			if oerr.Id == osincli.E_ACCESS_DENIED {
				return nil, false, nil
			}
		}
		return nil, false, err
	}

	identity := authapi.NewDefaultUserIdentityInfo(a.providerName, username)
	identity.Extra["token"] = token.AccessToken
	user, err := a.mapper.UserFor(identity)
	glog.V(4).Infof("Got userIdentityMapping: %#v", user)
	if err != nil {
		return nil, false, fmt.Errorf("Error creating or updating mapping for: %#v due to %v", identity, err)
	}

	return user, true, nil
}
