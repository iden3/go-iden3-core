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
//	idHex := c.Param("id")
//	id := common.HexToAddress(idHex)
//
//	var saveBackupMsg backupsrv.BackupData
//	err := c.BindJSON(&saveBackupMsg)
// 	if err != nil {
//      	Fail(c, "json parsing error", err)
//      	return
// 	}
//
//	version, err := backupservice.Save(id, saveBackupMsg)
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
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		fail(c, "error on id", err)
		return
	}

	data, err := backupservice.RecoverAll(id) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverAll", err)
		return
	}
	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverSinceVersion(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		fail(c, "error on id", err)
		return
	}
	versionstr := c.Param("version")
	version, err := strconv.ParseUint(versionstr, 10, 64)
	if err != nil {
		fail(c, "error on parsing version from url", err)
		return
	}

	data, err := backupservice.RecoverSinceVersion(id, version) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverSinceVersion", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverByType(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		fail(c, "error on id", err)
		return
	}
	dataType := c.Param("type")

	data, err := backupservice.RecoverByType(id, dataType) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverByType", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}

func handleRecoverSinceVersionByType(c *gin.Context) {
	idHex := c.Param("id")
	id, err := core.IDFromString(idHex)
	if err != nil {
		fail(c, "error on id", err)
		return
	}
	versionstr := c.Param("version")
	version, err := strconv.ParseUint(versionstr, 10, 64)
	if err != nil {
		fail(c, "error on parsing version from url", err)
		return
	}
	dataType := c.Param("type")

	data, err := backupservice.RecoverSinceVersionByType(id, version, dataType) // + proofs for authentication
	if err != nil {
		fail(c, "error on RecoverSinceVersionByType", err)
		return
	}

	c.JSON(200, gin.H{
		"backups": data,
	})
}
