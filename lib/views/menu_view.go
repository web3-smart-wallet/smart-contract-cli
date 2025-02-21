package views

import (
	"fmt"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
)

// MenuView renders the main menu page
func MenuView(choices []string, cursor int) string {
	s := constant.MenuPageTitle + "\n\n"

	for i, choice := range choices {
		cursorChar := constant.CursorInactive
		if cursor == i {
			cursorChar = constant.CursorActive
		}
		s += fmt.Sprintf("%s %s\n", cursorChar, choice)
	}
	s += "\n" + constant.MainMenuFooter + "\n"
	s += constant.ExitMessage + "\n"
	return s
} 