package views

import (
	"fmt"
	"strings"
	"time"

	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
)

// CheckTotalView 渲染查看 NFT 总量页面
func CheckTotalView(contracts []services.DeployedContract) string {
	var sb strings.Builder

	sb.WriteString(constant.DeployedContractPageTitle + "\n")
	sb.WriteString(string(constant.Separator) + "\n\n")

	if len(contracts) == 0 {
		sb.WriteString(constant.NoDeployedContract + "\n")
	} else {
		for i, contract := range contracts {
			sb.WriteString(fmt.Sprintf("合约 #%d:\n", i+1))
			sb.WriteString(fmt.Sprintf("地址: %s\n", contract.Address))
			sb.WriteString(fmt.Sprintf("部署时间: %s\n", contract.DeployTime.Format(time.RFC3339)))
			sb.WriteString(string(constant.Separator) + "\n")
		}
	}

	sb.WriteString("\n" + constant.BackToPrevious)
	sb.WriteString("\n" + constant.ExitMessage + "\n")

	return sb.String()
}
