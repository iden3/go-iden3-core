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

func handleInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":        "ok",
		"powdifficulty": backupservice.GetPoWDifficulty(),
	})
}

func handleSave(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)

	var saveBackupMsg backupsrv.BackupData
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
		fail(c, "error on RecoverAll", err)
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
		fail(c, "error on RecoverByTimestamp", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverByType(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	dataType := c.Param("type")

	data, err := backupservice.RecoverByType(idaddr, dataType) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverByType", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverByTimestampAndType(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr := common.HexToAddress(idaddrhex)
	timestampstr := c.Param("timestamp")
	timestamp, err := strconv.ParseUint(timestampstr, 10, 64)
	if err != nil {
		fail(c, "error on parsing timestamp from url", err)
		return
	}
	dataType := c.Param("type")

	data, err := backupservice.RecoverByTimestampAndType(idaddr, timestamp, dataType) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverByTimestampAndType", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}
