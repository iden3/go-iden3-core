package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/backupserver/config"
	"github.com/iden3/go-iden3/services/backupsrv"
	log "github.com/sirupsen/logrus"
)

var backupservice backupsrv.Service

func serveServiceApi() *http.Server {
	api := gin.Default()
	api.Use(cors.Default())
	serviceapi := api.Group("/api/v0.1")
	serviceapi.GET("/", handleInfo)
	serviceapi.POST("/:idaddr/save", handleSave)
	serviceapi.POST("/:idaddr/recover", handleRecover)
	//TODO get with specific version
	serviceapi.POST("/:idaddr/recover/version/:version", handleRecoverSinceVersion)
	serviceapi.POST("/:idaddr/recover/type/:type", handleRecoverByType)
	serviceapi.POST("/:idaddr/recover/version/:version/type/:type", handleRecoverSinceVersionByType)

	serviceapisrv := &http.Server{Addr: config.C.Server.ServiceApi, Handler: api}
	go func() {
		log.Info("API server at ", config.C.Server.ServiceApi)
		if err := serviceapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
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
