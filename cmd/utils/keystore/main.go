package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/iden3/go-iden3-crypto/babyjub"
	babykeystore "github.com/iden3/go-iden3/keystore"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Tool to create a key store\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s [opts] new/import/dump\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	path := flag.String("path", "", "keystore file path")
	skHex := flag.String("sk", "", "private key in hex to import")
	pass := flag.String("pass", "", "keystore password")
	light := flag.Bool("light", false, "Use light key derivation parameters")
	flag.Parse()
	if *path == "" || *pass == "" || len(flag.Args()) == 0 {
		usage()
		return
	}
	cmd := flag.Args()[0]

	params := babykeystore.StandardKeyStoreParams
	if *light {
		params = babykeystore.LightKeyStoreParams
	}

	storage := babykeystore.NewFileStorage(*path)
	ks, err := babykeystore.NewKeyStore(storage, params)
	if err != nil {
		panic(err)
	}
	switch cmd {
	case "new":
		pk, err := ks.NewKey([]byte(*pass))
		if err != nil {
			panic(err)
		}
		fmt.Println("Public key:", hex.EncodeToString(pk[:]))
	case "import":
		if *skHex == "" {
			usage()
			return
		}
		var sk babyjub.PrivateKey
		if _, err := hex.Decode(sk[:], []byte(*skHex)); err != nil {
			panic(err)
		}
		pk, err := ks.ImportKey(sk, []byte(*pass))
		if err != nil {
			panic(err)
		}
		fmt.Println("Public key:", hex.EncodeToString(pk[:]))
	case "dump":
		for i, pk := range ks.Keys() {
			fmt.Printf("%02d - pk: %v\n", i, hex.EncodeToString(pk[:]))
			sk, err := ks.ExportKey(&pk, []byte(*pass))
			if err != nil {
				panic(err)
			}
			fmt.Printf("     sk: %v\n", hex.EncodeToString(sk[:]))
		}
	default:
		usage()
		return
	}
}
