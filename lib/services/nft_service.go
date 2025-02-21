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
