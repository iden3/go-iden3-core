package endpoint

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/services/backupsrv"
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

	var saveBackupMsg backupsrv.SaveBackupMsg
	c.BindJSON(&saveBackupMsg)

	timestamp, err := backupservice.Save(idaddr, saveBackupMsg)
	if err != nil {
		fail(c, "error on SaveBackup", err)
		return
	}

	c.JSON(200, gin.H{
		"status":    "stored correctly",
		"timestamp": timestamp,
	})
}

func handleRecover(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)

	data, err := backupservice.RecoverAll(idaddr) // + proofs for authentication
	if err != nil {
		fail(c, "error on SaveBackup", err)
		return
	}
	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverByTimestamp(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	timestampstr := c.Param("timestamp")
	timestamp, err := strconv.ParseUint(timestampstr, 10, 64)
	if err != nil {
		fail(c, "error on parsing timestamp from url", err)
		return
	}

	data, err := backupservice.RecoverByTimestamp(idaddr, timestamp) // + proofs for authentication
	if err != nil {
		fail(c, "error on SaveBackup", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}
