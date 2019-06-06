package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	encryption "github.com/iden3/go-iden3/crypto/public-key-encryption"
	"io/ioutil"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Tool to generate asymmetric keys, encrypt and decrypt\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s [opts] gen/encrypt/decrypt\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	inPath := flag.String("in", "", "Input file")
	outPath := flag.String("out", "", "Output file")
	pkHex := flag.String("pk", "", "Public Key")
	kpHex := flag.String("kp", "", "Key Pair (contains Secret Key)")
	flag.Parse()
	if len(flag.Args()) == 0 {
		usage()
		return
	}

	cmd := flag.Args()[0]

	switch cmd {
	case "gen":
		// Create key pair
		kp := encryption.GenKP()
		fmt.Printf("Public Key: %s\n", hex.EncodeToString(kp.PublicKey.Bytes[:]))
		fmt.Printf("Key Pair: %s%s\n",
			hex.EncodeToString(kp.PublicKey.Bytes[:]), hex.EncodeToString(kp.SecretKey.Bytes[:]))
	case "encrypt":
		if *pkHex == "" || *inPath == "" || *outPath == "" {
			usage()
			return
		}
		println(inPath)
		println(outPath)
		pk, err := encryption.ImportBoxPublicKey(*pkHex)
		if err != nil {
			fmt.Println("Error importing public key:", err)
			return
		}
		data, err := ioutil.ReadFile(*inPath)
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}
		encData := encryption.Encrypt(pk, data)
		if err := ioutil.WriteFile(*outPath, encData, 0600); err != nil {
			fmt.Println("Error writing output file:", err)
			return
		}

	case "decrypt":
		if *kpHex == "" || *inPath == "" || *outPath == "" {
			usage()
			return
		}
		kp, err := encryption.ImportBoxKP(*kpHex)
		if err != nil {
			fmt.Println("Error importing key pair:", err)
			return
		}
		encData, err := ioutil.ReadFile(*inPath)
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}
		// Decrypt
		data, err := encryption.Decrypt(kp, encData)
		if err != nil {
			fmt.Println("Error decrypting input file:", err)
			return
		}
		if err := ioutil.WriteFile(*outPath, data, 0600); err != nil {
			fmt.Println("Error writing output file:", err)
			return
		}
	default:
		usage()
		return
	}
}
