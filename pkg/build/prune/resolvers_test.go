package prune

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

type mockResolver struct {
	builds []*buildapi.Build
	err    error
}

func (m *mockResolver) Resolve() ([]*buildapi.Build, error) {
	return m.builds, m.err
}

func TestMergeResolver(t *testing.T) {
	resolverA := &mockResolver{
		builds: []*buildapi.Build{
			mockBuild("a", "b", nil),
		},
	}
	resolverB := &mockResolver{
		builds: []*buildapi.Build{
			mockBuild("c", "d", nil),
		},
	}
	resolver := &mergeResolver{resolvers: []Resolver{resolverA, resolverB}}
	results, err := resolver.Resolve()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Unexpected results %v", results)
	}
	expectedNames := util.NewStringSet("b", "d")
	for _, build := range results {
		if !expectedNames.Has(build.Name) {
			t.Errorf("Unexpected name %v", build.Name)
		}
	}
}

func TestOrphanBuildResolver(t *testing.T) {
	activeBuildConfig := mockBuildConfig("a", "active-build-config")
	inactiveBuildConfig := mockBuildConfig("a", "inactive-build-config")

	buildConfigs := []*buildapi.BuildConfig{activeBuildConfig}
	builds := []*buildapi.Build{}

	expectedNames := util.StringSet{}
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

	for _, buildStatusOption := range buildStatusOptions {
		builds = append(builds, withStatus(mockBuild("a", string(buildStatusOption)+"-active", activeBuildConfig), buildStatusOption))
		builds = append(builds, withStatus(mockBuild("a", string(buildStatusOption)+"-inactive", inactiveBuildConfig), buildStatusOption))
		builds = append(builds, withStatus(mockBuild("a", string(buildStatusOption)+"-orphan", nil), buildStatusOption))
		if buildStatusFilterSet.Has(string(buildStatusOption)) {
			expectedNames.Insert(string(buildStatusOption) + "-inactive")
			expectedNames.Insert(string(buildStatusOption) + "-orphan")
		}
	}

	dataSet := NewDataSet(buildConfigs, builds)
	resolver := NewOrphanBuildResolver(dataSet, buildStatusFilter)
	results, err := resolver.Resolve()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	foundNames := util.StringSet{}
	for _, result := range results {
		foundNames.Insert(result.Name)
	}
	if len(foundNames) != len(expectedNames) || !expectedNames.HasAll(foundNames.List()...) {
		t.Errorf("expected %v, actual %v", expectedNames, foundNames)
	}
}

func TestPerBuildConfigResolver(t *testing.T) {
	buildStatusOptions := []buildapi.BuildStatus{
		buildapi.BuildStatusCancelled,
		buildapi.BuildStatusComplete,
		buildapi.BuildStatusError,
		buildapi.BuildStatusFailed,
		buildapi.BuildStatusNew,
		buildapi.BuildStatusPending,
		buildapi.BuildStatusRunning,
	}
	buildConfigs := []*buildapi.BuildConfig{
		mockBuildConfig("a", "build-config-1"),
		mockBuildConfig("b", "build-config-2"),
	}
	buildsPerStatus := 100
	builds := []*buildapi.Build{}
	for _, buildConfig := range buildConfigs {
		for _, buildStatusOption := range buildStatusOptions {
			for i := 0; i < buildsPerStatus; i++ {
				build := withStatus(mockBuild(buildConfig.Namespace, fmt.Sprintf("%v-%v-%v", buildConfig.Name, buildStatusOption, i), buildConfig), buildStatusOption)
				builds = append(builds, build)
			}
		}
	}

	now := util.Now()
	for i := range builds {
		creationTimestamp := util.NewTime(now.Time.Add(-1 * time.Duration(i) * time.Hour))
		builds[i].CreationTimestamp = creationTimestamp
	}

	// test number to keep at varying ranges
	for keep := 0; keep < buildsPerStatus*2; keep++ {
		dataSet := NewDataSet(buildConfigs, builds)

		expectedNames := util.StringSet{}
		buildCompleteStatusFilterSet := util.NewStringSet(string(buildapi.BuildStatusComplete))
		buildFailedStatusFilterSet := util.NewStringSet(string(buildapi.BuildStatusCancelled), string(buildapi.BuildStatusError), string(buildapi.BuildStatusFailed))

		for _, buildConfig := range buildConfigs {
			buildItems, err := dataSet.ListBuildsByBuildConfig(buildConfig)
			if err != nil {
				t.Errorf("Unexpected err %v", err)
			}
			completedBuilds, failedBuilds := []*buildapi.Build{}, []*buildapi.Build{}
			for _, build := range buildItems {
				if buildCompleteStatusFilterSet.Has(string(build.Status)) {
					completedBuilds = append(completedBuilds, build)
				} else if buildFailedStatusFilterSet.Has(string(build.Status)) {
					failedBuilds = append(failedBuilds, build)
				}
			}
			sort.Sort(sortableBuilds(completedBuilds))
			sort.Sort(sortableBuilds(failedBuilds))
			purgeCompleted := []*buildapi.Build{}
			purgeFailed := []*buildapi.Build{}
			if keep >= 0 && keep < len(completedBuilds) {
				purgeCompleted = completedBuilds[keep:]
			}
			if keep >= 0 && keep < len(failedBuilds) {
				purgeFailed = failedBuilds[keep:]
			}
			for _, build := range purgeCompleted {
				expectedNames.Insert(build.Name)
			}
			for _, build := range purgeFailed {
				expectedNames.Insert(build.Name)
			}
		}

		resolver := NewPerBuildConfigResolver(dataSet, keep, keep)
		results, err := resolver.Resolve()
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		foundNames := util.StringSet{}
		for _, result := range results {
			foundNames.Insert(result.Name)
		}
		if len(foundNames) != len(expectedNames) || !expectedNames.HasAll(foundNames.List()...) {
			expectedValues := expectedNames.List()
			actualValues := foundNames.List()
			sort.Strings(expectedValues)
			sort.Strings(actualValues)
			t.Errorf("keep %v\n, expected \n\t%v\n, actual \n\t%v\n", keep, expectedValues, actualValues)
		}
	}
}
