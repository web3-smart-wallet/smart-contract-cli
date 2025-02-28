package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// CheckTotalController handles the check total page logic
type CheckTotalController struct {
	nftService      *services.NftService
	contractAddress string
}

// NewCheckTotalController creates a new check total controller
func NewCheckTotalController(nftService *services.NftService, contractAddress string) *CheckTotalController {
	return &CheckTotalController{
		nftService:      nftService,
		contractAddress: contractAddress,
	}
}

// Update handles the check total page updates
func (c *CheckTotalController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {

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
func (c *CheckTotalController) View() string {
	return views.CheckTotalView()
}

func (c *CheckTotalController) Name() constant.Page {
	return constant.CheckTotalPage
}
