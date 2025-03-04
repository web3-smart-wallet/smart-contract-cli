package types

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
)

// PasswordControllerInterface defines the interface for password controllers
type ControllerInterface interface {
	// Update handles the password page updates
	Update(model AppModel, msg tea.Msg) (any, tea.Cmd)
	// View renders the password page
	View() string
	Name() constant.Page
}

// AppModel represents the main application model
type AppModel struct {
	CurrentPage    constant.Page
	Cursor         int
	ErrorMessage   string
	SuccessMessage string
	Loading        bool
	Logger         *Logger
	Controllers    map[constant.Page]ControllerInterface
}

type State struct {
	UploadWalletAddresses []string
	SelectedContract      string
	DeployStat            bool
	TokenURI              string
	NFTID                 string
	SendNFTStat           bool
}

var GlobalState = &State{}

// 上传地址
// model.state.uploadurl = "https://www.baidu.com"
type ContractChoice struct {
	Address    string
	DeployTime string
}
