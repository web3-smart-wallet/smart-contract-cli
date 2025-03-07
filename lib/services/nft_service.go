package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NFTServiceInterface interface {
	DeployContractWithABI(params DeployContractParams) (contractAddress string, err error)
}

// DeployContractParams contains all parameters needed for contract deployment
type DeployContractParams struct {
	// The compiled bytecode of the smart contract
	Bytecode string
	// The ABI of the constructor
	ConstructorABI string
	// Arguments to pass to the constructor
	ConstructorArgs []any
	// Optional gas limit (if not set, default 3000000 will be used)
	GasLimit uint64
	// Optional value to send with deployment (if not set, 0 will be used)
	Value *big.Int
	// 添加新字段
	InitialURI   string
	InitialOwner string
}

// ContractCallParams contains all parameters needed for contract function calls
type ContractCallParams struct {
	// The contract address to interact with
	ContractAddress string
	// The ABI of the contract
	ContractABI string
	// The name of the function to call
	FunctionName string
	// Arguments to pass to the function
	FunctionArgs []any
	// Optional gas limit (if not set, default 300000 will be used)
	GasLimit uint64
	// Optional value to send with the transaction (if not set, 0 will be used)
	Value *big.Int
}

type NftService struct {
	rpcUrl     string
	privateKey string
}

func NewNftService(rpcUrl string, privateKey string) *NftService {

	return &NftService{
		rpcUrl:     rpcUrl,
		privateKey: strings.TrimPrefix(privateKey, "0x"),
	}
}

// getKeyPair converts a private key string to ECDSA private key and corresponding public address
// Parameters:
//   - privateKeyStr: The private key in string format (without 0x prefix)
//
// Returns:
//   - *ecdsa.PrivateKey: The ECDSA private key
//   - common.Address: The Ethereum address derived from the public key
//   - error: Any error that occurred during conversion
func (s *NftService) getKeyPair() (privateKey *ecdsa.PrivateKey, fromAddress common.Address, err error) {
	// check if privatekey starts with 0x
	// Convert private key string to ECDSA private key
	privateKey, err = crypto.HexToECDSA(s.privateKey)
	if err != nil {
		return nil, common.Address{}, err
	}

	// Get the public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKey, fromAddress, nil
}

// DeployContractWithABI deploys smart contract to the blockchain with constructor arguments
// Parameters:
//   - params: DeployContractParams struct containing all deployment parameters
//
// Returns:
//   - contractAddress: The address where the contract was deployed
//   - error: Any error that occurred during deployment
func (s *NftService) DeployContractWithABI(params DeployContractParams) (contractAddress string, err error) {
	// Connect to the Ethereum client
	client, err := ethclient.Dial(s.rpcUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Get the key pair
	privateKey, fromAddress, err := s.getKeyPair()
	if err != nil {
		return "", err
	}

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// Decode the bytecode
	decodedBytecode := common.FromHex(params.Bytecode)

	// 获取部署者地址作为默认的 initialOwner（如果未指定）
	if params.InitialOwner == "" {
		params.InitialOwner = fromAddress.Hex()
	}

	// 设置构造函数参数
	params.ConstructorABI = `[{
		"inputs": [
			{"name": "initialOwner", "type": "address"},
			{"name": "newuri", "type": "string"}
		],
		"stateMutability": "nonpayable",
		"type": "constructor"
	}]`

	params.ConstructorArgs = []interface{}{
		common.HexToAddress(params.InitialOwner),
		params.InitialURI,
	}

	// If we have constructor arguments, encode them and append to bytecode
	if len(params.ConstructorArgs) > 0 {
		parsedABI, err := abi.JSON(strings.NewReader(params.ConstructorABI))
		if err != nil {
			return "", fmt.Errorf("failed to parse constructor ABI: %v", err)
		}

		// Pack the constructor arguments
		encodedArgs, err := parsedABI.Pack("", params.ConstructorArgs...)
		if err != nil {
			return "", fmt.Errorf("failed to encode constructor arguments: %v", err)
		}

		// Append encoded arguments to bytecode
		decodedBytecode = append(decodedBytecode, encodedArgs...)
	}

	// Set default gas limit if not provided
	gasLimit := params.GasLimit
	if gasLimit == 0 {
		gasLimit = 3000000
	}

	// Set default value if not provided
	value := params.Value
	if value == nil {
		value = big.NewInt(0)
	}

	// Create transaction data
	tx := types.NewContractCreation(nonce, value, gasLimit, gasPrice, decodedBytecode)

	// Get the chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return "", err
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		return "", err
	}

	// Return the contract address
	return receipt.ContractAddress.Hex(), nil
}

// CallContractFunction executes a function on a deployed smart contract
// Parameters:
//   - params: ContractCallParams struct containing all call parameters
//
// Returns:
//   - txHash: The transaction hash of the executed function call
//   - error: Any error that occurred during the function call
func (s *NftService) CallContractFunction(params ContractCallParams) (txHash string, err error) {
	// Connect to the Ethereum client
	client, err := ethclient.Dial(s.rpcUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Get the key pair
	privateKey, fromAddress, err := s.getKeyPair()
	if err != nil {
		return "", err
	}

	// Parse the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(params.ContractABI))
	if err != nil {
		return "", fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	// Pack the function data
	data, err := parsedABI.Pack(params.FunctionName, params.FunctionArgs...)
	if err != nil {
		return "", fmt.Errorf("failed to pack function data: %v", err)
	}

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// Set default gas limit if not provided
	gasLimit := params.GasLimit
	if gasLimit == 0 {
		gasLimit = 300000
	}

	// Set default value if not provided
	value := params.Value
	if value == nil {
		value = big.NewInt(0)
	}

	// Create transaction data
	contractAddress := common.HexToAddress(params.ContractAddress)
	tx := types.NewTransaction(
		nonce,
		contractAddress,
		value,
		gasLimit,
		gasPrice,
		data,
	)

	// Get the chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return "", err
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		return "", err
	}

	// Return the transaction hash
	return receipt.TxHash.Hex(), nil
}

// MintNFTToAddresses mints NFTs to multiple addresses
func (s *NftService) MintNFTToAddresses(contractAddr string, addresses []string, nftID string) (string, error) {
	// 将字符串地址转换为 common.Address 数组
	recipients := make([]common.Address, len(addresses))
	for i, addr := range addresses {
		recipients[i] = common.HexToAddress(addr)
	}

	// 将 nftID 转换为 big.Int
	// 去除可能的空格
	trimmedNftID := strings.TrimSpace(nftID)
	if trimmedNftID == "" {
		return "", fmt.Errorf("NFT ID不能为空")
	}

	tokenID, ok := new(big.Int).SetString(trimmedNftID, 10)
	if !ok {
		return "", fmt.Errorf("无效的 NFT ID: %s", trimmedNftID)
	}

	// 准备调用参数
	params := ContractCallParams{
		ContractAddress: contractAddr,
		ContractABI: `[{
			"inputs": [
				{"type": "address[]", "name": "accounts"},
				{"type": "uint256", "name": "ids"},
				{"type": "uint256", "name": "amounts"},
				{"type": "bytes", "name": "data"}
			],
			"name": "mintToMultple",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}]`,
		FunctionName: "mintToMultple",
		FunctionArgs: []interface{}{
			recipients,    // 直接传入 []common.Address
			tokenID,       // *big.Int
			big.NewInt(1), // *big.Int
			[]byte{},      // bytes
		},
		GasLimit: 500000, // 增加 gas limit 以处理多个地址
	}

	// 调用合约
	return s.CallContractFunction(params)
}

// SetURI sets the base URI for all tokens
func (s *NftService) SetURI(contractAddr string, newURI string) error {
	params := ContractCallParams{
		ContractAddress: contractAddr,
		ContractABI: `[{
			"inputs": [
				{"type": "string", "name": "newuri"}
			],
			"name": "setURI",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}]`,
		FunctionName: "setURI",
		FunctionArgs: []interface{}{
			newURI,
		},
		GasLimit: 300000,
	}

	// 调用合约
	_, err := s.CallContractFunction(params)
	return err
}
