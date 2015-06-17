package generator

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	"github.com/projectatomic/appinfra-next/pkg/build/generator"
	"github.com/projectatomic/appinfra-next/pkg/build/registry/clone"
)

// NewStorage creates a new storage object for build generation
func NewStorage(generator *generator.BuildGenerator) *CloneREST {
	return &CloneREST{generator: generator}
}

// CloneREST is a RESTStorage implementation for a BuildGenerator which supports only
// the Get operation (as the generator has no underlying storage object).
type CloneREST struct {
	generator *generator.BuildGenerator
}

// New creates a new build clone request
func (s *CloneREST) New() runtime.Object {
	return &buildapi.BuildRequest{}
}

// Create instantiates a new build from an existing build
func (s *CloneREST) Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error) {
	if err := rest.BeforeCreate(clone.Strategy, ctx, obj); err != nil {
		return nil, err
	}

	return s.generator.Clone(ctx, obj.(*buildapi.BuildRequest))
}
