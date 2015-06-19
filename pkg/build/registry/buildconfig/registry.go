package buildconfig

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/projectatomic/appinfra-next/pkg/build/api"
)

// Registry is an interface for things that know how to store BuildConfigs.
type Registry interface {
	// ListBuildConfigs obtains list of buildConfigs that match a selector.
	ListBuildConfigs(ctx kapi.Context, labels labels.Selector, field fields.Selector) (*api.BuildConfigList, error)
	// GetBuildConfig retrieves a specific buildConfig.
	GetBuildConfig(ctx kapi.Context, id string) (*api.BuildConfig, error)
	// CreateBuildConfig creates a new buildConfig.
	CreateBuildConfig(ctx kapi.Context, buildConfig *api.BuildConfig) error
	// UpdateBuildConfig updates a buildConfig.
	UpdateBuildConfig(ctx kapi.Context, buildConfig *api.BuildConfig) error
	// DeleteBuildConfig deletes a buildConfig.
	DeleteBuildConfig(ctx kapi.Context, id string) error
	// WatchBuildConfigs watches buildConfigs.
	WatchBuildConfigs(ctx kapi.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
}

// storage puts strong typing around storage calls
type storage struct {
	rest.StandardStorage
}

// NewRegistry returns a new Registry interface for the given Storage. Any mismatched
// types will panic.
func NewRegistry(s rest.StandardStorage) Registry {
	return &storage{s}
}

func (s *storage) ListBuildConfigs(ctx kapi.Context, label labels.Selector, field fields.Selector) (*api.BuildConfigList, error) {
	obj, err := s.List(ctx, label, field)
	if err != nil {
		return nil, err
	}
	return obj.(*api.BuildConfigList), nil
}

func (s *storage) WatchBuildConfigs(ctx kapi.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return s.Watch(ctx, label, field, resourceVersion)
}

func (s *storage) GetBuildConfig(ctx kapi.Context, name string) (*api.BuildConfig, error) {
	obj, err := s.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return obj.(*api.BuildConfig), nil
}

func (s *storage) CreateBuildConfig(ctx kapi.Context, build *api.BuildConfig) error {
	_, err := s.Create(ctx, build)
	return err
}

func (s *storage) UpdateBuildConfig(ctx kapi.Context, build *api.BuildConfig) error {
	_, _, err := s.Update(ctx, build)
	return err
}

func (s *storage) DeleteBuildConfig(ctx kapi.Context, name string) error {
	_, err := s.Delete(ctx, name, nil)
	return err
}
