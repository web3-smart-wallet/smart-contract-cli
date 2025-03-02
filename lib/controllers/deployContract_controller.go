package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// DeployContractController handles the deploy contract page logic
type DeployContractController struct {
	nftService *services.NftService
}

// NewDeployContractController creates a new deploy contract controller
func NewDeployContractController(nftService *services.NftService) *DeployContractController {
	return &DeployContractController{
		nftService: nftService,
	}
}

// Update handles the check total page updates
func (c *DeployContractController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyEsc:
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		}
	}

	return model, nil
}

// View renders the check total page
func (c *DeployContractController) View() string {
	return views.DeployContractView()
}

func (c *DeployContractController) Name() constant.Page {
	return constant.DeployContractPage
}
