package ipfailover

import (
	"io"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"

	deployapi "github.com/projectatomic/appinfra-next/pkg/deploy/api"
)

type IPFailoverConfiguratorPlugin interface {
	GetWatchPort() (int, error)
	GetSelector() (map[string]string, error)
	GetNamespace() (string, error)
	GetDeploymentConfig() (*deployapi.DeploymentConfig, error)
	Generate() (*kapi.List, error)
	Create(out io.Writer) error
}
