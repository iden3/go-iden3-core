package genericserver

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewServiceAPI(prefix string) (*gin.Engine, *gin.RouterGroup) {
	api := gin.Default()
	api.Use(cors.Default())

	serviceapi := api.Group(prefix)
	serviceapi.GET("/root", HandleGetRoot)

	return api, serviceapi
}

func NewAdminAPI(prefix string, stopch chan interface{}) (*gin.Engine, *gin.RouterGroup) {
	api := gin.Default()
	api.Use(cors.Default())
	adminapi := api.Group("/api/unstable")

	adminapi.POST("/stop", func(c *gin.Context) {
		// yeah, use curl -X POST http://<adminserver>/stop
		c.String(http.StatusOK, "got it, shutdowning server")
		stopch <- nil
	})

	adminapi.GET("/info", HandleInfo)
	adminapi.GET("/rawdump", HandleRawDump)
	adminapi.POST("/rawimport", HandleRawImport)
	adminapi.GET("/claimsdump", HandleClaimsDump)
	return api, adminapi
}

// https://golang.org/src/net/http/server.go?s=86961:87002#L3255
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func ListenAndServe(srv *http.Server, name string) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Infof("%s API is ready at %v", name, addr)
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
