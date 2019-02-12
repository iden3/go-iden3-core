package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/services/backupsrv"
)

func handleRegister(c *gin.Context) {
	var user backupsrv.User
	c.BindJSON(&user)

	err := backupservice.Register(user)
	if err != nil {
		fail(c, "error on Register", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

type backupMsg struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Backup   string `json:"backup" binding:"required"`
}

func handleBackupUpload(c *gin.Context) {
	var backupMsg backupMsg
	c.BindJSON(&backupMsg)

	user := backupsrv.User{
		Username: backupMsg.Username,
		Password: backupMsg.Password,
	}
	backupPacket := backupsrv.BackupPacket{
		Username: backupMsg.Username,
		Backup:   backupMsg.Backup,
	}

	err := backupservice.BackupUpload(user, backupPacket)
	if err != nil {
		fail(c, "error on BackupUpload", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func handleBackupDownload(c *gin.Context) {
	var user backupsrv.User
	c.BindJSON(&user)

	backupPacket, err := backupservice.BackupDownload(user)
	if err != nil {
		fail(c, "error on BackupDownload", err)
		return
	}

	c.JSON(http.StatusOK, backupPacket)
}
