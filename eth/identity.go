package eth

import (
	"bytes"
	"fmt"
	"math/big"

	common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type IdentityManager struct {
	client   Client
	deployer *Contract
	impl     *Contract
	proxy    *Contract
}
type Identity struct {
	operational common.Address
	relayer     common.Address
	recovery    common.Address
	revoke      common.Address
	impl        common.Address
}

func NewIdentityManager(client Client, deployer, impl *common.Address) *IdentityManager {
	deployerAbi, deployerCode := deployerContract.mustGet()
	proxyAbi, proxyCode := iden3proxyContract.mustGet()
	implAbi, implCode := iden3implContract.mustGet()

	return &IdentityManager{
		client:   client,
		deployer: NewContract(client, &deployerAbi, deployerCode, deployer),
		proxy:    NewContract(client, &proxyAbi, proxyCode, nil),
		impl:     NewContract(client, &implAbi, implCode, impl),
	}
}

func (i *IdentityManager) Initialized() bool {
	return i.deployer.Address() != nil && i.impl.Address() != nil
}

func (i *IdentityManager) Initialize() error {

	_, _, err := i.deployer.DeploySync()
	if err != nil {
		return err
	}
	_, _, err = i.impl.DeploySync()
	if err != nil {
		return err
	}

	log.Info("Deployer created at ", i.deployer.Address().Hex())
	log.Info("Implementation created at ", i.impl.Address().Hex())
	return nil
}

func (m *IdentityManager) codeAndAddress(id *Identity) (common.Address, []byte, error) {
	code, err := m.deployer.CreationBytes(
		id.operational,
		id.relayer,
		id.recovery,
		id.revoke,
		id.impl,
	)
	if err != nil {
		return common.Address{}, nil, err
	}
	return crypto.CreateAddress2(
		*m.deployer.Address(),
		common.BigToHash(big.NewInt(0)),
		code,
	), code, nil
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
	deployedcode, err := m.client.CodeAt(addr)
	if err != nil {
		return common.Address{}, err
	}
	if bytes.Compare(code, deployedcode) != 0 {
		return common.Address{}, fmt.Errorf("Bad deployed code")
	}

	return addr, nil
}
