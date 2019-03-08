package genericserver

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli"
)

// Claim
func CmdAddClaim(c *cli.Context) error {
	if err := MustRead(c); err != nil {
		return err
	}

	ks, acc := LoadKeyStore()
	client := LoadWeb3(ks, &acc)
	storage := LoadStorage()
	mt := LoadMerkele(storage)

	rootservice := LoadRootsService(client)
	claimservice := LoadClaimService(mt, rootservice, ks, acc)

	indexData := c.Args().Get(0)
	outData := c.Args().Get(1)

	var indexSlot [400 / 8]byte
	var dataSlot [496 / 8]byte
	if len(indexData) != len(indexSlot) || len(outData) != len(dataSlot) {
		return fmt.Errorf(
			"Length of indexSlot and dataSlot must be %v and %v respectively",
			len(indexSlot), len(dataSlot))
	}
	copy(indexSlot[:], indexData)
	copy(dataSlot[:], outData)
	claim := core.NewClaimBasic(indexSlot, dataSlot)
	fmt.Println("clam: " + common3.HexEncode(claim.Entry().Bytes()))

	err := claimservice.AddDirectClaim(*claim)
	if err != nil {
		return err
	}
	fmt.Print("root updated: " + mt.RootKey().Hex())

	mp, err := mt.GenerateProof(claim.Entry().HIndex(), nil)
	if err != nil {
		return err
	}
	fmt.Print("merkleproof: " + common3.HexEncode(mp.Bytes()))

	return nil
}

func CmdAddClaimsFromFile(c *cli.Context) error {
	if err := MustRead(c); err != nil {
		return err
	}
	// read config
	filepath := c.Args().Get(0)

	ks, acc := LoadKeyStore()
	client := LoadWeb3(ks, &acc)
	storage := LoadStorage()
	mt := LoadMerkele(storage)

	rootservice := LoadRootsService(client)

	fmt.Print("\n---\nimporting claims\n---\n\n")
	// csv file will have the following structure: indexData, noindexData
	csvFile, _ := os.Open(filepath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var err error
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		fmt.Println("importing claim with index: " + line[0] + ", outside index: " + line[1])

		var indexSlot [400 / 8]byte
		var dataSlot [496 / 8]byte
		if len(line[0]) != len(indexSlot) || len(line[1]) != len(dataSlot) {
			return fmt.Errorf(
				"Length of indexSlot and dataSlot must be %v and %v respectively",
				len(indexSlot), len(dataSlot))
		}
		copy(indexSlot[:], line[0])
		copy(dataSlot[:], line[1])
		claim := core.NewClaimBasic(indexSlot, dataSlot)
		// claim := core.NewGenericClaim("iden3.io", "generic", []byte(line[0]), []byte(line[1]))
		fmt.Println("clam: " + common3.HexEncode(claim.Entry().Bytes()) + "\n")

		// add claim to merkletree, without updating the root, that will be done on the end of the loop (csv file)
		err = mt.Add(claim.Entry())
		if err != nil {
			return err
		}
	}
	fmt.Print("\n---\ngenerating proofs\n---\n\n")
	// now, let's generate the proofs
	csvFile, _ = os.Open(filepath)
	reader = csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		fmt.Println("generating merkleproof of claim with index: " + line[0] + ", outside index: " + line[1])

		var indexSlot [400 / 8]byte
		var dataSlot [496 / 8]byte
		if len(line[0]) != len(indexSlot) || len(line[1]) != len(dataSlot) {
			return fmt.Errorf(
				"Length of indexSlot and dataSlot must be %v and %v respectively",
				len(indexSlot), len(dataSlot))
		}
		copy(indexSlot[:], line[0])
		copy(dataSlot[:], line[1])
		claim := core.NewClaimBasic(indexSlot, dataSlot)
		fmt.Println("clam: " + common3.HexEncode(claim.Entry().Bytes()))

		// the proofs better generate them once all claims are added
		mp, err := mt.GenerateProof(claim.Entry().HIndex(), nil)
		if err != nil {
			return err
		}
		fmt.Println("merkleproof: " + common3.HexEncode(mp.Bytes()) + "\n")
	}
	// update the root in the smart contract
	rootservice.SetRoot(*mt.RootKey())
	fmt.Println("merkletree root: " + mt.RootKey().Hex())

	return nil
}

// DB
func CmdDbRawDump(c *cli.Context) error {

	if err := MustRead(c); err != nil {
		return err
	}
	storage := LoadStorage()
	ldb := (storage.(*db.LevelDbStorage)).LevelDB()
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(hex.EncodeToString(iter.Key()), " ", hex.EncodeToString(iter.Value()))
	}
	iter.Release()
	return nil
}

func CmdDbIPFSexport(c *cli.Context) error {
	if err := MustRead(c); err != nil {
		return err
	}
	storage := LoadStorage()
	ldb := (storage.(*db.LevelDbStorage)).LevelDB()
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		sh := shell.NewShell("localhost:5001") // ipfs daemon IP:Port
		cid, err := sh.Add(bytes.NewReader(iter.Value()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
		fmt.Println("value of key "+common3.HexEncode(iter.Key())+" added, ipfs hash: ", cid)
	}
	iter.Release()
	return nil
}
