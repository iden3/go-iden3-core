package commands

import (
	"os"
	"os/signal"

	cfg "github.com/iden3/go-iden3/cmd/backupserver/config"
	"github.com/iden3/go-iden3/cmd/backupserver/endpoint"
	"github.com/urfave/cli"
)

func Start(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	// storage := cfg.LoadStorage()
	mongoservice := cfg.LoadMongoService()
	backupservice := cfg.LoadBackupService(mongoservice)

	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				os.Exit(1)
			}
		}
	}()

	endpoint.Serve(backupservice)
	// storage.Close()

	return nil
}
