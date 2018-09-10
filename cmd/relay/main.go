package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/cmd/relay/endpoint"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/web3"
)

func main() {
	config.MustRead(".", "config")

	storage, err := merkletree.NewLevelDbStorage("./db/")
	if err != nil {
		log.Fatal(err)
	}
	mt, err := merkletree.New(storage, 140)
	if err != nil {
		log.Fatal(err)
	}

	defer storage.Close()

	fmt.Println("mt.Root: " + mt.Root().Hex())

	// Ethereum
	err = web3srv.Open(config.C.Geth.URL, config.C.Server.PrivK)
	if err != nil {
		color.Red(err.Error())
		log.Fatal(err)
	}
	balance, err := web3srv.GetBalance(web3srv.Address)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"balance": balance.String(),
		"address": web3srv.Address.Hex(),
	}).Info("contract info")
	if balance.Int64() < 3000000 {
		color.Red("Not enough funds in the relay address")
		os.Exit(0)
	}
	endpoint.Serve(mt)
}
