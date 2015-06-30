package github

import (
	"testing"

	"github.com/projectatomic/atomic-enterprise/pkg/auth/oauth/external"
)

func TestGitHub(t *testing.T) {
	_ = external.Provider(NewProvider("github", "clientid", "clientsecret"))
}
