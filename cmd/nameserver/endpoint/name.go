package endpoint

import (
	"github.com/gin-gonic/gin"
	cfg "github.com/iden3/go-iden3/cmd/nameserver/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
)

func handleVinculateId(c *gin.Context) {
	var form struct {
		AssignName string `json:"assignName" binding:"required"`
	}
	var signedPacket core.SignedPacket
	signedPacket.Payload.Form = &form
	if err := c.BindJSON(&signedPacket); err != nil {
		fail(c, "BindJSON", err)
		return
	}
	if err := core.VerifySignedPacket(&signedPacket); err != nil {
		fail(c, "invalid signed packet", err)
		return
	}
	if signedPacket.Payload.Type != core.GENERICSIGV01 {
		fail(c, "invalid signed packet payload type", nil)
		return
	}

	claimAssignName, err := nameservice.VinculateId(form.AssignName, cfg.C.Domain, signedPacket.Header.Issuer)
	if err != nil {
		fail(c, "error name.VinculateId", err)
		return
	}

	// return claim with proofs
	proofClaimAssignName, err := claimservice.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"claimAssignName":      claimAssignName.Entry(),
		"name":                 form.AssignName,
		"idAddr":               signedPacket.Header.Issuer,
		"proofClaimAssignName": proofClaimAssignName,
	})
}
func handleClaimAssignNameResolv(c *gin.Context) {
	nameid := c.Param("name")

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
