package endpoint

import (
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/namesrv"
)

func handleVinculateId(c *gin.Context) {
	var vinculateIdMsg namesrv.VinculateIdMsg
	c.BindJSON(&vinculateIdMsg)
	claimAssignName, err := nameservice.VinculateId(vinculateIdMsg)
	if err != nil {
		fail(c, "error name.VinculateId", err)
	}

	// return claim with proofs
	proofClaimAssignName, err := claimservice.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"claimAssignName":      common3.HexEncode(claimAssignName.Entry().Bytes()),
		"name":                 vinculateIdMsg.Name,
		"idAddr":               claimAssignName.IdAddr,
		"proofClaimAssignName": proofClaimAssignName,
	})
}

func handleClaimAssignNameResolv(c *gin.Context) {
	nameid := c.Param("nameid")

	claimAssignName, err := nameservice.ResolvClaimAssignName(nameid)
	if err != nil {
		fail(c, "nameid not found in merkletree", err)
		return
	}

	proofClaimAssignName, err := claimservice.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"claim":                common3.HexEncode(claimAssignName.Entry().Bytes()),
		"idAddr":               claimAssignName.IdAddr,
		"proofClaimAssignName": proofClaimAssignName,
	})
}
