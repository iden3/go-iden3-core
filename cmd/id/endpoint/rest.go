package endpoint

import (
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}
	color.Red(msg)
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}

func handleVinculateID(c *gin.Context) {
	/*
		var vinculateIDMsg namesrv.VinculateIDMsg
		c.BindJSON(&vinculateIDMsg)
		privK, err := crypto.HexToECDSA(config.C.Relay.PrivK)
		if err != nil {
			fail(c, "error on parsing server.PrivK", err)
			return
		}
		assignNameClaim, err := namesrv.VinculateID(mt, vinculateIDMsg, config.C.ContractsAddress.Identities, privK)
		if err != nil {
			fail(c, "error name.VinculateID", err)
		}

		// return claim with proofs and signatures
		c.JSON(200, assignNameClaim)
	*/
}
func handleAssignNameClaimResolv(c *gin.Context) {
	/*
		nameid := c.Param("nameid")

		assignNameClaim, err := claimsrv.ResolvAssignNameClaim(mt, nameid, config.C.Namespace)
		if err != nil {
			fail(c, "nameid not found in merkletree", err)
			return
		}
		c.JSON(200, gin.H{
			"success": "ok",
			"claim":   assignNameClaim,
			"ethAddr":   assignNameClaim.EthAddr,
		})
	*/
}
