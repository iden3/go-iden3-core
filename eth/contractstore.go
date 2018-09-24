package eth

import "github.com/ethereum/go-ethereum/accounts/abi"

type ContractStore interface {
	Get(name string) (*abi.ABI, []byte, error)
}
