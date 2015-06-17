package scaler

import (
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/wait"

	"github.com/projectatomic/appinfra-next/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/deploy/util"
)

// DeploymentConfigScaler is a wrapper for the kubectl Scaler client
type DeploymentConfigScaler struct {
	c kubectl.ScalerClient
}

// ScalerFor returns the appropriate Scaler client depending on the provided
// kind of resource (Replication controllers and deploymentConfigs supported)
func ScalerFor(kind string, oc client.Interface, kc kclient.Interface) (kubectl.Scaler, error) {
	if kind != "DeploymentConfig" {
		return kubectl.ScalerFor(kind, kubectl.NewScalerClient(kc))

	}
	return &DeploymentConfigScaler{NewScalerClient(oc, kc)}, nil
}

// Scale updates a replication controller created by the DeploymentConfig with the provided namespace/name,
// to a new size, with optional precondition check (if preconditions is not nil),optional retries (if retry
//  is not nil), and then optionally waits for it's replica count to reach the new value (if wait is not nil).
func (scaler *DeploymentConfigScaler) Scale(namespace, name string, newSize uint, preconditions *kubectl.ScalePrecondition, retry, waitForReplicas *kubectl.RetryParams) error {
	if preconditions == nil {
		preconditions = &kubectl.ScalePrecondition{-1, ""}
	}
	if retry == nil {
		// Make it try only once, immediately
		retry = &kubectl.RetryParams{Interval: time.Millisecond, Timeout: time.Millisecond}
	}
	cond := kubectl.ScaleCondition(scaler, preconditions, namespace, name, newSize)
	if err := wait.Poll(retry.Interval, retry.Timeout, cond); err != nil {
		return err
	}
	if waitForReplicas != nil {
		rc := &kapi.ReplicationController{ObjectMeta: kapi.ObjectMeta{Namespace: namespace, Name: rcName}}
		return wait.Poll(waitForReplicas.Interval, waitForReplicas.Timeout,
			scaler.c.ControllerHasDesiredReplicas(rc))
	}
	return nil
}

// ScaleSimple does a simple one-shot attempt at scaling - not useful on it's own, but
// a necessary building block for Scale
func (scaler *DeploymentConfigScaler) ScaleSimple(namespace, name string, preconditions *kubectl.ScalePrecondition, newSize uint) (string, error) {
	const scaled = "scaled"
	controller, err := scaler.c.GetReplicationController(namespace, name)
	if err != nil {
		return "", kubectl.ControllerScaleError{kubectl.ControllerScaleGetFailure, "Unknown", err}
	}
	if preconditions != nil {
		if err := preconditions.Validate(controller); err != nil {
			return "", err
		}
	}
	controller.Spec.Replicas = int(newSize)
	// TODO: do retry on 409 errors here?
	if _, err := scaler.c.UpdateReplicationController(namespace, controller); err != nil {
		return "", kubectl.ControllerScaleError{kubectl.ControllerScaleUpdateFailure, controller.ResourceVersion, err}
	}
	// TODO: do a better job of printing objects here.
	return scaled, nil
}

// NewScalerClient returns a new Scaler client bundling both the OpenShift and
// Kubernetes clients
func NewScalerClient(oc client.Interface, kc kclient.Interface) kubectl.ScalerClient {
	return &realScalerClient{oc: oc, kc: kc}
}

// realScalerClient is a ScalerClient which uses an OpenShift and a Kube client.
type realScalerClient struct {
	oc client.Interface
	kc kclient.Interface
}

var rcName string

// GetReplicationController returns the most recent replication controller associated with the deploymentConfig
// with the provided namespace/name combination
func (c *realScalerClient) GetReplicationController(namespace, name string) (*kapi.ReplicationController, error) {
	dc, err := c.oc.DeploymentConfigs(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	rcName = util.LatestDeploymentNameForConfig(dc)
	return c.kc.ReplicationControllers(namespace).Get(rcName)
}

// UpdateReplicationController updates the provided replication controller
func (c *realScalerClient) UpdateReplicationController(namespace string, rc *kapi.ReplicationController) (*kapi.ReplicationController, error) {
	return c.kc.ReplicationControllers(namespace).Update(rc)
}

// ControllerHasDesiredReplicas checks whether the provided replication controller has the desired replicas
// number set
func (c *realScalerClient) ControllerHasDesiredReplicas(rc *kapi.ReplicationController) wait.ConditionFunc {
	return kclient.ControllerHasDesiredReplicas(c.kc, rc)
}
