package test

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	"github.com/projectatomic/appinfra-next/pkg/deploy/api"
)

type FakeDeploymentConfigStore struct {
	DeploymentConfig *api.DeploymentConfig
	Err              error
}

func NewFakeDeploymentConfigStore(deployment *api.DeploymentConfig) FakeDeploymentConfigStore {
	return FakeDeploymentConfigStore{DeploymentConfig: deployment}
}

func (s FakeDeploymentConfigStore) Add(obj interface{}) error {
	return s.Err
}

func (s FakeDeploymentConfigStore) Update(obj interface{}) error {
	return s.Err
}

func (s FakeDeploymentConfigStore) Delete(obj interface{}) error {
	return s.Err
}

func (s FakeDeploymentConfigStore) List() []interface{} {
	return []interface{}{s.DeploymentConfig}
}

func (s FakeDeploymentConfigStore) ContainedIDs() util.StringSet {
	return util.NewStringSet()
}

func (s FakeDeploymentConfigStore) Get(obj interface{}) (item interface{}, exists bool, err error) {
	return s.GetByKey("")
}

func (s FakeDeploymentConfigStore) GetByKey(id string) (item interface{}, exists bool, err error) {
	if s.Err != nil {
		return nil, false, err
	}
	if s.DeploymentConfig == nil {
		return nil, false, nil
	}

	return s.DeploymentConfig, true, nil
}

func (s FakeDeploymentConfigStore) Replace(list []interface{}) error {
	return nil
}
