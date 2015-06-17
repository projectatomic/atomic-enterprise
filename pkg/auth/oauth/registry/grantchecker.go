package registry

import (
	"fmt"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	"github.com/projectatomic/appinfra-next/pkg/auth/api"
	"github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthclientauthorization"
	"github.com/projectatomic/appinfra-next/pkg/oauth/scope"
)

type ClientAuthorizationGrantChecker struct {
	registry oauthclientauthorization.Registry
}

func NewClientAuthorizationGrantChecker(registry oauthclientauthorization.Registry) *ClientAuthorizationGrantChecker {
	return &ClientAuthorizationGrantChecker{registry}
}

func (c *ClientAuthorizationGrantChecker) HasAuthorizedClient(user user.Info, grant *api.Grant) (approved bool, err error) {
	id := c.registry.ClientAuthorizationName(user.GetName(), grant.Client.GetId())
	authorization, err := c.registry.GetClientAuthorization(kapi.NewContext(), id)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(authorization.UserUID) != 0 && authorization.UserUID != user.GetUID() {
		return false, fmt.Errorf("user %s UID %s does not match stored client authorization value for UID %s", user.GetName(), user.GetUID(), authorization.UserUID)
	}
	// TODO: improve this to allow the scope implementation to determine overlap
	if !scope.Covers(authorization.Scopes, scope.Split(grant.Scope)) {
		return false, nil
	}
	return true, nil
}
