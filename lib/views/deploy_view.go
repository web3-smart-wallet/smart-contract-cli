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
func DeployContractView() string {
	return `
Mint NFT 页面
--------------
这里是 Mint NFT 的界面

按 ESC 返回主菜单
按 Ctrl+C 退出程序
`
}
