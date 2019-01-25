package endpoint

//import (
//	"encoding/json"
//
//	jwt "github.com/dgrijalva/jwt-go"
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/websocket"
//	"github.com/iden3/go-iden3/cmd/centrauth/config"
//	"github.com/iden3/go-iden3/services/centrauthsrv"
//	log "github.com/sirupsen/logrus"
//)
//
//func fail(c *gin.Context, msg string, err error) {
//	if err != nil {
//		log.WithError(err).Error(msg)
//	} else {
//		log.Error(msg)
//	}
//	c.JSON(400, gin.H{
//		"error": msg,
//	})
//	return
//}

//func handleAuth(c *gin.Context) {
//
//	var authData centrauthsrv.AuthMsg
//	c.BindJSON(&authData)
//
//	err := centrauthsrv.Auth(authData)
//	if err != nil {
//		fail(c, "auth failed", err)
//		return
//	}
//
//	// generate new JWT
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id_address": authData.Address,
//	})
//	tokenString, err := token.SignedString([]byte(config.C.Jwtsecret))
//	if err != nil {
//		fail(c, "error generating token", err)
//		return
//	}
//
//	// send the token through websocket
//	var authTokenMsg centrauthsrv.AuthTokenMsg
//	authTokenMsg.Success = true
//	authTokenMsg.Token = tokenString
//	jResp, err := json.Marshal(authTokenMsg)
//	if err != nil {
//		fail(c, "error generating authToken json", err)
//		return
//	}
//	wsClients[authData.Challenge].conn.WriteMessage(websocket.TextMessage, jResp)
//	delete(wsClients, authData.Challenge)
//
//	c.JSON(200, gin.H{
//		"authenticated": true,
//	})
//}
