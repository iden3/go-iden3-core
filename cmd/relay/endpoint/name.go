package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/id/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/namesrv"
)

func handleVinculateID(c *gin.Context) {
	var vinculateIDMsg namesrv.VinculateIDMsg
	c.BindJSON(&vinculateIDMsg)
	assignNameClaim, err := nameservice.VinculateID(vinculateIDMsg)
	if err != nil {
		fail(c, "error name.VinculateID", err)
	}

	// return claim with proofs
	proofOfRelayClaim, err := claimservice.GetRelayClaimByHi(config.C.Namespace, assignNameClaim.Hi())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"assignNameClaim":   common3.BytesToHex(assignNameClaim.Bytes()),
		"name":              vinculateIDMsg.Name,
		"ethID":             assignNameClaim.EthID,
		"proofOfRelayClaim": proofOfRelayClaim.Hex(),
	})
}
func handleAssignNameClaimResolv(c *gin.Context) {
	nameid := c.Param("nameid")

	assignNameClaim, err := nameservice.ResolvAssignNameClaim(nameid, config.C.Namespace)
	if err != nil {
		fail(c, "nameid not found in merkletree", err)
		return
	}
	c.JSON(200, gin.H{
		"claim": common3.BytesToHex(assignNameClaim.Bytes()),
		"ethID": assignNameClaim.EthID,
	})
}
