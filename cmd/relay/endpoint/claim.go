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

// handleCommitNewIdRoot handles a request to set the root key of a user tree
// though a set root claim.
func handleCommitNewIdRoot(c *gin.Context) {
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
	rootBytes, err := common3.HexDecode(setRootMsg.Root)
	if err != nil {
		fail(c, "error on PostNewRoot parse root", err)
		return
	}
	var root merkletree.Hash
	copy(root[:], rootBytes[:32])

	// add the root through claimservice
	setRootClaim, err := claimservice.CommitNewIdRoot(idaddr, &setRootMsg.KSignPk.PublicKey,
		root, setRootMsg.Timestamp, setRootMsg.Signature)
	if err != nil {
		fail(c, "error on AddAuthorizeKSignClaim", err)
		return
	}

	// return claim with proofs
	proofRelayClaim, err := claimservice.GetClaimProofByHi(setRootClaim.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofRelayClaim,
	})
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	var bytesSignedMsg claimsrv.BytesSignedMsg
	c.BindJSON(&bytesSignedMsg)

	bytesValue, err := common3.HexDecode(bytesSignedMsg.ValueHex)
	if err != nil {
		fail(c, "error on parsing bytesSignedMsg.HexValue to bytes", err)
		return
	}

	// bytesValue to Element data
	var dataBytes [128]byte
	copy(dataBytes[:], bytesValue)
	data := merkletree.NewDataFromBytes(dataBytes)
	entry := merkletree.Entry{
		Data: *data,
	}

	claimValueMsg := claimsrv.ClaimValueMsg{
		ClaimValue: entry,
		Signature:  bytesSignedMsg.Signature,
		KSignPk:    bytesSignedMsg.KSignPk,
	}
	err = claimservice.AddUserIdClaim(idaddr, claimValueMsg)
	if err != nil {
		fail(c, "error on AddUserIdClaim", err)
		return
	}
	// return claim with proofs
	proofClaim, err := claimservice.GetClaimProofUserByHi(idaddr, entry.HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}

	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}

// handleGetIdRoot handles a request to query the root key of a user tree.
func handleGetIdRoot(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	idRoot, idRootProof, err := claimservice.GetIdRoot(idaddr)
	if err != nil {
		fail(c, "error on GetIdRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        claimservice.MT().RootKey().Hex(), // relay root
		"idRoot":      idRoot.Hex(),                      // user id root
		"proofIdRoot": common3.HexEncode(idRootProof),    // user id root proof in the relay merkletree
	})
	return
}

// handleGetClaimProofUserByHi handles the request to query the claim proof of
// a user claim (by hIndex).
func handleGetClaimProofUserByHi(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	idaddr := common.HexToAddress(idaddrhex)
	proofClaim, err := claimservice.GetClaimProofUserByHi(idaddr, hi)
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}

// handleGetClaimProofByHi handles the request to query the claim proof of a
// relay claim (by hIndex).
func handleGetClaimProofByHi(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofClaim, err := claimservice.GetClaimProofByHi(hi)
	if err != nil {
		fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}
