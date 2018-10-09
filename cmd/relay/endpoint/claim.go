package endpoint

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/claimsrv"

	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
)

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
	case common3.BytesToHex(core.DefaultTypeHash[:24]):
		claimDefault, err := core.ParseClaimDefaultBytes(bytesValue)
		if err != nil {
			fail(c, "error on parsing ClaimDefault bytes", err)
			return
		}

		claimValueMsg := claimsrv.ClaimValueMsg{
			claimDefault,
			bytesSignedMsg.SignatureHex,
			bytesSignedMsg.KSign,
			bytesSignedMsg.ProofOfKSignHex,
		}
		err = claimservice.AddUserIDClaim(config.C.Namespace, idaddr, claimValueMsg)
		if err != nil {
			fail(c, "error on AddUserIDClaim", err)
			return
		}
		// return claim with proofs
		proofOfClaim, err := claimservice.GetClaimByHi(config.C.Namespace, idaddr, claimDefault.Hi())
		if err != nil {
			fail(c, "error on GetClaimByHi", err)
			return
		}
		c.JSON(200, gin.H{
			"proofOfClaim": proofOfClaim.Hex(),
		})
		return

	case common3.BytesToHex(core.AssignNameTypeHash[:24]):
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
		proofOfClaim, err := claimservice.GetClaimByHi(config.C.Namespace, idaddr, assignNameClaim.Hi())
		if err != nil {
			fail(c, "error on GetClaimByHi", err)
			return
		}
		c.JSON(200, gin.H{
			"proofOfClaim": proofOfClaim.Hex(),
		})
		return

	case common3.BytesToHex(core.AuthorizeksignTypeHash[:24]):
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
		proofOfClaim, err := claimservice.GetClaimByHi(config.C.Namespace, idaddr, authorizeKSignClaim.Hi())
		if err != nil {
			fail(c, "error on GetClaimByHi", err)
			return
		}
		c.JSON(200, gin.H{
			"proofOfClaim": proofOfClaim.Hex(),
		})
		return

	case common3.BytesToHex(core.SetRootTypeHash[:24]):
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
	proofOfClaim, err := claimservice.GetClaimByHi(config.C.Namespace, idaddr, hi)
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofOfClaim": proofOfClaim.Hex(),
	})
	return
}
