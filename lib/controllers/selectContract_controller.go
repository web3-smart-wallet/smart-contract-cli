package controllers

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"
	views "github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// SelectContractController handles the select contract logic
type SelectContractController struct {
	contractService *services.ContractCompiler
	// model           *types.State
	choices []types.ContractChoice
	cursor  int
}

// NewSelectContractController creates a new select contract controller
func NewSelectContractController(contractService *services.ContractCompiler) *SelectContractController {
	// 从deployed_contracts.json读取已部署的合约地址
	contracts, err := contractService.GetDeployedContracts()
	if err != nil {
		return &SelectContractController{
			contractService: contractService,
			// model:           model,
			choices: []types.ContractChoice{},
			cursor:  0,
		}
	}

	// 修改这里，保存地址和部署时间
	choices := []types.ContractChoice{}
	for _, contract := range contracts {
		choices = append(choices, types.ContractChoice{
			Address:    contract.Address,
			DeployTime: contract.DeployTime.Format("2006-01-02 15:04:05"),
		})
	}

	// if len(choices) == 0 {
	// 	choices = []string{"暂无已部署的合约"}
	// }

	return &SelectContractController{
		contractService: contractService,
		choices:         choices,
		cursor:          0,
	}
}

// Update handles the menu page updates
func (c *SelectContractController) Update(model types.AppModel, msg tea.Msg) (interface{}, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := constant.KeyboardKey(msg.String())

		switch key {
		case constant.KeyUp:
			if c.cursor > 0 {
				c.cursor--
			}
		case constant.KeyDown:
			if c.cursor < len(c.choices)-1 {
				c.cursor++
			}
		case constant.KeyEnter:
			if len(c.choices) == 0 {
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: fmt.Errorf("请先部署合约后再进行空投操作")}
				}
			}
			// 只有当有已部署合约时才允许进入空投页面

			selectedContract := c.choices[c.cursor].Address
			types.GlobalState.SelectedContract = selectedContract

			// 获取选中合约的 tokenURI
			contracts, err := c.contractService.GetDeployedContracts()
			if err != nil {
				return model, func() tea.Msg {
					return types.ErrorMsg{Err: fmt.Errorf("获取合约信息失败: %v", err)}
				}
			}

			// 查找选中的合约并获取其 tokenURI
			for _, contract := range contracts {
				if contract.Address == selectedContract {
					types.GlobalState.TokenURI = contract.TokenURI
					break
				}
			}

			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.AirdropPage}
			}

		case constant.KeyEsc:
			return model, func() tea.Msg {
				return types.ChangePageMsg{Page: constant.MenuPage}
			}
		}
	}

	return model, nil
}

// View renders the menu page
func (c *SelectContractController) View() string {
	return views.SelectContractView(c.choices, c.cursor)
}

func (c *SelectContractController) Name() constant.Page {
	return constant.SelectContractPage
}
