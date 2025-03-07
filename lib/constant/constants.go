package constant

// PageState constants for different pages
type Page string

const (
	PasswordPage       Page = "PasswordPage"
	MenuPage           Page = "MenuPage"
	DeployPage         Page = "DeployPage"
	DeployContractPage Page = "DeployContractPage"
	AirdropPage        Page = "AirdropPage"
	UpLoadPage         Page = "UpLoadPage"
	ConfirmPage        Page = "ConfirmPage"
	CheckTotalPage     Page = "CheckTotalPage"
	SelectContractPage Page = "SelectContractPage"
)

// Common constants
// URL validation pattern

const URLPattern = `^(http|https)://[a-zA-Z0-9\-._~:/?#\[\]@!$&'()*+,;=]+$`
const MaxNFTIDLength = 10 // Maximum NFT ID length
const MaxURLLength = 255  // Maximum URL length

// Menu options
var (
	MainMenuChoices   = []string{"Deploy Contract", "AirDrop NFT", "Check Total NFT"}
	DeployMenuChoices = []string{"Deploy new Contract(ERC1155)", "Check Existing Contracts"}
)

// Input modes
const (
	NFTInputMode = "nft"
	URLInputMode = "url"
)

type KeyboardKey string

// UI Controls
const (
	CursorActive   KeyboardKey = ">"
	CursorInactive KeyboardKey = " "
	InputCursor    KeyboardKey = "_"
	Separator      KeyboardKey = "--------------"

	// Key commands
	KeyCtrlC     KeyboardKey = "ctrl+c"
	KeyEsc       KeyboardKey = "esc"
	KeyUp        KeyboardKey = "up"
	KeyDown      KeyboardKey = "down"
	KeyBackspace KeyboardKey = "backspace"
	KeyEnter     KeyboardKey = "enter"
)

// UI Messages
const (
	ErrorPrefix   = "❌ 错误: "
	SuccessPrefix = "✅ "

	// Page titles and headers
	PasswordPageTitle         = "请输入密码:"
	MenuPageTitle             = "请选择操作:"
	DeployPageTitle           = "部署合约页面"
	MintNFTPageTitle          = "Mint NFT 页面"
	AirdropPageTitle          = "空投 NFT 页面"
	UploadPageTitle           = "上传文件页面"
	CheckTotalPageTitle       = "查看 NFT 总量页面"
	DeployContractPageTitle   = "部署新合约"
	DeployedContractPageTitle = "已部署的合约页面"

	// Common UI elements
	MainMenuFooter    = "主菜单."
	ExitMessage       = "按 Ctrl+C 退出程序."
	BackToMenuMessage = "按 ESC 返回主菜单"
	EnterToContinue   = "按 Enter 继续"

	// Error messages
	EmptyNFTIDError    = "NFT ID 不能为空"
	LongNFTIDError     = "NFT ID 太长"
	EmptyURLError      = "URL 不能为空"
	LongURLError       = "URL 太长"
	InvalidURLError    = "无效的 URL 格式"
	WrongPasswordError = "密码错误"
	LoginSuccess       = "登录成功！"
	NoDeployedContract = "暂无已部署的合约"
)

// UI Messages - Additional
const (
	LoadingMessage     = "正在处理...\n\n"
	UnknownPageMessage = "未知页面\n"

	// Input prompts
	NFTInputPrompt = "请输入要空投的 NFT 编号：\n"
	URIPrompt      = "请输入 URI：\n"
	NFTIDLabel     = "NFT 编号: "

	// Navigation messages
	BackToPrevious   = "按 ESC 返回上一页"
	ReturnToNFTInput = "按 ESC 重新输入 NFT 编号"
	ConfirmAirdrop   = "按 Enter 确认空投"
)
