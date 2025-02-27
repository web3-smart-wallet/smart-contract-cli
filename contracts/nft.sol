// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.22;

import {ERC1155} from "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import {ERC1155Supply} from "@openzeppelin/contracts/token/ERC1155/extensions/ERC1155Supply.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is ERC1155, Ownable {
    uint16 public constant batchSize = 50;

    // 记录所有tokenId
    uint256[] private _tokenIds;
    mapping(uint256 => uint256) private _tokenIdToAmount;

    constructor(address initialOwner) ERC1155("") Ownable(initialOwner) {}

    function setURI(string memory newuri) public onlyOwner {
        _setURI(newuri);
    }

    
    function mint(address account, uint256 id, uint256 amount, bytes memory data)
        public
        onlyOwner
    {
        _mint(account, id, amount, data);
    }

    function mintToMultple(
        address[] memory accounts,
        uint256 ids,
        uint256 amounts,
        bytes memory data 
    ) public onlyOwner {
        require(accounts.length <= batchSize, "Batch size exceeds the limit");

        for (uint256 i = 0; i < accounts.length; i++) {
            _mint(accounts[i], ids, amounts, data);
        }
    }

    function mintBatch(address to, uint256[] memory ids, uint256[] memory amounts, bytes memory data)
        public  
        onlyOwner
    {
        _mintBatch(to, ids, amounts, data);
    }

}
