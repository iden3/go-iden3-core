package main

import (
	"fmt"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3/cmd/id/config"
	"github.com/iden3/go-iden3/cmd/id/endpoint"
	mtlib "github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/web3"
)

func main() {
	config.MustRead(".", "config")

	// MerkleTree leveldb
	storage, err := mtlib.NewLevelDbStorage("./db/")
	if err != nil {
		log.Fatal(err)
	}
	mt, err := mtlib.New(storage, 140)
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

	endpoint.Serve(mt)
}
