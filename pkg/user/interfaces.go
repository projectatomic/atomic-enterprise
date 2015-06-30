package user

import (
	"github.com/projectatomic/atomic-enterprise/pkg/user/api"
)

type Initializer interface {
	InitializeUser(identity *api.Identity, user *api.User) error
}
