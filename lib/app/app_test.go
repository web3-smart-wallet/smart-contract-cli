package app

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/suite"
)

type AppTestSuite struct {
	suite.Suite
}

func TestAppTestSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

func TestAppFlow(t *testing.T) {
	// Set up test environment variables
	t.Setenv("RPC_URL", "http://localhost:8545")
	t.Setenv("PRIVATE_KEY", "0xtest1234567890")
	t.Setenv("PASSWORD", "123")

	// Create a new test program
	tm := teatest.NewTestModel(t, initialModel())
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return strings.Contains(string(bts), "请输入密码:")
	}, teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*3))

	// enter password into password prompt
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("inc")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	// wait for error message
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return strings.Contains(string(bts), "错误: 密码错误")
	}, teatest.WithCheckInterval(time.Second),
		teatest.WithDuration(time.Second*3))

	tm = teatest.NewTestModel(t, initialModel())
	// enter password again
	// clear the input
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("123")})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	// wait for success message
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {

		return strings.Contains(string(bts), "Deploy Contract")
	}, teatest.WithCheckInterval(time.Second),
		teatest.WithDuration(time.Second*3))

}
