package imagestreammapping

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/projectatomic/appinfra-next/pkg/image/api"
)

// Registry is an interface for things that know how to store ImageStreamMapping objects.
type Registry interface {
	// CreateImageStreamMapping creates a new image stream mapping.
	CreateImageStreamMapping(ctx kapi.Context, mapping *api.ImageStreamMapping) (*kapi.Status, error)
}

// Storage is an interface for a standard REST Storage backend
type Storage interface {
	Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error)
}

// storage puts strong typing around storage calls
type storage struct {
	Storage
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s Storage) Registry {
	return &storage{s}
}

func (s *storage) CreateImageStreamMapping(ctx kapi.Context, mapping *api.ImageStreamMapping) (*kapi.Status, error) {
	obj, err := s.Create(ctx, mapping)
	if err != nil {
		return nil, err
	}
	return obj.(*kapi.Status), nil
}
