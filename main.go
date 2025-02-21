package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	constant "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"
	"github.com/web3-smart-wallet/smart-contract-cli/lib/views"
)

// Logger represents a custom logger that writes to both file and stdout
type Logger struct {
	file   *os.File
	logger *log.Logger
}

// NewLogger creates a new logger that writes to file only
func NewLogger() (*Logger, error) {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with timestamp in name
	timestamp := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logsDir, fmt.Sprintf("app_%s.log", timestamp))

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create logger that only writes to file
	logger := log.New(file, "", log.Ldate|log.Ltime)

	return &Logger{
		file:   file,
		logger: logger,
	}, nil
}

// Log writes a message to the log
func (l *Logger) Log(level, message string) {
	l.logger.Printf("[%s] %s", level, message)
}

// Close closes the log file
func (l *Logger) Close() error {
	return l.file.Close()
}

// 添加页面状态常量
const (
	passwordPage = iota // 添加密码页面作为第一个状态
	menuPage
	deployPage
	deployContractPage
	airdropPage
	upLoadPage
	checkTotalPage
)

// 在文件开头添加自定义消息类型
type airdropMsg struct {
	nftID  string
	nftURL string
}

// 添加错误消息类型
type errorMsg struct {
	err error
}

// 添加成功消息类型
type successMsg struct {
	message string
}

type model struct {
	choices        []string // 菜单选项
	cursor         int      // 当前光标位置
	selected       int      // 当前选中的选项
	currentPage    int      // 当前页面状态
	deployChoices  []string // 部署合约选项
	deployCursor   int      // 部署合约光标位置
	nftInput       string   // 输入框内容
	graphURL       string   // Graph URL输入内容
	inputMode      string   // 输入模式：'nft' 或 'url'
	inputCursor    int      // 输入框光标位置
	password       string   // 用户输入的密码
	authenticated  bool     // 验证状态
	errorMessage   string   // 错误消息
	successMessage string   // 成功消息
	loading        bool     // 加载状态
	logger         *Logger
}

func initialModel() model {
	logger, err := NewLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	return model{
		choices:        constant.MainMenuChoices,
		cursor:         0,
		selected:       0,
		currentPage:    constant.PasswordPage,
		deployChoices:  constant.DeployMenuChoices,
		deployCursor:   0,
		nftInput:       "",
		graphURL:       "",
		inputMode:      "nft",
		inputCursor:    0,
		password:       "",
		authenticated:  false,
		errorMessage:   "",
		successMessage: "",
		loading:        false,
		logger:         logger,
	}
}

// logger

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case errorMsg:
		m.errorMessage = msg.err.Error()
		m.logger.Log("ERROR", msg.err.Error())
		return m, nil

	case successMsg:
		m.successMessage = msg.message
		m.logger.Log("INFO", msg.message)
		return m, nil

	case airdropMsg:
		// 处理空投消息，跳转到上传页面
		m.currentPage = upLoadPage
		return m, nil

	case tea.KeyMsg:
		m.errorMessage = ""
		m.successMessage = ""
		key := constant.KeyboardKey(msg.String())
		switch key {
		case constant.KeyCtrlC:
			return m, tea.Quit
		case constant.KeyEsc:
			if m.currentPage != constant.MenuPage {
				if m.currentPage == constant.AirdropPage && m.inputMode == constant.URLInputMode {
					m.graphURL = ""
					m.inputMode = constant.NFTInputMode
				} else if m.currentPage == constant.UpLoadPage {
					m.currentPage = constant.AirdropPage
				} else {
					m.currentPage = constant.MenuPage
					m.nftInput = ""
				}
			}
		case constant.KeyUp:
			if m.currentPage == constant.MenuPage && m.cursor > 0 {
				m.cursor--
			}
			if m.currentPage == constant.DeployPage && m.deployCursor > 0 {
				m.deployCursor--
			}
		case constant.KeyDown:
			if m.currentPage == constant.MenuPage && m.cursor < len(m.choices)-1 {
				m.cursor++
			}
			if m.currentPage == constant.DeployPage && m.deployCursor < len(m.deployChoices)-1 {
				m.deployCursor++
			}
		case constant.KeyBackspace:
			if m.currentPage == constant.AirdropPage {
				if m.inputMode == constant.NFTInputMode && len(m.nftInput) > 0 {
					m.nftInput = m.nftInput[:len(m.nftInput)-1]

				} else if m.inputMode == constant.URLInputMode && len(m.graphURL) > 0 {
					m.graphURL = m.graphURL[:len(m.graphURL)-1]
				}
			} else if m.currentPage == constant.PasswordPage && len(m.password) > 0 {
				m.password = m.password[:len(m.password)-1]
			}
		case constant.KeyEnter:
			if m.currentPage == constant.PasswordPage {
				if constant.Password == m.password {
					m.authenticated = true
					m.currentPage = constant.MenuPage
					m.password = ""
					return m, func() tea.Msg {
						return successMsg{message: constant.LoginSuccess}
					}
				}
				return m, func() tea.Msg {
					return errorMsg{err: fmt.Errorf(constant.WrongPasswordError)}
				}
			} else if m.currentPage == constant.MenuPage {
				// 根据选择切换到对应页面
				switch m.cursor {
				case 0:
					m.currentPage = deployPage
				case 1:
					m.currentPage = airdropPage
					m.nftInput = ""
				case 2:
					m.currentPage = checkTotalPage
				}
			} else if m.currentPage == airdropPage {
				if m.inputMode == constant.NFTInputMode {
					if len(m.nftInput) == 0 {
						return m, func() tea.Msg {
							return errorMsg{err: fmt.Errorf(constant.EmptyNFTIDError)}
						}
					}
					if len(m.nftInput) > constant.MaxNFTIDLength {
						return m, func() tea.Msg {
							return errorMsg{err: fmt.Errorf(constant.LongNFTIDError)}
						}
					}
					m.inputMode = constant.URLInputMode
					return m, nil
				}

				if m.inputMode == constant.URLInputMode {
					m.logger.Log("INFO", fmt.Sprintf("graphURL: %s", m.graphURL))
					if len(m.graphURL) == 0 {
						return m, func() tea.Msg {
							return errorMsg{err: fmt.Errorf(constant.EmptyURLError)}
						}
					} else if len(m.graphURL) > constant.MaxURLLength {
						return m, func() tea.Msg {
							return errorMsg{err: fmt.Errorf(constant.LongURLError)}
						}
					}
					matched, _ := regexp.MatchString(constant.URLPattern, m.graphURL)
					if !matched {
						return m, func() tea.Msg {
							return errorMsg{err: fmt.Errorf(constant.InvalidURLError)}
						}
					} else {
						return m, func() tea.Msg {
							return airdropMsg{nftID: m.nftInput, nftURL: m.graphURL}
						}
					}
				}
			} else if m.currentPage == upLoadPage {
				// 上传文件
				m.currentPage = constant.MenuPage
				return m, nil

			} else if m.currentPage == deployPage {
				switch m.deployCursor {
				case 0:
					m.currentPage = deployContractPage
					// case 1:
					// 	m.currentPage = checkTotalPage
				}
			}

		default:
			if m.currentPage == airdropPage {
				if m.inputMode == constant.NFTInputMode {
					// 只接受数字输入
					if _, err := strconv.Atoi(msg.String()); err == nil {
						m.nftInput += msg.String()
					}
				} else if m.inputMode == constant.URLInputMode {

					m.graphURL += msg.String()
				}
			} else if m.currentPage == constant.PasswordPage {
				m.password += msg.String()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	if m.errorMessage != "" {
		s.WriteString(fmt.Sprintf("\n%s%s\n\n", constant.ErrorPrefix, m.errorMessage))
	}

	if m.successMessage != "" {
		s.WriteString(fmt.Sprintf("\n%s%s\n\n", constant.SuccessPrefix, m.successMessage))
	}

	if m.loading {
		s.WriteString(constant.LoadingMessage)
	}

	switch m.currentPage {
	case constant.PasswordPage:
		view := views.PasswordView(m.password)
		s.WriteString(view)
	case constant.MenuPage:
		view := views.MenuView(m.choices, m.cursor)
		s.WriteString(view)
	case constant.DeployPage:
		view := views.DeployView(m.deployChoices, m.deployCursor)
		s.WriteString(view)
	case constant.DeployContractPage:
		view := views.DeployContractView()
		s.WriteString(view)
	case constant.AirdropPage:
		view := views.AirdropView(m.inputMode, m.nftInput, m.graphURL)
		s.WriteString(view)
	case constant.UpLoadPage:
		view := views.UploadView()
		s.WriteString(view)
	case constant.CheckTotalPage:
		view := views.CheckTotalView()
		s.WriteString(view)
	default:
		return constant.UnknownPageMessage
	}

	return s.String()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, please make sure it exists")
		os.Exit(1)
	}

	constant.Password = os.Getenv("PASSWORD")
	// check if password is set
	if len(constant.Password) == 0 {
		fmt.Println("Password is not set, please set it in the .env file")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Close the logger when the program exits
	if model, ok := m.(model); ok {
		if err := model.logger.Close(); err != nil {
			fmt.Printf("Error closing logger: %v\n", err)
		}
	}
}
