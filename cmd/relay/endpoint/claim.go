package endpoint

import (
	"bytes"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/claimsrv"

	"github.com/iden3/go-iden3/merkletree"
)

// handleCommitNewIdRoot handles a request to set the root key of a user tree
// though a set root claim.
func handleCommitNewIdRoot(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}

	var setRootMsg claimsrv.SetRootMsg
	err = c.BindJSON(&setRootMsg)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	idMsg, err := core.IDFromString(setRootMsg.Id)
	if err != nil {
		genericserver.Fail(c, "error on parse setRootMsg.Id", err)
		return
	}

	// make sure that the given id from the post url matches with the id from the post data
	if !bytes.Equal(id.Bytes(), idMsg.Bytes()) {
		genericserver.Fail(c, "error on PostNewRoot, id not match", errors.New("PostNewRoot id not match"))
		return
	}
	rootBytes, err := common3.HexDecode(setRootMsg.Root)
	if err != nil {
		genericserver.Fail(c, "error on PostNewRoot parse root", err)
		return
	}
	var root merkletree.Hash
	copy(root[:], rootBytes[:32])

	// add the root through genericserver.Claimservice
	setRootClaim, err := genericserver.Claimservice.CommitNewIdRoot(id, &setRootMsg.KSignPk.PublicKey,
		root, setRootMsg.Timestamp, setRootMsg.Signature)
	if err != nil {
		genericserver.Fail(c, "error on AddAuthorizeKSignClaim", err)
		return
	}

	// return claim with proofs
	proofRelayClaim, err := genericserver.Claimservice.GetClaimProofByHi(setRootClaim.Entry().HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofRelayClaim,
	})
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}
	var bytesSignedMsg claimsrv.BytesSignedMsg
	err = c.BindJSON(&bytesSignedMsg)
	if err != nil {
		genericserver.Fail(c, "json parsing error", err)
		return
	}

	bytesValue, err := common3.HexDecode(bytesSignedMsg.ValueHex)
	if err != nil {
		genericserver.Fail(c, "error on parsing bytesSignedMsg.HexValue to bytes", err)
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
	err = genericserver.Claimservice.AddUserIdClaim(id, claimValueMsg)
	if err != nil {
		genericserver.Fail(c, "error on AddUserIdClaim", err)
		return
	}
	// return claim with proofs
	proofClaim, err := genericserver.Claimservice.GetClaimProofUserByHi(id, entry.HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetClaimByHi", err)
		return
	}

	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}

// handleGetIdRoot handles a request to query the root key of a user tree.
func handleGetIdRoot(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}
	idRoot, idRootProof, err := genericserver.Claimservice.GetIdRoot(id)
	if err != nil {
		genericserver.Fail(c, "error on GetIdRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        genericserver.Claimservice.MT().RootKey().Hex(), // relay root
		"idRoot":      idRoot.Hex(),                                    // user id root
		"proofIdRoot": common3.HexEncode(idRootProof),                  // user id root proof in the relay merkletree
	})
	return
}

// handleGetClaimProofUserByHi handles the request to query the claim proof of
// a user claim (by hIndex).
func handleGetClaimProofUserByHi(c *gin.Context) {
	idHex := c.Param("id")
	hiHex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hiHex)
	if err != nil {
		genericserver.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	id, err := core.IDFromString(idHex)
	if err != nil {
		genericserver.Fail(c, "error on parse id", err)
		return
	}
	proofClaim, err := genericserver.Claimservice.GetClaimProofUserByHi(id, hi)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimByHi", err)
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
		genericserver.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofClaim, err := genericserver.Claimservice.GetClaimProofByHi(hi)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofClaim,
	})
	return
}
