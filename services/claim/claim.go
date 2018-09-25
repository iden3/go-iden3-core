package claimsrv

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/web3"
	"github.com/iden3/go-iden3/utils"
)

// GetNextVersion returns the next version of a claim, given a Hash(index)
func GetNextVersion(mt *merkletree.MerkleTree, hi merkletree.Hash) (uint32, error) {
	var version uint32
	// loop until found a nextversion that don't exist
	for {
		// merkletree.GetValueInPos(hi)
		b, err := mt.GetValueInPos(hi)
		if err != nil {
			return 0, err
		}
		// if value not exist, return version 0
		if bytes.Equal(merkletree.EmptyNodeValue[:], b) {
			break
		}
		// get version bytes
		versionBytes := b[60:64]
		version = core.EthBytesToUint32(versionBytes)
		version++

		// get claim with version+1 from the merkletree
		versionBytes, err = core.Uint32ToEthBytes(version)
		if err != nil {
			return 0, err
		}
		copy(b[60:64], versionBytes)
		value, err := core.ParseValueFromBytes(b)
		if err != nil {
			return 0, err
		}
		hi = merkletree.HashBytes(value.Bytes()[:value.IndexLength()])
	}
	// return version
	return version, nil
}

// GetNonRevocationProof returns the next version Hi (that don't exist in the tree, it's value is Empty) with merkleproof and root
func GetNonRevocationProof(mt *merkletree.MerkleTree, hi merkletree.Hash) (ProofOfTreeLeaf, error) {
	var value merkletree.Value
	// get claim value in bytes
	b, err := mt.GetValueInPos(hi)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	if bytes.Equal(b, merkletree.EmptyNodeValue[:]) {
		return ProofOfTreeLeaf{}, errors.New("Hi not found in the merkle tree")
	}

	nextVersion, err := GetNextVersion(mt, hi)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}

	// get claim with version+1 from the merkletree
	nextVersionBytes, err := core.Uint32ToEthBytes(nextVersion)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	copy(b[60:64], nextVersionBytes)

	value, err = core.ParseValueFromBytes(b)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	mp, err := mt.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	nonRevocationProof := ProofOfTreeLeaf{
		Leaf:  merkletree.EmptyNodeValue[:],
		Hi:    merkletree.HashBytes(value.Bytes()[:value.IndexLength()]),
		Proof: mp,
		Root:  mt.Root(),
	}
	return nonRevocationProof, nil
}

// AddAssignNameClaim adds AssignNameClaim into the merkletree, updates the root in the smart contract, and returns the merkle proof of the claim in the merkletree
func AddAssignNameClaim(mt *merkletree.MerkleTree, assignNameClaim core.AssignNameClaim, identitiesContractHex string, relayPrivK *ecdsa.PrivateKey) (merkletree.Hash, []byte, []byte, error) {
	// TODO maybe it have no sense to verify a signature of the relay in this claim, if to call this function is passed already the relayPrivK
	// ---
	// verify signature
	// signature, err := common3.HexToBytes(assignNameClaimMsg.Signature)
	// if err != nil {
	// 	return merkletree.Hash{}, []byte{}, []byte{}, err
	// }
	// msgHash := assignNameClaimMsg.AssignNameClaim.Ht()
	// if !utils.VerifySig(assignNameClaimMsg.AssignNameClaim.EthID, signature, msgHash[:]) {
	// 	return merkletree.Hash{}, []byte{}, []byte{}, errors.New("signature can not be verified")
	// }
	// ---

	// get next version of the claim
	version, err := GetNextVersion(mt, assignNameClaim.Hi())
	if err != nil {
		return merkletree.Hash{}, []byte{}, []byte{}, err
	}
	assignNameClaim.BaseIndex.Version = version

	// add AssignNameClaim to the Relay's merkletree
	// err := mt.Add(assignNameClaimMsg.AssignNameClaim)
	err = mt.Add(assignNameClaim)
	if err != nil {
		log.Fatal(err)
		return merkletree.Hash{}, []byte{}, []byte{}, err
	}

	// update relay's root in smart contract
	contractAddress := common.HexToAddress(identitiesContractHex)
	err = web3srv.AddRoot(mt.Root(), contractAddress)
	if err != nil {
		return merkletree.Hash{}, []byte{}, []byte{}, err
	}

	// generate proofs mp
	mp, err := mt.GenerateProof(assignNameClaim.Hi())
	if err != nil {
		return merkletree.Hash{}, []byte{}, []byte{}, err
	}
	// TODO add proof of non revocated

	// sign root
	sig, err := utils.Sign(mt.Root(), relayPrivK)
	if err != nil {
		return merkletree.Hash{}, []byte{}, []byte{}, err
	}
	return mt.Root(), mp, sig, nil
}

// ResolvAssignNameClaim returns the AssignNameClaim from the merkletree, given a nameid and a namespace
func ResolvAssignNameClaim(mt *merkletree.MerkleTree, nameid, namespace string) (core.AssignNameClaim, error) {
	// get name and domain
	s := strings.Split(nameid, "@")
	name := s[0]
	domain := s[1]

	// build the AssignNameClaim Partial with the given data of the Index
	nameHash := merkletree.HashBytes([]byte(name))
	domainHash := merkletree.HashBytes([]byte(domain))
	// domainHash := merkletree.HashBytes([]byte(domain))
	claimPartial := core.NewAssignNameClaim(namespace, nameHash, domainHash, common.Address{})
	version, err := GetNextVersion(mt, claimPartial.Hi())
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	claimPartial.BaseIndex.Version = version - 1
	// get the complete AssignNameClaim in that merkle tree position
	claimInPosBytes, err := mt.GetValueInPos(claimPartial.Hi())
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	if bytes.Equal(claimInPosBytes, merkletree.EmptyNodeValue[:]) {
		return core.AssignNameClaim{}, errors.New("not found")
	}
	assignNameClaim, err := core.ParseAssignNameClaimBytes(claimInPosBytes)
	return assignNameClaim, err
}

// AddAuthorizeKSignClaim adds AuthorizeKSignClaim into the ID's merkletree, and adds the ID's merkle root into the Relay's merkletree inside a SetRootClaim. Returns the merkle proof of both Claims
func AddAuthorizeKSignClaim(mt *merkletree.MerkleTree, ethID common.Address, authorizeKSignClaimMsg AuthorizeKSignClaimMsg, identitiesContractHex string) ([]byte, []byte, error) {
	// TODO this function is almost the same than AddUserIDClaim() function. Maybe delete this one and keep the generic one
	// verify signature of the AuthorizeKSignClaim
	signature, err := common3.HexToBytes(authorizeKSignClaimMsg.Signature)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	msgHash := utils.EthHash(authorizeKSignClaimMsg.AuthorizeKSignClaim.Bytes())
	signature[64] -= 27
	if !utils.VerifySig(ethID, signature, msgHash[:]) {
		return []byte{}, []byte{}, errors.New("signature can not be verified")
	}

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethID.Bytes())
	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// add AuthorizeKSignClaim into the User's ID Merkle Tree
	err = userMT.Add(authorizeKSignClaimMsg.AuthorizeKSignClaim)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// create new SetRootClaim
	setRootClaim := core.NewSetRootClaim("iden3.io", ethID, userMT.Root())
	// get next version of the claim
	version, err := GetNextVersion(mt, setRootClaim.Hi())
	if err != nil {
		return []byte{}, []byte{}, err
	}
	setRootClaim.BaseIndex.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = mt.Add(setRootClaim)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// update Relay's Root in the Smart Contract
	contractAddress := common.HexToAddress(identitiesContractHex)
	err = web3srv.AddRoot(mt.Root(), contractAddress)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	// return the proof of the AuthorizeKSignClaim, and the proof of the SetRootClaim
	claimProof, err := userMT.GenerateProof(authorizeKSignClaimMsg.AuthorizeKSignClaim.Hi())
	if err != nil {
		return []byte{}, []byte{}, err
	}
	idRootProof, err := mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return claimProof, idRootProof, nil
}

// AddUserIDClaim adds a claim into the ID's merkle tree, and with the ID's root, creates a new SetRootClaim and adds it to the Relay's merkletree
func AddUserIDClaim(mt *merkletree.MerkleTree, namespace string, ethID common.Address, claimValueMsg ClaimValueMsg, identitiesContractHex string) ([]byte, []byte, error) {
	// verify signature
	signature, err := common3.HexToBytes(claimValueMsg.Signature)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	// msgHash := merkletree.HashBytes(claimValueMsg.ClaimValue.Bytes())
	msgHash := utils.EthHash(claimValueMsg.ClaimValue.Bytes())
	signature[64] -= 27
	if !utils.VerifySig(ethID, signature, msgHash[:]) {
		return []byte{}, []byte{}, errors.New("signature can not be verified")
	}

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethID.Bytes())
	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// add claim in User ID Merkle Tree
	err = userMT.Add(claimValueMsg.ClaimValue)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// setRootClaim of the user in the Relay Merkle Tree
	// create new SetRootClaim
	setRootClaim := core.NewSetRootClaim(namespace, ethID, userMT.Root())
	version, err := GetNextVersion(mt, setRootClaim.Hi())
	if err != nil {
		return []byte{}, []byte{}, err
	}
	setRootClaim.BaseIndex.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = mt.Add(setRootClaim)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	// update Relay Root in Smart Contract
	contractAddress := common.HexToAddress(identitiesContractHex)
	err = web3srv.AddRoot(mt.Root(), contractAddress)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	// return proof of SetRoot and original Claim
	hiClaimValue := merkletree.HashBytes(claimValueMsg.ClaimValue.Bytes()[:claimValueMsg.ClaimValue.IndexLength()])
	claimProof, err := userMT.GenerateProof(hiClaimValue)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	idRootProof, err := mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return claimProof, idRootProof, nil
}

// GetIDRoot returns the root of an ID tree, and the proof of that Root ID tree in the Relay Merkle Tree
func GetIDRoot(mt *merkletree.MerkleTree, ethID common.Address) (merkletree.Hash, []byte, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethID.Bytes())
	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	// build SetRootClaim of the user id
	setRootClaim := core.NewSetRootClaim("iden3.io", ethID, userMT.Root())
	version, err := GetNextVersion(mt, setRootClaim.Hi())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	setRootClaim.BaseIndex.Version = version - 1
	// get proof of SetRootProof in the Relay tree
	idRootProof, err := mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	return userMT.Root(), idRootProof, nil
}

// GetClaimByHi given a Hash(index) (Hi) and an ID, returns the Claim in that Hi position inside the ID's merkletree, and the SetRootClaim with the ID's root in the Relay's merkletree
func GetClaimByHi(mt *merkletree.MerkleTree, namespace string, ethID common.Address, hi merkletree.Hash) (ProofOfTreeLeaf, ProofOfTreeLeaf, ProofOfTreeLeaf, ProofOfTreeLeaf, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}

	// get the value in the hi position
	valueBytes, err := userMT.GetValueInPos(hi)
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}
	value, err := core.ParseValueFromBytes(valueBytes)
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}

	// get the proof of the value in the User ID Tree
	idProof, err := userMT.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}

	claimProof := ProofOfTreeLeaf{
		Leaf:  valueBytes,
		Hi:    merkletree.HashBytes(value.Bytes()[:value.IndexLength()]),
		Proof: idProof,
		Root:  userMT.Root(),
	}

	// build SetRootClaim
	setRootClaim := core.NewSetRootClaim(namespace, ethID, userMT.Root())
	version, err := GetNextVersion(mt, setRootClaim.Hi())
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}
	setRootClaim.BaseIndex.Version = version - 1
	// get the proof of the SetRootClaim in the Relay Tree
	relayProof, err := mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}
	setRootClaimProof := ProofOfTreeLeaf{
		Leaf:  setRootClaim.Bytes(),
		Hi:    setRootClaim.Hi(),
		Proof: relayProof,
		Root:  mt.Root(),
	}
	// get non revocation proofs of the claim
	claimNonRevocationProof, err := GetNonRevocationProof(userMT, merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}
	setRootClaimNonRevocationProof, err := GetNonRevocationProof(mt, setRootClaim.Hi())
	if err != nil {
		return ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, ProofOfTreeLeaf{}, err
	}

	return claimProof, setRootClaimProof, claimNonRevocationProof, setRootClaimNonRevocationProof, nil
}
