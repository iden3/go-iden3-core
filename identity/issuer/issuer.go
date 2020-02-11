package issuer

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	ErrIdenPubOnChainNil      = fmt.Errorf("IdenPubOnChain is nil")
	ErrIdenStatePendingNotNil = fmt.Errorf("Update of the published IdenState is pending")
)

var (
	dbPrefixClaimsTree     = []byte("treeclaims:")
	dbPrefixRevocationTree = []byte("treerevocation:")
	dbPrefixRootsTree      = []byte("treeroots:")
	dbPrefixIdenStateList  = []byte("idenstates:")
	dbKeyConfig            = []byte("config")
	dbKeyKOp               = []byte("kop")
	dbKeyId                = []byte("id")
	dbKeyNonceIdx          = []byte("nonceidx")
	dbKeyIdenStateOnChain  = []byte("idenstateonchain")
	dbKeyIdenStatePending  = []byte("idenstatepending")
	dbKeySetStateEthTx     = []byte("setstateethtx")
)

// Config allows configuring the creation of an Issuer.
type Config struct {
	MaxLevelsClaimsTree     int
	MaxLevelsRevocationTree int
	MaxLevelsRootsTree      int
}

// ConfigDefault is a default configuration for the Issuer.
var ConfigDefault = Config{MaxLevelsClaimsTree: 140, MaxLevelsRevocationTree: 140, MaxLevelsRootsTree: 140}

// IdenStateTreeRoots is the set of the three roots of each Identity Merkle Tree.
type IdenStateTreeRoots struct {
	ClaimsRoot      *merkletree.Hash
	RevocationsRoot *merkletree.Hash
	RootsRoot       *merkletree.Hash
}

// TODO: Add mutex!

// Issuer is an identity that issues claims
type Issuer struct {
	rw            *sync.RWMutex
	storage       db.Storage
	id            *core.ID
	claimsMt      *merkletree.MerkleTree
	revocationsMt *merkletree.MerkleTree
	rootsMt       *merkletree.MerkleTree
	// idenPubOnChain can be nil if the identity doesn't connect to the blockchain.
	idenPubOnChain idenpubonchain.IdenPubOnChainer
	keyStore       *keystore.KeyStore
	kOpComp        *babyjub.PublicKeyComp
	nonceGen       *UniqueNonceGen
	idenStateList  *StorageList
	// idenStateOnChain is the last known identity state checked to be in
	// the Smart Contract.  idenStateOnChain can be nil if the identity
	// doesn't connect to the blockchain.
	_idenStateOnChain *merkletree.Hash
	// idenStatePending is a newly calculated identity state that is being
	// published in the Smart Contract but the transaction to publish it is
	// still pending.
	_idenStatePending *merkletree.Hash
	_setStateEthTx    *types.Transaction
	cfg               Config
}

func (is *Issuer) idenStateOnChain() *merkletree.Hash {
	return is._idenStateOnChain
}

func (is *Issuer) setIdenStateOnChain(tx db.Tx, v *merkletree.Hash) {
	if v == nil {
		tx.Put(dbKeyIdenStateOnChain, merkletree.HashZero[:])
	} else {
		tx.Put(dbKeyIdenStateOnChain, v[:])
	}
	is._idenStateOnChain = v
}

func (is *Issuer) loadIdenStateOnChain() error {
	b, err := is.storage.Get(dbKeyIdenStateOnChain)
	if err != nil {
		return err
	}
	var v merkletree.Hash
	copy(v[:], b)
	if v.Equals(&merkletree.HashZero) {
		is._idenStateOnChain = nil
	} else {
		is._idenStateOnChain = &v
	}
	return nil
}

func (is *Issuer) idenStatePending() *merkletree.Hash {
	return is._idenStatePending
}

func (is *Issuer) setIdenStatePending(tx db.Tx, v *merkletree.Hash) {
	if v == nil {
		tx.Put(dbKeyIdenStatePending, merkletree.HashZero[:])
	} else {
		tx.Put(dbKeyIdenStatePending, v[:])
	}
	is._idenStatePending = v
}

func (is *Issuer) loadIdenStatePending() error {
	b, err := is.storage.Get(dbKeyIdenStatePending)
	if err == db.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}
	var v merkletree.Hash
	copy(v[:], b)
	if v.Equals(&merkletree.HashZero) {
		is._idenStatePending = nil
	} else {
		is._idenStatePending = &v
	}
	return nil
}

func (is *Issuer) setStateEthTx() *types.Transaction {
	return is._setStateEthTx
}

func (is *Issuer) setSetStateEthTx(tx db.Tx, v *types.Transaction) error {
	vJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}
	tx.Put(dbKeySetStateEthTx, vJSON)
	return nil
}

func (is *Issuer) loadSetStateEthTx() error {
	vJSON, err := is.storage.Get(dbKeySetStateEthTx)
	if err == db.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}
	var v types.Transaction
	if err := json.Unmarshal(vJSON, &v); err != nil {
		return err
	}
	is._setStateEthTx = &v
	return nil
}

// loadMTs loads the three identity merkle trees from the storage using the configuration.
func loadMTs(cfg *Config, storage db.Storage) (*merkletree.MerkleTree, *merkletree.MerkleTree, *merkletree.MerkleTree, error) {
	cltStorage := storage.WithPrefix(dbPrefixClaimsTree)
	retStorage := storage.WithPrefix(dbPrefixRevocationTree)
	rotStorage := storage.WithPrefix(dbPrefixRootsTree)

	clt, err := merkletree.NewMerkleTree(cltStorage, cfg.MaxLevelsClaimsTree)
	if err != nil {
		return nil, nil, nil, err
	}
	ret, err := merkletree.NewMerkleTree(retStorage, cfg.MaxLevelsRevocationTree)
	if err != nil {
		return nil, nil, nil, err
	}
	rot, err := merkletree.NewMerkleTree(rotStorage, cfg.MaxLevelsRootsTree)
	if err != nil {
		return nil, nil, nil, err
	}

	return clt, ret, rot, nil
}

// New creates a new Issuer, creating a new genesis ID and initializes the storages.
func New(cfg Config, kOpComp *babyjub.PublicKeyComp, extraGenesisClaims []merkletree.Entrier, storage db.Storage, keyStore *keystore.KeyStore, idenPubOnChain idenpubonchain.IdenPubOnChainer) (*Issuer, error) {
	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, err
	}

	tx, err := storage.NewTx()

	if err != nil {
		return nil, err
	}

	// Initialize the UniqueNonceGen to generate revocation nonces for claims.
	nonceGen := NewUniqueNonceGen(NewStorageValue(dbKeyNonceIdx))
	nonceGen.Init(tx)

	// Create the Claim to authorize the Operational Key (kOp)
	kOp, err := kOpComp.Decompress()
	if err != nil {
		return nil, err
	}
	nonce, err := nonceGen.Next(tx)
	if err != nil {
		return nil, err
	}
	claimKOp := claims.NewClaimAuthorizeKSignBabyJub(kOp, nonce)
	id, _, err := genesis.CalculateIdGenesisMT(clt, rot, claimKOp, extraGenesisClaims)
	if err != nil {
		return nil, err
	}

	tx.Put(dbKeyId, id[:])
	tx.Put(dbKeyKOp, kOpComp[:])

	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	tx.Put(dbKeyConfig, cfgJSON)

	idenStateList := NewStorageList(dbPrefixIdenStateList)

	is := Issuer{
		rw:             &sync.RWMutex{},
		id:             id,
		claimsMt:       clt,
		revocationsMt:  ret,
		rootsMt:        rot,
		idenPubOnChain: idenPubOnChain,
		// idenStateWriter: idenStateWriter,
		keyStore:      keyStore,
		kOpComp:       kOpComp,
		storage:       storage,
		nonceGen:      nonceGen,
		idenStateList: idenStateList,
		cfg:           cfg,
	}

	// Initalize the history of idenStates
	idenState, idenStateTreeRoots := is.state()
	idenStateList.Init(tx)

	if err := idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return nil, err
	}

	// Initialize IdenStateOnChain and IdenStatePending to nil (writes to storage).
	is.setIdenStateOnChain(tx, nil)
	is.setIdenStatePending(tx, nil)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &is, nil

}

// Load creates an Issuer by loading a previously created Issuer (with New).
func Load(storage db.Storage, keyStore *keystore.KeyStore, idenPubOnChain idenpubonchain.IdenPubOnChainer) (*Issuer, error) {
	var cfg Config
	cfgJSON, err := storage.Get(dbKeyConfig)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return nil, err
	}

	kOpCompBytes, err := storage.Get(dbKeyConfig)
	if err != nil {
		return nil, err
	}
	var kOpComp babyjub.PublicKeyComp
	copy(kOpComp[:], kOpCompBytes)

	var id core.ID
	idBytes, err := storage.Get(dbKeyId)
	if err != nil {
		return nil, err
	}
	copy(id[:], idBytes)

	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, err
	}

	nonceGen := NewUniqueNonceGen(NewStorageValue(dbKeyNonceIdx))
	idenStateList := NewStorageList(dbPrefixIdenStateList)

	is := Issuer{
		rw:             &sync.RWMutex{},
		id:             &id,
		claimsMt:       clt,
		revocationsMt:  ret,
		rootsMt:        rot,
		idenPubOnChain: idenPubOnChain,
		keyStore:       keyStore,
		kOpComp:        &kOpComp,
		storage:        storage,
		nonceGen:       nonceGen,
		idenStateList:  idenStateList,
		cfg:            cfg,
	}

	if err := is.loadIdenStateOnChain(); err != nil {
		return nil, err
	}
	if err := is.loadIdenStatePending(); err != nil {
		return nil, err
	}

	if err := is.SyncIdenStatePublic(); err != nil {
		if err != ErrIdenPubOnChainNil {
			return nil, err
		}
	}
	return &is, nil
}

// state returns the current Identity State and the three merkle tree roots.
func (is *Issuer) state() (*merkletree.Hash, IdenStateTreeRoots) {
	is.rw.RLock()
	defer is.rw.RUnlock()
	clr, rer, ror := is.claimsMt.RootKey(), is.revocationsMt.RootKey(), is.rootsMt.RootKey()
	idenState := core.IdenState(clr, rer, ror)
	return idenState, IdenStateTreeRoots{
		ClaimsRoot:      clr,
		RevocationsRoot: rer,
		RootsRoot:       ror,
	}
}

// ID returns the Issuer ID (Identity ID).
func (is *Issuer) ID() *core.ID {
	return is.id
}

// SyncIdenStatePublic updates the IdenStateOnChain and IdenStatePending from
// the values in the Smart Contract.
func (is *Issuer) SyncIdenStatePublic() error {
	tx, err := is.storage.NewTx()
	if err != nil {
		return err
	}
	if err := is.syncIdenStatePublic(tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (is *Issuer) syncIdenStatePublic(tx db.Tx) error {
	is.rw.Lock()
	defer is.rw.Unlock()
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	idenStateData, err := is.idenPubOnChain.GetState(is.id)
	if err != nil {
		return err
	}
	if is.idenStateOnChain() == nil {
		// If the IdenState is not in the blockchain, the result will be HashZero
		if idenStateData.IdenState.Equals(&merkletree.HashZero) {
			return nil
		}
		is.setIdenStateOnChain(tx, idenStateData.IdenState)
		return nil
	}
	if is.idenStatePending() == nil {
		// If there's no IdenState pending to be set on chain, the
		// obtained one must be the idenStateOnChain
		if !idenStateData.IdenState.Equals(is.idenStateOnChain()) {
			return fmt.Errorf("Fatal error: Identity State in the Smart Contract (%v)"+
				" doesn't match the expected OnChain one (%v).",
				idenStateData.IdenState, is.idenStateOnChain())
		}
		return nil
	}
	// If there's an IdenState pending to be set on chain, the
	// obtained one can be:

	// a. the idenStateOnchan (in this case, we still have an
	// IdenState pending to be set on chain).
	if idenStateData.IdenState.Equals(is.idenStateOnChain()) {
		return nil
	}

	// b. the idenStatePending (in this case, we no longer have an
	// IdenState pending and it becomes the idenStateOnChain).
	if idenStateData.IdenState.Equals(is.idenStatePending()) {
		is.setIdenStatePending(tx, nil)
		is.setIdenStateOnChain(tx, idenStateData.IdenState)
		return nil
	}

	// c. Neither the idenStatePending nor the idenStateOnchain
	// (unexpected result).
	return fmt.Errorf("Fatal error: Identity State in the Smart Contract (%v)"+
		" doesn't match the Pending one (%v) nor the OnChain one (%v).",
		idenStateData.IdenState, is.idenStatePending(), is.idenStateOnChain())
}

// GenCredentialExistence generates an existence credential (claim + proof of
// existence) of an issued claim.  The result contains all data necessary to
// validate the credential against the Identity State found in the blockchain.
func (is *Issuer) GenCredentialExistence(claim merkletree.Entrier) (*proof.CredentialExistence, error) {
	// idenState, err := m.idenStateWriter.GetIdenState(is.id)
	// if err != nil {
	// 	return nil, err
	// }
	// var clr *merkletree.Hash // TODO
	// clt, err := m.claimsMt.Snapshot(clr)
	// if err != nil {
	// 	return nil, err
	// }
	// proofClaim, err := proof.GetClaimProofByHi(mt, hi)
	// if err != nil {
	// 	return nil, err
	// }
	// proofClaim.ID = is.id
	// proofClaim.BlockN = idenState.BlockN
	// proofClaim.BlockTimestamp = idenState.BlockTs
	// return proofClaim, nil
	return nil, fmt.Errorf("TODO")
}

// IssueClaim adds a new claim to the Claims Merkle Tree of the Issuer.  The
// Identity State is not updated.
func (is *Issuer) IssueClaim(claim merkletree.Entrier) error {
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	err := is.claimsMt.AddClaim(claim)
	if err != nil {
		return err
	}
	return nil
}

// getIdenStateByIdx gets identity state and identity state tree roots of the
// Issuer from the stored list at index idx.
func (is *Issuer) getIdenStateByIdx(tx db.Tx, idx uint32) (*merkletree.Hash, *IdenStateTreeRoots, error) {
	var idenStateTreeRoots IdenStateTreeRoots
	idenStateBytes, err := is.idenStateList.GetByIdx(tx, idx, &idenStateTreeRoots)
	if err != nil {
		return nil, nil, err
	}
	var idenState merkletree.Hash
	copy(idenState[:], idenStateBytes)
	return &idenState, &idenStateTreeRoots, nil
}

// getIdenStateTreeRoots gets the identity state tree roots of the Issuer from
// the stored list by identity state.
// func (is *Issuer) getIdenStateTreeRoots(tx db.Tx, idenState *merkletree.Hash) (*IdenStateTreeRoots, error) {
// 	var idenStateTreeRoots IdenStateTreeRoots
// 	if err := is.idenStateList.Get(tx, idenState[:], &idenStateTreeRoots); err != nil {
// 		return nil, err
// 	}
// 	return &idenStateTreeRoots, nil
// }

// PublishState calculates the current Issuer identity state, and if it's
// different than the last one, it publishes in in the blockchain.
func (is *Issuer) PublishState() error {
	is.rw.Lock()
	defer is.rw.Unlock()
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	if is.idenStatePending() != nil {
		return ErrIdenStatePendingNotNil
	}
	idenState, idenStateTreeRoots := is.state()

	tx, err := is.storage.NewTx()
	if err != nil {
		return err
	}
	idenStateListLen, err := is.idenStateList.Length(tx)
	if err != nil {
		return err
	}
	idenStateLast, _, err := is.getIdenStateByIdx(tx, idenStateListLen-1)
	if err != nil {
		return err
	}

	if idenState == idenStateLast {
		// IdenState hasn't changed, there's no need to do anything!
		return nil
	}

	if err := is.idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return err
	}

	// Publish the State in the Smart Contract.
	ethTx, err := is.idenPubOnChain.SetState(is.id, idenState, nil, nil, nil)
	if err != nil {
		return err
	}

	if err := is.setSetStateEthTx(tx, ethTx); err != nil {
		return err
	}

	is.setIdenStatePending(tx, idenState)

	if err := tx.Commit(); err != nil {
		return err
	}
	return fmt.Errorf("TODO")
}

// RevokeClaim revokes an already issued claim.
func (is *Issuer) RevokeClaim(claim merkletree.Entrier) error {
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	data, err := is.claimsMt.GetDataByIndex(claim.Entry().HIndex())
	if err != nil {
		return err
	}
	nonce := claims.GetRevocationNonce(&merkletree.Entry{Data: *data})

	if err := claims.AddLeafRevocationsTree(is.revocationsMt, nonce, 0xffffffff); err != nil {
		return err
	}
	return nil
}

// UpdateClaim allows updating the value of an already issued claim.
func (is *Issuer) UpdateClaim(hIndex *merkletree.Hash, value []merkletree.ElemBytes) error {
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	return fmt.Errorf("TODO")
}

// Sign signs a message by the kOp of the issuer.
func (is *Issuer) Sign(string) (string, error) {
	return "", fmt.Errorf("TODO")
}

// Sign signs a binary message by the kOp of the issuer.
func (is *Issuer) SignBinary(string) (string, error) {
	return "", fmt.Errorf("TODO")
}
