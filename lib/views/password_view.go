package views

import (
	"strings"
)

// PasswordView renders the password input page
func PasswordView(password string) string {
	s := "请输入密码:\n\n"
	s += "> " + strings.Repeat("*", len(password))
	if len(password) == 0 {
		s += "_"
	}
	s += "\n\n按 Enter 确认\n"
	s += "按 Ctrl+C 退出程序\n"
	return s
}
