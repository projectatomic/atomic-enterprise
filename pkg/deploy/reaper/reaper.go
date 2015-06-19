package reaper

import (
	"fmt"
	"time"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kerrors "github.com/GoogleCloudPlatform/kubernetes/pkg/api/errors"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/kubectl"
	"github.com/golang/glog"

	"github.com/projectatomic/appinfra-next/pkg/client"
	"github.com/projectatomic/appinfra-next/pkg/deploy/util"
)

// ReaperFor returns the appropriate Reaper client depending on the provided
// kind of resource (Replication controllers, pods, services, and deploymentConfigs
// supported)
func ReaperFor(kind string, oc *client.Client, kc *kclient.Client) (kubectl.Reaper, error) {
	if kind != "DeploymentConfig" {
		return kubectl.ReaperFor(kind, kc)
	}
	return &DeploymentConfigReaper{oc: oc, kc: kc, pollInterval: kubectl.Interval, timeout: kubectl.Timeout}, nil
}

// DeploymentConfigReaper implements the Reaper interface for deploymentConfigs
type DeploymentConfigReaper struct {
	oc                    client.Interface
	kc                    kclient.Interface
	pollInterval, timeout time.Duration
}

// Stop scales a replication controller via its deployment configuration down to
// zero replicas, waits for all of them to get deleted and then deletes both the
// replication controller and its deployment configuration.
func (reaper *DeploymentConfigReaper) Stop(namespace, name string, gracePeriod *kapi.DeleteOptions) (string, error) {
	// If the config is already deleted, it may still have associated
	// deployments which didn't get cleaned up during prior calls to Stop. If
	// the config can't be found, still make an attempt to clean up the
	// deployments.
	//
	// It's important to delete the config first to avoid an undesirable side
	// effect which can cause the deployment to be re-triggered upon the
	// config's deletion. See https://github.com/projectatomic/appinfra-next/issues/2721
	// for more details.
	err := reaper.oc.DeploymentConfigs(namespace).Delete(name)
	configNotFound := kerrors.IsNotFound(err)
	if err != nil && !configNotFound {
		return "", err
	}

	// Clean up deployments related to the config.
	rcList, err := reaper.kc.ReplicationControllers(namespace).List(util.ConfigSelector(name))
	if err != nil {
		return "", err
	}
	rcReaper, err := kubectl.ReaperFor("ReplicationController", reaper.kc)
	if err != nil {
		return "", err
	}

	// If there is neither a config nor any deployments, we can return NotFound.
	deployments := rcList.Items
	if configNotFound && len(deployments) == 0 {
		return "", kerrors.NewNotFound("DeploymentConfig", name)
	}
	for _, rc := range deployments {
		if _, err = rcReaper.Stop(rc.Namespace, rc.Name, gracePeriod); err != nil {
			// Better not error out here...
			glog.Infof("Cannot delete ReplicationController %s/%s: %v", rc.Namespace, rc.Name, err)
		}
	}

	return fmt.Sprintf("%s stopped", name), nil
}
