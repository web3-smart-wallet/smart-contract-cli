package views

import (
	"fmt"

	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	types "github.com/web3-smart-wallet/smart-contract-cli/lib/types"
)

func SelectContractView(choices []types.ContractChoice, cursor int) string {
	s := "选择要操作的合约地址\n"
	s += string(constant.Separator) + "\n\n"

	if len(choices) == 0 {
		s += "暂无已部署的合约,请先部署合约后再进行空投操作\n"

	} else {

		for i, choice := range choices {
			cursorChar := " "
			if cursor == i {
				cursorChar = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursorChar, choice.Address)
			s += fmt.Sprintf("  部署时间: %s\n\n", choice.DeployTime)
		}
	}

	s += "\n" + constant.BackToMenuMessage + "\n"
	s += constant.ExitMessage + "\n"

	return s
}

// AirdropView renders the airdrop page
func AirdropView(inputMode, nftInput, uri string) string {
	s := constant.AirdropPageTitle + "\n"
	s += fmt.Sprintf("合约地址: %s\n", types.GlobalState.SelectedContract)
	s += fmt.Sprintf("当前TokenURI: %s\n", types.GlobalState.TokenURI)
	s += string(constant.Separator) + "\n\n"

	if inputMode == constant.NFTInputMode {
		s += constant.NFTInputPrompt
		s += fmt.Sprintf("> %s", nftInput)
		if len(nftInput) == 0 {
			s += string(constant.InputCursor)
		}
		s += "\n\n" + constant.EnterToContinue
		s += "\n" + constant.BackToPrevious + "\n"
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
