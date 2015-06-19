// +build integration,!no-etcd

package integration

import (
	"io/ioutil"
	"os"
	"testing"

	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	"github.com/projectatomic/appinfra-next/pkg/client"
	configapi "github.com/projectatomic/appinfra-next/pkg/cmd/server/api"
	"github.com/projectatomic/appinfra-next/pkg/cmd/util/tokencmd"
	testutil "github.com/projectatomic/appinfra-next/test/util"
)

func TestHTPasswd(t *testing.T) {
	htpasswdFile, err := ioutil.TempFile("", "test.htpasswd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(htpasswdFile.Name())

	masterOptions, err := testutil.DefaultMasterOptions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	masterOptions.OAuthConfig.IdentityProviders[0] = configapi.IdentityProvider{
		Name:            "htpasswd",
		UseAsChallenger: true,
		UseAsLogin:      true,
		Provider: runtime.EmbeddedObject{
			&configapi.HTPasswdPasswordIdentityProvider{
				File: htpasswdFile.Name(),
			},
		},
	}

	clusterAdminKubeConfig, err := testutil.StartConfiguredMaster(masterOptions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	clientConfig, err := testutil.GetClusterAdminClientConfig(clusterAdminKubeConfig)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Use the server and CA info
	anonConfig := kclient.Config{}
	anonConfig.Host = clientConfig.Host
	anonConfig.CAFile = clientConfig.CAFile
	anonConfig.CAData = clientConfig.CAData

	// Make sure we can't authenticate
	if _, err := tokencmd.RequestToken(&anonConfig, nil, "username", "password"); err == nil {
		t.Error("Expected error, got none")
	}

	// Update the htpasswd file with output of `htpasswd -n -b username password`
	userpass := "username:$apr1$4Ci5I8yc$85R9vc4fOgzAULsldiUuv."
	ioutil.WriteFile(htpasswdFile.Name(), []byte(userpass), os.FileMode(0600))

	// Make sure we can get a token
	accessToken, err := tokencmd.RequestToken(&anonConfig, nil, "username", "password")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(accessToken) == 0 {
		t.Errorf("Expected access token, got none")
	}

	// Make sure we can use the token, and it represents who we expect
	userConfig := anonConfig
	userConfig.BearerToken = accessToken
	userClient, err := client.New(&userConfig)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	user, err := userClient.Users().Get("~")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if user.Name != "username" {
		t.Fatalf("Expected username as the user, got %v", user)
	}
}
