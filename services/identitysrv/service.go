package identitysrv

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/utils"

	log "github.com/sirupsen/logrus"
)

const CompressedPkLen = 33

type Service interface {
	Initialized() bool
	AddressOf(id *Identity) (common.Address, error)
	Deploy(id *Identity) (common.Address, *types.Transaction, error)
	IsDeployed(idaddr common.Address) (bool, error)
	Info(idaddr common.Address) (*Info, error)
	Forward(idaddr common.Address, ksignpk *ecdsa.PublicKey, to common.Address, data []byte, value *big.Int, gas uint64, sig []byte) (common.Hash, error)
	Add(id *Identity) (*core.ProofClaim, error)
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
	Operational   common.Address   `json:"operational"`
	OperationalPk *utils.PublicKey `json:"operationalPk"`
	Relayer       common.Address   `json:"relayer"`
	Recoverer     common.Address   `json:"recoverer"`
	Revokator     common.Address   `json:"revokator"`
	Impl          common.Address   `json:"impl"`
}

func (i *Identity) Encode() []byte {
	var b bytes.Buffer
	b.Write(i.Operational[:])
	b.Write(crypto.CompressPubkey(&i.OperationalPk.PublicKey))
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
	var operationalPkComp [CompressedPkLen]byte
	if _, err := b.Read(operationalPkComp[:]); err != nil {
		return err
	}
	if pk, err := crypto.DecompressPubkey(operationalPkComp[:]); err != nil {
		return err
	} else {
		i.OperationalPk = &utils.PublicKey{PublicKey: *pk}
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

func (is *ServiceImpl) Initialized() bool {
	return is.deployer.Address() != nil && is.impl.Address() != nil
}

func (is *ServiceImpl) codeAndAddress(id *Identity) (common.Address, []byte, error) {
	code, err := is.proxy.CreationBytes(
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
		*is.deployer.Address(),
		common.BigToHash(big.NewInt(0)),
		code,
	)

	return addr, code, nil
}

// AddressOf returns the address of the smart contract given the identity data
// of a user.
func (is *ServiceImpl) AddressOf(id *Identity) (common.Address, error) {
	addr, _, err := is.codeAndAddress(id)
	return addr, err
}

// IsDeployed checks if the users's smart contract is deployed in the blockchain.
func (is *ServiceImpl) IsDeployed(idaddr common.Address) (bool, error) {
	deployedcode, err := is.deployer.Client().CodeAt(idaddr)
	if err != nil {
		return false, err
	}
	if len(deployedcode) == 0 {
		return false, nil
	}
	return true, nil
}

// Deploy deploys the user's smart contract in the blockchain.
func (is *ServiceImpl) Deploy(id *Identity) (common.Address, *types.Transaction, error) {

	addr, code, err := is.codeAndAddress(id)
	if err != nil {
		return common.Address{}, nil, err
	}
	tx, err := is.deployer.SendTransaction(nil, 0, "create", code)
	if err != nil {
		return common.Address{}, nil, err
	}
	return addr, tx, nil
}

func (is *ServiceImpl) Info(idaddr common.Address) (*Info, error) {

	var info Info

	code, err := is.impl.Client().CodeAt(idaddr)
	if err != nil {
		return nil, err
	}
	if code == nil || len(code) == 0 {
		return nil, nil
	}

	info.Codehash = sha256.Sum256(code)

	if err := is.impl.At(&idaddr).Call(&info, "info"); err != nil {
		return nil, err
	}
	if err := is.impl.At(&idaddr).Call(&info.LastNonce, "lastNonce"); err != nil {
		return nil, err
	}
	return &info, nil

}

func (is *ServiceImpl) Forward(
	idaddr common.Address,
	ksignpk *ecdsa.PublicKey,
	to common.Address,
	data []byte,
	value *big.Int,
	gas uint64,
	sig []byte,
) (common.Hash, error) {

	ksignclaim := core.NewClaimAuthorizeKSignSecp256k1(ksignpk)
	proof, err := is.cs.GetClaimProofUserByHiOld(idaddr, *ksignclaim.Entry().HIndex())
	if err != nil {
		log.Warn("Error retieving proof ", err)
		return common.Hash{}, err
	}

	auth := packAuth(
		proof.ClaimProof.Leaf,
		proof.ClaimProof.Root[:],
		proof.ClaimProof.Proof,
		proof.ClaimNonRevocationProof.Proof,

		proof.SetRootClaimProof.Leaf,
		proof.SetRootClaimProof.Root[:],
		proof.SetRootClaimProof.Proof,
		proof.SetRootClaimNonRevocationProof.Proof,

		proof.Date, proof.Signature,
	)
	proxy := is.impl.At(&idaddr)

	tx, err := proxy.SendTransaction(
		big.NewInt(0), 4000000,
		"forward",
		to, data, value, big.NewInt(int64(gas)), sig, auth,
	)
	if err == nil {
		_, err = proxy.Client().WaitReceipt(tx.Hash())
		return tx.Hash(), err
	}

	return common.Hash{}, err

}

// Add creates a merkle tree of a new user in the relay, given the identity
// data of the user.
func (is *ServiceImpl) Add(id *Identity) (*core.ProofClaim, error) {
	var err error

	idaddr, _, err := is.codeAndAddress(id)
	if err != nil {
		return nil, err
	}

	if _, err := is.sto.Get(idaddr[:]); err == nil {
		return nil, fmt.Errorf("the identity %v with id %+v already exists in the Relay", idaddr, *id)
	}

	tx, err := is.sto.NewTx()
	if err != nil {
		return nil, err
	}

	// store identity
	tx.Put(idaddr[:], id.Encode())
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	claim := core.NewClaimAuthorizeKSignSecp256k1(&id.OperationalPk.PublicKey)
	err = is.cs.AddClaimAuthorizeKSignSecp256k1First(idaddr, *claim)
	if err != nil {
		return nil, err
	}

	return is.cs.GetClaimProofUserByHi(idaddr, claim.Entry().HIndex())
}

func (is *ServiceImpl) List(limit int) ([]common.Address, error) {

	kvs, err := is.sto.List(limit)
	if err != nil {
		return nil, err
	}
	addrs := make([]common.Address, 0, len(kvs))
	for _, e := range kvs {
		var addr common.Address
		copy(addr[:], e.K)
		addrs = append(addrs, addr)
	}
	return addrs, err
}

func (is *ServiceImpl) Get(idaddr common.Address) (*Identity, error) {

	data, err := is.sto.Get(idaddr[:])
	if err != nil {
		return nil, err
	}
	var id Identity
	err = id.Decode(data)
	return &id, err
}

func (is *ServiceImpl) DeployerAddr() *common.Address {
	return is.deployer.Address()
}

func (is *ServiceImpl) ImplAddr() *common.Address {
	return is.impl.Address()
}

func packAuth(
	kclaimBytes, kclaimRoot, kclaimExistenceProof, kclaimNonNextExistenceProof []byte,
	rclaimBytes, rclaimRoot, rclaimExistenceProof, rclaimNonNextExistenceProof []byte,
	rclaimSigDate int64,
	rclaimSigRSV []byte) []byte {

	var b bytes.Buffer
	writeBytes := func(v []byte) {
		var vlen [2]byte
		binary.BigEndian.PutUint16(vlen[:], uint16(len(v)))
		b.Write(vlen[:])
		b.Write(v)
	}
	writeUint64 := func(v uint64) {
		var val [8]byte
		binary.BigEndian.PutUint64(val[:], v)
		b.Write(val[:])
	}

	writeBytes(kclaimBytes)
	b.Write(kclaimRoot)
	writeBytes(kclaimExistenceProof)
	writeBytes(kclaimNonNextExistenceProof)
	writeBytes(rclaimBytes)
	b.Write(rclaimRoot)
	writeBytes(rclaimExistenceProof)
	writeBytes(rclaimNonNextExistenceProof)

	writeUint64(uint64(rclaimSigDate))
	writeBytes(rclaimSigRSV)

	return b.Bytes()
}
