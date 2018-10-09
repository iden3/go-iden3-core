package identitysrv

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"math"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type Service interface {
	Initialized() bool
	AddressOf(id *Identity) (common.Address, error)
	Deploy(id *Identity) (common.Address, *types.Transaction, error)
	IsDeployed(idaddr common.Address) (bool, error)
	Info(idaddr common.Address) (*Info, error)
	Forward(idaddr common.Address, to common.Address, data []byte, value *big.Int, gas *big.Int, sig []byte, auth []byte) (common.Hash, error)
	Add(id *Identity) error
	List(limit int) ([]common.Address, error)
	Get(idaddr common.Address) (*Identity, error)
	DeployerAddr() *common.Address
	ImplAddr() *common.Address
}

type ServiceImpl struct {
	deployer *eth.Contract
	impl     *eth.Contract
	proxy    *eth.Contract
	cs       claimsrv.Service
	sto      db.Storage
}

type Identity struct {
	Operational common.Address
	Relayer     common.Address
	Recoverer   common.Address
	Revokator   common.Address
	Impl        common.Address
}

func (i *Identity) Encode() []byte {
	var b bytes.Buffer
	b.Write(i.Operational[:])
	b.Write(i.Relayer[:])
	b.Write(i.Recoverer[:])
	b.Write(i.Revokator[:])
	b.Write(i.Impl[:])
	return b.Bytes()
}
func (i *Identity) Decode(encoded []byte) error {
	b := bytes.NewBuffer(encoded)
	if _, err := b.Read(i.Operational[:]); err != nil {
		return err
	}
	if _, err := b.Read(i.Relayer[:]); err != nil {
		return err
	}
	if _, err := b.Read(i.Recoverer[:]); err != nil {
		return err
	}
	if _, err := b.Read(i.Revokator[:]); err != nil {
		return err
	}
	if _, err := b.Read(i.Impl[:]); err != nil {
		return err
	}
	return nil
}

type Info struct {
	Codehash      common.Hash
	Impl          common.Address
	Recoverer     common.Address
	RecovererProp common.Address
	Revoker       common.Address
	Relay         common.Address
	LastNonce     *big.Int
}

func New(deployer, impl, proxy *eth.Contract, cs claimsrv.Service, sto db.Storage) *ServiceImpl {
	return &ServiceImpl{
		deployer: deployer,
		proxy:    proxy,
		impl:     impl,
		cs:       cs,
		sto:      sto,
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

func (m *ServiceImpl) Deploy(id *Identity) (common.Address, *types.Transaction, error) {

	addr, code, err := m.codeAndAddress(id)
	if err != nil {
		return common.Address{}, nil, err
	}
	tx, err := m.deployer.SendTransaction(nil, 0, "create", code)
	if err != nil {
		return common.Address{}, nil, err
	}
	return addr, tx, nil
}

func (s *ServiceImpl) Info(idaddr common.Address) (*Info, error) {

	var info Info
	
	code, err := s.impl.Client().CodeAt(idaddr)
	if err != nil {
		return nil, err
	}
	if code == nil || len(code)==0 {
		return nil, nil
	}

	info.Codehash = sha256.Sum256(code)

	if err := s.impl.At(&idaddr).Call(&info, "info"); err != nil {
		return nil, err
	}
	if err := s.impl.At(&idaddr).Call(&info.LastNonce, "lastNonce"); err != nil {
		return nil, err
	}
	return &info, nil

}

func (s *ServiceImpl) Forward(
	idaddr common.Address,
	to common.Address,
	data []byte,
	value *big.Int,
	gas *big.Int,
	sig []byte,
	auth []byte,
) (common.Hash, error) {

	proxy := s.impl.At(&idaddr)
	tx, err := proxy.SendTransaction(
		big.NewInt(0), 0,
		"forward",
		to, data, value, gas, sig, auth,
	)
	if err == nil {
		_, err = proxy.Client().WaitReceipt(tx.Hash())
	}

	return tx.Hash(), err
}

func (s *ServiceImpl) Add(id *Identity) error {

	var err error

	idaddr, _, err := s.codeAndAddress(id)
	if err != nil {
		return err
	}

	tx, err := s.sto.NewTx()
	if err != nil {
		return err
	}

	// store identity
	tx.Put(idaddr[:], id.Encode())
	if err = tx.Commit(); err != nil {
		return err
	}

	claim := core.NewOperationalKSignClaim("iden3.io", id.Operational, 0, math.MaxUint64)
	return s.cs.AddAuthorizeKSignClaimFirst(idaddr,claim)
}

func (m *ServiceImpl) List(limit int) ([]common.Address, error) {

	kvs,err := m.sto.List(limit)
	if err != nil {
		return nil, err
	}
	addrs := make([]common.Address,0,len(kvs))
	for _,e := range kvs {
		var addr common.Address
		copy(addr[:],e.K)
		addrs = append(addrs,addr)
	}
	return addrs, err
}

func (m *ServiceImpl) Get(idaddr common.Address) (*Identity, error) {

	data, err := m.sto.Get(idaddr[:])
	if err != nil {
		return nil,err
	}
	var id Identity
	err = id.Decode(data)
	return &id,err
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
