package endpoint

import (
	"strconv"

	"github.com/gin-gonic/gin"
	//"github.com/iden3/go-iden3/services/backupsrv"
	"github.com/iden3/go-iden3/core"
)

func handleInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":        "ok",
		"powdifficulty": backupservice.GetPoWDifficulty(),
	})
}

// TODO: Redo
//func handleSave(c *gin.Context) {
//	idaddrhex := c.Param("idaddr")
//	idaddr := common.HexToAddress(idaddrhex)
//
//	var saveBackupMsg backupsrv.BackupData
//	c.BindJSON(&saveBackupMsg)
//
//	version, err := backupservice.Save(idaddr, saveBackupMsg)
//	if err != nil {
//		fail(c, "error on SaveBackup", err)
//		return
//	}
//
//	c.JSON(200, gin.H{
//		"status":  "stored correctly",
//		"version": version,
//	})
//}

func handleRecover(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idAddr, err := core.IDFromString(idaddrhex)
	if err != nil {
		fail(c, "error on idAddr", err)
		return
	}

	data, err := backupservice.RecoverAll(idAddr) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverAll", err)
		return
	}
	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverSinceVersion(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr, err := core.IDFromString(idaddrhex)
	if err != nil {
		fail(c, "error on idAddr", err)
		return
	}
	versionstr := c.Param("version")
	version, err := strconv.ParseUint(versionstr, 10, 64)
	if err != nil {
		fail(c, "error on parsing version from url", err)
		return
	}

	data, err := backupservice.RecoverSinceVersion(idaddr, version) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverSinceVersion", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverByType(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr, err := core.IDFromString(idaddrhex)
	if err != nil {
		fail(c, "error on idAddr", err)
		return
	}
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

func handleRecoverSinceVersionByType(c *gin.Context) {
	idaddrhex := c.Param("idaddr")
	idaddr, err := core.IDFromString(idaddrhex)
	if err != nil {
		fail(c, "error on idAddr", err)
		return
	}
	versionstr := c.Param("version")
	version, err := strconv.ParseUint(versionstr, 10, 64)
	if err != nil {
		fail(c, "error on parsing version from url", err)
		return
	}
	dataType := c.Param("type")

	data, err := backupservice.RecoverSinceVersionByType(idaddr, version, dataType) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverSinceVersionByType", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}
