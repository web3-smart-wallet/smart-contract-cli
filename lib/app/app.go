package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/controllers"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/models"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/pages/password"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/services"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/types"

	tea "github.com/charmbracelet/bubbletea"
)

// Create a local type that embeds the imported type
type LocalModel struct {
	types.AppModel
	types.State
}

func initialModel() LocalModel {
	logger, err := types.NewLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// Get configuration from environment variables
	rpcURL := os.Getenv("RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")
	userPassword := os.Getenv("PASSWORD")

	if rpcURL == "" || privateKey == "" {
		fmt.Println("RPC_URL or PRIVATE_KEY is not set, please set it in the .env file")
		os.Exit(1)
	}

	if userPassword == "" {
		fmt.Println("PASSWORD is not set, please set it in the .env file")
		os.Exit(1)
	}

	// Create shared services
	nftService := services.NewNftService(rpcURL, privateKey)
	passwordService := password.NewService(userPassword)
	contractService := services.NewContractCompiler("./artifacts")

	// Create shared models
	airdropModel := models.NewAirdropModel()
	deployContractModel := models.NewDeployContractModel()

	// Create controllers
	passwordController := controllers.NewPasswordController(passwordService)
	menuController := controllers.NewMenuController(constant.MainMenuChoices)
	deployController := controllers.NewDeployController(constant.DeployMenuChoices)

	deployContractController := controllers.NewDeployContractController(nftService, contractService, deployContractModel)
	selectContractController := controllers.NewSelectContractController(contractService)
	airdropController := controllers.NewAirdropController(airdropModel, nftService)
	uploadController := controllers.NewUploadController()
	confirmController := controllers.NewConfirmController(nftService)
	checkController := controllers.NewCheckTotalController(nftService, contractService)

	// 添加合约选择控制器

	return LocalModel{
		AppModel: types.AppModel{
			CurrentPage:    constant.PasswordPage,
			Cursor:         0,
			ErrorMessage:   "",
			SuccessMessage: "",
			Loading:        false,
			Logger:         logger,
			Controllers: map[constant.Page]types.ControllerInterface{
				constant.PasswordPage:       passwordController,
				constant.MenuPage:           menuController,
				constant.DeployPage:         deployController,
				constant.DeployContractPage: deployContractController,
				constant.SelectContractPage: selectContractController,
				constant.AirdropPage:        airdropController,
				constant.UpLoadPage:         uploadController,
				constant.ConfirmPage:        confirmController,
				constant.CheckTotalPage:     checkController,
			},
		},
		State: types.State{
			UploadWalletAddresses: []string{},
			SelectedContract:      "", // 添加选中的合约地址
		},
	}
}

// Define methods on the local type
func (m LocalModel) Init() tea.Cmd {
	return nil
}

func (m LocalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case types.ErrorMsg:
		m.AppModel.ErrorMessage = msg.Err.Error()
		m.AppModel.Logger.Log("ERROR", msg.Err.Error())
		return m, nil

	case types.SuccessMsg:
		m.AppModel.SuccessMessage = msg.Message
		m.AppModel.Logger.Log("INFO", msg.Message)
		return m, nil

	case types.ChangePageMsg:
		m.AppModel.CurrentPage = msg.Page
		m.AppModel.Cursor = 0 // Reset cursor when changing pages
		return m, nil

	case tea.KeyMsg:
		m.AppModel.ErrorMessage = ""
		m.AppModel.SuccessMessage = ""
		key := constant.KeyboardKey(msg.String())

		// Global key handlers
		if key == constant.KeyCtrlC {
			return m, tea.Quit
		}

		// Page-specific updates
		var cmd tea.Cmd
		var result any

		controller := m.AppModel.Controllers[m.AppModel.CurrentPage]
		result, cmd = controller.Update(m.AppModel, msg)
		// Add type assertion to convert interface{} back to AppModel
		if result != nil {
			m.AppModel = result.(types.AppModel)
		}

		return m, cmd
	}

	return m, nil
}

func (m LocalModel) View() string {
	var s strings.Builder

	if m.AppModel.ErrorMessage != "" {
		s.WriteString(fmt.Sprintf("\n%s%s\n\n", constant.ErrorPrefix, m.AppModel.ErrorMessage))
	}

	if m.AppModel.SuccessMessage != "" {
		s.WriteString(fmt.Sprintf("\n%s%s\n\n", constant.SuccessPrefix, m.AppModel.SuccessMessage))
	}

	if m.AppModel.Loading {
		s.WriteString(constant.LoadingMessage)
	}

	controller := m.AppModel.Controllers[m.AppModel.CurrentPage]
	s.WriteString(controller.View())

	return s.String()
}

// Run the application
func Run() {
	_ = godotenv.Load() // ignore error since it's not required

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Close the logger when the program exits
	if model, ok := m.(LocalModel); ok {
		if err := model.Logger.Close(); err != nil {
			fmt.Printf("Error closing logger: %v\n", err)
		}
	}
}
