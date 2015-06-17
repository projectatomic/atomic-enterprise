package builder

import (
	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
)

// getBuildEnvVars returns a map with the environment variables that should be added
// to the built image
func getBuildEnvVars(build *buildapi.Build) map[string]string {
	envVars := map[string]string{
		"OPENSHIFT_BUILD_NAME":      build.Name,
		"OPENSHIFT_BUILD_NAMESPACE": build.Namespace,
		"OPENSHIFT_BUILD_SOURCE":    build.Parameters.Source.Git.URI,
	}
	if build.Parameters.Source.Git.Ref != "" {
		envVars["OPENSHIFT_BUILD_REFERENCE"] = build.Parameters.Source.Git.Ref
	}
	if build.Parameters.Revision != nil &&
		build.Parameters.Revision.Git != nil &&
		build.Parameters.Revision.Git.Commit != "" {
		envVars["OPENSHIFT_BUILD_COMMIT"] = build.Parameters.Revision.Git.Commit
	}
	if build.Parameters.Strategy.Type == buildapi.SourceBuildStrategyType {
		userEnv := build.Parameters.Strategy.SourceStrategy.Env
		for _, v := range userEnv {
			envVars[v.Name] = v.Value
		}
	}
	return envVars
}
