package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// DeployController handles the deploy page logic
type DeployController struct {
	choices []string
	cursor  int
}

// NewDeployController creates a new deploy controller
func NewDeployController(choices []string) *DeployController {
	return &DeployController{
		choices: choices,
		cursor:  0,
	}
}

// Update handles the deploy page updates
func (c *DeployController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyUp:
			if c.cursor > 0 {
				c.cursor--
			}
		case constant.KeyDown:
			if c.cursor < len(c.choices)-1 {
				c.cursor++
			}
		case constant.KeyEnter:
			// Navigate to the selected page
			if c.cursor == 0 {
				return model, func() tea.Msg {
					return types.ChangePageMsg{Page: constant.DeployContractPage}
				}
			}
		case constant.KeyEsc:
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		}
	}

	return model, nil
}

// View renders the deploy page
func (c *DeployController) View() string {
	return views.DeployView(c.choices, c.cursor)
}

func (c *DeployController) Name() constant.Page {
	return constant.DeployPage
}
