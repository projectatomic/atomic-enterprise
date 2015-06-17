package prune

import (
	"sort"
	"testing"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
)

type mockPruneRecorder struct {
	set util.StringSet
	err error
}

func (m *mockPruneRecorder) Handler(deployment *kapi.ReplicationController) error {
	m.set.Insert(deployment.Name)
	return m.err
}

func (m *mockPruneRecorder) Verify(t *testing.T, expected util.StringSet) {
	if len(m.set) != len(expected) || !m.set.HasAll(expected.List()...) {
		expectedValues := expected.List()
		actualValues := m.set.List()
		sort.Strings(expectedValues)
		sort.Strings(actualValues)
		t.Errorf("expected \n\t%v\n, actual \n\t%v\n", expectedValues, actualValues)
	}
}

func TestPruneTask(t *testing.T) {
	deploymentStatusOptions := []deployapi.DeploymentStatus{
		deployapi.DeploymentStatusComplete,
		deployapi.DeploymentStatusFailed,
		deployapi.DeploymentStatusNew,
		deployapi.DeploymentStatusPending,
		deployapi.DeploymentStatusRunning,
	}
	deploymentStatusFilter := []deployapi.DeploymentStatus{
		deployapi.DeploymentStatusComplete,
		deployapi.DeploymentStatusFailed,
	}
	deploymentStatusFilterSet := util.StringSet{}
	for _, deploymentStatus := range deploymentStatusFilter {
		deploymentStatusFilterSet.Insert(string(deploymentStatus))
	}

	for _, orphans := range []bool{true, false} {
		for _, deploymentStatusOption := range deploymentStatusOptions {
			keepYoungerThan := time.Hour

			now := util.Now()
			old := util.NewTime(now.Time.Add(-1 * keepYoungerThan))

			deploymentConfigs := []*deployapi.DeploymentConfig{}
			deployments := []*kapi.ReplicationController{}

			deploymentConfig := mockDeploymentConfig("a", "deployment-config")
			deploymentConfigs = append(deploymentConfigs, deploymentConfig)

			deployments = append(deployments, withCreated(withStatus(mockDeployment("a", "build-1", deploymentConfig), deploymentStatusOption), now))
			deployments = append(deployments, withCreated(withStatus(mockDeployment("a", "build-2", deploymentConfig), deploymentStatusOption), old))
			deployments = append(deployments, withSize(withCreated(withStatus(mockDeployment("a", "build-3-with-replicas", deploymentConfig), deploymentStatusOption), old), 4))
			deployments = append(deployments, withCreated(withStatus(mockDeployment("a", "orphan-build-1", nil), deploymentStatusOption), now))
			deployments = append(deployments, withCreated(withStatus(mockDeployment("a", "orphan-build-2", nil), deploymentStatusOption), old))
			deployments = append(deployments, withSize(withCreated(withStatus(mockDeployment("a", "orphan-build-3-with-replicas", nil), deploymentStatusOption), old), 4))

			keepComplete := 1
			keepFailed := 1
			expectedValues := util.StringSet{}
			filter := &andFilter{
				filterPredicates: []FilterPredicate{
					FilterDeploymentsPredicate,
					FilterZeroReplicaSize,
					NewFilterBeforePredicate(keepYoungerThan),
				},
			}
			dataSet := NewDataSet(deploymentConfigs, filter.Filter(deployments))
			resolver := NewPerDeploymentConfigResolver(dataSet, keepComplete, keepFailed)
			if orphans {
				resolver = &mergeResolver{
					resolvers: []Resolver{resolver, NewOrphanDeploymentResolver(dataSet, deploymentStatusFilter)},
				}
			}
			expectedDeployments, err := resolver.Resolve()
			for _, item := range expectedDeployments {
				expectedValues.Insert(item.Name)
			}

			recorder := &mockPruneRecorder{set: util.StringSet{}}
			task := NewPruneTasker(deploymentConfigs, deployments, keepYoungerThan, orphans, keepComplete, keepFailed, recorder.Handler)
			err = task.PruneTask()
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			recorder.Verify(t, expectedValues)
		}
	}

}
