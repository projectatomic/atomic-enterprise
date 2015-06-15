package generator

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/build/generator"
	"github.com/projectatomic/appinfra-next/pkg/build/registry/instantiate"
)

// NewStorage creates a new storage object for build generation
func NewStorage(generator *generator.BuildGenerator) *InstantiateREST {
	return &InstantiateREST{generator: generator}
}

// InstantiateREST is a RESTStorage implementation for a BuildGenerator which supports only
// the Get operation (as the generator has no underlying storage object).
type InstantiateREST struct {
	generator *generator.BuildGenerator
}

// New creates a new build generation request
func (s *InstantiateREST) New() runtime.Object {
	return &buildapi.BuildRequest{}
}

// Create instantiates a new build from a build configuration
func (s *InstantiateREST) Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error) {
	if err := rest.BeforeCreate(instantiate.Strategy, ctx, obj); err != nil {
		return nil, err
	}

	return s.generator.Instantiate(ctx, obj.(*buildapi.BuildRequest))
}
