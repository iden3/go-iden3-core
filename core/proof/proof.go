package proof

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	// common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/merkletree"
	// "github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	ErrRevokedClaim = errors.New("the claim is revoked: the next version exists")
)

// ProofClaimPartial is a proof of existence and non-existence of a claim in
// a single tree (only one level).
type ProofClaimPartial struct {
	Mtp0 *merkletree.Proof `json:"mtp0" binding:"required"`
	Mtp1 *merkletree.Proof `json:"mtp1" binding:"required"`
	Root *merkletree.Hash  `json:"root" binding:"required"`
}

// Returns either (false, nil) or (true, error)
func (p *ProofClaimPartial) Verify(claim *merkletree.Entry) (bool, error) {
	// Proof of existence verification
	if !p.Mtp0.Existence {
		return false, fmt.Errorf("Mtp0 is a non-existence proof")
	}
	if !merkletree.VerifyProof(p.Root, p.Mtp0, claim.HIndex(), claim.HValue()) {
		return false, fmt.Errorf("Mtp0 doesn't match with the root")
	}

	// Proof of non-existence of next version (revocation) verification
	if p.Mtp1.Existence {
		return false, fmt.Errorf("Mtp1 is an existence proof")
	}
	// Make a copy of the claim and increase the version
	claimNext := claim.Clone()
	claimType, claimVer := claims.GetClaimTypeVersionFromData(&claimNext.Data)
	claims.SetClaimTypeVersionInData(&claimNext.Data, claimType, claimVer+1)
	if !merkletree.VerifyProof(p.Root, p.Mtp1, claimNext.HIndex(), claimNext.HValue()) {
		return false, fmt.Errorf("Mtp1 doesn't match with the root")
	}
	return true, nil
}

func (pcp *ProofClaimPartial) String() string {
	buf := bytes.NewBufferString("ProofClaimPartial:\n")
	fmt.Fprintf(buf, "mtp0: %v\n", pcp.Mtp0)
	fmt.Fprintf(buf, "mtp0: %v\n", pcp.Mtp1)
	fmt.Fprintf(buf, "root: %v", pcp.Root)
	return buf.String()
}

// RelayAux is auxiliary data used to check a proof when the identity publishes
// roots via a Relay.
type RelayAux struct {
	// Version is SetRootClaim.Version
	Version uint32 `json:"version" binding:"required"`
	// Era is SetRootClaim.Era
	Era uint32 `json:"era" binding:"required"`
	// Proof is the ProofClaimPartial of the SetRootClaim in the Relay MT
	Proof ProofClaimPartial `json:"proof" binding:"required"`

	// RelayID is the ID of the Relay authorized by the Identity
	RelayID *core.ID `json:"relayId" binding:"required"`

	// GenesisRoot is the Genesis Root of the Identity
	// GenesisRoot *merkletree.Hash `json:"genesisRoot" binding:"required"`
	// MtpClaimAuthRelay is the mtp of the Identity Genesis ClaimAuthService authorizing the Relay
	MtpClaimAuthRelay *merkletree.Proof `json:"mtpAuthRelay" binding:"required"`
}

func (ra *RelayAux) ProofClaimAuthRelay(id *core.ID) *ProofClaimGenesis {
	return &ProofClaimGenesis{
		Mtp: ra.MtpClaimAuthRelay,
		Id:  id,
	}
}

func (ra *RelayAux) ClaimAuthRelay() *merkletree.Entry {
	return claims.NewClaimAuthorizeService(claims.ServiceTypeRelay, ra.RelayID.String(), "", "").Entry()
}

type RootData struct {
	BlockN         uint64
	BlockTimestamp int64
	Root           *merkletree.Hash
}

// ProofClaim is a complete proof of a claim that includes all the proofs of
// existence and non-existence for mutliple levels from the claim of a tree to
// the signed root of possibly another tree whose root binding:"required".
type ProofClaim struct {
	Claim          *merkletree.Entry `json:"claim" binding:"required"`
	ID             *core.ID          `json:"id" binding:"required"`
	BlockN         uint64            `json:"blockN" binding:"required"`
	BlockTimestamp int64             `json:"blockTS" binding:"required"`
	Proof          ProofClaimPartial `json:"proof" binding:"required,dive"`

	RelayAux *RelayAux `json:"relayAux"`
}

func (pc *ProofClaim) String() string {
	buf := bytes.NewBufferString("ProofClaim:\n")
	fmt.Fprintf(buf, "blockTS: %v\n", time.Unix(pc.BlockTimestamp, 0))
	fmt.Fprintf(buf, "blockN: %v\n", pc.BlockN)
	fmt.Fprintf(buf, "claim: %v\n", pc.Claim)
	if pc.RelayAux != nil {
		fmt.Fprintf(buf, "relayAux: %v", *pc.RelayAux)
	}
	return buf.String()
}

// PublishedData returns the id of the root publisher with the corresponding
// block number and block timestamp linked to publish the root.
func (pc *ProofClaim) PublishedData() (*core.ID, uint64, int64) {
	var publisherID *core.ID
	if pc.RelayAux != nil {
		publisherID = pc.RelayAux.RelayID
	} else {
		publisherID = pc.ID
	}
	return publisherID, pc.BlockN, pc.BlockTimestamp

}

// CheckProofClaim checks the claim proofs from the bottom to the top are valid and not revoked, and that the top root is signed by relayAddr.
// Returns either (false, nil) or (true, error)
// TODO Check id-root in the blockchain!
func (pc *ProofClaim) Verify(publishedRoot *merkletree.Hash) (bool, error) {
	var publisherClaim *merkletree.Entry
	if pc.RelayAux != nil {
		relayAux := pc.RelayAux
		// Verify that the identity has authorized the relay ID in a genesis claim
		proofClaimAuthRelay := relayAux.ProofClaimAuthRelay(pc.ID)
		if ok, err := proofClaimAuthRelay.Verify(relayAux.ClaimAuthRelay()); !ok {
			return false, fmt.Errorf("verification of ProofClaim.RelayAux.ProofClaimAuthRelay failed: %v", err)
		}
		// Verify that the claim is under the identity MT
		if ok, err := relayAux.Proof.Verify(pc.Claim); !ok {
			return false, fmt.Errorf("verification of ProofClaim.RelayAux.Proof failed: %v", err)
		}
		// Construct setRootClaim from the identity root
		setRootClaim, err := claims.NewClaimSetRootKey(pc.ID, relayAux.Proof.Root)
		if err != nil {
			return false, err
		}
		setRootClaim.Version = relayAux.Version
		setRootClaim.Era = relayAux.Era
		publisherClaim = setRootClaim.Entry()
	} else {
		publisherClaim = pc.Claim
	}

	// Verify that the publisherClaim is in the publisher MT
	if ok, err := pc.Proof.Verify(publisherClaim); !ok {
		return false, fmt.Errorf("verification of ProofClaim.Proof failed: %v", err)
	}

	// Verify that the root matches with the published root passed as argument
	if !pc.Proof.Root.Equals(publishedRoot) {
		return false, fmt.Errorf("ProofClaim root doesn't match the expected published root")
	}

	return true, nil
}

func VerifyGenesisMTProof(id *core.ID, proof *merkletree.Proof, hIndex, hValue *merkletree.Hash) (bool, error) {
	clr, err := merkletree.RootFromProof(proof, hIndex, hValue)
	if err != nil {
		return false, err
	}
	idenState := core.IdenState(clr, &merkletree.HashZero, clr)

	if eq := bytes.Equal(id[:], core.IdGenesisFromIdenState(idenState)[:]); !eq {
		return false, fmt.Errorf("calclated root doesn't match proof root")
	}
	return true, nil
}

// ProofClaimGenesis is a proof that a claim belongs to the genesis tree of an
// Id.
type ProofClaimGenesis struct {
	Mtp *merkletree.Proof `json:"mtp" binding:"required"`
	Id  *core.ID          `json:"id" binding:"required"`
}

// Verify that the claim belongs to the genesis tree with the specified root
// which was used to generate the Id.
func (p *ProofClaimGenesis) Verify(claim *merkletree.Entry) (bool, error) {
	if !p.Mtp.Existence {
		return false, fmt.Errorf("Mtp is a non-existence proof")
	}
	if ok, err := VerifyGenesisMTProof(p.Id, p.Mtp, claim.HIndex(), claim.HValue()); !ok {
		return false, fmt.Errorf("Mtp doesn't match with the genesis Id: %v", err)
	}
	return true, nil
}

// GetNonRevocationMTProof is a helper function to return a proof of non
// existence of the following version of a given claim (leafData).  If the
// following version exists, an error is returned.
func GetNonRevocationMTProof(mt *merkletree.MerkleTree, leafData *merkletree.Data, hi *merkletree.Hash) (*merkletree.Proof, error) {
	claimType, claimVersion := claims.GetClaimTypeVersionFromData(leafData)

	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], leafData[:])
	claims.SetClaimTypeVersionInData(leafDataCpy, claimType, claimVersion+1)
	entry := merkletree.Entry{
		Data: *leafDataCpy,
	}
	proof, err := mt.GenerateProof(entry.HIndex(), nil)
	if err != nil {
		return nil, err
	}
	if proof.Existence {
		return nil, ErrRevokedClaim
	}
	return proof, nil
}

// GetClaimProofByHi given a Hash(index) (Hi) and an id, returns the Claim
// in that Hi position inside the User merkletree, it's proof of existence and
// of non-revocation, and the proof of existence and of non-revocation for the
// set root claim in the relay tree, all in the form of a ProofClaim.  The
// result is not yet signed and has no timestamp.
func GetClaimProofByHi(mt *merkletree.MerkleTree, hi *merkletree.Hash) (*ProofClaim, error) {
	// get the value in the hi position
	leafData, err := mt.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}

	// get the MT proof of existence of the claim and the non-existence of
	// the claim's next version in the Relay Tree
	mtpExist, err := mt.GenerateProof(hi, nil)
	if err != nil {
		return nil, err
	}
	mtpNonExist, err := GetNonRevocationMTProof(mt, leafData, hi)
	if err != nil {
		return nil, err
	}

	rootKey := mt.RootKey()

	proofClaimPartial := ProofClaimPartial{
		Mtp0: mtpExist,
		Mtp1: mtpNonExist,
		Root: rootKey,
	}
	proofClaim := ProofClaim{
		Claim:          &merkletree.Entry{Data: *leafData},
		ID:             nil,
		BlockN:         0,
		BlockTimestamp: 0,
		Proof:          proofClaimPartial,
		RelayAux:       nil,
	}

	return &proofClaim, nil
}

type PredicateProof struct {
	LeafEntry               *merkletree.Entry
	MtpNonExistInOldRoot    *merkletree.Proof
	MtpExist                *merkletree.Proof
	MtpNonExistNextVersion  *merkletree.Proof
	MtpExistPreviousVersion *merkletree.Proof
	OldRoot                 *merkletree.Hash
	Root                    *merkletree.Hash
}

func GetPreviousVersionEntry(entry *merkletree.Entry) (*merkletree.Entry, error) {
	claimType, claimVer := claims.GetClaimTypeVersionFromData(&entry.Data)
	if claimVer == 0 {
		return nil, errors.New("claim is in version 0, can not exist a previous version")
	}
	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], entry.Data[:])
	claims.SetClaimTypeVersionInData(leafDataCpy, claimType, claimVer-1)
	entry1 := merkletree.Entry{
		Data: *leafDataCpy,
	}
	return &entry1, nil
}
func GetNextVersionEntry(entry *merkletree.Entry) *merkletree.Entry {
	claimType, claimVer := claims.GetClaimTypeVersionFromData(&entry.Data)
	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], entry.Data[:])
	claims.SetClaimTypeVersionInData(leafDataCpy, claimType, claimVer+1)
	entry1 := merkletree.Entry{
		Data: *leafDataCpy,
	}
	return &entry1
}

// GetPredicateProof, ϕ_min
// checks that:
// - 0: tree is updated incrementally
// 	- claim position was empty in oldRoot
// - 1: claim is added correctly
// 	- claim position contains the claim in currentRoot
// - 2: claim is not revocated
// 	- claim (version+1) is empty in currentRoot
// in case that the claim version != 0:
// - 3: claim is at the expected version
//	- claim (version-1) exist in oldRoot
// - 4: current version is incremental from the last one
//	- siblings of check_0 are inside siblings of check_1
//
// *TODO The output format will depend on the zkSnark inputs format (not specified yet)
func GetPredicateProof(mt *merkletree.MerkleTree, oldRoot, hi *merkletree.Hash) (*PredicateProof, error) {
	// proof_0: that claim position was empty in oldRoot
	mtpNonExistInOldRoot, err := mt.GenerateProof(hi, oldRoot)
	if err != nil {
		return nil, err
	}

	// proof_1: that claim position contains the claim in newRoot
	mtpExist, err := mt.GenerateProof(hi, nil)
	if err != nil {
		return nil, err
	}

	// proof_2: that claim (version+1) is empty in newRoot
	// get the value in the hi position
	leafData, err := mt.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}
	mtpNonExistNextVersion, err := GetNonRevocationMTProof(mt, leafData, hi)
	if err != nil {
		return nil, err
	}

	leafDataCpy := &merkletree.Data{}
	copy(leafDataCpy[:], leafData[:])
	entry := merkletree.Entry{
		Data: *leafDataCpy,
	}

	// checks 3 and 4 are necessary if the claim.Version != 0
	var mtpExistPrevVersion *merkletree.Proof
	_, v := claims.GetClaimTypeVersion(&entry)
	if v != 0 {
		// version is not 0 (v!=0), so we need to provide proof_3 and proof_4

		// proof_3: that claim (version-1) is not empty in oldRoot
		entryPrevVersion, err := GetPreviousVersionEntry(&entry)
		if err != nil {
			return nil, err
		}
		mtpExistPrevVersion, err = mt.GenerateProof(entryPrevVersion.HIndex(), oldRoot)
		if err != nil {
			return nil, err
		}

		// proof_4: that siblings of check_0 are inside of siblings of check_1
		// this proof don't needs more additional data
		// is something that the verifier needs to verify with the data from the other proofs
	}

	predicateProof := &PredicateProof{
		LeafEntry:               &entry,
		MtpNonExistInOldRoot:    mtpNonExistInOldRoot,   // proof_0
		MtpExist:                mtpExist,               // proof_1
		MtpNonExistNextVersion:  mtpNonExistNextVersion, // proof_2
		MtpExistPreviousVersion: mtpExistPrevVersion,    // proof_3
		OldRoot:                 oldRoot,
		Root:                    mt.RootKey(),
	}
	return predicateProof, nil
}

// VerifyPredicateProof, ϕ_min
// checks that:
// - 0: tree is updated incrementally
// 	- claim position was empty in oldRoot
// - 1: claim is added correctly
// 	- claim position contains the claim in currentRoot
// - 2: claim is not revocated
// 	- claim (version+1) is empty in currentRoot
// in case that the claim version != 0:
// - 3: claim is at the expected version
//	- claim (version-1) exist in oldRoot
// - 4: current version is incremental from the last one
//	- siblings of check_0 are inside of check_1
//
// *TODO The input format will depend on the zkSnark inputs format (not specified yet)
func VerifyPredicateProof(p *PredicateProof) bool {
	if bytes.Equal(p.Root.Bytes(), p.OldRoot.Bytes()) {
		return false
	}

	// check_0: that claim position was empty in oldRoot
	if p.MtpNonExistInOldRoot.Existence {
		// should be a proof of non existence, if not, verification fails
		return false
	}
	if !merkletree.VerifyProof(p.OldRoot, p.MtpNonExistInOldRoot, p.LeafEntry.HIndex(), p.LeafEntry.HValue()) {
		return false
	}

	// check_1: that claim position contains the claim in currentRoot
	if !p.MtpExist.Existence {
		// should be a proof of existence, if not, verification fails
		return false
	}
	if !merkletree.VerifyProof(p.Root, p.MtpExist, p.LeafEntry.HIndex(), p.LeafEntry.HValue()) {
		return false
	}

	// check_2: that claim (version+1) is empty in currentRoot
	if p.MtpNonExistNextVersion.Existence {
		// should be a proof of non existence, if not, verification fails
		return false
	}
	entry1 := GetNextVersionEntry(p.LeafEntry)
	if !merkletree.VerifyProof(p.Root, p.MtpNonExistNextVersion, entry1.HIndex(), entry1.HValue()) {
		return false
	}

	// checks 3 and 4 are necessary if the claim.Version != 0
	_, v := claims.GetClaimTypeVersion(p.LeafEntry)
	if v == 0 {
		// if version == 0, return true expected checks have passed
		return true
	}

	// check_3: that claim (version-1) is not empty in oldRoot
	if !p.MtpExistPreviousVersion.Existence {
		// should be a proof of existence, if not, verification fails
		return false
	}
	entryPrevVersion, err := GetPreviousVersionEntry(p.LeafEntry)
	if err != nil {
		// if err!=nil means that there is no previous version possible, as the current version is 0
		return false
	}
	if !merkletree.VerifyProof(p.OldRoot, p.MtpExistPreviousVersion, entryPrevVersion.HIndex(), entryPrevVersion.HValue()) {
		return false
	}

	// check_4: check that siblings of check_0 are inside the siblings of check_1
	// p.MtpNonExistInOldRoot.Siblings == p.MtpExist.Siblings[:len(p.MtpNonExistInOldRoot.Siblings)]
	for i := 0; i < len(p.MtpNonExistInOldRoot.Siblings); i++ {
		if !bytes.Equal(p.MtpNonExistInOldRoot.Siblings[i].Bytes(), p.MtpExist.Siblings[i].Bytes()) {
			return false
		}
	}

	return true
}

type IdenState struct {
	BlockTs int64
	BlockN  uint64
	Value   *merkletree.Hash
}

type CredentialExistence struct {
	Id        *core.ID
	IdenState IdenState
	MtpClaim  *merkletree.Proof
	Claim     *merkletree.Entry
	RevRoot   *merkletree.Hash
	RooRoot   *merkletree.Hash
	IdPub     string
}

type CredentialValidity struct {
	CredentialExistence CredentialExistence
	IdenState           IdenState
	MtpNotNonce         *merkletree.Proof
	ClaRoot             *merkletree.Hash
	RooRoot             *merkletree.Hash
}
