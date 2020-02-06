package idenmanager

// DEPRECATED in favour of identity/issuer

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/components/idenmanager/messages"
	"github.com/iden3/go-iden3-core/components/idensigner"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/core/proof"
	crypto3 "github.com/iden3/go-iden3-core/crypto"
	"github.com/iden3/go-iden3-core/db"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"

	// "github.com/iden3/go-iden3-core/services/idenstatewriter"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

var PREFIX_MERKLETREE = []byte("merkletree")

var (
	ErrNotFound = errors.New("value not found")
)

type IdenManager struct {
	id *core.ID
	mt *merkletree.MerkleTree
	// idenStateWriter idenstatewriter.IdenStateWriter
	signer idensigner.IdenSigner
}

func New(id *core.ID, mt *merkletree.MerkleTree,
	signer idensigner.IdenSigner) *IdenManager {
	return &IdenManager{id, mt, signer}
}

// ID returns the id.
func (m *IdenManager) ID() *core.ID {
	return m.id
}

// MT returns the merkle tree.
func (m *IdenManager) MT() *merkletree.MerkleTree {
	return m.mt
}

// IdenStateWriter returns the IdenStateWriter
// func (m *IdenManager) IdenStateWriter() idenstatewriter.IdenStateWriter {
// 	return m.idenStateWriter
// }

// CheckSetRootParams checks the params corresponding to the SetRoot0Req.
// 1. Check that id and ProofClaimAuthKOp.Id match
// 2. Parse ProofClaimAuthKOp.Claim to get kOp
// 3. Verify that sig(kOp, oldRoot+newRoot) == signature
// 4. Verify ProofClaimAuthKOp
func CheckSetRootParams(id *core.ID, setRootReq messages.SetRoot0Req) (bool, error) {

	//// 1. Check that id and ProofClaimAuthKOp.Id match
	if !id.Equal(setRootReq.ProofClaimAuthKOp.Id) {
		return false, fmt.Errorf("id and ProofClaimAuthKOp.Id don't match")
	}

	//// 2. Parse ClaimAuthKOp to get kOp
	claim, err := claims.NewClaimFromEntry(setRootReq.ClaimAuthKOp)
	if err != nil {
		return false, fmt.Errorf("Error parsing ClaimAuthKOp: %v", err)
	}
	claimAuthorizeKSign, ok := claim.(*claims.ClaimAuthorizeKSignBabyJub)
	if !ok {
		return false, fmt.Errorf("Invalid claim type in ClaimAuthKOp.Claim," +
			"expected ClaimAuthorizeKSignBabyJub")
	}

	//// 3. Verify that sig(kOp, oldRoot+newRoot) == signature
	kSignComp := claimAuthorizeKSign.PublicKeyComp()
	msg := append(setRootReq.OldRoot[:], setRootReq.NewRoot[:]...)
	// check the signature with PrefixMinorUpdate
	if ok, err := babykeystore.VerifySignature(kSignComp, setRootReq.Signature, msg, setRootReq.Date, babykeystore.PrefixMinorUpdate); !ok {
		return false, fmt.Errorf("root signature doesn't match with kOp from ClaimAuthKOp: %v", err)
	}
	//// 4. Verify ProofClaimAuthKOp
	if ok, err := setRootReq.ProofClaimAuthKOp.Verify(setRootReq.ClaimAuthKOp); !ok {
		return false, fmt.Errorf("Verification of ProofClaimAuthKOp failed: %v", err)
	}
	return true, nil
}

func (m *IdenManager) UpdateSetRootClaim(id *core.ID, setRootReq messages.SetRoot0Req) (*claims.ClaimSetRootKey, error) {
	ok, err := CheckSetRootParams(id, setRootReq)
	if err != nil || !ok {
		return nil, fmt.Errorf("SetRoot params verification not passed, " + err.Error())
	}

	// Add new SetRootClaim with id -> setRootReq.NewRoot
	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey, err := claims.NewClaimSetRootKey(id, setRootReq.NewRoot)
	if err != nil {
		return nil, err
	}
	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	if err = m.mt.AddClaim(claimSetRootKey); err != nil {
		return nil, err
	}

	// update Relay Root in Smart Contract
	// m.idenStateWriter.SetRoot(*m.mt.RootKey())

	return claimSetRootKey, nil
}

// SetNewIdRoot checks that the data is valid and performs a claim in the Relay merkletree setting the new Root of the emiting Id
func (m *IdenManager) CommitNewIdRoot(id core.ID, kSignPk *ecdsa.PublicKey, root merkletree.Hash,
	timestamp int64, signature *crypto3.SignatureEthMsg) (*claims.ClaimSetRootKey, error) {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := m.mt.Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIddb(userMT, kSignPk) {
		return nil, errors.New("can not verify the KSign")
	}
	// in the future the user merkletree will be in the client side, and this step will be a check of the ProofKSign

	// check data timestamp
	verified := utils.VerifyTimestamp(timestamp, 30000) //needs to be from last 30 seconds
	if !verified {
		return nil, errors.New("timestamp too old")
	}
	// check signature with id
	// whee data signed is id+root+timestamp
	timestampBytes := common3.Uint64ToEthBytes(uint64(timestamp))
	// signature of id+root+timestamp, only valid if is from last X seconds
	var msg []byte
	msg = append(msg, id.Bytes()...)
	msg = append(msg, root.Bytes()...)
	msg = append(msg, timestampBytes...)

	if !crypto3.VerifySigEthMsg(crypto.PubkeyToAddress(*kSignPk), signature, msg) {
		return nil, errors.New("signature can not be verified")
	}

	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey, err := claims.NewClaimSetRootKey(&id, &root)
	if err != nil {
		return nil, err
	}
	// entry := claimSetRootKey.Entry()
	// version, err := GetNextVersion(m.mt, entry.HIndex())
	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = m.mt.AddClaim(claimSetRootKey)
	if err != nil {
		return nil, err
	}

	// update Relay Root in Smart Contract
	// m.idenStateWriter.SetRoot(*m.mt.RootKey())

	return claimSetRootKey, nil
}

// AddClaimAuthorizeKSignSecp256k1First adds ClaimAuthorizeKSignSecp256k1 into
// the Id's merkletree, and adds the Id's merkle root into the Relay's
// merkletree inside a ClaimSetRootKey. Returns the merkle proof of both Claims
func (m *IdenManager) AddClaimAuthorizeKSignSecp256k1First(id core.ID,
	claimAuthorizeKSignSecp256k1 claims.ClaimAuthorizeKSignSecp256k1) error {

	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := m.mt.Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return err
	}

	// add ClaimAuthorizeKSign into the User's Id Merkle Tree
	err = userMT.AddClaim(&claimAuthorizeKSignSecp256k1)
	if err != nil {
		return err
	}

	// create new ClaimSetRootKey
	claimSetRootKey, err := claims.NewClaimSetRootKey(&id, userMT.RootKey())
	if err != nil {
		return err
	}

	// get next version of the claim
	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = m.mt.AddClaim(claimSetRootKey)
	if err != nil {
		return err
	}

	// update Relay's Root in the Smart Contract
	// m.idenStateWriter.SetRoot(*m.mt.RootKey())

	return nil
}

// AddUserIdClaim adds a claim into the Id's merkle tree, and with the Id's root, creates a new ClaimSetRootKey and adds it to the Relay's merkletree
func (m *IdenManager) AddUserIdClaim(id *core.ID, claimValueMsg messages.ClaimValueReq) error {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := m.mt.Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIddb(userMT, &claimValueMsg.KSignPk.PublicKey) {
		return errors.New("can not verify the KSign")
	}

	// verify signature with KSign
	if !crypto3.VerifySigEthMsg(crypto.PubkeyToAddress(claimValueMsg.KSignPk.PublicKey),
		claimValueMsg.Signature, claimValueMsg.ClaimValue.Bytes()) {
		return errors.New("signature can not be verified")
	}

	// add claim in User Id Merkle Tree
	err = userMT.AddEntry(&claimValueMsg.ClaimValue)
	if err != nil {
		return err
	}

	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey, err := claims.NewClaimSetRootKey(id, userMT.RootKey())
	if err != nil {
		return err
	}
	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = m.mt.AddClaim(claimSetRootKey)
	if err != nil {
		return err
	}

	// update Relay Root in Smart Contract
	// m.idenStateWriter.SetRoot(*m.mt.RootKey())

	return nil
}

// AddClaim adds a claim directly to the Relay merkletree
func (m *IdenManager) AddClaim(claim merkletree.Entrier) error {
	err := m.mt.AddClaim(claim)
	if err != nil {
		return err
	}
	// m.idenStateWriter.SetRoot(*m.mt.RootKey())
	return nil
}

// GetIdRoot returns the root of an Id tree, and the proof of that Root Id tree in the Relay Merkle Tree
func (m *IdenManager) GetIdRoot(id *core.ID) (merkletree.Hash, []byte, error) {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := m.mt.Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}

	// build ClaimSetRootKey of the user id
	claimSetRootKey, err := claims.NewClaimSetRootKey(id, userMT.RootKey())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}

	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	claimSetRootKey.Version = version - 1

	// get proof of SetRootProof in the Relay tree
	idRootProof, err := m.mt.GenerateProof(claimSetRootKey.Entry().HIndex(), nil)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	return *userMT.RootKey(), idRootProof.Bytes(), nil
}

// GetSetRootClaim returns the last SetRootKey claim corresponding to an id
// with a proof to the root in the blockchain.
func (m *IdenManager) GetSetRootClaim(id *core.ID) (*proof.ProofClaim, error) {
	// build ClaimSetRootKey of the user id
	claimSetRootKey, err := claims.NewClaimSetRootKey(id, &merkletree.Hash{})
	if err != nil {
		return nil, err
	}

	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimSetRootKey.Version = version - 1

	return m.GetClaimProofByHiBlockchain(claimSetRootKey.Entry().HIndex())
}

// TODO: Remove this
// GetClaimProofUserByHiOld given a Hash(index) (Hi) and an Id, returns the Claim in that Hi position inside the Id's merkletree, and the ClaimSetRootKey with the Id's root in the Relay's merkletree
func (m *IdenManager) GetClaimProofUserByHiOld(id *core.ID, hi *merkletree.Hash) (*messages.ProofClaimUserRes, error) {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := m.mt.Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, err
	}

	// get the value in the hi position
	// valueBytes, err := userMT.GetValueInPos(hi)
	leafData, err := userMT.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}
	// if bytes.Equal(valueBytes, merkletree.EmptyNodeValue[:]) {
	//         return nil, ErrNotFound
	// }

	// value, err := core.ParseValueFromBytes(valueBytes)
	// if err != nil {
	//         return nil, err
	// }

	// get the proof of the value in the User Id Tree
	// idProof, err := userMT.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]), nil)
	idProof, err := userMT.GenerateProof(hi, nil)
	if err != nil {
		return nil, err
	}

	leafBytes := leafData.Bytes()
	claimProof := messages.ProofTreeLeaf{
		Leaf:  leafBytes[:],
		Proof: idProof.Bytes(),
		Root:  *userMT.RootKey(),
	}

	// build ClaimSetRootKey
	claimSetRootKey, err := claims.NewClaimSetRootKey(id, userMT.RootKey())
	if err != nil {
		return nil, err
	}
	version, err := GetNextVersion(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimSetRootKey.Version = version - 1

	// get the proof of the ClaimSetRootKey in the Relay Tree
	relayProof, err := m.mt.GenerateProof(claimSetRootKey.Entry().HIndex(), nil)
	if err != nil {
		return nil, err
	}
	claimSetRootKeyProof := messages.ProofTreeLeaf{
		Leaf:  claimSetRootKey.Entry().Bytes(),
		Proof: relayProof.Bytes(),
		Root:  *m.mt.RootKey(),
	}

	// get non revocation proofs of the claim
	claimNonRevocationProof, err := getNonRevocationProof(userMT, hi)
	if err != nil {
		return nil, err
	}
	claimSetRootKeyNonRevocationProof, err := getNonRevocationProof(m.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}

	// sign root + date
	sig, date, err := m.signer.SignEthMsgDate(claimSetRootKeyProof.Root[:])
	if err != nil {
		return nil, err
	}

	proofClaim := messages.ProofClaimUserRes{
		ClaimProof:                     claimProof,
		SetRootClaimProof:              *claimNonRevocationProof,
		ClaimNonRevocationProof:        claimSetRootKeyProof,
		SetRootClaimNonRevocationProof: *claimSetRootKeyNonRevocationProof,
		Date:                           date,
		Signature:                      sig[:],
	}
	return &proofClaim, nil
}

// GetClaimProofByHiBlockchain given a Hash(index) (Hi), returns the Claim in that Hi
// position inside the Relay merkletree, and it's proof of existence and of
// non-revocated, all in the form of a ProofClaim, using a root that is
// published in the blockchain.  The result is signed (with
// a timestamp) by the service.
func (m *IdenManager) GetClaimProofByHiBlockchain(hi *merkletree.Hash) (*proof.ProofClaim, error) {
	// rootData, err := m.idenStateWriter.GetRoot(m.id)
	rootData, err := &proof.RootData{}, fmt.Errorf("DEPRECATED")
	if err != nil {
		return nil, err
	}
	mt, err := m.mt.Snapshot(rootData.Root)
	if err != nil {
		return nil, err
	}
	fmt.Printf("> G %+v\n", rootData.Root.String())
	fmt.Printf("> B %+v\n", mt.RootKey().String())
	proofClaim, err := proof.GetClaimProofByHi(mt, hi)
	if err != nil {
		return nil, err
	}

	proofClaim.ID = m.id
	proofClaim.BlockN = rootData.BlockN
	proofClaim.BlockTimestamp = rootData.BlockTimestamp

	return proofClaim, nil
}

// CreateIdGenesis initializes the id MerkleTree with the given the kop, kdisable,
// kreenable and kupdateRoots public keys. Where the id is calculated a MerkleTree containing
// that initial data, calculated in the function CalculateIdGenesis()
func (m *IdenManager) CreateIdGenesis(kop *babyjub.PublicKey, kdis, kreen, kupdateRoot common.Address) (*core.ID, *proof.ProofClaim, error) {

	id, proofClaims, err := genesis.CalculateIdGenesisFrom4Keys(kop, kdis, kreen, kupdateRoot)
	if err != nil {
		return nil, nil, err
	}

	// add the claims into the storage merkletree of that identity
	stoUserId := m.MT().Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, nil, err
	}

	proofClaimsList := []proof.ProofClaim{proofClaims.KOp, proofClaims.KDis,
		proofClaims.KReen, proofClaims.KUpdateRoot}
	for _, proofClaim := range proofClaimsList {
		err = userMT.AddEntry(proofClaim.Claim)
		if err != nil {
			return nil, nil, err
		}
	}

	// create new ClaimSetRootKey
	claimSetRootKey, err := claims.NewClaimSetRootKey(id, userMT.RootKey())
	if err != nil {
		return nil, nil, err
	}

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = m.MT().AddClaim(claimSetRootKey)
	if err != nil {
		return nil, nil, err
	}

	// update Relay's Root in the Smart Contract
	// m.idenStateWriter.SetRoot(*m.MT().RootKey())

	return id, &proofClaims.KOp, nil
}

// getNonRevocationProof returns the next version Hi (that don't exist in the tree, it's value is Empty) with merkleproof and root
func getNonRevocationProof(mt *merkletree.MerkleTree, hi *merkletree.Hash) (*messages.ProofTreeLeaf, error) {
	// var value merkletree.Value

	// get claim value in bytes
	leafData, err := mt.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}

	claimType, _ := claims.GetClaimTypeVersionFromData(leafData)
	nextVersion, err := GetNextVersion(mt, hi)
	if err != nil {
		return nil, err
	}

	claims.SetClaimTypeVersionInData(leafData, claimType, nextVersion)

	entry := merkletree.Entry{
		Data: *leafData,
	}
	mp, err := mt.GenerateProof(entry.HIndex(), nil)
	if err != nil {
		return nil, err
	}
	leafBytes := entry.Bytes()
	nonRevocationProof := messages.ProofTreeLeaf{
		Leaf:  leafBytes[:],
		Proof: mp.Bytes(),
		Root:  *mt.RootKey(),
	}
	return &nonRevocationProof, nil
}

// GetNextVersion returns the next version of a claim, given a Hash(index)
func GetNextVersion(mt *merkletree.MerkleTree, hi *merkletree.Hash) (uint32, error) {
	var claimType claims.ClaimType
	var version uint32

	// loop until we find a nextversion that don't exist
	for {
		leafData, err := mt.GetDataByIndex(hi)
		if err == merkletree.ErrEntryIndexNotFound {
			return version, nil
		} else if err != nil {
			return 0, err
		}
		claimType, version = claims.GetClaimTypeVersionFromData(leafData)
		version++

		claims.SetClaimTypeVersionInData(leafData, claimType, version)

		entry := merkletree.Entry{
			Data: *leafData,
		}
		hi = entry.HIndex()
	}
}

// NewMerkleTreeUser creates a new user merkle tree by using an storage with
// the user addres prefix.
func NewMerkleTreeUser(id *core.ID, storage db.Storage, levels int) (*merkletree.MerkleTree, error) {
	stoUserId := storage.WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)
	if userMT, err := merkletree.NewMerkleTree(stoUserId, levels); err != nil {
		return nil, err
	} else {
		return userMT, nil
	}
}
