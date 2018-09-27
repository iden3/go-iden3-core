package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/cmd/relay/endpoint"
)

func main() {

	config.MustRead(".", "config")
	ks, acc := config.LoadKeyStore()
	client := config.LoadWeb3(ks, &acc)
	mt := config.LoadMerkele()
	defer mt.Storage().Close()

	rootservice := config.LoadRootsService(client)
	claimservice := config.LoadClaimService(mt, rootservice, ks, acc)

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

	endpoint.Serve(rootservice, claimservice)
}
