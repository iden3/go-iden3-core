package idenpubonchain

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/eth/contracts"
)

type StateDepsAddresses struct {
	Verifier *common.Address
}

type DeployStateResult struct {
	Verifier eth.ContractData
	State    eth.ContractData
}

func DeployState(client *eth.Client, depsAddrs *StateDepsAddresses) (DeployStateResult, error) {
	if depsAddrs == nil {
		depsAddrs = &StateDepsAddresses{}
	}
	var result DeployStateResult

	if depsAddrs.Verifier == nil {
		if contract, err := client.Deploy("Verifier",
			func(c *ethclient.Client, auth *bind.TransactOpts) (common.Address,
				*types.Transaction, interface{}, error) {
				return contracts.DeployVerifier(auth, c)
			}); err != nil {
			return result, err
		} else {
			result.Verifier = contract
		}
	} else {
		result.Verifier.Address = *depsAddrs.Verifier
	}

	if contract, err := client.Deploy("State",
		func(c *ethclient.Client, auth *bind.TransactOpts) (common.Address,
			*types.Transaction, interface{}, error) {
			return contracts.DeployState(auth, c, result.Verifier.Address)
		}); err != nil {
		return result, err
	} else {
		result.State = contract
	}

	return result, nil
}
