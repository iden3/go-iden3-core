package endpoint

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
)

func handleVinculateId(c *gin.Context) {
	var signedPacket signedpacketsrv.SignedPacket
	if err := c.BindJSON(&signedPacket); err != nil {
		genericserver.Fail(c, "BindJSON", err)
		return
	}
	if err := signedPacketVerifier.VerifySignedPacketGeneric(&signedPacket); err != nil {
		genericserver.Fail(c, "invalid signed packet", err)
		return
	}
	form := signedPacket.Payload.Form.(map[string]string)

	claimAssignName, err := nameService.VinculateId(form["assignName"], genericserver.C.Domain,
		signedPacket.Header.Issuer)
	if err != nil {
		genericserver.Fail(c, "error name.VinculateId", err)
		return
	}

	// return claim with proofs
	proofClaimAssignName, err := claimService.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"ethName":         fmt.Sprintf("%v@%v", form["assignName"], genericserver.C.Domain),
		"proofAssignName": proofClaimAssignName,
	})
}

func handleClaimAssignNameResolv(c *gin.Context) {
	nameid := c.Param("name")

	claimAssignName, err := nameService.ResolvClaimAssignName(nameid)
	if err != nil {
		genericserver.Fail(c, "nameid not found in merkletree", err)
		return
	}

	proofClaimAssignName, err := claimService.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetClaimByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"id":              claimAssignName.Id,
		"proofAssignName": proofClaimAssignName,
	})
}
