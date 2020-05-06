package local

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	zkparsers "github.com/iden3/go-circom-prover-verifier/parsers"
	zktypes "github.com/iden3/go-circom-prover-verifier/types"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree"
	// cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
)

type IdenStateHistory struct {
	IdenStates []*proof.IdenStateData
	ByTime     map[int64]*proof.IdenStateData
	ByBlock    map[uint64]*proof.IdenStateData
}

func NewIdenStateHistory() *IdenStateHistory {
	return &IdenStateHistory{
		IdenStates: make([]*proof.IdenStateData, 0),
		ByTime:     make(map[int64]*proof.IdenStateData),
		ByBlock:    make(map[uint64]*proof.IdenStateData),
	}
}

func (h *IdenStateHistory) Add(idenStateData *proof.IdenStateData) {
	h.IdenStates = append(h.IdenStates, idenStateData)
	h.ByTime[idenStateData.BlockTs] = idenStateData
	h.ByBlock[idenStateData.BlockN] = idenStateData
}

type IdIdenStateData struct {
	Id            *core.ID
	IdenStateData *proof.IdenStateData
}

// IdenPubOnChain is an implementation of the IdenPubOnnChainer that instead of
// interacting with the blockchain has a local copy of the identities states.
// All writes (Init and Set) are set to a pending queue, and written into the
// internal state once the Sync function is called.
type IdenPubOnChain struct {
	rw             sync.RWMutex
	idenStatesData map[core.ID]*IdenStateHistory
	pendingInit    []*IdIdenStateData
	pendingSet     []*IdIdenStateData
	timeNow        func() time.Time
	blockNow       func() uint64
	verifyingKey   *zktypes.Vk
}

// New creates a new IdenPubOnChain
func New(timeNow func() time.Time, blockNow func() uint64, verifyingKey *zktypes.Vk) *IdenPubOnChain {
	return &IdenPubOnChain{
		idenStatesData: make(map[core.ID]*IdenStateHistory),
		pendingInit:    make([]*IdIdenStateData, 0),
		pendingSet:     make([]*IdIdenStateData, 0),
		timeNow:        timeNow,
		blockNow:       blockNow,
		verifyingKey:   verifyingKey,
	}
}

// Sync all the pending writes to the local state.
func (ip *IdenPubOnChain) Sync() {
	ip.rw.Lock()
	defer ip.rw.Unlock()
	for _, idIdenStateData := range ip.pendingInit {
		idenStatesData := NewIdenStateHistory()
		idenStatesData.Add(idIdenStateData.IdenStateData)
		ip.idenStatesData[*idIdenStateData.Id] = idenStatesData
	}
	for _, idIdenStateData := range ip.pendingSet {
		ip.idenStatesData[*idIdenStateData.Id].Add(idIdenStateData.IdenStateData)
	}
	ip.pendingInit = make([]*IdIdenStateData, 0)
	ip.pendingSet = make([]*IdIdenStateData, 0)
}

// GetState returns the Identity State Data of the given ID from the IdenStates Smart Contract.
func (ip *IdenPubOnChain) GetState(id *core.ID) (*proof.IdenStateData, error) {
	ip.rw.RLock()
	defer ip.rw.RUnlock()
	idenStatesData, ok := ip.idenStatesData[*id]
	if !ok {
		return nil, idenpubonchain.ErrIdenNotOnChain
	}
	return idenStatesData.IdenStates[len(idenStatesData.IdenStates)-1], nil
}

// GetStateByBlock returns the Identity State Data of the given ID published at
// queryBlockN from the IdenStates Smart Contract.
func (ip *IdenPubOnChain) GetStateByBlock(id *core.ID, queryBlockN uint64) (*proof.IdenStateData, error) {
	ip.rw.RLock()
	defer ip.rw.RUnlock()
	idenStatesData, ok := ip.idenStatesData[*id]
	if !ok {
		return nil, idenpubonchain.ErrIdenNotOnChain
	}
	idenState, ok := idenStatesData.ByBlock[queryBlockN]
	if !ok {
		return nil, idenpubonchain.ErrIdenByBlockNotFound
	}
	return idenState, nil
}

// GetStateByTime returns the Identity State Data of the given ID published at
// queryBlockTs from the IdenStates Smart Contract.
func (ip *IdenPubOnChain) GetStateByTime(id *core.ID, queryBlockTs int64) (*proof.IdenStateData, error) {
	ip.rw.RLock()
	defer ip.rw.RUnlock()
	idenStatesData, ok := ip.idenStatesData[*id]
	if !ok {
		return nil, idenpubonchain.ErrIdenNotOnChain
	}
	idenState, ok := idenStatesData.ByTime[queryBlockTs]
	if !ok {
		return nil, idenpubonchain.ErrIdenByTimeNotFound
	}
	return idenState, nil
}

// SetState updates the Identity State of the given ID in the IdenStates Smart Contract.
func (ip *IdenPubOnChain) SetState(id *core.ID, newState *merkletree.Hash,
	zkProof *zktypes.Proof) (*types.Transaction, error) {
	ip.rw.Lock()
	defer ip.rw.Unlock()
	idenStatesData, ok := ip.idenStatesData[*id]
	if !ok {
		return nil, idenpubonchain.ErrIdenNotOnChain
	}
	oldState := idenStatesData.IdenStates[len(idenStatesData.IdenStates)-1].IdenState
	if !ip.verifyZKP(zkProof, id, oldState, newState) {
		return nil, fmt.Errorf("zkproof verification failed")
	}
	idenState := proof.IdenStateData{
		BlockN:    ip.blockNow(),
		BlockTs:   ip.timeNow().Unix(),
		IdenState: newState,
	}
	ip.pendingSet = append(ip.pendingSet, &IdIdenStateData{Id: id, IdenStateData: &idenState})
	return &types.Transaction{}, nil
}

func g1ToBigInts(g1 *bn256.G1) [2]*big.Int {
	numBytes := 256 / 8
	bs := g1.Marshal()
	x := new(big.Int).SetBytes(bs[:numBytes])
	y := new(big.Int).SetBytes(bs[numBytes:])
	return [2]*big.Int{x, y}
}

func g2ToBigInts(g2 *bn256.G2) [2][2]*big.Int {
	numBytes := 256 / 8
	bs := g2.Marshal()
	xx := new(big.Int).SetBytes(bs[0*numBytes : 1*numBytes])
	xy := new(big.Int).SetBytes(bs[1*numBytes : 2*numBytes])
	yx := new(big.Int).SetBytes(bs[2*numBytes : 3*numBytes])
	yy := new(big.Int).SetBytes(bs[3*numBytes : 4*numBytes])
	// return [2][2]*big.Int{[2]*big.Int{xy, xx}, [2]*big.Int{yy, yx}}
	return [2][2]*big.Int{[2]*big.Int{xx, xy}, [2]*big.Int{yx, yy}}
}

/*
func g1ToBigInts(g1 *bn256.G1) [2]*big.Int {
	numBytes := 256 / 8
	bs := g1.Marshal()
	x := new(big.Int).SetBytes(cryptoUtils.SwapEndianness(bs[:numBytes]))
	y := new(big.Int).SetBytes(cryptoUtils.SwapEndianness(bs[numBytes:]))
	return [2]*big.Int{x, y}
}

func g2ToBigInts(g2 *bn256.G2) [2][2]*big.Int {
	numBytes := 256 / 8
	bs := g2.Marshal()
	xx := new(big.Int).SetBytes(cryptoUtils.SwapEndianness(bs[0*numBytes : 1*numBytes]))
	xy := new(big.Int).SetBytes(cryptoUtils.SwapEndianness(bs[1*numBytes : 2*numBytes]))
	yx := new(big.Int).SetBytes(cryptoUtils.SwapEndianness(bs[2*numBytes : 3*numBytes]))
	yy := new(big.Int).SetBytes(cryptoUtils.SwapEndianness(bs[3*numBytes : 4*numBytes]))
	println("DBG swapEndian, swap xy")
	return [2][2]*big.Int{[2]*big.Int{xy, xx}, [2]*big.Int{yy, yx}}
}
*/

func proofToBigInts(proof *zktypes.Proof) (a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) {
	a = g1ToBigInts(proof.A)
	b = g2ToBigInts(proof.B)
	c = g1ToBigInts(proof.C)
	return a, b, c
}

// InitState initializes the first Identity State of the given ID in the IdenStates Smart Contract.
func (ip *IdenPubOnChain) InitState(id *core.ID, genesisState,
	newState *merkletree.Hash, zkProof *zktypes.Proof) (*types.Transaction, error) {
	ip.rw.Lock()
	defer ip.rw.Unlock()
	_, ok := ip.idenStatesData[*id]
	if ok {
		return nil, fmt.Errorf("identity already exists on chain")
	}

	fmt.Printf(`    "id": "%v",
`,
		id.BigInt().String())
	fmt.Printf(`    "genesisState": "%v",
`,
		genesisState.BigInt().String())
	fmt.Printf(`    "newState": "%v",
`,
		newState.BigInt().String())
	proofA, proofB, proofC := proofToBigInts(zkProof)
	fmt.Printf(`    "a": ["%v",
			    "%v"],
`,
		proofA[0], proofA[1])
	fmt.Printf(`    "b": [
    			    ["%v",
			     "%v"],
			    ["%v",
			     "%v"]],
`,
		proofB[0][0], proofB[0][1], proofB[1][0], proofB[1][1])
	fmt.Printf(`    "c": ["%v",
			    "%v"]
`,
		proofC[0], proofC[1])
	proofJSON, err := zkparsers.ProofToJson(zkProof)
	if err != nil {
		panic(err)
	}
	fmt.Println("proof:", string(proofJSON))
	if !ip.verifyZKP(zkProof, id, genesisState, newState) {
		return nil, fmt.Errorf("zkproof verification failed")
	}
	idenState := proof.IdenStateData{
		BlockN:    ip.blockNow(),
		BlockTs:   ip.timeNow().Unix(),
		IdenState: newState,
	}
	ip.pendingInit = append(ip.pendingInit, &IdIdenStateData{Id: id, IdenStateData: &idenState})
	return types.NewTransaction(0, common.Address{}, nil, 0, nil,
		new(big.Int).SetUint64(ip.blockNow()).Bytes()), nil
}

// TxConfirmBlocks returns the number of confirmed blocks of transaction tx.
func (ip *IdenPubOnChain) TxConfirmBlocks(tx *types.Transaction) (*big.Int, error) {
	blockNumber := new(big.Int).SetBytes(tx.Data())
	currentBlock := new(big.Int).SetUint64(ip.blockNow())
	return currentBlock.Sub(currentBlock, blockNumber), nil
}

func (ip *IdenPubOnChain) verifyZKP(zkProof *zktypes.Proof,
	id *core.ID, oldState, newState *merkletree.Hash) bool {
	var idElem merkletree.ElemBytes
	copy(idElem[:], id[:])
	publicSignals := []*big.Int{idElem.BigInt(), oldState.BigInt(), newState.BigInt()}
	ok := verifier.Verify(ip.verifyingKey, zkProof, publicSignals)
	return ok
}
