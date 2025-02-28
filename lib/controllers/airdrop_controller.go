package controllers

import (
	"fmt"
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/models"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// AirdropController handles the airdrop page logic
type AirdropController struct {
	model      *models.AirdropModel
	nftService *services.NftService
}

// NewAirdropController creates a new airdrop controller
func NewAirdropController(model *models.AirdropModel, nftService *services.NftService) *AirdropController {
	return &AirdropController{
		model:      model,
		nftService: nftService,
	}
}

// Update handles the airdrop page updates
func (c *AirdropController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyBackspace:
			if c.model.InputMode == constant.NFTInputMode && len(c.model.NFTInput) > 0 {
				c.model.NFTInput = c.model.NFTInput[:len(c.model.NFTInput)-1]
			} else if c.model.InputMode == constant.URLInputMode && len(c.model.URI) > 0 {
				c.model.URI = c.model.URI[:len(c.model.URI)-1]
			}
		case constant.KeyEnter:
			if c.model.InputMode == constant.NFTInputMode {
				if len(c.model.NFTInput) == 0 {
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: fmt.Errorf(constant.EmptyNFTIDError)}
					}
				}
				if len(c.model.NFTInput) > constant.MaxNFTIDLength {
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: fmt.Errorf(constant.LongNFTIDError)}
					}
				}
				c.model.InputMode = constant.URLInputMode
				return model, nil
			}

			if c.model.InputMode == constant.URLInputMode {
				if len(c.model.URI) == 0 {
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: fmt.Errorf(constant.EmptyURLError)}
					}
				} else if len(c.model.URI) > constant.MaxURLLength {
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: fmt.Errorf(constant.LongURLError)}
					}
				}
				matched, _ := regexp.MatchString(constant.URLPattern, c.model.URI)
				if !matched {
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: fmt.Errorf(constant.InvalidURLError)}
					}
				} else {
					return model, func() tea.Msg {
						return types.ChangePageMsg{Page: constant.UpLoadPage}
					}
				}
			}
		case constant.KeyEsc:
			if c.model.InputMode == constant.URLInputMode {
				c.model.URI = ""
				c.model.InputMode = constant.NFTInputMode
				return model, nil
			}
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		default:
			if c.model.InputMode == constant.NFTInputMode {
				// Only accept numeric input
				if _, err := strconv.Atoi(msg.String()); err == nil {
					c.model.NFTInput += msg.String()
				}
			} else if c.model.InputMode == constant.URLInputMode {
				c.model.URI += msg.String()
			}
		}
	}

	return model, nil
}

// View renders the airdrop page
func (c *AirdropController) View() string {
	return views.AirdropView(c.model.InputMode, c.model.NFTInput, c.model.URI)
}

func (c *AirdropController) Name() constant.Page {
	return constant.AirdropPage
}
