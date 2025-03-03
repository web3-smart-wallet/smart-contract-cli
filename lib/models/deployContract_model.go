package models

import "github.com/web3-smart-wallet/smart-contract-cli/lib/services"

// deployContract represents the data for the deployContract page
type DeployContractModel struct {
	// ShowError bool
	URI                 string
	AvailableContracts  []services.AvailableContract
	SelectedContract    int  // Index of the selected contract, -1 if none selected
	IsSelectingContract bool // Whether we're in contract selection mode
}

// NewDeployContractModel creates a new deployContract model
func NewDeployContractModel() *DeployContractModel {
	return &DeployContractModel{
		// ShowError: false,
		URI:                 "",
		AvailableContracts:  []services.AvailableContract{},
		SelectedContract:    -1,
		IsSelectingContract: true,
	}
}
