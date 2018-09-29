package identitysrv

import (
	"bytes"
	"encoding/binary"
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
	IsDeployed(idaddr common.Address) (bool, error)
	Info(idaddr common.Address) (*Info, error)
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
	Recoverer   common.Address
	Revokator   common.Address
	Impl        common.Address
}

type Info struct {
	Impl          common.Address
	Recoverer     common.Address
	RecovererProp common.Address
	Revoker       common.Address
	Relay         common.Address
	LastNonce     *big.Int
}

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

func (s *ServiceImpl) codeAndAddress(id *Identity) (common.Address, []byte, error) {
	code, err := s.proxy.CreationBytes(
		id.Operational,
		id.Relayer,
		id.Recoverer,
		id.Revokator,
		id.Impl,
	)
	if err != nil {
		return common.Address{}, nil, err
	}
	addr := crypto.CreateAddress2(
		*s.deployer.Address(),
		common.BigToHash(big.NewInt(0)),
		code,
	)

	return addr, code, nil
}

func (m *ServiceImpl) AddressOf(id *Identity) (common.Address, error) {
	addr, _, err := m.codeAndAddress(id)
	return addr, err
}

func (m *ServiceImpl) IsDeployed(idaddr common.Address) (bool, error) {
	deployedcode, err := m.deployer.Client().CodeAt(idaddr)
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

func (s *ServiceImpl) Info(idaddr common.Address) (*Info, error) {

	var info Info
	if err := s.impl.At(&idaddr).Call(&info, "info"); err != nil {
		return nil, err
	}
	if err := s.impl.At(&idaddr).Call(&info.LastNonce, "lastNonce"); err != nil {
		return nil, err
	}
	return &info, nil

}

func (s *ServiceImpl) Forward(
	id *Identity,
	to common.Address,
	data []byte,
	value *big.Int,
	gas *big.Int,
	sig []byte,
	auth []byte,
) (common.Hash, error) {

	addr, _, err := s.codeAndAddress(id)
	if err != nil {
		return common.Hash{}, err
	}

	proxy := s.impl.At(&addr)
	tx, _, err := proxy.SendTransactionSync(
		big.NewInt(0), 0,
		"forward",
		to, data, value, gas, sig, auth,
	)

	return tx.Hash(), err
}

func (s *ServiceImpl) DeployerAddr() *common.Address {
	return s.deployer.Address()
}

func (s *ServiceImpl) ImplAddr() *common.Address {
	return s.impl.Address()
}

func PackAuth(
	kclaimBytes, kclaimRoot, kclaimExistenceProof, kclaimNonNextExistenceProof []byte,
	rclaimBytes, rclaimRoot, rclaimExistenceProof, rclaimNonNextExistenceProof []byte,
	rclaimSigDate uint64,
	rclaimSigR, rclaimSigS []byte, rclaimSigV uint8) []byte {

	var b bytes.Buffer
	b.Write(kclaimBytes)
	b.Write(kclaimRoot)
	b.Write(kclaimExistenceProof)
	b.Write(kclaimNonNextExistenceProof)
	b.Write(rclaimBytes)
	b.Write(rclaimRoot)
	b.Write(rclaimExistenceProof)
	b.Write(rclaimNonNextExistenceProof)

	var v [4]byte
	binary.LittleEndian.PutUint64(v[:], rclaimSigDate)
	b.Write(v[:])

	b.Write(rclaimSigR)
	b.Write(rclaimSigS)
	b.Write([]byte{rclaimSigV})

	return b.Bytes()
}
