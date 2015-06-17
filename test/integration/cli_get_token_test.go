// +build integration,!no-etcd

package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/master"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/tools/etcdtest"

	// for osinserver setup.
	"github.com/projectatomic/appinfra-next/pkg/api/latest"
	"github.com/projectatomic/appinfra-next/pkg/auth/authenticator/challenger/passwordchallenger"
	"github.com/projectatomic/appinfra-next/pkg/auth/authenticator/password/allowanypassword"
	"github.com/projectatomic/appinfra-next/pkg/auth/authenticator/request/basicauthrequest"
	oauthhandlers "github.com/projectatomic/appinfra-next/pkg/auth/oauth/handlers"
	oauthregistry "github.com/projectatomic/appinfra-next/pkg/auth/oauth/registry"
	"github.com/projectatomic/appinfra-next/pkg/auth/userregistry/identitymapper"
	"github.com/projectatomic/appinfra-next/pkg/cmd/server/origin"
	accesstokenregistry "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthaccesstoken"
	accesstokenetcd "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthaccesstoken/etcd"
	authorizetokenregistry "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthauthorizetoken"
	authorizetokenetcd "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthauthorizetoken/etcd"
	clientregistry "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthclient"
	clientetcd "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthclient/etcd"
	clientauthregistry "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthclientauthorization"
	clientauthetcd "github.com/projectatomic/appinfra-next/pkg/oauth/registry/oauthclientauthorization/etcd"
	"github.com/projectatomic/appinfra-next/pkg/oauth/server/osinserver"
	"github.com/projectatomic/appinfra-next/pkg/oauth/server/osinserver/registrystorage"
	identityregistry "github.com/projectatomic/appinfra-next/pkg/user/registry/identity"
	identityetcd "github.com/projectatomic/appinfra-next/pkg/user/registry/identity/etcd"
	userregistry "github.com/projectatomic/appinfra-next/pkg/user/registry/user"
	useretcd "github.com/projectatomic/appinfra-next/pkg/user/registry/user/etcd"

	"github.com/projectatomic/appinfra-next/pkg/cmd/util/clientcmd"
	"github.com/projectatomic/appinfra-next/pkg/cmd/util/tokencmd"
	testutil "github.com/projectatomic/appinfra-next/test/util"
)

func init() {
	testutil.RequireEtcd()
}

func TestCLIGetToken(t *testing.T) {
	testutil.DeleteAllEtcdKeys()

	// setup
	etcdClient := testutil.NewEtcdClient()
	etcdHelper, _ := master.NewEtcdHelper(etcdClient, latest.Version, etcdtest.PathPrefix())

	accessTokenStorage := accesstokenetcd.NewREST(etcdHelper)
	accessTokenRegistry := accesstokenregistry.NewRegistry(accessTokenStorage)
	authorizeTokenStorage := authorizetokenetcd.NewREST(etcdHelper)
	authorizeTokenRegistry := authorizetokenregistry.NewRegistry(authorizeTokenStorage)
	clientStorage := clientetcd.NewREST(etcdHelper)
	clientRegistry := clientregistry.NewRegistry(clientStorage)
	clientAuthStorage := clientauthetcd.NewREST(etcdHelper)
	clientAuthRegistry := clientauthregistry.NewRegistry(clientAuthStorage)

	userStorage := useretcd.NewREST(etcdHelper)
	userRegistry := userregistry.NewRegistry(userStorage)
	identityStorage := identityetcd.NewREST(etcdHelper)
	identityRegistry := identityregistry.NewRegistry(identityStorage)

	identityMapper := identitymapper.NewAlwaysCreateUserIdentityToUserMapper(identityRegistry, userRegistry)

	authRequestHandler := basicauthrequest.NewBasicAuthAuthentication(allowanypassword.New("get-token-test", identityMapper), true)
	authHandler := oauthhandlers.NewUnionAuthenticationHandler(
		map[string]oauthhandlers.AuthenticationChallenger{"login": passwordchallenger.NewBasicAuthChallenger("openshift")}, nil, nil)

	storage := registrystorage.New(accessTokenRegistry, authorizeTokenRegistry, clientRegistry, oauthregistry.NewUserConversion())
	config := osinserver.NewDefaultServerConfig()

	grantChecker := oauthregistry.NewClientAuthorizationGrantChecker(clientAuthRegistry)
	grantHandler := oauthhandlers.NewAutoGrant()

	server := osinserver.New(
		config,
		storage,
		osinserver.AuthorizeHandlers{
			oauthhandlers.NewAuthorizeAuthenticator(
				authRequestHandler,
				authHandler,
				oauthhandlers.EmptyError{},
			),
			oauthhandlers.NewGrantCheck(
				grantChecker,
				grantHandler,
				oauthhandlers.EmptyError{},
			),
		},
		osinserver.AccessHandlers{
			oauthhandlers.NewDenyAccessAuthenticator(),
		},
		osinserver.NewDefaultErrorHandler(),
	)
	mux := http.NewServeMux()
	server.Install(mux, origin.OpenShiftOAuthAPIPrefix)
	oauthServer := httptest.NewServer(http.Handler(mux))
	defer oauthServer.Close()
	t.Logf("oauth server is on %v\n", oauthServer.URL)

	// create the default oauth clients with redirects to our server
	origin.CreateOrUpdateDefaultOAuthClients(oauthServer.URL, []string{oauthServer.URL}, clientRegistry)

	flags := pflag.NewFlagSet("test-flags", pflag.ContinueOnError)
	clientCfg := clientcmd.NewConfig()
	clientCfg.Bind(flags)
	flags.Parse(strings.Split("--master="+oauthServer.URL, " "))

	reader := bytes.NewBufferString("user\npass")

	accessToken, err := tokencmd.RequestToken(clientCfg.OpenShiftConfig(), reader, "", "")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(accessToken) == 0 {
		t.Error("Expected accessToken, but did not get one")
	}

	// lets see if this access token is any good
	token, err := accessTokenRegistry.GetAccessToken(kapi.NewContext(), accessToken)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if token.UserName != "user" {
		t.Errorf("Expected token for \"user\", but got: %#v", token)
	}
}
