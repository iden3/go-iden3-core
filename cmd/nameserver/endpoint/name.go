package endpoint

import (
	"fmt"

	"github.com/gin-gonic/gin"
	cfg "github.com/iden3/go-iden3/cmd/nameserver/config"
	"github.com/iden3/go-iden3/core"
)

func handleVinculateId(c *gin.Context) {
	var signedPacket core.SignedPacket
	if err := c.BindJSON(&signedPacket); err != nil {
		fail(c, "BindJSON", err)
		return
	}
	if err := core.VerifySignedPacketGeneric(&signedPacket); err != nil {
		fail(c, "invalid generic signed packet", err)
		return
	}
	form := signedPacket.Payload.Form.(map[string]string)

	claimAssignName, err := nameservice.VinculateId(form["assignName"], cfg.C.Domain,
		signedPacket.Header.Issuer)
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
		"ethName":         fmt.Sprintf("%v@%v", form["assignName"], cfg.C.Domain),
		"proofAssignName": proofClaimAssignName,
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
		"idAddr":          claimAssignName.IdAddr,
		"proofAssignName": proofClaimAssignName,
	})
}
