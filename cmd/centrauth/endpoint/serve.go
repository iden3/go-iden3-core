package endpoint

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/centrauth/config"
)

func Serve() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/auth", handleAuth)
	r.GET("/ws/:id", handleWs)
	r.Run(config.C.Server.ServiceApi)
}
