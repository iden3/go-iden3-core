package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/middleware/iden-assert-auth"
	"github.com/iden3/go-iden3/services/discoverysrv"
	"github.com/iden3/go-iden3/services/nameresolversrv"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
)

func handleGetHello(c *gin.Context) {
	user := auth.GetUser(c)
	c.JSON(200, gin.H{
		"idAddr":  common3.HexEncode(user.IdAddr[:]),
		"ethName": user.EthName,
		"text":    "Hello World.",
	})
}

func main() {
	nonceDb := core.NewNonceDb()
	domain := "test.eth"

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	nameResolverService, err := nameresolversrv.New("/tmp/go-iden3/names.json")
	if err != nil {
		log.Fatal(err)
	}
	discoveryservice, err := discoverysrv.New("/tmp/go-iden3/identitites.json")
	if err != nil {
		log.Fatal(err)
	}
	signedpacketservice := signedpacketsrv.New(discoveryservice, nameResolverService)
	authapi, err := auth.AddAuthMiddleware(&r.RouterGroup, domain, nonceDb, []byte("password"),
		signedpacketservice)
	if err != nil {
		log.Fatal(err)
	}

	authapi.GET("/hello", handleGetHello)

	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}
