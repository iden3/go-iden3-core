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
	PoseidonCircomlib *common.Address
	Poseidon          *common.Address
	EddsaBabyJub      *common.Address
}

type DeployStateResult struct {
	PoseidonCircomlib eth.ContractData
	Poseidon          eth.ContractData
	EddsaBabyJub      eth.ContractData
	State             eth.ContractData
}

func DeployState(client *eth.Client, depsAddrs *StateDepsAddresses) (DeployStateResult, error) {
	if depsAddrs == nil {
		depsAddrs = &StateDepsAddresses{}
	}
	var result DeployStateResult

	if depsAddrs.PoseidonCircomlib == nil {
		if contract, err := client.Deploy("PoseidonCircomlib",
			func(c *ethclient.Client, auth *bind.TransactOpts) (common.Address,
				*types.Transaction, interface{}, error) {
				return contracts.DeployPoseidonUnit(auth, c)
			}); err != nil {
			return result, err
		} else {
			result.PoseidonCircomlib = contract
		}
	} else {
		result.PoseidonCircomlib.Address = *depsAddrs.PoseidonCircomlib
	}

	if depsAddrs.Poseidon == nil {
		if contract, err := client.Deploy("Poseidon",
			func(c *ethclient.Client, auth *bind.TransactOpts) (common.Address,
				*types.Transaction, interface{}, error) {
				return contracts.DeployPoseidon(auth, c, result.PoseidonCircomlib.Address)
			}); err != nil {
			return result, err
		} else {
			result.Poseidon = contract
		}
	} else {
		result.Poseidon.Address = *depsAddrs.Poseidon
	}

	if depsAddrs.EddsaBabyJub == nil {
		if contract, err := client.Deploy("EddsaBabyJub",
			func(c *ethclient.Client, auth *bind.TransactOpts) (common.Address,
				*types.Transaction, interface{}, error) {
				return contracts.DeployEddsaBabyJubJub(auth, c, result.PoseidonCircomlib.Address)
			}); err != nil {
			return result, err
		} else {
			result.EddsaBabyJub = contract
		}
	} else {
		result.EddsaBabyJub.Address = *depsAddrs.EddsaBabyJub
	}

	if contract, err := client.Deploy("State",
		func(c *ethclient.Client, auth *bind.TransactOpts) (common.Address,
			*types.Transaction, interface{}, error) {
			return contracts.DeployState(auth, c, result.PoseidonCircomlib.Address, result.EddsaBabyJub.Address)
		}); err != nil {
		return result, err
	} else {
		result.State = contract
	}

	return result, nil
}
