package configchange

import (
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"

	api "github.com/projectatomic/appinfra-next/pkg/api/latest"
	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
	deployapitest "github.com/projectatomic/appinfra-next/pkg/deploy/api/test"
	deployutil "github.com/projectatomic/appinfra-next/pkg/deploy/util"
)

// TestHandle_newConfigNoTriggers ensures that a change to a config with no
// triggers doesn't result in a new config version bump.
func TestHandle_newConfigNoTriggers(t *testing.T) {
	controller := &DeploymentConfigChangeController{
		decodeConfig: func(deployment *kapi.ReplicationController) (*deployapi.DeploymentConfig, error) {
			return deployutil.DecodeDeploymentConfig(deployment, api.Codec)
		},
		changeStrategy: &changeStrategyImpl{
			generateDeploymentConfigFunc: func(namespace, name string) (*deployapi.DeploymentConfig, error) {
				t.Fatalf("unexpected generation of deploymentConfig")
				return nil, nil
			},
			updateDeploymentConfigFunc: func(namespace string, config *deployapi.DeploymentConfig) (*deployapi.DeploymentConfig, error) {
				t.Fatalf("unexpected update of deploymentConfig")
				return config, nil
			},
		},
	}

	config := deployapitest.OkDeploymentConfig(1)
	config.Triggers = []deployapi.DeploymentTriggerPolicy{}
	err := controller.Handle(config)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestHandle_newConfigTriggers ensures that the creation of a new config
// (with version 0) with a config change trigger results in a version bump and
// cause update for initial deployment.
func TestHandle_newConfigTriggers(t *testing.T) {
	var updated *deployapi.DeploymentConfig

	controller := &DeploymentConfigChangeController{
		decodeConfig: func(deployment *kapi.ReplicationController) (*deployapi.DeploymentConfig, error) {
			return deployutil.DecodeDeploymentConfig(deployment, api.Codec)
		},
		changeStrategy: &changeStrategyImpl{
			generateDeploymentConfigFunc: func(namespace, name string) (*deployapi.DeploymentConfig, error) {
				return deployapitest.OkDeploymentConfig(1), nil
			},
			updateDeploymentConfigFunc: func(namespace string, config *deployapi.DeploymentConfig) (*deployapi.DeploymentConfig, error) {
				updated = config
				return config, nil
			},
		},
	}

	config := deployapitest.OkDeploymentConfig(0)
	config.Triggers = []deployapi.DeploymentTriggerPolicy{deployapitest.OkConfigChangeTrigger()}
	err := controller.Handle(config)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated == nil {
		t.Fatalf("expected config to be updated")
	}

	if e, a := 1, updated.LatestVersion; e != a {
		t.Fatalf("expected update to latestversion=%d, got %d", e, a)
	}

	if updated.Details == nil {
		t.Fatalf("expected config change details to be set")
	} else if updated.Details.Causes == nil {
		t.Fatalf("expected config change causes to be set")
	} else if updated.Details.Causes[0].Type != deployapi.DeploymentTriggerOnConfigChange {
		t.Fatalf("expected config change cause to be set to config change trigger, got %s", updated.Details.Causes[0].Type)
	}
}

// TestHandle_changeWithTemplateDiff ensures that a pod template change to a
// config with a config change trigger results in a version bump and cause
// update.
func TestHandle_changeWithTemplateDiff(t *testing.T) {
	var updated *deployapi.DeploymentConfig

	controller := &DeploymentConfigChangeController{
		decodeConfig: func(deployment *kapi.ReplicationController) (*deployapi.DeploymentConfig, error) {
			return deployutil.DecodeDeploymentConfig(deployment, api.Codec)
		},
		changeStrategy: &changeStrategyImpl{
			generateDeploymentConfigFunc: func(namespace, name string) (*deployapi.DeploymentConfig, error) {
				return deployapitest.OkDeploymentConfig(2), nil
			},
			updateDeploymentConfigFunc: func(namespace string, config *deployapi.DeploymentConfig) (*deployapi.DeploymentConfig, error) {
				updated = config
				return config, nil
			},
			getDeploymentFunc: func(namespace, name string) (*kapi.ReplicationController, error) {
				deployment, _ := deployutil.MakeDeployment(deployapitest.OkDeploymentConfig(1), kapi.Codec)
				return deployment, nil
			},
		},
	}

	config := deployapitest.OkDeploymentConfig(1)
	config.Triggers = []deployapi.DeploymentTriggerPolicy{deployapitest.OkConfigChangeTrigger()}
	config.Template.ControllerTemplate.Template.Spec.Containers[1].Name = "modified"
	err := controller.Handle(config)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated == nil {
		t.Fatalf("expected config to be updated")
	}

	if e, a := 2, updated.LatestVersion; e != a {
		t.Fatalf("expected update to latestversion=%d, got %d", e, a)
	}

	if updated.Details == nil {
		t.Fatalf("expected config change details to be set")
	} else if updated.Details.Causes == nil {
		t.Fatalf("expected config change causes to be set")
	} else if updated.Details.Causes[0].Type != deployapi.DeploymentTriggerOnConfigChange {
		t.Fatalf("expected config change cause to be set to config change trigger, got %s", updated.Details.Causes[0].Type)
	}
}

// TestHandle_changeWithoutTemplateDiff ensures that an updated config with no
// pod template diff results in the config version remaining the same.
func TestHandle_changeWithoutTemplateDiff(t *testing.T) {
	config := deployapitest.OkDeploymentConfig(1)
	config.Triggers = []deployapi.DeploymentTriggerPolicy{deployapitest.OkConfigChangeTrigger()}

	generated := false
	updated := false

	controller := &DeploymentConfigChangeController{
		decodeConfig: func(deployment *kapi.ReplicationController) (*deployapi.DeploymentConfig, error) {
			return deployutil.DecodeDeploymentConfig(deployment, api.Codec)
		},
		changeStrategy: &changeStrategyImpl{
			generateDeploymentConfigFunc: func(namespace, name string) (*deployapi.DeploymentConfig, error) {
				generated = true
				return config, nil
			},
			updateDeploymentConfigFunc: func(namespace string, config *deployapi.DeploymentConfig) (*deployapi.DeploymentConfig, error) {
				updated = true
				return config, nil
			},
			getDeploymentFunc: func(namespace, name string) (*kapi.ReplicationController, error) {
				deployment, _ := deployutil.MakeDeployment(deployapitest.OkDeploymentConfig(1), kapi.Codec)
				return deployment, nil
			},
		},
	}

	err := controller.Handle(config)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if generated {
		t.Error("Unexpected generation of deploymentConfig")
	}

	if updated {
		t.Error("Unexpected update of deploymentConfig")
	}
}
