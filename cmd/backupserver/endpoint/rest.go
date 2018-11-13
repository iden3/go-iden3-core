package endpoint

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
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

func handleSave(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)

	err := backupservice.SaveBackup(idaddr) // + proofs for authentication
	if err != nil {
		fail(c, "error on SaveBackup", err)
		return
	}

	c.JSON(200, gin.H{
		"backup": "dev",
	})
}

func handleRecover(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)

	err := backupservice.RecoverBackup(idaddr) // + proofs for authentication
	if err != nil {
		fail(c, "error on SaveBackup", err)
		return
	}

	c.JSON(200, gin.H{
		"backup": "dev",
	})
}
