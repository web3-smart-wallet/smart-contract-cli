package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/pages/password"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
)

// PasswordController handles the password page logic
type PasswordController struct {
	service *password.Service
	input   string
}

// NewPasswordController creates a new password controller
func NewPasswordController(service *password.Service) *PasswordController {
	return &PasswordController{
		service: service,
		input:   "",
	}
}

// Update handles the password page updates
func (c *PasswordController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyBackspace:
			if len(c.input) > 0 {
				c.input = c.input[:len(c.input)-1]
			}
		case constant.KeyEnter:
			err := c.service.VerifyPassword(c.input)
			if err != nil {
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: err}
				}
			}

			// Password verified, move to menu page
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		default:
			c.input += msg.String()
		}
	}

	return model, nil
}

// View renders the password page
func (c *PasswordController) View() string {
	return password.View(c.input)
}

func (c *PasswordController) Name() constant.Page {
	return constant.PasswordPage
}
