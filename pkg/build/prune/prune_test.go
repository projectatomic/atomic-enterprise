package prune

import (
	"sort"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

type mockPruneRecorder struct {
	set util.StringSet
	err error
}

func (m *mockPruneRecorder) Handler(build *buildapi.Build) error {
	m.set.Insert(build.Name)
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
	buildStatusOptions := []buildapi.BuildStatus{
		buildapi.BuildStatusCancelled,
		buildapi.BuildStatusComplete,
		buildapi.BuildStatusError,
		buildapi.BuildStatusFailed,
		buildapi.BuildStatusNew,
		buildapi.BuildStatusPending,
		buildapi.BuildStatusRunning,
	}
	buildStatusFilter := []buildapi.BuildStatus{
		buildapi.BuildStatusCancelled,
		buildapi.BuildStatusComplete,
		buildapi.BuildStatusError,
		buildapi.BuildStatusFailed,
	}
	buildStatusFilterSet := util.StringSet{}
	for _, buildStatus := range buildStatusFilter {
		buildStatusFilterSet.Insert(string(buildStatus))
	}

	for _, orphans := range []bool{true, false} {
		for _, buildStatusOption := range buildStatusOptions {
			keepYoungerThan := time.Hour

			now := util.Now()
			old := util.NewTime(now.Time.Add(-1 * keepYoungerThan))

			buildConfigs := []*buildapi.BuildConfig{}
			builds := []*buildapi.Build{}

			buildConfig := mockBuildConfig("a", "build-config")
			buildConfigs = append(buildConfigs, buildConfig)

			builds = append(builds, withCreated(withStatus(mockBuild("a", "build-1", buildConfig), buildStatusOption), now))
			builds = append(builds, withCreated(withStatus(mockBuild("a", "build-2", buildConfig), buildStatusOption), old))
			builds = append(builds, withCreated(withStatus(mockBuild("a", "orphan-build-1", nil), buildStatusOption), now))
			builds = append(builds, withCreated(withStatus(mockBuild("a", "orphan-build-2", nil), buildStatusOption), old))

			keepComplete := 1
			keepFailed := 1
			expectedValues := util.StringSet{}
			filter := &andFilter{
				filterPredicates: []FilterPredicate{NewFilterBeforePredicate(keepYoungerThan)},
			}
			dataSet := NewDataSet(buildConfigs, filter.Filter(builds))
			resolver := NewPerBuildConfigResolver(dataSet, keepComplete, keepFailed)
			if orphans {
				resolver = &mergeResolver{
					resolvers: []Resolver{resolver, NewOrphanBuildResolver(dataSet, buildStatusFilter)},
				}
			}
			expectedBuilds, err := resolver.Resolve()
			for _, build := range expectedBuilds {
				expectedValues.Insert(build.Name)
			}

			recorder := &mockPruneRecorder{set: util.StringSet{}}
			task := NewPruneTasker(buildConfigs, builds, keepYoungerThan, orphans, keepComplete, keepFailed, recorder.Handler)
			err = task.PruneTask()
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			recorder.Verify(t, expectedValues)
		}
	}

}
