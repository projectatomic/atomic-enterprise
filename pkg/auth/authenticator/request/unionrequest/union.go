package unionrequest

import (
	"net/http"

	kerrors "github.com/GoogleCloudPlatform/kubernetes/pkg/util/errors"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/auth/user"
	"github.com/projectatomic/appinfra-next/pkg/auth/authenticator"
)

// TODO remove this in favor of kubernetes types

type Authenticator struct {
	Handlers    []authenticator.Request
	FailOnError bool
}

// NewUnionAuthentication returns a request authenticator that validates credentials using a chain of authenticator.Request objects
func NewUnionAuthentication(authRequestHandlers ...authenticator.Request) authenticator.Request {
	return &Authenticator{Handlers: authRequestHandlers}
}

// AuthenticateRequest authenticates the request using a chain of authenticator.Request objects.  The first
// success returns that identity.  Errors are only returned if no matches are found.
func (authHandler *Authenticator) AuthenticateRequest(req *http.Request) (user.Info, bool, error) {
	errors := []error{}
	for _, currAuthRequestHandler := range authHandler.Handlers {
		info, ok, err := currAuthRequestHandler.AuthenticateRequest(req)
		if err == nil && ok {
			return info, ok, err
		}
		if err != nil {
			if authHandler.FailOnError {
				return nil, false, err
			}
			errors = append(errors, err)
		}
	}

	return nil, false, kerrors.NewAggregate(errors)
}
