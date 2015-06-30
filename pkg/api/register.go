package api

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	_ "github.com/projectatomic/atomic-enterprise/pkg/authorization/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/build/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/deploy/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/image/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/oauth/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/project/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/route/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/sdn/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/template/api"
	_ "github.com/projectatomic/atomic-enterprise/pkg/user/api"
)

// Codec is the identity codec for this package - it can only convert itself
// to itself.
var Codec = runtime.CodecFor(api.Scheme, "")

func init() {
	api.Scheme.AddKnownTypes("")
}
