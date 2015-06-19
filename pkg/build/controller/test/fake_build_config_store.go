package test

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

type FakeBuildConfigStore struct {
	Build *buildapi.BuildConfig
	Err   error
}

func NewFakeBuildConfigStore(build *buildapi.BuildConfig) FakeBuildConfigStore {
	return FakeBuildConfigStore{Build: build}
}

func (s FakeBuildConfigStore) Add(obj interface{}) error {
	return s.Err
}

func (s FakeBuildConfigStore) Update(obj interface{}) error {
	return s.Err
}

func (s FakeBuildConfigStore) Delete(obj interface{}) error {
	return s.Err
}

func (s FakeBuildConfigStore) List() []interface{} {
	return []interface{}{s.Build}
}

func (s FakeBuildConfigStore) ListKeys() []string {
	return []string{"config"}
}

func (s FakeBuildConfigStore) ContainedIDs() util.StringSet {
	return util.NewStringSet()
}

func (s FakeBuildConfigStore) Get(obj interface{}) (item interface{}, exists bool, err error) {
	return s.GetByKey("")
}

func (s FakeBuildConfigStore) GetByKey(id string) (item interface{}, exists bool, err error) {
	if s.Err != nil {
		return nil, false, err
	}
	if s.Build == nil {
		return nil, false, nil
	}

	return s.Build, true, nil
}

func (s FakeBuildConfigStore) Replace(list []interface{}) error {
	return nil
}
