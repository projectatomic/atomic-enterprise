package v1

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/projectatomic/appinfra-next/pkg/cmd/server/api"
)

var Codec = runtime.CodecFor(api.Scheme, "v1")

func init() {
	api.Scheme.AddKnownTypes("v1",
		&MasterConfig{},
		&NodeConfig{},
		&SessionSecrets{},

		&IdentityProvider{},
		&BasicAuthPasswordIdentityProvider{},
		&AllowAllPasswordIdentityProvider{},
		&DenyAllPasswordIdentityProvider{},
		&HTPasswdPasswordIdentityProvider{},
		&RequestHeaderIdentityProvider{},
		&GitHubIdentityProvider{},
		&GoogleIdentityProvider{},
		&OpenIDIdentityProvider{},
		&GrantConfig{},
	)
}

func (*IdentityProvider) IsAnAPIObject()                  {}
func (*BasicAuthPasswordIdentityProvider) IsAnAPIObject() {}
func (*AllowAllPasswordIdentityProvider) IsAnAPIObject()  {}
func (*DenyAllPasswordIdentityProvider) IsAnAPIObject()   {}
func (*HTPasswdPasswordIdentityProvider) IsAnAPIObject()  {}
func (*RequestHeaderIdentityProvider) IsAnAPIObject()     {}
func (*GitHubIdentityProvider) IsAnAPIObject()            {}
func (*GoogleIdentityProvider) IsAnAPIObject()            {}
func (*OpenIDIdentityProvider) IsAnAPIObject()            {}
func (*GrantConfig) IsAnAPIObject()                       {}

func (*MasterConfig) IsAnAPIObject()   {}
func (*NodeConfig) IsAnAPIObject()     {}
func (*SessionSecrets) IsAnAPIObject() {}
