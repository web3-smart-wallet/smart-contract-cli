package types

import "github.com/web3-smart-wallet/smart-contract-cli/lib/constant"

type ErrorMsg struct {
	Err error
}

type SuccessMsg struct {
	Message string
}

type ChangePageMsg struct {
	Page constant.Page
}
