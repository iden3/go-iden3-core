package endpoint

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"

	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
)

func handleCommitNewIDRoot(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)

	var setRootMsg claimsrv.SetRootMsg
	c.BindJSON(&setRootMsg)

	idaddrMsg := common.HexToAddress(setRootMsg.IdAddr)
	kSign := common.HexToAddress(setRootMsg.KSign)

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

	// add the root throught claimservice
	setRootClaim, err := claimservice.CommitNewIDRoot(idaddr, kSign, root, setRootMsg.Timestamp, signature)
	if err != nil {
		fail(c, "error on AddAuthorizeKSignClaim", err)
		return
	}

	// return claim with proofs
	proofOfRelayClaim, err := claimservice.GetRelayClaimByHi(setRootClaim.Hi())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofOfClaim": proofOfRelayClaim.Hex(),
	})
}

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
	typeBytes := bytesValue[32:56]

	switch common3.BytesToHex(typeBytes) {
	case common3.BytesToHex(core.DefaultType):
		claimDefault, err := core.ParseGenericClaimBytes(bytesValue)
		if err != nil {
			fail(c, "error on parsing GenericClaim bytes", err)
			return
		}

		claimValueMsg := claimsrv.ClaimValueMsg{
			claimDefault,
			bytesSignedMsg.SignatureHex,
			bytesSignedMsg.KSign,
		}
		err = claimservice.AddUserIDClaim(idaddr, claimValueMsg)
		if err != nil {
			fail(c, "error on AddUserIDClaim", err)
			return
		}
		// return claim with proofs
		proofOfClaim, err := claimservice.GetClaimByHi(idaddr, claimDefault.Hi())
		if err != nil {
			fail(c, "error on GetClaimByHi", err)
			return
		}
		c.JSON(200, gin.H{
			"proofOfClaim": proofOfClaim.Hex(),
		})
		return

	case common3.BytesToHex(core.AssignNameType):
		assignNameClaim, err := core.ParseAssignNameClaimBytes(bytesValue)
		if err != nil {
			fail(c, "error on parsing AssignNameClaim bytes", err)
			return
		}

		err = claimservice.AddAssignNameClaim(assignNameClaim)
		if err != nil {
			fail(c, "error on AddAssignNameClaim", err)
			return
		}

		// return claim with proofs
		proofOfClaim, err := claimservice.GetClaimByHi(idaddr, assignNameClaim.Hi())
		if err != nil {
			fail(c, "error on GetClaimByHi", err)
			return
		}
		c.JSON(200, gin.H{
			"proofOfClaim": proofOfClaim.Hex(),
		})
		return

	case common3.BytesToHex(core.AuthorizeksignType):
		authorizeKSignClaim, err := core.ParseAuthorizeKSignClaimBytes(bytesValue)
		if err != nil {
			fail(c, "error on parsing AuthorizeKSignClaim bytes", err)
			return
		}
		authorizeKSignClaimMsg := claimsrv.AuthorizeKSignClaimMsg{
			authorizeKSignClaim,
			bytesSignedMsg.SignatureHex,
			bytesSignedMsg.KSign,
		}
		err = claimservice.AddAuthorizeKSignClaim(idaddr, authorizeKSignClaimMsg)
		if err != nil {
			fail(c, "error on AddAuthorizeKSignClaim", err)
			return
		}
		// return claim with proofs
		proofOfClaim, err := claimservice.GetClaimByHi(idaddr, authorizeKSignClaim.Hi())
		if err != nil {
			fail(c, "error on GetClaimByHi", err)
			return
		}
		c.JSON(200, gin.H{
			"proofOfClaim": proofOfClaim.Hex(),
		})
		return

	case common3.BytesToHex(core.SetRootType):
		break

	default:
		fail(c, "type not found", errors.New("claim type not found"))
	}
}

func handleGetIDRoot(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	idRoot, idRootProof, err := claimservice.GetIDRoot(idaddr)
	if err != nil {
		fail(c, "error on GetIDRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        claimservice.MT().Root().Hex(),  // relay root
		"idRoot":      idRoot.Hex(),                    // user id root
		"idRootProof": common3.BytesToHex(idRootProof), // user id root proof in the relay merkletree
	})
	return
}

func handleGetClaimByHi(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	hihex := c.Param("hi")
	hiBytes, err := common3.HexToBytes(hihex)
	if err != nil {
		fail(c, "error on HexToBytes of Hi", err)
		return
	}
	var hi merkletree.Hash
	copy(hi[:], hiBytes)
	idaddr := common.HexToAddress(idaddrhex)
	proofOfClaim, err := claimservice.GetClaimByHi(idaddr, hi)
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofOfClaim": proofOfClaim.Hex(),
	})
	return
}
