package endpoint

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"

	"github.com/iden3/go-iden3/merkletree"
)

// handleCommitNewIDRoot handles a request to set the root key of a user tree
// though a set root claim.
func handleCommitNewIDRoot(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)

	var setRootMsg claimsrv.SetRootMsg
	c.BindJSON(&setRootMsg)

	idaddrMsg := common.HexToAddress(setRootMsg.IdAddr)

	// make sure that the given idaddr from the post url matches with the idaddr from the post data
	if !bytes.Equal(idaddr.Bytes(), idaddrMsg.Bytes()) {
		fail(c, "error on PostNewRoot, idaddr not match", errors.New("PostNewRoot idaddr not match"))
		return
	}
	// get signature from setRootClaimMsg
	signature, err := common3.HexToBytes(setRootMsg.Signature)
	if err != nil {
		fail(c, "error on PostNewRoot parse signature", err)
		return
	}
	rootBytes, err := common3.HexToBytes(setRootMsg.Root)
	if err != nil {
		fail(c, "error on PostNewRoot parse root", err)
		return
	}
	var root merkletree.Hash
	copy(root[:], rootBytes[:32])

	// add the root through claimservice
	setRootClaim, err := claimservice.CommitNewIDRoot(idaddr, &setRootMsg.KSignPk.PublicKey,
		root, setRootMsg.Timestamp, signature)
	if err != nil {
		fail(c, "error on AddAuthorizeKSignClaim", err)
		return
	}

	// return claim with proofs
	proofOfRelayClaim, err := claimservice.GetClaimProofByHi(setRootClaim.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofOfClaim": proofOfRelayClaim,
	})
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	var bytesSignedMsg claimsrv.BytesSignedMsg
	c.BindJSON(&bytesSignedMsg)

	bytesValue, err := common3.HexToBytes(bytesSignedMsg.ValueHex)
	if err != nil {
		fail(c, "error on parsing bytesSignedMsg.HexValue to bytes", err)
		return
	}

	// bytesValue to Element data
	var dataBytes [128]byte
	copy(dataBytes[:], bytesValue)
	data := merkletree.BytesToData(dataBytes)
	entry := merkletree.Entry{
		Data: *data,
	}

	claimValueMsg := claimsrv.ClaimValueMsg{
		ClaimValue: entry,
		Signature:  bytesSignedMsg.SignatureHex,
		KSignPk:    bytesSignedMsg.KSignPk,
	}
	err = claimservice.AddUserIDClaim(idaddr, claimValueMsg)
	if err != nil {
		fail(c, "error on AddUserIDClaim", err)
		return
	}
	// return claim with proofs
	proofOfClaim, err := claimservice.GetClaimProofUserByHi(idaddr, entry.HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}

	c.JSON(200, gin.H{
		"proofOfClaim": proofOfClaim,
	})
	return
}

// handleGetIDRoot handles a request to query the root key of a user tree.
func handleGetIDRoot(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	idRoot, idRootProof, err := claimservice.GetIDRoot(idaddr)
	if err != nil {
		fail(c, "error on GetIDRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        claimservice.MT().RootKey().Hex(), // relay root
		"idRoot":      idRoot.Hex(),                      // user id root
		"idRootProof": common3.BytesToHex(idRootProof),   // user id root proof in the relay merkletree
	})
	return
}

// handleGetClaimProofUserByHi handles the request to query the claim proof of
// a user claim (by hIndex).
func handleGetClaimProofUserByHi(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	hihex := c.Param("hi")
	hiBytes, err := common3.HexToBytes(hihex)
	if err != nil {
		fail(c, "error on HexToBytes of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	idaddr := common.HexToAddress(idaddrhex)
	proofOfClaim, err := claimservice.GetClaimProofUserByHi(idaddr, hi)
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofOfClaim": proofOfClaim,
	})
	return
}

// handleGetClaimProofByHi handles the request to query the claim proof of a
// relay claim (by hIndex).
func handleGetClaimProofByHi(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexToBytes(hihex)
	if err != nil {
		fail(c, "error on HexToBytes of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofOfClaim, err := claimservice.GetClaimProofByHi(hi)
	if err != nil {
		fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofOfClaim": proofOfClaim,
	})
	return
}
