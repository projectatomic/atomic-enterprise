package validation

import (
	"github.com/projectatomic/appinfra-next/pkg/cmd/server/api"
)

func ValidateAllInOneConfig(master *api.MasterConfig, node *api.NodeConfig) ValidationResults {
	validationResults := ValidationResults{}

	validationResults.Append(ValidateMasterConfig(master).Prefix("masterConfig"))

	validationResults.AddErrors(ValidateNodeConfig(node).Prefix("nodeConfig")...)

	// Validation between the configs

	return validationResults
}
