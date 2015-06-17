package strategy

import (
	"errors"
	"fmt"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/golang/glog"

	buildapi "github.com/projectatomic/appinfra-next/pkg/build/api"
	buildutil "github.com/projectatomic/appinfra-next/pkg/build/util"
)

// CustomBuildStrategy creates a build using a custom builder image.
type CustomBuildStrategy struct {
	// Codec is the codec to use for encoding the output pod.
	// IMPORTANT: This may break backwards compatibility when
	// it changes.
	Codec runtime.Codec
}

// CreateBuildPod creates the pod to be used for the Custom build
func (bs *CustomBuildStrategy) CreateBuildPod(build *buildapi.Build) (*kapi.Pod, error) {
	data, err := bs.Codec.Encode(build)
	if err != nil {
		return nil, fmt.Errorf("failed to encode the build: %v", err)
	}

	strategy := build.Parameters.Strategy.CustomStrategy
	containerEnv := []kapi.EnvVar{{Name: "BUILD", Value: string(data)}}

	if build.Parameters.Source.Git != nil {
		containerEnv = append(containerEnv, kapi.EnvVar{
			Name: "SOURCE_REPOSITORY", Value: build.Parameters.Source.Git.URI,
		})
	}

	if strategy == nil || strategy.From == nil || len(strategy.From.Name) == 0 {
		return nil, errors.New("CustomBuildStrategy cannot be executed without image")
	}

	if len(strategy.Env) > 0 {
		containerEnv = append(containerEnv, strategy.Env...)
	}

	if strategy.ExposeDockerSocket {
		glog.V(2).Infof("ExposeDockerSocket is enabled for %s build", build.Name)
		containerEnv = append(containerEnv, kapi.EnvVar{Name: "DOCKER_SOCKET", Value: dockerSocketPath})
	}

	privileged := true
	pod := &kapi.Pod{
		ObjectMeta: kapi.ObjectMeta{
			Name:      buildutil.GetBuildPodName(build),
			Namespace: build.Namespace,
			Labels:    getPodLabels(build),
		},
		Spec: kapi.PodSpec{
			ServiceAccount: build.Parameters.ServiceAccount,
			Containers: []kapi.Container{
				{
					Name:  "custom-build",
					Image: strategy.From.Name,
					Env:   containerEnv,
					// TODO: run unprivileged https://github.com/projectatomic/appinfra-next/issues/662
					SecurityContext: &kapi.SecurityContext{
						Privileged: &privileged,
					},
				},
			},
			RestartPolicy: kapi.RestartPolicyNever,
		},
	}

	if err := setupBuildEnv(build, pod); err != nil {
		return nil, err
	}

	pod.Spec.Containers[0].ImagePullPolicy = kapi.PullIfNotPresent
	pod.Spec.Containers[0].Resources = build.Parameters.Resources

	if strategy.ExposeDockerSocket {
		setupDockerSocket(pod)
		setupDockerSecrets(pod, build.Parameters.Output.PushSecret, strategy.PullSecret)
	}
	setupSourceSecrets(pod, build.Parameters.Source.SourceSecret)
	return pod, nil
}
