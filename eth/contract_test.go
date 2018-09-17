package eth

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/stretchr/testify/assert"
)

/*

pragma solidity ^0.4.24;

contract Counterfactual {
    uint256 public iv;
    constructor (uint256 _iv) public {
        iv = _iv;
    }
}
*/

const bytecodestr = "608060405234801561001057600080fd5b506040516020806100cc83398101604052516000556099806100336000396000f300608060405260043610603e5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416635f64335881146043575b600080fd5b348015604e57600080fd5b5060556067565b60408051918252519081900360200190f35b600054815600a165627a7a72305820b300b91cec4bc3007d5feb0f1f8b9a465e21a75e18f1b014ea1ac7756d42a3800029"

const abistr = `
[
	{
		"constant": true,
		"inputs": [],
		"name": "iv",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "_iv",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]
`

func TestConterfactual(t *testing.T) {

	client := new(ClientMock)
	client.On("NetworkID").Return(big.NewInt(4), nil)

	abiobj, err := abi.JSON(strings.NewReader(abistr))
	assert.Nil(t, err)
	bytecode, err := hex.DecodeString(bytecodestr)
	assert.Nil(t, err)

	contract := NewContract(client, &abiobj, bytecode, nil)

	gasPrice := big.NewInt(0)
	gasPrice.SetString("4000000000", 10) // 4 gigawei
	gasLimit := uint64(2000000)

	creatoraddraddr, contractaddr, _, err := contract.Conterfactual(gasLimit, gasPrice, big.NewInt(1001))
	assert.Nil(t, err)

	assert.Equal(t, "0xa8FF5dD292C3f9790007d112559820f7ceB0D020", creatoraddraddr.Hex())
	assert.Equal(t, "0x9bc449f23f16bF48D95AE6F6BE43b900cc514034", contractaddr.Hex())
}
