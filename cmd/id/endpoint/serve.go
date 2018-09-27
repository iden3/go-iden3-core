package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/services/claimsrv"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Origin", "*")
		c.Writer.Header().Add("X-Requested-With", "*")
		c.Next()
	}
}

func Serve(cs *claimsrv.Service) {
	/*
		claimservice = cs

		r := gin.Default()
		r.Use(corsMiddleware())
		r.POST("/vinculateid", handleVinculateID)
		r.GET("/identities/resolv/:nameid", handleAssignNameClaimResolv)
		r.Run(":" + config.C.Server.Port)
	*/
}
