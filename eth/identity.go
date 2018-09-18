package eth

import (
	"math/big"

	common "github.com/ethereum/go-ethereum/common"
)

type IdentityManager struct {
	client             Client
	iden3impl          *Contract
	iden3delegateproxy *Contract
}
type Identity struct {
	gasLimit     uint64
	gasPrice     *big.Int
	operationals []common.Address
	relayer      common.Address
	recovery     common.Address
	impl         common.Address
}

func (i *IdentityManager) DeployBase() {
	_, _, err := i.iden3impl.DeploySync()
}

func (m *IdentityManager) Deploy(id *Identity) {

	creator, contract, rawtx, err := m.iden3delegateproxy.Conterfactual(
		id.gasLimit, id.gasPrice,
		id.operationals,
		id.relayer,
		id.recovery,
		id.impl,
	)

	// top up creator
	tx, r, err := m.client.SendTransactionSync(creator, big.NewInt(299), 21000, []byte{})

	// send the raw tx
	m.client.SendRawTx(rawtx)
}

func (m *IdentityManager) AddressOf(id *Identity) common.Address {
	_, contract, _, _ := m.iden3delegateproxy.Conterfactual(
		id.gasLimit, id.gasPrice,
		id.operationals,
		id.relayer,
		id.recovery,
		id.impl,
	)
	return contract
}
