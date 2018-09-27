package identitysrv

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3/eth"

	log "github.com/sirupsen/logrus"
)

type Service interface {
	Initialized() bool
	AddressOf(id *Identity) (common.Address, error)
	Deploy(id *Identity) (common.Address, error)
	Forward(id *Identity, to common.Address, data []byte, value *big.Int, gas *big.Int, sig []byte, auth []byte) (common.Hash, error)
	DeployerAddr() *common.Address
	ImplAddr() *common.Address
}

type ServiceImpl struct {
	deployer *eth.Contract
	impl     *eth.Contract
	proxy    *eth.Contract
}

type Identity struct {
	Operational common.Address
	Relayer     common.Address
	Recovery    common.Address
	Revoke      common.Address
	Impl        common.Address
}

const (
	deployerContract = "deployer"
	proxyContract    = "proxy"
	implContract     = "impl"
)

func New(deployer, impl, proxy *eth.Contract) *ServiceImpl {
	return &ServiceImpl{
		deployer: deployer,
		proxy:    proxy,
		impl:     impl,
	}
}

func (i *ServiceImpl) Initialized() bool {
	return i.deployer.Address() != nil && i.impl.Address() != nil
}

func (m *ServiceImpl) codeAndAddress(id *Identity) (common.Address, []byte, error) {
	code, err := m.proxy.CreationBytes(
		id.Operational,
		id.Relayer,
		id.Recovery,
		id.Revoke,
		id.Impl,
	)
	if err != nil {
		return common.Address{}, nil, err
	}
	addr := crypto.CreateAddress2(
		*m.deployer.Address(),
		common.BigToHash(big.NewInt(0)),
		code,
	)
	log.Info("caller=", m.deployer.Address().Hex(),
		" salt=", common.BigToHash(big.NewInt(0)).Hex(),
		" code=", hex.EncodeToString(code),
		" adress=", addr.Hex())

	return addr, code, nil
}

func (m *ServiceImpl) AddressOf(id *Identity) (common.Address, error) {
	addr, _, err := m.codeAndAddress(id)
	return addr, err
}

func (m *ServiceImpl) Deployed(id *Identity) (bool, error) {
	addr, _, err := m.codeAndAddress(id)
	if err != nil {
		return false, err
	}
	deployedcode, err := m.deployer.Client().CodeAt(addr)
	if err != nil {
		return false, err
	}
	if len(deployedcode) == 0 {
		return false, nil
	}
	return true, nil
}

func (m *ServiceImpl) Deploy(id *Identity) (common.Address, error) {

	addr, code, err := m.codeAndAddress(id)
	if err != nil {
		return common.Address{}, err
	}
	_, _, err = m.deployer.SendTransactionSync(nil, 0, "create", code)
	if err != nil {
		return common.Address{}, err
	}
	log.Info("Deployed identity at ", addr.Hex())
	return addr, nil
}

func (m *ServiceImpl) Ping(id *Identity) error {
	addr, err := m.AddressOf(id)
	if err != nil {
		return err
	}
	var pong string
	err = m.impl.At(&addr).Call(&pong, "ping")
	if err != nil {
		return err
	}
	if pong != "pong" {
		return fmt.Errorf("Not returned pong (%v instead)", pong)
	}
	return nil
}

func (m *ServiceImpl) Forward(
	id *Identity,
	to common.Address,
	data []byte,
	value *big.Int,
	gas *big.Int,
	sig []byte,
	auth []byte,
) (common.Hash, error) {

	addr, _, err := m.codeAndAddress(id)
	if err != nil {
		return common.Hash{}, err
	}

	proxy := m.impl.At(&addr)
	tx, _, err := proxy.SendTransactionSync(
		big.NewInt(0), 0,
		"forward",
		to, data, value, gas, sig, auth,
	)

	return tx.Hash(), err
}

func (m *ServiceImpl) DeployerAddr() *common.Address {
	return m.deployer.Address()
}
func (m *ServiceImpl) ImplAddr() *common.Address {
	return m.impl.Address()
}
