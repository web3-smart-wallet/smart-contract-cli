package controllers

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// UploadController handles the upload page logic
type UploadController struct {
	filePath        string
	uploadError     error
	walletAddresses []string
}

// NewUploadController creates a new upload controller
func NewUploadController() *UploadController {
	return &UploadController{
		filePath:        "./addresses.txt",
		uploadError:     nil,
		walletAddresses: []string{},
	}
}

// Update handles the upload page updates
func (c *UploadController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyEnter:
			// Read and parse the file
			addresses, err := c.parseWalletAddresses(c.filePath)
			if err != nil {
				c.uploadError = err
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: err}
				}
			}
			c.walletAddresses = addresses

			// Pass the addresses to the confirm controller
			// m.confirmController.SetWalletAddresses(addresses)
			// m.confirmController.SetNFTID(m.airdropController.model.NFTInput)
			// m.confirmController.SetURI(m.airdropController.model.URI)

			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.ConfirmPage}
			}
		case constant.KeyEsc:
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.AirdropPage}
			}
		}
	}

	return model, nil
}

// View renders the upload page
func (c *UploadController) View() string {
	var errorStr string
	if c.uploadError != nil {
		errorStr = c.uploadError.Error()
	}
	return views.UploadView(c.filePath, errorStr)
}

// parseWalletAddresses reads and validates wallet addresses from a file
func (c *UploadController) parseWalletAddresses(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// 清空之前的地址列表，避免重复
	types.GlobalState.UploadWalletAddresses = []string{}

	// Ethereum address regex
	ethAddressRegex := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !ethAddressRegex.MatchString(line) {
			return nil, fmt.Errorf("第 %d 行包含无效的以太坊地址: %s", i+1, line)
		}

		types.GlobalState.UploadWalletAddresses = append(types.GlobalState.UploadWalletAddresses, line)
	}

	if len(types.GlobalState.UploadWalletAddresses) == 0 {
		return nil, fmt.Errorf("文件中没有找到有效的钱包地址")
	}

	return types.GlobalState.UploadWalletAddresses, nil
}

func (c *UploadController) Name() constant.Page {
	return constant.UpLoadPage
}
