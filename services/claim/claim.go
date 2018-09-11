package claimsrv

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/web3"
	"github.com/iden3/go-iden3/utils"
)

// GetNextVersion returns the next version of a claim, given a Hash(index)
func GetNextVersion(mt *merkletree.MerkleTree, hi merkletree.Hash) (uint32, error) {

	// merkletree.GetValueInPos(hi)
	b, err := mt.GetValueInPos(hi)
	if err != nil {
		return 0, err
	}
	// if value not exist, return version 0
	if bytes.Equal(merkletree.EmptyNodeValue[:], b) {
		return 0, nil
	}
	// get version bytes
	versionBytes := b[64:68]
	version := common3.BytesToUint32(versionBytes)
	version++
	// return version
	return version, nil
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
	assignNameClaim.Version = version

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
	mp, err := mt.GenerateProof(assignNameClaim)
	if err != nil {
		return merkletree.Hash{}, []byte{}, []byte{}, err
	}
	// TODO add proof of non revocated

	// sign root
	// privK, err := crypto.HexToECDSA(relayPrivK)
	// if err != nil {
	// 	return merkletree.Hash{}, []byte{}, []byte{}, err
	// }
	sig, err := crypto.Sign(mt.Root().Bytes(), relayPrivK)
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
	// verify signature of the AuthorizeKSignClaim
	signature, err := common3.HexToBytes(authorizeKSignClaimMsg.Signature)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	msgHash := authorizeKSignClaimMsg.AuthorizeKSignClaim.Ht()
	fmt.Println(common3.BytesToHex(signature))
	fmt.Println(ethID.Hex())
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
	setRootClaim := core.NewSetRootClaim("iden3.io", authorizeKSignClaimMsg.AuthorizeKSignClaim.ExtraIndex.KeyToAuthorize, userMT.Root())
	// get next version of the claim
	version, err := GetNextVersion(mt, setRootClaim.Hi())
	if err != nil {
		return []byte{}, []byte{}, err
	}
	setRootClaim.Version = version
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
	claimProof, err := userMT.GenerateProof(authorizeKSignClaimMsg.AuthorizeKSignClaim)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	idRootProof, err := mt.GenerateProof(setRootClaim)
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
	msgHash := merkletree.HashBytes(claimValueMsg.ClaimValue.Bytes())
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
	setRootClaim.BaseIndex.Version++ // TODO autoincrement
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
	claimProof, err := userMT.GenerateProof(claimValueMsg.ClaimValue)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	idRootProof, err := mt.GenerateProof(setRootClaim)
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
	// get proof of SetRootProof in the Relay tree
	idRootProof, err := mt.GenerateProof(setRootClaim)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	return userMT.Root(), idRootProof, nil
}

// GetClaimByHi given a Hash(index) (Hi) and an ID, returns the Claim in that Hi position inside the ID's merkletree, and the SetRootClaim with the ID's root in the Relay's merkletree
func GetClaimByHi(mt *merkletree.MerkleTree, namespace string, ethID common.Address, hi merkletree.Hash) (merkletree.Value, []byte, merkletree.Hash, merkletree.Value, []byte, merkletree.Hash, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := mt.Storage().WithPrefix(ethID.Bytes())
	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return nil, []byte{}, merkletree.Hash{}, nil, []byte{}, merkletree.Hash{}, err
	}
	// get the value in the hi position
	valueBytes, err := userMT.GetValueInPos(hi)
	if err != nil {
		return nil, []byte{}, merkletree.Hash{}, nil, []byte{}, merkletree.Hash{}, err
	}
	value, err := core.ParseValueFromBytes(valueBytes)
	if err != nil {
		return nil, []byte{}, merkletree.Hash{}, nil, []byte{}, merkletree.Hash{}, err
	}
	// get the proof of the value in the User ID Tree
	idProof, err := userMT.GenerateProof(value)
	if err != nil {
		return nil, []byte{}, merkletree.Hash{}, nil, []byte{}, merkletree.Hash{}, err
	}
	// build SetRootClaim
	setRootClaim := core.NewSetRootClaim(namespace, ethID, userMT.Root())
	setRootClaim.BaseIndex.Version++ // TODO autoincrement version
	// get the proof of the SetRootClaim in the Relay Tree
	relayProof, err := mt.GenerateProof(setRootClaim)
	if err != nil {
		return nil, []byte{}, merkletree.Hash{}, nil, []byte{}, merkletree.Hash{}, err
	}

	return value, idProof, userMT.Root(), setRootClaim, relayProof, mt.Root(), nil

}
