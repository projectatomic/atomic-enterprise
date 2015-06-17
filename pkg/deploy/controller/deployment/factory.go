package deployment

import (
	"fmt"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/cache"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/record"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	kutil "github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"

	controller "github.com/projectatomic/appinfra-next/pkg/controller"
	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
	deployutil "github.com/projectatomic/appinfra-next/pkg/deploy/util"
)

// DeploymentControllerFactory can create a DeploymentController that creates
// deployer pods in a configurable way.
type DeploymentControllerFactory struct {
	// KubeClient is a Kubernetes client.
	KubeClient kclient.Interface
	// Codec is used for encoding/decoding.
	Codec runtime.Codec
	// ServiceAccount is the service account name to run deployer pods as
	ServiceAccount string
	// Environment is a set of environment which should be injected into all deployer pod containers.
	Environment []kapi.EnvVar
	// DeployerImage specifies which Docker image can support the default strategies.
	DeployerImage string
}

// Create creates a DeploymentController.
func (factory *DeploymentControllerFactory) Create() controller.RunnableController {
	deploymentLW := &deployutil.ListWatcherImpl{
		// TODO: Investigate specifying annotation field selectors to fetch only 'deployments'
		// Currently field selectors are not supported for replication controllers
		ListFunc: func() (runtime.Object, error) {
			return factory.KubeClient.ReplicationControllers(kapi.NamespaceAll).List(labels.Everything())
		},
		WatchFunc: func(resourceVersion string) (watch.Interface, error) {
			return factory.KubeClient.ReplicationControllers(kapi.NamespaceAll).Watch(labels.Everything(), fields.Everything(), resourceVersion)
		},
	}
	deploymentQueue := cache.NewFIFO(cache.MetaNamespaceKeyFunc)
	cache.NewReflector(deploymentLW, &kapi.ReplicationController{}, deploymentQueue, 2*time.Minute).Run()

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(factory.KubeClient.Events(""))

	deployController := &DeploymentController{
		serviceAccount: factory.ServiceAccount,
		deploymentClient: &deploymentClientImpl{
			getDeploymentFunc: func(namespace, name string) (*kapi.ReplicationController, error) {
				return factory.KubeClient.ReplicationControllers(namespace).Get(name)
			},
			updateDeploymentFunc: func(namespace string, deployment *kapi.ReplicationController) (*kapi.ReplicationController, error) {
				return factory.KubeClient.ReplicationControllers(namespace).Update(deployment)
			},
		},
		podClient: &podClientImpl{
			getPodFunc: func(namespace, name string) (*kapi.Pod, error) {
				return factory.KubeClient.Pods(namespace).Get(name)
			},
			createPodFunc: func(namespace string, pod *kapi.Pod) (*kapi.Pod, error) {
				return factory.KubeClient.Pods(namespace).Create(pod)
			},
			deletePodFunc: func(namespace, name string) error {
				return factory.KubeClient.Pods(namespace).Delete(name, nil)
			},
			updatePodFunc: func(namespace string, pod *kapi.Pod) (*kapi.Pod, error) {
				return factory.KubeClient.Pods(namespace).Update(pod)
			},
			// Find deployer pods using the label they should all have which
			// correlates them to the named deployment.
			getDeployerPodsForFunc: func(namespace, name string) ([]kapi.Pod, error) {
				labelSel, err := labels.Parse(fmt.Sprintf("%s=%s", deployapi.DeployerPodForDeploymentLabel, name))
				if err != nil {
					return []kapi.Pod{}, err
				}
				pods, err := factory.KubeClient.Pods(namespace).List(labelSel, fields.Everything())
				if err != nil {
					return []kapi.Pod{}, err
				}
				return pods.Items, nil
			},
		},
		makeContainer: func(strategy *deployapi.DeploymentStrategy) (*kapi.Container, error) {
			return factory.makeContainer(strategy)
		},
		decodeConfig: func(deployment *kapi.ReplicationController) (*deployapi.DeploymentConfig, error) {
			return deployutil.DecodeDeploymentConfig(deployment, factory.Codec)
		},
		recorder: eventBroadcaster.NewRecorder(kapi.EventSource{Component: "deployer"}),
	}

	return &controller.RetryController{
		Queue: deploymentQueue,
		RetryManager: controller.NewQueueRetryManager(
			deploymentQueue,
			cache.MetaNamespaceKeyFunc,
			func(obj interface{}, err error, retries controller.Retry) bool {
				if _, isFatal := err.(fatalError); isFatal {
					kutil.HandleError(err)
					return false
				}
				if retries.Count > 1 {
					return false
				}
				return true
			},
			kutil.NewTokenBucketRateLimiter(1, 10),
		),
		Handle: func(obj interface{}) error {
			deployment := obj.(*kapi.ReplicationController)
			return deployController.Handle(deployment)
		},
	}
}

// makeContainer creates containers in the following way:
//
//   1. For the Recreate and Rolling strategies, strategy, use the factory's
//      DeployerImage as the container image, and the factory's Environment
//      as the container environment.
//   2. For all Custom strategy, use the strategy's image for the container
//      image, and use the combination of the factory's Environment and the
//      strategy's environment as the container environment.
//
// An error is returned if the deployment strategy type is not supported.
func (factory *DeploymentControllerFactory) makeContainer(strategy *deployapi.DeploymentStrategy) (*kapi.Container, error) {
	// Set default environment values
	environment := []kapi.EnvVar{}
	for _, env := range factory.Environment {
		environment = append(environment, env)
	}

	// Every strategy type should be handled here.
	switch strategy.Type {
	case deployapi.DeploymentStrategyTypeRecreate, deployapi.DeploymentStrategyTypeRolling:
		// Use the factory-configured image.
		return &kapi.Container{
			Image: factory.DeployerImage,
			Env:   environment,
		}, nil
	case deployapi.DeploymentStrategyTypeCustom:
		// Use user-defined values from the strategy input.
		for _, env := range strategy.CustomParams.Environment {
			environment = append(environment, env)
		}
		return &kapi.Container{
			Image: strategy.CustomParams.Image,
			Env:   environment,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported deployment strategy type: %s", strategy.Type)
	}
}
