package imagestream

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	"github.com/projectatomic/appinfra-next/pkg/image/api"
)

// Registry is an interface for things that know how to store ImageStream objects.
type Registry interface {
	// ListImageStreams obtains a list of image streams that match a selector.
	ListImageStreams(ctx kapi.Context, selector labels.Selector) (*api.ImageStreamList, error)
	// GetImageStream retrieves a specific image stream.
	GetImageStream(ctx kapi.Context, id string) (*api.ImageStream, error)
	// CreateImageStream creates a new image stream.
	CreateImageStream(ctx kapi.Context, repo *api.ImageStream) (*api.ImageStream, error)
	// UpdateImageStream updates an image stream.
	UpdateImageStream(ctx kapi.Context, repo *api.ImageStream) (*api.ImageStream, error)
	// UpdateImageStream updates an image stream's status.
	UpdateImageStreamStatus(ctx kapi.Context, repo *api.ImageStream) (*api.ImageStream, error)
	// DeleteImageStream deletes an image stream.
	DeleteImageStream(ctx kapi.Context, id string) (*kapi.Status, error)
	// WatchImageStreams watches for new/changed/deleted image streams.
	WatchImageStreams(ctx kapi.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
}

// Storage is an interface for a standard REST Storage backend
type Storage interface {
	rest.GracefulDeleter
	rest.Lister
	rest.Getter
	rest.Watcher

	Create(ctx kapi.Context, obj runtime.Object) (runtime.Object, error)
	Update(ctx kapi.Context, obj runtime.Object) (runtime.Object, bool, error)
}

// storage puts strong typing around storage calls
type storage struct {
	Storage
	status rest.Updater
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s Storage, status rest.Updater) Registry {
	return &storage{s, status}
}

func (s *storage) ListImageStreams(ctx kapi.Context, label labels.Selector) (*api.ImageStreamList, error) {
	obj, err := s.List(ctx, label, fields.Everything())
	if err != nil {
		return nil, err
	}
	return obj.(*api.ImageStreamList), nil
}

func (s *storage) GetImageStream(ctx kapi.Context, imageStreamID string) (*api.ImageStream, error) {
	obj, err := s.Get(ctx, imageStreamID)
	if err != nil {
		return nil, err
	}
	return obj.(*api.ImageStream), nil
}

func (s *storage) CreateImageStream(ctx kapi.Context, imageStream *api.ImageStream) (*api.ImageStream, error) {
	obj, err := s.Create(ctx, imageStream)
	if err != nil {
		return nil, err
	}
	return obj.(*api.ImageStream), nil
}

func (s *storage) UpdateImageStream(ctx kapi.Context, imageStream *api.ImageStream) (*api.ImageStream, error) {
	obj, _, err := s.Update(ctx, imageStream)
	if err != nil {
		return nil, err
	}
	return obj.(*api.ImageStream), nil
}

func (s *storage) UpdateImageStreamStatus(ctx kapi.Context, imageStream *api.ImageStream) (*api.ImageStream, error) {
	obj, _, err := s.status.Update(ctx, imageStream)
	if err != nil {
		return nil, err
	}
	return obj.(*api.ImageStream), nil
}

func (s *storage) DeleteImageStream(ctx kapi.Context, imageStreamID string) (*kapi.Status, error) {
	obj, err := s.Delete(ctx, imageStreamID, nil)
	if err != nil {
		return nil, err
	}
	return obj.(*kapi.Status), nil
}

func (s *storage) WatchImageStreams(ctx kapi.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return s.Watch(ctx, label, field, resourceVersion)
}
