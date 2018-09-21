package eth

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

/*
   nightly.2018.9.20+commit.2150aea3.Emscripten.clang
   pragma solidity ^0.4.25;

   contract Deployer {
       event Created(address addr);
       function create(bytes memory _code) public {
           address addr;
           assembly {
               addr := create2(0,add(_code,0x20),mload(_code),0)
           }
           require(addr!=address(0x0));
           emit Created(addr);
       }
   }
*/

type contractspec struct {
	abistr  string
	codestr string
}

func (c *contractspec) mustGet() (abi.ABI, []byte) {
	abi, err := abi.JSON(strings.NewReader(c.abistr))
	if err != nil {
		panic(err)
	}
	code, err := hex.DecodeString(c.codestr)
	if err != nil {
		panic(err)
	}
	return abi, code
}

var deployerContract = contractspec{
	abistr:  `[{"constant": false,"inputs": [{"name": "_code","type": "bytes"}],"name": "create","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"anonymous": false,"inputs": [{"indexed": false,"name": "addr","type": "address"}],"name": "Created","type": "event"}]`,
	codestr: "0x608060405234801561001057600080fd5b506101eb806100206000396000f3fe608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063cf5ba53f14610046575b600080fd5b34801561005257600080fd5b5061010c6004803603602081101561006957600080fd5b810190808035906020019064010000000081111561008657600080fd5b82018360208201111561009857600080fd5b803590602001918460018302840111640100000000831117156100ba57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061010e565b005b6000808251602084016000f59050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561015857600080fd5b7f1449abf21e49fd025f33495e77f7b1461caefdd3d4bb646424a3f445c4576a5b81604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390a1505056fea165627a7a72305820c2f5af3c0faf620540eaff3fa36329f67c7b3b8dce8b0db18a6c04a97c3b8f590029",
}

var iden3implContract = contractspec{
	abistr:  "",
	codestr: "0x",
}

var iden3proxyContract = contractspec{
	abistr:  "",
	codestr: "0x",
}
