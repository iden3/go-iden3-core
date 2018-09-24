package eth

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

type IdentityManager struct {
	client   Client
	deployer *Contract
	impl     *Contract
	proxy    *Contract
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

func NewIdentityManager(client Client, store ContractStore, deployer, impl *common.Address) (*IdentityManager, error) {

	deployerAbi, deployerCode, err := store.Get(deployerContract)
	if err != nil {
		return nil, err
	}
	proxyAbi, proxyCode, err := store.Get(proxyContract)
	if err != nil {
		return nil, err
	}
	implAbi, implCode, err := store.Get(implContract)
	if err != nil {
		return nil, err
	}

	return &IdentityManager{
		client:   client,
		deployer: NewContract(client, deployerAbi, deployerCode, deployer),
		proxy:    NewContract(client, proxyAbi, proxyCode, nil),
		impl:     NewContract(client, implAbi, implCode, impl),
	}, err
}

func (i *IdentityManager) Initialized() bool {
	return i.deployer.Address() != nil && i.impl.Address() != nil
}

func (i *IdentityManager) Initialize() error {

	log.Info("Deploying deployer")
	_, _, err := i.deployer.DeploySync()
	if err != nil {
		return err
	}
	log.Info("Deploying implementation")
	_, _, err = i.impl.DeploySync()
	if err != nil {
		return err
	}

	log.Info("Deployer created at ", i.deployer.Address().Hex())
	log.Info("Implementation created at ", i.impl.Address().Hex())
	return nil
}

func (m *IdentityManager) codeAndAddress(id *Identity) (common.Address, []byte, error) {
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

func (m *IdentityManager) AddressOf(id *Identity) (common.Address, error) {
	addr, _, err := m.codeAndAddress(id)
	return addr, err
}

func (m *IdentityManager) Deployed(id *Identity) (bool, error) {
	addr, code, err := m.codeAndAddress(id)
	if err != nil {
		return false, err
	}
	deployedcode, err := m.client.CodeAt(addr)
	if err != nil {
		return false, err
	}
	if len(deployedcode) == 0 {
		return false, nil
	}
	if bytes.Compare(code, deployedcode) != 0 {
		return false, fmt.Errorf("Bad deployed code")
	}
	return true, nil
}

func (m *IdentityManager) Deploy(id *Identity) (common.Address, error) {

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

func (m *IdentityManager) Ping(id *Identity) error {
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

func (m *IdentityManager) Forward(
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

func (m *IdentityManager) DeployerAddr() *common.Address {
	return m.deployer.Address()
}
func (m *IdentityManager) ImplAddr() *common.Address {
	return m.impl.Address()
}
