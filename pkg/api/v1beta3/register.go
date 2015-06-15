package v1beta3

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	_ "github.com/projectatomic/appinfra-next/pkg/authorization/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/build/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/deploy/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/image/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/oauth/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/project/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/route/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/sdn/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/template/api/v1beta3"
	_ "github.com/projectatomic/appinfra-next/pkg/user/api/v1beta3"
)

// Codec encodes internal objects to the v1beta3 scheme
var Codec = runtime.CodecFor(api.Scheme, "v1beta3")

func init() {
	api.Scheme.AddKnownTypes("v1beta3")
}
