package v1

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	_ "github.com/projectatomic/atomic-enterprise/pkg/authorization/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/build/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/deploy/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/image/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/oauth/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/project/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/route/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/sdn/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/template/api/v1"
	_ "github.com/projectatomic/atomic-enterprise/pkg/user/api/v1"
)

// Codec encodes internal objects to the v1 scheme
var Codec = runtime.CodecFor(api.Scheme, "v1")

func init() {
}
