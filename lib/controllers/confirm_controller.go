package controllers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// ConfirmController handles the confirmation page logic
type ConfirmController struct {
	nftService      *services.NftService
	contractAddress string
	walletAddresses []string
	nftID           string
	uri             string
}

// NewConfirmController creates a new confirm controller
func NewConfirmController(nftService *services.NftService, contractAddress string) *ConfirmController {
	return &ConfirmController{
		nftService:      nftService,
		contractAddress: contractAddress,
		walletAddresses: []string{},
		nftID:           "",
		uri:             "",
	}
}

// SetWalletAddresses sets the wallet addresses
func (c *ConfirmController) SetWalletAddresses(addresses []string) {
	c.walletAddresses = addresses
}

// SetNFTID sets the NFT ID
func (c *ConfirmController) SetNFTID(nftID string) {
	c.nftID = nftID
}

// SetURI sets the URI
func (c *ConfirmController) SetURI(uri string) {
	c.uri = uri
}

// Update handles the confirm page updates
func (c *ConfirmController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyEnter:
			// Set the URI first
			if err := c.nftService.SetURI(c.contractAddress, c.uri); err != nil {
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: fmt.Errorf("设置 URI 失败: %v", err)}
				}
			}

			// Send NFTs to addresses
			txHash, err := c.nftService.MintNFTToAddresses(c.contractAddress, c.walletAddresses, c.nftID)
			if err != nil {
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: fmt.Errorf("发送 NFT 失败: %v", err)}
				}
			}

			model.Logger.Log("INFO", fmt.Sprintf("NFT 发送成功，交易哈希: %s", txHash))

			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		case constant.KeyEsc:
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		}
	}

	return model, nil
}

// View renders the confirm page
func (c *ConfirmController) View() string {
	return views.ConfirmView(types.GlobalState.UploadWalletAddresses)
}

func (c *ConfirmController) Name() constant.Page {
	return constant.ConfirmPage
}
