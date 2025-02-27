package views

import (
	"fmt"

	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
)

// AirdropView renders the airdrop page
func AirdropView(inputMode, nftInput, uri string) string {
	s := constant.AirdropPageTitle + "\n"
	s += string(constant.Separator) + "\n\n"

	if inputMode == constant.NFTInputMode {
		s += constant.NFTInputPrompt
		s += fmt.Sprintf("> %s", nftInput)
		if len(nftInput) == 0 {
			s += string(constant.InputCursor)
		}
		s += "\n\n" + constant.EnterToContinue
		s += "\n" + constant.BackToMenuMessage + "\n"
		s += constant.ExitMessage + "\n"
		// s += constant.ExitMessage + "\n"
	} else {
		s += fmt.Sprintf("%s%s\n\n", constant.NFTIDLabel, nftInput)
		s += constant.URIPrompt
		s += fmt.Sprintf("> %s", uri)
		if len(uri) == 0 {
			s += string(constant.InputCursor)
		}
		s += "\n\n" + constant.ConfirmAirdrop
		s += "\n" + constant.ReturnToNFTInput + "\n"
		s += constant.ExitMessage + "\n"
	}

	return s
}
