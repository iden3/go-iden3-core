package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/services/centrauthsrv"
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

	var authData centrauthsrv.AuthMsg
	c.BindJSON(&authData)

	err := centrauthsrv.Auth(authData)
	if err != nil {
		fail(c, "auth failed", err)
		return
	}

	c.JSON(200, gin.H{
		"authenticated": true,
	})
}
