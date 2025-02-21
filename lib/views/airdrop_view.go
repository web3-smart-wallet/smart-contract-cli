package views

import (
	"fmt"

	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
)

// AirdropView renders the airdrop page
func AirdropView(inputMode, nftInput, graphURL string) string {
	s := constant.AirdropPageTitle + "\n"
	s += string(constant.Separator) + "\n\n"

	if inputMode == constant.NFTInputMode {
		s += constant.NFTInputPrompt
		if len(nftInput) == 0 {
			s += string(constant.InputCursor)
		}
		s += "\n\n" + constant.EnterToContinue
		s += "\n" + constant.BackToMenuMessage + "\n"
		s += constant.ExitMessage + "\n"
		s += constant.ExitMessage + "\n"
	} else {
		s += fmt.Sprintf("%s%s\n\n", constant.NFTIDLabel, nftInput)
		s += constant.GraphURLPrompt
		s += fmt.Sprintf("> %s", graphURL)
		if len(graphURL) == 0 {
			s += string(constant.InputCursor)
		}
		s += "\n\n" + constant.ConfirmAirdrop
		s += "\n" + constant.ReturnToNFTInput + "\n"
		s += constant.ExitMessage + "\n"
	}

	return s
}
