package views

import (
	"fmt"

	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
)

// DeployView renders the deploy menu page
func DeployView(deployChoices []string, deployCursor int) string {
	s := constant.DeployPageTitle + "\n" + string(constant.Separator) + "\n"

	for i, choice := range deployChoices {
		cursor := constant.CursorInactive
		if deployCursor == i {
			cursor = constant.CursorActive
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n" + constant.BackToMenuMessage + "\n"
	s += constant.ExitMessage + "\n"
	return s
}

// DeployContractView renders the deploy contract page
func DeployContractView(uri string) string {
	s := constant.DeployContractPageTitle + "\n" + string(constant.Separator) + "\n\n"

	s += constant.URIPrompt + "\n"

	// if len(uri) > 0 {
	// 	s += "\n"
	// }

	// s += string(constant.InputCursor)

	s += fmt.Sprintf("> %s", uri)
	if len(uri) == 0 {
		s += string(constant.InputCursor)
	}

	// if showError {
	// 	s += "\n\n" + constant.ErrorPrefix + constant.InvalidURLError
	// }

	s += "\n\n" + constant.BackToMenuMessage + "\n"
	s += constant.ExitMessage

	return s
}
