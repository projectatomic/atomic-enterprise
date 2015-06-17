package origin

import (
	configapi "github.com/projectatomic/appinfra-next/pkg/cmd/server/api"
)

// AssetConfig defines the required parameters for starting the OpenShift master
type AssetConfig struct {
	Options configapi.AssetConfig

	// TODO: possibly change to point to MasterConfig's version
	OpenshiftEnabled bool
}

// BuildAssetConfig returns a new AssetConfig
func BuildAssetConfig(options configapi.MasterConfig) (*AssetConfig, error) {
	return &AssetConfig{
		Options:          *options.AssetConfig,
		OpenshiftEnabled: options.OpenshiftEnabled,
	}, nil
}
