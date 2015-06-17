package authenticator

import (
	"net/http"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	"github.com/projectatomic/appinfra-next/pkg/auth/api"
)

type Token interface {
	AuthenticateToken(token string) (user.Info, bool, error)
}

type Request interface {
	AuthenticateRequest(req *http.Request) (user.Info, bool, error)
}

type Password interface {
	AuthenticatePassword(user, password string) (user.Info, bool, error)
}

type Assertion interface {
	AuthenticateAssertion(assertionType, data string) (user.Info, bool, error)
}

type Client interface {
	AuthenticateClient(client api.Client) (user.Info, bool, error)
}

type RequestFunc func(req *http.Request) (user.Info, bool, error)

func (f RequestFunc) AuthenticateRequest(req *http.Request) (user.Info, bool, error) {
	return f(req)
}
