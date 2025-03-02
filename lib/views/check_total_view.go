package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
)

// CheckTotalView 渲染查看 NFT 总量页面
func CheckTotalView(contracts []services.DeployedContract) string {
	var sb strings.Builder

	sb.WriteString("\n已部署的合约信息\n")
	sb.WriteString("----------------\n\n")

	if len(contracts) == 0 {
		sb.WriteString("目前还没有已部署的合约\n")
	} else {
		for i, contract := range contracts {
			sb.WriteString(fmt.Sprintf("合约 #%d:\n", i+1))
			sb.WriteString(fmt.Sprintf("地址: %s\n", contract.Address))
			sb.WriteString(fmt.Sprintf("部署时间: %s\n", contract.DeployTime.Format(time.RFC3339)))
			sb.WriteString("----------------\n")
		}
	}

	sb.WriteString("\n按 ESC 返回主菜单\n")
	sb.WriteString("按 Ctrl+C 退出程序\n")

	return sb.String()
}
