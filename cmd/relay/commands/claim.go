package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/iden3/go-iden3/cmd/genericserver"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var ClaimCommands = []cli.Command{
	{
		Name:  "claim",
		Usage: "claim add",
		Subcommands: []cli.Command{{
			Name:   "add",
			Usage:  "claim add",
			Action: cmdAddClaim,
		}},
	},
	{
		Name:  "claims",
		Usage: "claims import from file",
		Subcommands: []cli.Command{{
			Name:   "fromfile",
			Usage:  "import claims from file",
			Action: cmdAddClaimsFromFile,
		}},
	},
}

func cmdAddClaim(c *cli.Context) error {
	if err := genericserver.MustRead(c); err != nil {
		return err
	}

	ks, acc := genericserver.LoadKeyStore()
	client := genericserver.LoadWeb3(ks, &acc)
	storage := genericserver.LoadStorage()
	mt := genericserver.LoadMerkele(storage)

	rootService := genericserver.LoadRootsService(client)
	claimService := genericserver.LoadClaimService(mt, rootService, ks, acc)

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

	err := claimService.AddDirectClaim(*claim)
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

func cmdAddClaimsFromFile(c *cli.Context) error {
	if err := genericserver.MustRead(c); err != nil {
		return err
	}
	// read config
	filepath := c.Args().Get(0)

	ks, acc := genericserver.LoadKeyStore()
	client := genericserver.LoadWeb3(ks, &acc)
	storage := genericserver.LoadStorage()
	mt := genericserver.LoadMerkele(storage)

	rootService := genericserver.LoadRootsService(client)

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
	rootService.SetRoot(*mt.RootKey())
	fmt.Println("merkletree root: " + mt.RootKey().Hex())

	return nil
}
