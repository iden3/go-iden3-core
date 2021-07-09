package issuer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-core/components/idenpuboffchain"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/eth"
	"github.com/iden3/go-iden3-core/keystore"
	zkutils "github.com/iden3/go-iden3-core/utils/zk"
	"github.com/iden3/go-merkletree-sql"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"

	"github.com/iden3/go-circom-prover-verifier/prover"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	witnesscalc "github.com/iden3/go-circom-witnesscalc"

	log "github.com/sirupsen/logrus"
)

var (
	ErrIdenGenesisOnly                    = fmt.Errorf("identity is genesis only")
	ErrIdenPubOnChainNil                  = fmt.Errorf("idenPubOnChain is nil")
	ErrIdenStateSNARKPathsNil             = fmt.Errorf("idenStateZkProofConf is nil")
	ErrEthClientNil                       = fmt.Errorf("ethClient is nil")
	ErrIdenPubOffChainWriterNil           = fmt.Errorf("idenPubOffChainWriter is nil")
	ErrIdenStatePendingNotNil             = fmt.Errorf("update of the published IdenState is pending")
	ErrIdenStateOnChainZero               = fmt.Errorf("no IdenState known to be on chain")
	ErrClaimNotFoundStateOnChain          = fmt.Errorf("claim not found under the on chain identity state")
	ErrClaimNotFoundClaimsTree            = fmt.Errorf("claim not found in the claims tree: the claim hasn't been issued")
	ErrClaimNotYetInOnChainState          = fmt.Errorf("claim has been issued but is not yet under a published on chain identity state")
	ErrFailedVerifyZkProofIdenStateUpdate = fmt.Errorf("failed verifing generated zk proof of identity state update")
)

var (
	dbPrefixClaimsTree        = []byte("treeclaims:")
	dbPrefixRevocationTree    = []byte("treerevocation:")
	dbPrefixRootsTree         = []byte("treeroots:")
	dbPrefixIdenStateList     = []byte("idenstates:")
	dbKeyConfig               = []byte("config")
	dbKeyKOp                  = []byte("kop")
	dbKeyClaimKOpHi           = []byte("claimkophi")
	dbKeyGenesisClaimKOpMtp   = []byte("genclaimkopmtp")
	dbKeyGenesisClaimTreeRoot = []byte("genclr")
	dbKeyId                   = []byte("id")
	dbKeyNonceIdx             = []byte("nonceidx")
	// dbKeyIdenStateOnChain     = []byte("idenstateonchain")
	dbKeyIdenStateDataOnChain       = []byte("idenstatedataonchain")
	dbKeyIdenStatePending           = []byte("idenstatepending")
	dbKeyIdenStatePendingTransacted = []byte("idenstatependingtxed")
	dbKeyEthTxSetState              = []byte("ethtxsetstate")
	dbKeyEthTxInitState             = []byte("ethtxinitstate")
)

var (
	SigPrefixSetState = []byte("setstate:")
)

// ConfigDefault is a default configuration for the Issuer.
var ConfigDefault = Config{MaxLevelsClaimsTree: 140, MaxLevelsRevocationTree: 140, MaxLevelsRootsTree: 140, GenesisOnly: false, ConfirmBlocks: 3}

// Config allows configuring the creation of an Issuer.
type Config struct {
	MaxLevelsClaimsTree     int
	MaxLevelsRevocationTree int
	MaxLevelsRootsTree      int
	GenesisOnly             bool
	ConfirmBlocks           uint64
}

// IdenStateZkProofConf are the paths to the SNARK related files required to
// generate an identity state update zkSNARK proof.
type IdenStateZkProofConf struct {
	Levels int
	Files  zkutils.ZkFiles
}

// IdenStateTreeRoots is the set of the three roots of each Identity Merkle Tree.
type IdenStateTreeRoots struct {
	ClaimsTreeRoot      *merkletree.Hash
	RevocationsTreeRoot *merkletree.Hash
	RootsTreeRoot       *merkletree.Hash
}

// Issuer is an identity that issues claims
type Issuer struct {
	rw              *sync.RWMutex
	storage         db.Storage
	id              *core.ID
	claimsTree      *merkletree.MerkleTree
	revocationsTree *merkletree.MerkleTree
	rootsTree       *merkletree.MerkleTree
	// idenPubOnChain can be nil if the identity doesn't connect to the blockchain.
	idenPubOnChain idenpubonchain.IdenPubOnChainer
	// idenPubOffChainWriter can be nil if the identity doesn't ever update
	// it's state after genesis.
	idenPubOffChainWriter idenpuboffchain.IdenPubOffChainWriter
	keyStore              *keystore.KeyStore
	kOpComp               *babyjub.PublicKeyComp
	nonceGen              *UniqueNonceGen
	idenStateList         *db.StorageList
	// _idenStateOnChain     *merkletree.Hash
	// idenStateDataOnChain is the last known identity state checked to be
	// in the Smart Contract.
	_idenStateDataOnChain *proof.IdenStateData
	// idenStatePending is a newly calculated identity state that is being
	// published in the Smart Contract but the transaction to publish it is
	// still pending.
	_idenStatePending           *merkletree.Hash
	_idenStatePendingTransacted bool
	_ethTxSetState              *types.Transaction
	_ethTxInitState             *types.Transaction
	idenStateZkProofConf        *IdenStateZkProofConf
	cfg                         Config
}

//
// Persistence setters and getters
//

func (is *Issuer) idenStateDataOnChain() *proof.IdenStateData { return is._idenStateDataOnChain }

func (is *Issuer) setIdenStateDataOnChain(tx db.Tx, v *proof.IdenStateData) error {
	is._idenStateDataOnChain = v
	return db.StoreJSON(tx, dbKeyIdenStateDataOnChain, v)
}

func (is *Issuer) loadIdenStateDataOnChain() error {
	is._idenStateDataOnChain = &proof.IdenStateData{}
	return db.LoadJSON(is.storage, dbKeyIdenStateDataOnChain, is._idenStateDataOnChain)
}

func (is *Issuer) IdenStateOnChain() *merkletree.Hash {
	return is._idenStateDataOnChain.IdenState
}

func (is *Issuer) idenStateOnChain() *merkletree.Hash {
	return is._idenStateDataOnChain.IdenState
}

func (is *Issuer) IdenStatePending() (*merkletree.Hash, bool) {
	return is._idenStatePending, is._idenStatePendingTransacted
}

func (is *Issuer) idenStatePending() (*merkletree.Hash, bool) {
	return is._idenStatePending, is._idenStatePendingTransacted
}

func byte2bool(v byte) bool {
	if v == 0 {
		return false
	} else {
		return true
	}
}

func bool2byte(v bool) byte {
	if v {
		return 1
	} else {
		return 0
	}
}

func (is *Issuer) setIdenStatePending(tx db.Tx, v *merkletree.Hash, transacted bool) {
	is._idenStatePending = v
	is._idenStatePendingTransacted = transacted
	tx.Put(dbKeyIdenStatePending, v[:])
	tx.Put(dbKeyIdenStatePendingTransacted, []byte{bool2byte(transacted)})
}

func (is *Issuer) loadIdenStatePending() error {
	b, err := is.storage.Get(dbKeyIdenStatePending)
	if err != nil {
		return err
	}
	var v merkletree.Hash
	copy(v[:], b)
	is._idenStatePending = &v

	transactedBytes, err := is.storage.Get(dbKeyIdenStatePendingTransacted)
	if err != nil {
		return err
	}
	is._idenStatePendingTransacted = byte2bool(transactedBytes[0])
	return nil
}

func (is *Issuer) ethTxSetState() *types.Transaction { return is._ethTxSetState }

func (is *Issuer) setEthTxSetState(tx db.Tx, v *types.Transaction) error {
	is._ethTxSetState = v
	return db.StoreJSON(tx, dbKeyEthTxSetState, v)
}

func (is *Issuer) loadEthTxSetState() error {
	is._ethTxSetState = &types.Transaction{}
	return db.LoadJSON(is.storage, dbKeyEthTxSetState, &is._ethTxSetState)
}

func (is *Issuer) ethTxInitState() *types.Transaction { return is._ethTxInitState }

func (is *Issuer) setEthTxInitState(tx db.Tx, v *types.Transaction) error {
	is._ethTxInitState = v
	return db.StoreJSON(tx, dbKeyEthTxInitState, v)
}

func (is *Issuer) loadEthTxInitState() error {
	is._ethTxInitState = &types.Transaction{}
	return db.LoadJSON(is.storage, dbKeyEthTxInitState, &is._ethTxInitState)
}

// loadMTs loads the three identity merkle trees from the storage using the configuration.
func loadMTs(cfg *Config, storage db.Storage) (*merkletree.MerkleTree, *merkletree.MerkleTree,
	*merkletree.MerkleTree, error) {
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

// Create a new Issuer, creating a new genesis ID and initializes the
// storages.  The extraGenesisClaims metadata's are updated.
func Create(cfg Config, kOpComp *babyjub.PublicKeyComp, extraGenesisClaims []claims.Claimer,
	storage db.Storage, keyStore *keystore.KeyStore) (*core.ID, error) {
	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, err
	}

	tx, err := storage.NewTx()

	if err != nil {
		return nil, err
	}

	// Initialize the UniqueNonceGen to generate revocation nonces for claims.
	nonceGen := NewUniqueNonceGen(db.NewStorageValue(dbKeyNonceIdx))
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
	claimKOp := claims.NewClaimKeyBabyJub(kOp, claims.BabyJubKeyTypeAuthorizeKSign)
	claimKOp.Metadata().RevNonce = nonce
	extraGenesisClaimsEntriers := make([]merkletree.Entrier, len(extraGenesisClaims))
	for i, claim := range extraGenesisClaims {
		nonce, err := nonceGen.Next(tx)
		if err != nil {
			return nil, err
		}
		claim.Metadata().RevNonce = nonce
		extraGenesisClaimsEntriers[i] = claim
	}
	id, err := genesis.CalculateIdGenesisMT(clt, rot, claimKOp, extraGenesisClaimsEntriers)
	if err != nil {
		return nil, err
	}
	claimKOpHi, err := claimKOp.Entry().HIndex()
	if err != nil {
		return nil, err
	}
	claimKOpMtp, _, err := clt.GenerateProof(claimKOpHi.BigInt(), nil)
	if err != nil {
		return nil, err
	}

	tx.Put(dbKeyId, id[:])
	tx.Put(dbKeyKOp, kOpComp[:])
	tx.Put(dbKeyClaimKOpHi, claimKOpHi[:])
	if err := db.StoreJSON(tx, dbKeyGenesisClaimKOpMtp, claimKOpMtp); err != nil {
		return nil, err
	}
	if err := db.StoreJSON(tx, dbKeyGenesisClaimTreeRoot, clt.Root()); err != nil {
		return nil, err
	}

	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	tx.Put(dbKeyConfig, cfgJSON)

	idenStateList := db.NewStorageList(dbPrefixIdenStateList)

	is := Issuer{
		rw:                    &sync.RWMutex{},
		id:                    id,
		claimsTree:            clt,
		revocationsTree:       ret,
		rootsTree:             rot,
		idenPubOnChain:        nil,
		idenPubOffChainWriter: nil,
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

	// Initialize IdenStateDataOnChain and IdenStatePending to zero (writes to storage).
	if err := is.setIdenStateDataOnChain(tx, &proof.IdenStateData{IdenState: &merkletree.HashZero}); err != nil {
		return nil, err
	}
	is.setIdenStatePending(tx, &merkletree.HashZero, false)
	if err := is.setEthTxInitState(tx, nil); err != nil {
		return nil, err
	}
	if err := is.setEthTxSetState(tx, nil); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return is.id, nil
}

// Load creates an Issuer by loading a previously created Issuer (with New).
func Load(storage db.Storage, keyStore *keystore.KeyStore,
	idenPubOnChain idenpubonchain.IdenPubOnChainer,
	idenStateZkProofConf *IdenStateZkProofConf,
	idenPubOffChainWriter idenpuboffchain.IdenPubOffChainWriter) (*Issuer, error) {
	var cfg Config
	cfgJSON, err := storage.Get(dbKeyConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting config from storage: %w", err)
	}
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return nil, err
	}
	if !cfg.GenesisOnly {
		if idenPubOnChain == nil {
			return nil, ErrIdenPubOnChainNil
		}
		if idenStateZkProofConf == nil {
			return nil, ErrIdenStateSNARKPathsNil
		}
		if err := idenStateZkProofConf.Files.LoadAll(); err != nil {
			return nil, fmt.Errorf("error loading zk files: %w", err)
		}
		if idenPubOffChainWriter == nil {
			return nil, ErrIdenPubOffChainWriterNil
		}
	}

	kOpCompBytes, err := storage.Get(dbKeyKOp)
	if err != nil {
		return nil, fmt.Errorf("error getting kop from storage: %w", err)
	}
	var kOpComp babyjub.PublicKeyComp
	copy(kOpComp[:], kOpCompBytes)

	var id core.ID
	idBytes, err := storage.Get(dbKeyId)
	if err != nil {
		return nil, fmt.Errorf("error getting id from storage: %w", err)
	}
	copy(id[:], idBytes)

	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, fmt.Errorf("error loading merkle trees from storage: %w", err)
	}

	nonceGen := NewUniqueNonceGen(db.NewStorageValue(dbKeyNonceIdx))
	idenStateList := db.NewStorageList(dbPrefixIdenStateList)

	is := Issuer{
		rw:                    &sync.RWMutex{},
		id:                    &id,
		claimsTree:            clt,
		revocationsTree:       ret,
		rootsTree:             rot,
		idenPubOnChain:        idenPubOnChain,
		idenPubOffChainWriter: idenPubOffChainWriter,
		keyStore:              keyStore,
		kOpComp:               &kOpComp,
		storage:               storage,
		nonceGen:              nonceGen,
		idenStateList:         idenStateList,
		idenStateZkProofConf:  idenStateZkProofConf,
		cfg:                   cfg,
	}

	if err := is.loadIdenStateDataOnChain(); err != nil {
		return nil, err
	}
	if err := is.loadIdenStatePending(); err != nil {
		return nil, err
	}
	if err := is.loadEthTxInitState(); err != nil {
		return nil, err
	}
	if err := is.loadEthTxSetState(); err != nil {
		return nil, err
	}

	if !is.cfg.GenesisOnly {
		if err := is.SyncIdenStatePublic(); err != nil {
			return nil, fmt.Errorf("error syncing idenstate from smart contract: %w", err)
		}
	}
	return &is, nil
}

// state returns the current Identity State and the three merkle tree roots.
func (is *Issuer) state() (*merkletree.Hash, IdenStateTreeRoots) {
	clr, rer, ror := is.claimsTree.Root(), is.revocationsTree.Root(), is.rootsTree.Root()
	idenState := core.IdenState(clr, rer, ror)
	return idenState, IdenStateTreeRoots{
		ClaimsTreeRoot:      clr,
		RevocationsTreeRoot: rer,
		RootsTreeRoot:       ror,
	}
}

// State calculates and returns the current Identity State and the three merkle tree roots.
func (is *Issuer) State() (*merkletree.Hash, IdenStateTreeRoots) {
	is.rw.RLock()
	defer is.rw.RUnlock()
	return is.state()
}

// StateDataOnChain returns the last known IdentityState Data known to be on chain.
func (is *Issuer) StateDataOnChain() *proof.IdenStateData {
	is.rw.RLock()
	defer is.rw.RUnlock()
	return is.idenStateDataOnChain()
}

// ID returns the Issuer ID (Identity ID).
func (is *Issuer) ID() *core.ID {
	return is.id
}

// KeyOperational returns the identity's operational key.
func (is *Issuer) KeyOperational() *babyjub.PublicKeyComp {
	return is.kOpComp
}

// SyncIdenStatePublic updates the IdenStateOnChain and IdenStatePending from
// the values in the Smart Contract.
func (is *Issuer) SyncIdenStatePublic() error {
	if is.cfg.GenesisOnly {
		return ErrIdenGenesisOnly
	}
	is.rw.Lock()
	defer is.rw.Unlock()
	// If there's a pending state, check that the ethereum Tx was
	// succsefully and only call GetState when the number of confirmed
	// blocks is equal or higher than is.cfg.ConfirmBlocks
	idenStatePending, transacted := is.idenStatePending()
	// (C)(idenStatePending: X, transacted: true)
	if !idenStatePending.Equals(&merkletree.HashZero) && transacted {
		var ethTx *types.Transaction
		// If idenStateOnChain is zero, the pending state was caused by
		// InitState.  Otherwise it was a regular SetState.
		if is.idenStateOnChain().Equals(&merkletree.HashZero) {
			ethTx = is.ethTxInitState()
		} else {
			ethTx = is.ethTxSetState()
		}
		confirmBlocks, err := is.idenPubOnChain.TxConfirmBlocks(ethTx)
		if err == eth.ErrReceiptNotReceived {
			return nil
		} else if err != nil {
			return fmt.Errorf("TxConfirmBlocks: %w", err)
		}
		log.WithField("tx", ethTx.Hash().Hex()).
			WithField("TxConfirmBlocks", confirmBlocks).
			WithField("cfg.ConfirmBlocks", is.cfg.ConfirmBlocks).
			Debug("State Update Tx")
		if confirmBlocks.Cmp(new(big.Int).SetUint64(is.cfg.ConfirmBlocks)) == -1 {
			return nil
		}
	}

	idenStateData, err := is.idenPubOnChain.GetState(is.id)
	if err == idenpubonchain.ErrIdenNotOnChain {
		idenStateData = &proof.IdenStateData{
			IdenState: &merkletree.HashZero,
		}
	} else if err != nil {
		return fmt.Errorf("error calling idenstates smart contract getState: %w", err)
	}

	// (A)(idenStatePending: 0, transacted: false)
	if idenStatePending.Equals(&merkletree.HashZero) && !transacted {
		// If there's no IdenState pending to be set on chain, the
		// obtained one must be the idenStateOnChain (Zero for genesis
		// / empty in the smart contract).
		if idenStateData.IdenState.Equals(is.idenStateOnChain()) {
			return nil
		}

		return fmt.Errorf("Fatal error: Identity State in the Smart Contract (%v)"+
			" doesn't match the expected OnChain one (%v).",
			idenStateData.IdenState, is.idenStateOnChain())
	}
	// If there's an IdenState pending to be set on chain, the
	// obtained one can be:

	// a. the idenStateOnchan (in this case, we still have an
	// IdenState pending to be set on chain).
	if idenStateData.IdenState.Equals(is.idenStateOnChain()) {
		return nil
	}

	// b. the idenStatePending (in this case, we no longer have an
	// IdenState pending and it becomes the idenStateOnChain, so we update
	// the sync state).
	// (C)(idenStatePending: X, transacted: true)
	if idenStateData.IdenState.Equals(idenStatePending) {
		tx, err := is.storage.NewTx()
		if err != nil {
			return err
		}
		is.setIdenStatePending(tx, &merkletree.HashZero, false)
		if err := is.setIdenStateDataOnChain(tx, idenStateData); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}

	// c. Neither the idenStatePending nor the idenStateOnchain
	// (unexpected result).
	return fmt.Errorf("Fatal error: Identity State in the Smart Contract (%v)"+
		" doesn't match the Pending one (%v) nor the OnChain one (%v).",
		idenStateData.IdenState, idenStatePending, is.idenStateOnChain())
}

// IssueClaim adds a new claim to the Claims Merkle Tree of the Issuer.  The
// Identity State is not updated.  The claim metadata is updated if the issue
// is successfull.
func (is *Issuer) IssueClaim(claim claims.Claimer) error {
	if is.cfg.GenesisOnly {
		return ErrIdenGenesisOnly
	}
	is.rw.Lock()
	defer is.rw.Unlock()
	tx, err := is.storage.NewTx()
	if err != nil {
		return err
	}
	nonce, err := is.nonceGen.Next(tx)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	claim.Metadata().RevNonce = nonce
	err = is.claimsTree.AddEntry(claim.Entry())
	if err != nil {
		return err
	}
	return nil
}

// getIdenStateByIdx gets identity state and identity state tree roots of the
// Issuer from the stored list at index idx.
func (is *Issuer) getIdenStateByIdx(tx db.Tx, idx int64) (*merkletree.Hash, *IdenStateTreeRoots, error) {
	idxAbs := uint32(idx)
	if idx < 0 {
		idenStateListLen, err := is.idenStateList.Length(tx)
		if err != nil {
			return nil, nil, err
		}
		if int64(idenStateListLen) < -idx {
			return nil, nil, fmt.Errorf("idenStateListLen (%v) < -Idx (%v)", idenStateListLen, -idx)
		}
		idxAbs = idenStateListLen - uint32(-idx)
	}
	var idenStateTreeRoots IdenStateTreeRoots
	idenStateBytes, err := is.idenStateList.GetByIdx(tx, idxAbs, &idenStateTreeRoots)
	if err != nil {
		return nil, nil, err
	}
	var idenState merkletree.Hash
	copy(idenState[:], idenStateBytes)
	return &idenState, &idenStateTreeRoots, nil
}

// getIdenStateTreeRoots gets the identity state tree roots of the Issuer from
// the stored list by identity state.
func (is *Issuer) getIdenStateTreeRoots(tx db.Tx, idenState *merkletree.Hash) (*IdenStateTreeRoots, error) {
	var idenStateTreeRoots IdenStateTreeRoots
	if err := is.idenStateList.Get(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return nil, err
	}
	return &idenStateTreeRoots, nil
}

// idenStatePending state graph:
// -> (A)(idenStatePending: 0, transacted: false) -> (B)(idenStatePending: X, transacted: false)
//                     ^\ (C)(idenStatePending: X, transacted: true) </

// PublishState calculates the current Issuer identity state, and if it's
// different than the last one, it publishes in in the blockchain.
func (is *Issuer) PublishState() error {
	if is.cfg.GenesisOnly {
		return ErrIdenGenesisOnly
	}
	is.rw.Lock()
	defer is.rw.Unlock()
	idenStatePending, transacted := is.idenStatePending()
	// (C)(idenStatePending: X, transacted: true)
	if !idenStatePending.Equals(&merkletree.HashZero) && transacted {
		return ErrIdenStatePendingNotNil
	}

	idenState, idenStateTreeRoots := is.state()

	tx0, err := is.storage.NewTx() // Read only Tx
	defer tx0.Close()
	if err != nil {
		return err
	}
	idenStateLast, idenStateTreeRootsLast, err := is.getIdenStateByIdx(tx0, -1)
	if err != nil {
		return err
	}

	// (A)(idenStatePending: 0, transacted: false) && idenState != idenStateLast
	if idenStatePending.Equals(&merkletree.HashZero) && !transacted {
		if idenState.Equals(idenStateLast) {
			// IdenState hasn't changed, there's no need to do
			// anything!
			return nil
		}

		// idenState != idenStateLast

		// If the ClaimsTreeRoot has changed (claims have been added), add the
		// ClaimsTreeRoot to the RootsTree.
		if !idenStateTreeRoots.ClaimsTreeRoot.Equals(idenStateTreeRootsLast.ClaimsTreeRoot) {
			if err := claims.AddLeafRootsTree(is.rootsTree, idenStateTreeRoots.ClaimsTreeRoot); err != nil {
				return err
			}
			idenState, idenStateTreeRoots = is.state()
		}

		tx, err := is.storage.NewTx()
		if err != nil {
			return err
		}

		if err := is.idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
			return err
		}

		is.setIdenStatePending(tx, idenState, false)

		if err := tx.Commit(); err != nil {
			return err
		}
	} else {
		idenStateLast, idenStateTreeRootsLast, err = is.getIdenStateByIdx(tx0, -2)
		if err != nil {
			return err
		}
	}

	// (B)(idenStatePending: X, transacted: false)

	zkProofOut, err := is.GenZkProofIdenStateUpdate(idenStateLast, idenState)
	if err != nil {
		return err
	}

	tx, err := is.storage.NewTx()
	if err != nil {
		return err
	}

	if is.idenStateOnChain().Equals(&merkletree.HashZero) {
		// Identity State not present in the Smart Contract. First time
		// publishing it.
		ethTx, err := is.idenPubOnChain.InitState(is.id, idenStateLast, idenState, &zkProofOut.Proof)
		if err != nil {
			return fmt.Errorf("error calling idenstates smart contract initState: %w", err)
		}

		if err := is.setEthTxInitState(tx, ethTx); err != nil {
			return err
		}
	} else {
		// Identity State already present in the Smart Contract.
		// Update it.
		ethTx, err := is.idenPubOnChain.SetState(is.id, idenState, &zkProofOut.Proof)
		if err != nil {
			return fmt.Errorf("error calling idenstates smart contract setState: %w", err)
		}

		if err := is.setEthTxSetState(tx, ethTx); err != nil {
			return err
		}
	}
	is.setIdenStatePending(tx, idenState, true)

	if err := tx.Commit(); err != nil {
		return err
	}

	publicData := idenpuboffchain.PublicData{
		IdenState:           idenState,
		ClaimsTreeRoot:      idenStateTreeRoots.ClaimsTreeRoot,
		RevocationsTreeRoot: idenStateTreeRoots.RevocationsTreeRoot,
		RevocationsTree:     is.revocationsTree,
		RootsTreeRoot:       idenStateTreeRoots.RootsTreeRoot,
		RootsTree:           is.rootsTree,
	}

	// finally, Publish the Public Off Chain identity data
	if err := is.idenPubOffChainWriter.Publish(is.id, &publicData); err != nil {
		return err
	}
	return nil
}

// RevokeClaim revokes an already issued claim.
func (is *Issuer) RevokeClaim(claim merkletree.Entrier) error {
	if is.cfg.GenesisOnly {
		return ErrIdenGenesisOnly
	}
	is.rw.Lock()
	defer is.rw.Unlock()

	hi, err := claim.Entry().HIndex()
	if err != nil {
		return err
	}
	data, err := is.claimsTree.GetDataByIndex(hi)
	if err != nil {
		return err
	}
	nonce := claims.GetRevocationNonce(&merkletree.Entry{Data: *data})

	if err := claims.AddLeafRevocationsTree(is.revocationsTree, nonce, 0xffffffff); err != nil {
		return err
	}
	return nil
}

// UpdateClaim allows updating the value of an already issued claim.
func (is *Issuer) UpdateClaim(hIndex *merkletree.Hash, value []merkletree.ElemBytes) error {
	if is.cfg.GenesisOnly {
		return ErrIdenGenesisOnly
	}
	return fmt.Errorf("TODO")
}

// Sign signs a message by the kOp of the issuer.
func (is *Issuer) Sign(string) (string, error) {
	return "", fmt.Errorf("TODO")
}

// SignBinary signs a binary message by the kOp of the issuer.
func (is *Issuer) SignBinary(prefix, msg []byte) (*babyjub.SignatureComp, error) {
	return is.keyStore.SignRaw(is.kOpComp, append(prefix, msg...))
}

// SignState signs the Identity State transition (oldState+newState) by the kOp of the issuer.
func (is *Issuer) SignState(oldState, newState *merkletree.Hash) (*babyjub.SignatureComp, error) {
	var prefix31 [31]byte
	copy(prefix31[:], SigPrefixSetState)
	prefixBigInt := new(big.Int)
	utils.SetBigIntFromLEBytes(prefixBigInt, prefix31[:])

	toHash := []*big.Int{prefixBigInt, oldState.BigInt(), newState.BigInt(), big.NewInt(0), big.NewInt(0), big.NewInt(0)}

	return is.SignElems(toHash)
}

// SignElems signs a [poseidon.T]*big.Int of elements in *big.Int format
func (is *Issuer) SignElems(toHash []*big.Int) (*babyjub.SignatureComp, error) {
	e, err := poseidon.Hash(toHash)
	if err != nil {
		return nil, err
	}
	return is.keyStore.SignElem(is.kOpComp, e)
}

func generateExistenceMTProof(mt *merkletree.MerkleTree, hi, root *merkletree.Hash) (*merkletree.Proof, error) {
	mtp, _, err := mt.GenerateProof(hi.BigInt(), root)
	if err != nil {
		return nil, err
	}
	if !mtp.Existence {
		return nil, ErrClaimNotFoundStateOnChain
	}
	return mtp, nil
}

// GenCredentialExistence generates an existence credential (claim + proof of
// existence) of an issued claim.  The result contains all data necessary to
// validate the credential against the Identity State found in the blockchain.
// For now, there are no genesis credentials.
func (is *Issuer) GenCredentialExistence(claim merkletree.Entrier) (*proof.CredentialExistence, error) {
	// TODO: Once a genesis credential is implemented, figure out what to
	// return here.  Maybe this function will error and there will be a
	// "GenCredentialExistenceGenesis".  Maybe this function will be able
	// to return a credential even when there's no state on chain.
	if is.cfg.GenesisOnly {
		return nil, ErrIdenGenesisOnly
	}
	tx, err := is.storage.NewTx()
	if err != nil {
		return nil, err
	}
	is.rw.RLock()
	defer is.rw.RUnlock()
	idenStateData := is.idenStateDataOnChain()
	if idenStateData.IdenState.Equals(&merkletree.HashZero) {
		return nil, ErrIdenStateOnChainZero
	}
	idenStateTreeRoots, err := is.getIdenStateTreeRoots(tx, idenStateData.IdenState)
	if err != nil {
		return nil, err
	}
	claimEntry := claim.Entry()
	hi, err := claimEntry.HIndex()
	if err != nil {
		return nil, err
	}
	mtpExist, err := generateExistenceMTProof(is.claimsTree, hi,
		idenStateTreeRoots.ClaimsTreeRoot)
	if err != nil {
		// We were unable to generate a proof from the claims tree
		// associated with the on chain identity state.  Check if the
		// claim exists in the current claims tree.
		if err := is.claimsTree.EntryExists(claimEntry, nil); err != nil {
			return nil, ErrClaimNotFoundClaimsTree
		} else {
			return nil, ErrClaimNotYetInOnChainState
		}
	} else {
		// We were able to generate a proof from the claims tree with
		// the HIndex.  Check the HValue is also valid!
		if err := is.claimsTree.EntryExists(claimEntry, idenStateTreeRoots.ClaimsTreeRoot); err != nil {
			return nil, ErrClaimNotFoundClaimsTree
		}
	}
	return &proof.CredentialExistence{
		Id:                  is.id,
		IdenStateData:       *idenStateData,
		MtpClaim:            mtpExist,
		Claim:               claimEntry,
		RevocationsTreeRoot: idenStateTreeRoots.RevocationsTreeRoot,
		RootsTreeRoot:       idenStateTreeRoots.RootsTreeRoot,
		IdenPubUrl:          is.idenPubOffChainWriter.Url(),
	}, nil
}

type IdOwnershipGenesisInputs struct {
	Id             *big.Int
	PrivateKey     *big.Int
	MtpSiblings    []*big.Int
	ClaimsTreeRoot *big.Int
	// RevTreeRoot    *big.Int
	// RootTreeRoot   *big.Int
}

func (is *Issuer) GenIdOwnershipGenesisInputs(levels int) (*IdOwnershipGenesisInputs, error) {
	sk, err := is.keyStore.ExportKey(is.kOpComp)
	if err != nil {
		return nil, err
	}

	var mtp merkletree.Proof
	err = db.LoadJSON(is.storage, dbKeyGenesisClaimKOpMtp, &mtp)
	if err != nil {
		return nil, err
	}
	siblings := mtp.AllSiblingsCircom(levels)
	if len(siblings) != levels+1 {
		return nil, fmt.Errorf("number of mtp siblings in genesis ClaimTree (%v) "+
			"is higher than requested levels (%v)", len(siblings), levels+1)
	}

	var genesisClaimTreeRoot merkletree.Hash
	err = db.LoadJSON(is.storage, dbKeyGenesisClaimTreeRoot, &genesisClaimTreeRoot)
	if err != nil {
		return nil, err
	}
	return &IdOwnershipGenesisInputs{
		Id:             is.id.BigInt(),
		PrivateKey:     (*big.Int)(sk.Scalar()),
		MtpSiblings:    siblings,
		ClaimsTreeRoot: genesisClaimTreeRoot.BigInt(),
		// RevTreeRoot    :
		// RootTreeRoot   :
	}, nil
}

func (is *Issuer) GenZkProofIdenStateUpdate(oldIdState, newIdState *merkletree.Hash) (*zkutils.ZkProofOut, error) {
	pk, err := is.idenStateZkProofConf.Files.ProvingKey()
	if err != nil {
		return nil, fmt.Errorf("error loading zk pk: %w", err)
	}
	vk, err := is.idenStateZkProofConf.Files.VerificationKey()
	if err != nil {
		return nil, fmt.Errorf("error loading zk vk: %w", err)
	}

	idOwnershipInputs, err := is.GenIdOwnershipGenesisInputs(is.idenStateZkProofConf.Levels)
	if err != nil {
		return nil, fmt.Errorf("error generating idOwnership inputs: %w", err)
	}

	inputs := make(map[string]interface{})

	inputs["id"] = idOwnershipInputs.Id
	inputs["oldIdState"] = oldIdState.BigInt()
	inputs["userPrivateKey"] = idOwnershipInputs.PrivateKey
	inputs["siblings"] = idOwnershipInputs.MtpSiblings
	inputs["claimsTreeRoot"] = idOwnershipInputs.ClaimsTreeRoot
	inputs["newIdState"] = newIdState.BigInt()

	witnessCalcWASM, err := is.idenStateZkProofConf.Files.WitnessCalcWASM()
	if err != nil {
		return nil, fmt.Errorf("error loading zk witnessCalc WASM: %w", err)
	}
	wit, err := witnesscalc.CalculateWitnessBinWASM(witnessCalcWASM, inputs)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	proof, pubSignals, err := prover.GenerateProof(pk, wit)
	if err != nil {
		return nil, err
	}
	// Verify zk proof
	if !verifier.Verify(vk, proof, pubSignals) {
		return nil, ErrFailedVerifyZkProofIdenStateUpdate
	}

	log.WithField("elapsed", time.Since(start)).Debug("Proof generated")
	return &zkutils.ZkProofOut{Proof: *proof, PubSignals: pubSignals}, nil
}

// TODO: Create an Admin struct that exposes the following:
// - The 3 Merle Trees
// - RawDump(f func(key, value string))
// - RawImport(raw map[string]string) (int, error)
// - ClaimsDump() map[string]string
// The return and input types are open to change.  They are based on the old
// components/idenadminutils/idenadminutils.go
