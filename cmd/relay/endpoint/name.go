package endpoint

import (
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/namesrv"
)

func handleVinculateID(c *gin.Context) {
	var vinculateIDMsg namesrv.VinculateIDMsg
	c.BindJSON(&vinculateIDMsg)
	claimAssignName, err := nameservice.VinculateID(vinculateIDMsg)
	if err != nil {
		fail(c, "error name.VinculateID", err)
	}

	// return claim with proofs
	proofOfClaimAssignName, err := claimservice.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"claimAssignName":        common3.BytesToHex(claimAssignName.Entry().Bytes()),
		"name":                   vinculateIDMsg.Name,
		"ethAddr":                claimAssignName.EthAddr,
		"proofOfClaimAssignName": proofOfClaimAssignName,
	})
}
func handleClaimAssignNameResolv(c *gin.Context) {
	nameid := c.Param("nameid")

	claimAssignName, err := nameservice.ResolvClaimAssignName(nameid)
	if err != nil {
		fail(c, "nameid not found in merkletree", err)
		return
	}

	proofOfClaimAssignName, err := claimservice.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"claim":                  common3.BytesToHex(claimAssignName.Entry().Bytes()),
		"ethAddr":                claimAssignName.EthAddr,
		"proofOfClaimAssignName": proofOfClaimAssignName,
	})
}
