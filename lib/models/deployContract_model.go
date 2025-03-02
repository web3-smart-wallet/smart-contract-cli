package models

// deployContract represents the data for the deployContract page
type DeployContractModel struct {
	// ShowError bool
	URI string
}

// NewDeployContractModel creates a new deployContract model
func NewDeployContractModel() *DeployContractModel {
	return &DeployContractModel{
		// ShowError: false,
		URI: "",
	}
}
