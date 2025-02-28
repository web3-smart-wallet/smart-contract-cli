package models

// AirdropModel represents the data for the airdrop page
type AirdropModel struct {
	NFTInput   string
	URI        string
	InputMode  string
}

// NewAirdropModel creates a new airdrop model
func NewAirdropModel() *AirdropModel {
	return &AirdropModel{
		NFTInput:  "",
		URI:       "",
		InputMode: "nft",
	}
} 