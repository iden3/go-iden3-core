package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/rootsrv"
)

var claimservice claimsrv.Service
var rootservice rootsrv.Service

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Origin", "*")
		c.Writer.Header().Add("X-Requested-With", "*")
		c.Next()
	}
}

func Serve(rs rootsrv.Service, cs claimsrv.Service) {

	claimservice = cs
	rootservice = rs

	r := gin.Default()
	r.Use(corsMiddleware())
	r.GET("/root", handleGetRoot)
	r.POST("/claim/:idaddr", handlePostClaim)
	r.GET("/claim/:idaddr/root", handleGetIDRoot)
	r.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)
	r.Run(config.C.Server.ServiceApi)
}
