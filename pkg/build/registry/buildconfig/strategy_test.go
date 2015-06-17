package buildconfig

import (
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

func TestBuildConfigStrategy(t *testing.T) {
	ctx := kapi.NewDefaultContext()
	if !Strategy.NamespaceScoped() {
		t.Errorf("BuildConfig is namespace scoped")
	}
	if Strategy.AllowCreateOnUpdate() {
		t.Errorf("BuildConfig should not allow create on update")
	}
	buildConfig := &buildapi.BuildConfig{
		ObjectMeta: kapi.ObjectMeta{Name: "config-id", Namespace: "namespace"},
		Parameters: buildapi.BuildParameters{
			Source: buildapi.BuildSource{
				Type: buildapi.BuildSourceGit,
				Git: &buildapi.GitBuildSource{
					URI: "http://github.com/my/repository",
				},
				ContextDir: "context",
			},
			Strategy: buildapi.BuildStrategy{
				Type:           buildapi.DockerBuildStrategyType,
				DockerStrategy: &buildapi.DockerBuildStrategy{},
			},
			Output: buildapi.BuildOutput{
				DockerImageReference: "repository/data",
			},
		},
	}
	Strategy.PrepareForCreate(buildConfig)
	errs := Strategy.Validate(ctx, buildConfig)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %v", errs)
	}

	buildConfig.ResourceVersion = "foo"
	errs = Strategy.ValidateUpdate(ctx, buildConfig, buildConfig)
	if len(errs) != 0 {
		t.Errorf("Unexpected error validating %v", errs)
	}
	invalidBuildConfig := &buildapi.BuildConfig{}
	errs = Strategy.Validate(ctx, invalidBuildConfig)
	if len(errs) == 0 {
		t.Errorf("Expected error validating")
	}
}
