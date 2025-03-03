package controllers

import (
	"fmt"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/models"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// DeployContractController handles the deploy contract page logic
type DeployContractController struct {
	nftService       *services.NftService
	contractCompiler *services.ContractCompiler
	model            *models.DeployContractModel
}

// NewDeployContractController creates a new deploy contract controller
func NewDeployContractController(
	nftService *services.NftService,
	contractCompiler *services.ContractCompiler,
	model *models.DeployContractModel,

) *DeployContractController {
	controller := &DeployContractController{
		nftService:       nftService,
		contractCompiler: contractCompiler,
		model:            model,
	}

	// Load available contracts
	contracts, err := contractCompiler.GetAvailableContracts()
	if err != nil {
		// Log error but don't fail
		fmt.Printf("Warning: Failed to load available contracts: %v\n", err)
	}
	model.AvailableContracts = contracts

	return controller
}

// Update handles the deploy contract page updates
func (c *DeployContractController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyEsc:
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.DeployPage}
			}
		case constant.KeyEnter:
			if c.model.IsSelectingContract {
				if c.model.SelectedContract >= 0 && c.model.SelectedContract < len(c.model.AvailableContracts) {
					// Move to URI input after contract selection
					c.model.IsSelectingContract = false
					return model, nil
				}
				return model, nil
			}

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
			}

			selectedContract := c.model.AvailableContracts[c.model.SelectedContract]

			// Create deployment parameters
			params := services.DeployContractParams{
				Bytecode:   selectedContract.Bytecode,
				InitialURI: c.model.URI,
				GasLimit:   3000000,
			}

			// Deploy contract
			contractAddr, err := c.nftService.DeployContractWithABI(params)
			if err != nil {
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: err}
				}
			}

			// Save contract info
			err = c.contractCompiler.SaveDeployedContract(contractAddr, c.model.URI, selectedContract.ABI)
			if err != nil {
				model.Logger.Log("ERROR", fmt.Sprintf("保存合约信息失败: %v", err))
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: fmt.Errorf("save contract info failed")}
				}
			}

			// Set success message
			model.SuccessMessage = fmt.Sprintf("合约部署成功！地址: %s", contractAddr)

			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}

		case constant.KeyBackspace:
			if !c.model.IsSelectingContract && len(c.model.URI) > 0 {
				c.model.URI = c.model.URI[:len(c.model.URI)-1]
			}

		case constant.KeyUp:
			if c.model.IsSelectingContract && c.model.SelectedContract > 0 {
				c.model.SelectedContract--
			}

		case constant.KeyDown:
			if c.model.IsSelectingContract && c.model.SelectedContract < len(c.model.AvailableContracts)-1 {
				c.model.SelectedContract++
			}

		default:
			if !c.model.IsSelectingContract && len(msg.String()) == 1 {
				c.model.URI += msg.String()
			}
		}
	}

	return model, nil
}

// View renders the deploy contract page
func (c *DeployContractController) View() string {
	return views.DeployContractView(c.model)
}

func (c *DeployContractController) Name() constant.Page {
	return constant.DeployContractPage
}
