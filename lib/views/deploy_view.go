package views

import (
	"fmt"

	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/models"
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
func DeployContractView(model *models.DeployContractModel) string {
	s := constant.DeployContractPageTitle + "\n" + string(constant.Separator) + "\n\n"

	if model.IsSelectingContract {
		if len(model.AvailableContracts) == 0 {
			s += "没有检测到可部署的合约。\n"
			s += "请确保在 contracts/ 目录下有编译好的合约JSON文件。\n"
		} else {
			s += "你想部署哪一个合约？（检测到有以下可以部署的合约）\n\n"

			for i, contract := range model.AvailableContracts {
				cursor := constant.CursorInactive
				if model.SelectedContract == i {
					cursor = constant.CursorActive
				}
				s += fmt.Sprintf("%s %s(%s)\n", cursor, contract.ContractName, contract.FilePath)
			}
		}
	} else {
		s += constant.URIPrompt + "\n"
		s += fmt.Sprintf("> %s", model.URI)
		if len(model.URI) == 0 {
			s += string(constant.InputCursor)
		}
	}

	s += "\n\n" + constant.BackToMenuMessage + "\n"
	s += constant.ExitMessage

	return s
}
