package test

import (
	"sync"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

type BuildRegistry struct {
	Err            error
	Builds         *buildapi.BuildList
	Build          *buildapi.Build
	DeletedBuildID string
	sync.Mutex
}

func (r *BuildRegistry) ListBuilds(ctx kapi.Context, labels labels.Selector, field fields.Selector) (*buildapi.BuildList, error) {
	r.Lock()
	defer r.Unlock()
	return r.Builds, r.Err
}

func (r *BuildRegistry) GetBuild(ctx kapi.Context, id string) (*buildapi.Build, error) {
	r.Lock()
	defer r.Unlock()
	return r.Build, r.Err
}

func (r *BuildRegistry) CreateBuild(ctx kapi.Context, build *buildapi.Build) error {
	r.Lock()
	defer r.Unlock()
	r.Build = build
	return r.Err
}

func (r *BuildRegistry) UpdateBuild(ctx kapi.Context, build *buildapi.Build) error {
	r.Lock()
	defer r.Unlock()
	r.Build = build
	return r.Err
}

func (r *BuildRegistry) DeleteBuild(ctx kapi.Context, id string) error {
	r.Lock()
	defer r.Unlock()
	r.DeletedBuildID = id
	r.Build = nil
	return r.Err
}

func (r *BuildRegistry) WatchBuilds(ctx kapi.Context, label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return nil, r.Err
}
