package api

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	_ "github.com/projectatomic/appinfra-next/pkg/authorization/api"
	_ "github.com/projectatomic/appinfra-next/pkg/build/api"
	_ "github.com/projectatomic/appinfra-next/pkg/deploy/api"
	_ "github.com/projectatomic/appinfra-next/pkg/image/api"
	_ "github.com/projectatomic/appinfra-next/pkg/oauth/api"
	_ "github.com/projectatomic/appinfra-next/pkg/project/api"
	_ "github.com/projectatomic/appinfra-next/pkg/route/api"
	_ "github.com/projectatomic/appinfra-next/pkg/sdn/api"
	_ "github.com/projectatomic/appinfra-next/pkg/template/api"
	_ "github.com/projectatomic/appinfra-next/pkg/user/api"
)

// Codec is the identity codec for this package - it can only convert itself
// to itself.
var Codec = runtime.CodecFor(api.Scheme, "")

func init() {
	api.Scheme.AddKnownTypes("")
}
