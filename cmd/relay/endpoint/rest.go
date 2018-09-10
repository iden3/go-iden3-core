package endpoint

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claim"
	"github.com/iden3/go-iden3/services/web3"
	log "github.com/sirupsen/logrus"
)

func fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}

func handleGetRoot(c *gin.Context) {
	// get the contract data
	contractAddress := common.HexToAddress(config.C.ContractsAddress.Identities)
	root, err := web3srv.GetRoot(contractAddress)
	if err != nil {
		fail(c, "error contract.GetRoot(contractAddress)", err)
		return
	}
	c.JSON(200, gin.H{
		"root":         mt.Root().Hex(),
		"contractRoot": common3.BytesToHex(root[:]),
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
	typeBytes := bytesValue[32:64]
	switch common3.BytesToHex(typeBytes) {
	case merkletree.HashBytes([]byte("default")).Hex():
		break
	case merkletree.HashBytes([]byte("assignname")).Hex():
		assignNameClaim, err := core.ParseAssignNameClaimBytes(bytesValue)
		if err != nil {
			fail(c, "error on parsing AssignNameClaim bytes", err)
			return
		}

		privK, err := crypto.HexToECDSA(config.C.Server.PrivK)
		if err != nil {
			fail(c, "error on parsing server.PrivK", err)
			return
		}
		_, mp, sig, err := claimsrv.AddAssignNameClaim(mt, assignNameClaim, config.C.ContractsAddress.Identities, privK)
		if err != nil {
			fail(c, "error on AddAssignNameClaim", err)
			return
		}
		// return claim with proofs and signatures
		c.JSON(200, gin.H{
			"sig":  sig,
			"root": mt.Root().Hex(),
			"mp":   mp,
		})
		return
	case merkletree.HashBytes([]byte("authorizeksign")).Hex():
		authorizeKSignClaim, err := core.ParseAuthorizeKSignClaimBytes(bytesValue)
		if err != nil {
			fail(c, "error on parsing AuthorizeKSignClaim bytes", err)
			return
		}
		authorizeKSignClaimMsg := claimsrv.AuthorizeKSignClaimMsg{
			authorizeKSignClaim,
			bytesSignedMsg.SignatureHex,
		}
		claimProof, idRootProof, err := claimsrv.AddAuthorizeKSignClaim(mt, idaddr, authorizeKSignClaimMsg, config.C.ContractsAddress.Identities)
		if err != nil {
			fail(c, "error on AddAuthorizeKSignClaim", err)
			return
		}
		// return claim with proofs and signatures
		c.JSON(200, gin.H{
			"claimProof":  claimProof,
			"root":        mt.Root().Hex(),
			"idRootProof": idRootProof,
		})
		return
	case merkletree.HashBytes([]byte("setroot")).Hex():
		break
	default:
		fail(c, "type not found", errors.New("claim type not found"))
	}
}

func handleGetIDRoot(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	idRoot, idRootProof, err := claimsrv.GetIDRoot(mt, idaddr)
	if err != nil {
		fail(c, "error on GetIDRoot", err)
		return
	}
	c.JSON(200, gin.H{
		"root":        mt.Root().Hex(), // relay root
		"idRoot":      idRoot,          // user id root
		"idRootProof": idRootProof,     // user id root proof in the relay merkletree
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
	claim, idProof, idRoot, setClaimRoot, relayProof, relayRoot, err := claimsrv.GetClaimByHi(mt, config.C.Namespace, idaddr, hi)
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"claim":        common3.BytesToHex(claim.Bytes()),
		"idProof":      common3.BytesToHex(idProof),
		"setClaimRoot": common3.BytesToHex(setClaimRoot.Bytes()),
		"idRoot":       idRoot.Hex(),
		"relayProof":   common3.BytesToHex(relayProof),
		"relayRoot":    relayRoot.Hex(),
	})
	return
}
