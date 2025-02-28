package controllers

import (
	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// MenuController handles the main menu logic
type MenuController struct {
	choices []string
	cursor  int
}

// NewMenuController creates a new menu controller
func NewMenuController(choices []string) *MenuController {
	return &MenuController{
		choices: choices,
		cursor:  0,
	}
}

// Update handles the menu page updates
func (c *MenuController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {
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
			var nextPage constant.Page
			switch c.cursor {
			case 0:
				nextPage = constant.DeployPage
			case 1:
				nextPage = constant.AirdropPage
			case 2:
				nextPage = constant.CheckTotalPage
			}

			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: nextPage}
			}
		}
	}

	return model, nil
}

// View renders the menu page
func (c *MenuController) View() string {
	return views.MenuView(c.choices, c.cursor)
}

func (c *MenuController) Name() constant.Page {
	return constant.MenuPage
}
