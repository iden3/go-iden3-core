package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/backupserver/config"
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/services/backupsrv"
	log "github.com/sirupsen/logrus"
)

var backupservice backupsrv.Service

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

func serveServiceApi() *http.Server {
	api := gin.Default()
	api.Use(cors.Default())
	serviceapi := api.Group("/api/unstable")
	serviceapi.GET("/", handleInfo)

	// BACKUP SERVICE
	serviceapi.POST("/register", handleRegister)
	serviceapi.POST("/backup/upload", handleBackupUpload)
	serviceapi.POST("/backup/download", handleBackupDownload)

	// SYNCHRONIZATION SERVICE
	//serviceapi.POST("/:id/save", handleSave) // TODO: Redo
	serviceapi.POST("/folder/:id/recover", handleRecover)
	//TODO get with specific version
	serviceapi.POST("/folder/:id/recover/version/:version", handleRecoverSinceVersion)
	serviceapi.POST("/folder/:id/recover/type/:type", handleRecoverByType)
	serviceapi.POST("/folder/:id/recover/version/:version/type/:type", handleRecoverSinceVersionByType)

	serviceapisrv := &http.Server{Addr: config.C.Server.ServiceApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return serviceapisrv
}

func Serve(bs backupsrv.Service) {
	backupservice = bs

	stopch := make(chan interface{})

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				stopch <- nil
			}
		}
	}()

	serviceapisrv := serveServiceApi()
	// wait until shutdown signal
	<-stopch
	log.Info("Shutdown Server ...")
	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}
}
