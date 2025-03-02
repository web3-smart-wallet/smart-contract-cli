package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ContractArtifact struct {
	Bytecode string `json:"bytecode"`
	Abi      string `json:"abi"`
}

type DeployedContract struct {
	Address    string    `json:"address"`
	Abi        string    `json:"abi"`
	DeployTime time.Time `json:"deploy_time"`
}

type DeployedContracts struct {
	Contracts []DeployedContract `json:"contracts"`
}

type ContractCompiler struct {
	artifactsPath string
}

func NewContractCompiler(artifactsPath string) *ContractCompiler {
	return &ContractCompiler{
		artifactsPath: artifactsPath,
	}
}

// GetContractBytecode 从编译后的合约文件中读取字节码
func (c *ContractCompiler) GetContractBytecode() (string, string, error) {
	// 更新合约文件路径
	artifactPath := filepath.Join(c.artifactsPath, "contracts", "nft.sol", "MyToken.json")

	data, err := os.ReadFile(artifactPath)
	if err != nil {
		return "", "", fmt.Errorf("读取合约文件失败: %v", err)
	}

	var artifact ContractArtifact
	if err := json.Unmarshal(data, &artifact); err != nil {
		return "", "", fmt.Errorf("解析合约文件失败: %v", err)
	}

	if artifact.Bytecode == "" {
		return "", "", fmt.Errorf("合约字节码为空")
	}

	return artifact.Bytecode, artifact.Abi, nil
}

// SaveDeployedContract 保存已部署的合约信息
func (c *ContractCompiler) SaveDeployedContract(address string, abi string) error {
	deployedFile := "deployed_contracts.json"
	var deployedContracts DeployedContracts

	// 如果文件存在，读取现有内容
	if _, err := os.Stat(deployedFile); err == nil {
		data, err := os.ReadFile(deployedFile)
		if err != nil {
			return fmt.Errorf("读取已部署合约文件失败: %v", err)
		}

		if err := json.Unmarshal(data, &deployedContracts); err != nil {
			return fmt.Errorf("解析已部署合约文件失败: %v", err)
		}
	}

	// 添加新部署的合约信息
	newContract := DeployedContract{
		Address:    address,
		Abi:        abi,
		DeployTime: time.Now(),
	}
	deployedContracts.Contracts = append(deployedContracts.Contracts, newContract)

	// 保存到JSON文件
	data, err := json.MarshalIndent(deployedContracts, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化合约信息失败: %v", err)
	}

	if err := os.WriteFile(deployedFile, data, 0644); err != nil {
		return fmt.Errorf("保存合约信息失败: %v", err)
	}

	return nil
}

// GetDeployedContracts 获取所有已部署的合约信息
func (c *ContractCompiler) GetDeployedContracts() ([]DeployedContract, error) {
	deployedFile := "deployed_contracts.json"

	// 如果文件不存在，返回空列表
	if _, err := os.Stat(deployedFile); os.IsNotExist(err) {
		return []DeployedContract{}, nil
	}

	data, err := os.ReadFile(deployedFile)
	if err != nil {
		return nil, fmt.Errorf("读取已部署合约文件失败: %v", err)
	}

	var deployedContracts DeployedContracts
	if err := json.Unmarshal(data, &deployedContracts); err != nil {
		return nil, fmt.Errorf("解析已部署合约文件失败: %v", err)
	}

	return deployedContracts.Contracts, nil
}

// GetLatestDeployedContract 获取最新部署的合约信息
func (c *ContractCompiler) GetLatestDeployedContract() (*DeployedContract, error) {
	contracts, err := c.GetDeployedContracts()
	if err != nil {
		return nil, err
	}

	if len(contracts) == 0 {
		return nil, fmt.Errorf("没有找到已部署的合约")
	}

	// 返回最后一个部署的合约
	return &contracts[len(contracts)-1], nil
}
