package commands

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/cmd/relay/endpoint"
)

func Start(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)
	mt := cfg.LoadMerkele()
	defer mt.Storage().Close()

	rootservice := cfg.LoadRootsService(client)
	claimservice := cfg.LoadClaimService(mt, rootservice, ks, acc)

	// Check for founds
	balance, err := client.BalanceAt(acc.Address)
	if err != nil {
		log.Panic(err)
	}
	log.WithFields(log.Fields{
		"balance": balance.String(),
		"address": acc.Address.Hex(),
	}).Info("Account balance retrieved")
	if balance.Int64() < 3000000 {
		log.Panic("Not enough funds in the relay address")
	}
	rootservice.Start()

	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				rootservice.StopAndJoin()
				os.Exit(1)
			}
		}
	}()
	endpoint.Serve(rootservice, claimservice)
	return nil
}
