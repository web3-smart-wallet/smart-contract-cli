package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
)

type CurrentTime struct {
	CurrentTime time.Time
}

// UploadView renders the upload file page
func UploadView(filePath string, errorMsg string) string {
	var sb strings.Builder

	sb.WriteString("\n=== 文件上传页面 ===\n\n")
	sb.WriteString("请将钱包地址列表保存在 ./addresses.txt 文件中\n")
	sb.WriteString("每行一个地址\n\n")
	sb.WriteString("按 Enter 读取文件\n")
	sb.WriteString("按 ESC 返回上一页\n")

	if errorMsg != "" {
		sb.WriteString(fmt.Sprintf("\n错误: %s\n", errorMsg))
	}

	return sb.String()
}

// ConfirmView renders the confirmation page
func ConfirmView(addresses []string) string {
	var sb strings.Builder
	if types.GlobalState.SendNFTStat {
		sb.WriteString(fmt.Sprintf("NFT发送时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
		sb.WriteString("=== 确认发送 NFT ===\n\n")
	} else {
		sb.WriteString("\n=== 确认发送 NFT ===\n\n")
	}
	sb.WriteString(fmt.Sprintf("即将向 %d 个地址发送 NFT\n\n", len(addresses)))

	// 显示前5个地址作为预览
	if len(addresses) > 0 {
		sb.WriteString("地址预览：\n")
		previewCount := min(5, len(addresses))
		for i := 0; i < previewCount; i++ {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, addresses[i]))
		}
		if len(addresses) > 5 {
			sb.WriteString("...\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("按 Enter 确认发送\n")
	sb.WriteString("按 ESC 取消操作\n")

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
