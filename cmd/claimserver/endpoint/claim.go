package endpoint

import (

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/utils"

	"github.com/iden3/go-iden3/merkletree"
)

type claimData struct {
	IdAddr  	common.Address
	Data      string
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	var m claimData
	c.BindJSON(&m)

	data, err := common3.HexDecode(m.Data)
	if err != nil {
		genericserver.Fail(c, "error on handlePostClaim", err)
		return
	}

	var nameHash [256 / 8]byte
	hash := utils.HashBytes([]byte(data))
	copy(nameHash[:], hash[len(hash)-256/8:])
	claim := core.NewClaimLinkObjectIdentity(0, 5, 0, m.IdAddr, nameHash)

	err = genericserver.Claimservice.AddLinkObjectClaim(*claim)
	if err != nil {
		genericserver.Fail(c, "error on AddLinkObjectClaim", err)
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
	return
}

// handleGetClaimProofByHi handles the request to query the claim proof of a
// server claim (by hIndex).
func handleGetClaimProofByHi(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		genericserver.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofOfClaim, err := genericserver.Claimservice.GetClaimProofByHi(hi)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofOfClaim,
	})
	return
}
