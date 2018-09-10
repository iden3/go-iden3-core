package namesrv

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/cmd/id/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claim"
	"github.com/iden3/go-iden3/utils"
)

// VinculateID creates an adds a AssignNameClaim vinculating a name and an address, into the merkletree
func VinculateID(mt *merkletree.MerkleTree, vinculateIDMsg VinculateIDMsg, identitiesContractHex string, relayPrivK *ecdsa.PrivateKey) (core.AssignNameClaim, error) {
	// TODO calculate EthID from the AssignNameClaim.RawIdentityTx
	ethID := common.HexToAddress(vinculateIDMsg.Msg.EthID) // tmp

	// verify vinculateIDMsg.Msg signature with EthID
	// sigBytes := hexutil.MustDecode(vinculateIDMsg.MsgSignature)
	sigBytes, err := common3.HexToBytes(vinculateIDMsg.MsgSignature)
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	msgHash := vinculateIDMsg.MsgHash()
	verified := utils.VerifySig(ethID, sigBytes, msgHash[:])
	if !verified {
		return core.AssignNameClaim{}, errors.New("signature can not be verified")
	}
	// add AssignNameClaim to merkle tree
	namespaceHash := merkletree.HashBytes([]byte(config.C.Namespace))
	nameHash := merkletree.HashBytes([]byte(vinculateIDMsg.Msg.Name))
	domainHash := merkletree.HashBytes([]byte(config.C.Domain))
	assignNameClaim := core.NewAssignNameClaim(namespaceHash, nameHash, domainHash, ethID)
	// signature, err := utils.Sign(assignNameClaim.Ht(), relayPrivK)
	// if err != nil {
	// 	return core.AssignNameClaim{}, errors.New("error signing")
	// }
	// signatureHex := common3.BytesToHex(signature)
	// assignNameClaimMsg := claimsrv.AssignNameClaimMsg{
	// 	assignNameClaim,
	// 	signatureHex,
	// }
	// send the AssignNameClaim to the service/relay
	_, _, _, err = claimsrv.AddAssignNameClaim(mt, assignNameClaim, identitiesContractHex, relayPrivK)
	return assignNameClaim, err
}
