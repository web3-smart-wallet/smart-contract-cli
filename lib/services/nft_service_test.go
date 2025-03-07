package services

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type NftServiceTestSuite struct {
	suite.Suite
	NftService *NftService
}

const (
	HardhatNodeUrl    = "http://localhost:8545"
	HardhatPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	HardhatAddress    = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	AbiPath           = "../../artifacts/contracts/nft.sol/MyToken.json"
)

type AbiInfoFile struct {
	Abi      any    `json:"abi"`
	Bytecode string `json:"bytecode"`
}

// run hardhat node before running the tests
func TestNftServiceSuite(t *testing.T) {
	suite.Run(t, new(NftServiceTestSuite))
}

func (s *NftServiceTestSuite) SetupTest() {
	s.NftService = NewNftService(HardhatNodeUrl, HardhatPrivateKey)
}

func (s *NftServiceTestSuite) TestDeployContract() {
	// read bytecode from file
	abiInfoFile, err := os.ReadFile(AbiPath)
	if err != nil {
		s.T().Fatalf("failed to read bytecode: %v", err)
	}

	var abiInfo AbiInfoFile
	err = json.Unmarshal(abiInfoFile, &abiInfo)
	if err != nil {
		s.T().Fatalf("failed to unmarshal abi info: %v", err)
	}

	abiBytes, err := json.Marshal(abiInfo.Abi)
	if err != nil {
		s.T().Fatalf("failed to marshal abi: %v", err)
	}

	// Convert the address string to common.Address
	initialOwner := common.HexToAddress(HardhatAddress)

	contractAddress, err := s.NftService.DeployContractWithABI(
		DeployContractParams{
			Bytecode:       string(abiInfo.Bytecode),
			ConstructorABI: string(abiBytes),
			ConstructorArgs: []any{
				initialOwner,
			},
		},
	)
	if err != nil {
		s.T().Fatalf("failed to deploy contract: %v", err)
	}
	s.Require().NotNil(contractAddress)
	s.Require().NotEmpty(contractAddress)
}

func (s *NftServiceTestSuite) TestSetURI() {
	// First deploy the contract
	abiInfoFile, err := os.ReadFile(AbiPath)
	if err != nil {
		s.T().Fatalf("failed to read bytecode: %v", err)
	}

	var abiInfo AbiInfoFile
	err = json.Unmarshal(abiInfoFile, &abiInfo)
	if err != nil {
		s.T().Fatalf("failed to unmarshal abi info: %v", err)
	}

	abiBytes, err := json.Marshal(abiInfo.Abi)
	if err != nil {
		s.T().Fatalf("failed to marshal abi: %v", err)
	}

	initialOwner := common.HexToAddress(HardhatAddress)

	// Deploy the contract
	contractAddress, err := s.NftService.DeployContractWithABI(
		DeployContractParams{
			Bytecode:       string(abiInfo.Bytecode),
			ConstructorABI: string(abiBytes),
			ConstructorArgs: []any{
				initialOwner,
			},
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(contractAddress)

	// Call setURI function
	newURI := "https://api.example.com/token/{id}"
	txHash, err := s.NftService.CallContractFunction(
		ContractCallParams{
			ContractAddress: contractAddress,
			ContractABI:     string(abiBytes),
			FunctionName:    "setURI",
			FunctionArgs:    []any{newURI},
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(txHash)
}

func (s *NftServiceTestSuite) TestMintToMultiple() {
	// ... existing code for reading ABI and deploying contract ...
	abiInfoFile, err := os.ReadFile(AbiPath)
	if err != nil {
		s.T().Fatalf("failed to read bytecode: %v", err)
	}

	var abiInfo AbiInfoFile
	err = json.Unmarshal(abiInfoFile, &abiInfo)
	if err != nil {
		s.T().Fatalf("failed to unmarshal abi info: %v", err)
	}

	abiBytes, err := json.Marshal(abiInfo.Abi)
	if err != nil {
		s.T().Fatalf("failed to marshal abi: %v", err)
	}

	initialOwner := common.HexToAddress(HardhatAddress)

	// Deploy the contract
	contractAddress, err := s.NftService.DeployContractWithABI(
		DeployContractParams{
			Bytecode:       string(abiInfo.Bytecode),
			ConstructorABI: string(abiBytes),
			ConstructorArgs: []any{
				initialOwner,
			},
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(contractAddress)

	// 准备测试数据：创建多个接收地址
	recipients := []string{
		"0x70997970C51812dc3A010C7d01b50e0d17dc79C8", // Hardhat 账号 #1
		"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", // Hardhat 账号 #2
		"0x90F79bf6EB2c4f870365E785982E1f101E93b906", // Hardhat 账号 #3
	}

	// 将字符串地址转换为 common.Address 数组
	addresses := make([]common.Address, len(recipients))
	for i, addr := range recipients {
		addresses[i] = common.HexToAddress(addr)
	}

	// 调用 mintToMultple 函数
	tokenId := big.NewInt(1) // NFT ID
	amount := big.NewInt(1)  // 每个地址接收的数量
	data := []byte{}         // 额外数据（空）

	txHash, err := s.NftService.CallContractFunction(
		ContractCallParams{
			ContractAddress: contractAddress,
			ContractABI:     string(abiBytes),
			FunctionName:    "mintToMultple",
			FunctionArgs:    []any{addresses, tokenId, amount, data},
		},
	)
	s.Require().NoError(err)
	s.Require().NotEmpty(txHash)
}
