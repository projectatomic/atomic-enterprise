package registry

import (
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/client"
	oclient "github.com/projectatomic/appinfra-next/pkg/oauth/client"
)

type OAuthAccessTokenSource interface {
	AuthenticatePassword(username, password string) (string, bool, error)
}

type Authenticator struct {
	token OAuthAccessTokenSource
	host  string
	rt    http.RoundTripper
}

func New(token OAuthAccessTokenSource, host string, rt http.RoundTripper) *Authenticator {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &Authenticator{token, host, rt}
}

func (a *Authenticator) AuthenticatePassword(username, password string) (user.Info, bool, error) {
	token, ok, err := a.token.AuthenticatePassword(username, password)
	if !ok || err != nil {
		return nil, false, err
	}

	auth := oclient.OAuthWrapper{a.rt, token}

	client, err := client.New(&kclient.Config{Transport: auth, Host: a.host})
	if err != nil {
		return nil, false, err
	}
	u, err := client.Users().Get("~")
	if err != nil {
		return nil, false, err
	}

	info := &user.DefaultInfo{
		Name: u.Name,
		UID:  string(u.UID),
	}

	return info, true, nil
}
