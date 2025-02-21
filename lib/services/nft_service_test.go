package services

import (
	"encoding/json"
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
