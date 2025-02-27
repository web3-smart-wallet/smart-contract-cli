package views

import (
	"fmt"
	"strings"
)

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

// 添加确认页面的视图函数
func ConfirmView(addressCount int) string {
	var sb strings.Builder

	sb.WriteString("\n=== 确认发送 NFT ===\n\n")
	sb.WriteString(fmt.Sprintf("即将向 %d 个地址发送 NFT\n\n", addressCount))
	sb.WriteString("按 Enter 确认发送\n")
	sb.WriteString("按 ESC 取消操作\n")

	return sb.String()
}
