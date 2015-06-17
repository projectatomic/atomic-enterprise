package user

import (
	"github.com/projectatomic/appinfra-next/pkg/user/api"
)

type Initializer interface {
	InitializeUser(identity *api.Identity, user *api.User) error
}
