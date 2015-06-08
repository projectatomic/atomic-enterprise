package origin

import (
	configapi "github.com/openshift/origin/pkg/cmd/server/api"
)

// MasterConfig defines the required parameters for starting the OpenShift master
type AssetConfig struct {
	Options configapi.AssetConfig

	// TODO: possibly change to point to MasterConfig's version
	OpenshiftEnabled bool
}

func BuildAssetConfig(options configapi.MasterConfig) (*AssetConfig, error) {
	return &AssetConfig{
		Options:          *options.AssetConfig,
		OpenshiftEnabled: options.OpenshiftEnabled,
	}, nil
}
