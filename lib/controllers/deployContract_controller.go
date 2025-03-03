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
	return &DeployContractController{
		nftService:       nftService,
		contractCompiler: contractCompiler,
		model:            model,
	}
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

				// 获取合约字节码和ABI
				bytecode, abi, err := c.contractCompiler.GetContractBytecode()
				if err != nil {
					// c.model.ShowError = true
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: err}
					}
				}

				// 创建部署参数
				params := services.DeployContractParams{
					Bytecode:   bytecode,
					InitialURI: c.model.URI,
					GasLimit:   3000000,
				}

				// 部署合约
				contractAddr, err := c.nftService.DeployContractWithABI(params)
				if err != nil {
					// c.model.ShowError = true
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: err}
					}
				}

				// 保存合约信息
				err = c.contractCompiler.SaveDeployedContract(contractAddr, c.model.URI, abi)
				if err != nil {
					// 记录错误但不中断流程
					model.Logger.Log("ERROR", fmt.Sprintf("保存合约信息失败: %v", err))
					return model, func() tea.Msg {
						return types.ErrorMsg{Err: fmt.Errorf("save contract info failed")}
					}
				}

				// 设置成功消息
				model.SuccessMessage = fmt.Sprintf("合约部署成功！地址: %s", contractAddr)

				return model, func() tea.Msg {
					return types.ChangePageMsg{Page: constant.MenuPage}

				}
			}

		case constant.KeyBackspace:
			if len(c.model.URI) > 0 {
				c.model.URI = c.model.URI[:len(c.model.URI)-1]
				// c.model.ShowError = false
			}

		default:
			if len(msg.String()) == 1 {
				c.model.URI += msg.String()
				// c.model.ShowError = false
			}
		}
	}

	return model, nil
}

// validateURI checks if the URI is valid
// func (c *DeployContractController) validateURI() bool {
// 	if len(c.model.URI) == 0 {
// 		return false
// 	}
// 	matched, _ := regexp.MatchString(constant.URLPattern, c.model.URI)
// 	return matched
// }

// View renders the deploy contract page
func (c *DeployContractController) View() string {
	return views.DeployContractView(c.model.URI)
}

func (c *DeployContractController) Name() constant.Page {
	return constant.DeployContractPage
}
