package endpoint

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/cauth/challenge"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/utils"
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

func handleAuth(c *gin.Context) {

	var authData AuthMsg
	c.BindJSON(&authData)
	err := challenge.VerifyTimestamp(authData.Challenge)
	if err != nil {
		fail(c, "verifyTimestamp of Challenge failed", err)
		return
	}
	addr := common.HexToAddress(authData.Address)

	sigBytes, err := common3.HexToBytes(authData.Signature)
	if err != nil {
		fail(c, "parsing signature failed", err)
		return
	}
	challengeBytes, err := common3.HexToBytes(authData.Challenge)
	if err != nil {
		fail(c, "parsing signature failed", err)
		return
	}
	verified := utils.VerifySig(addr, sigBytes, challengeBytes)
	if !verified {
		fail(c, "verifyChallenge failed", err)
		return
	}

	c.JSON(200, gin.H{
		"ok": "ok",
	})
}
