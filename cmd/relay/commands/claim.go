package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	cfg "github.com/iden3/go-iden3/cmd/relay/config"
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
	if err := cfg.MustRead(c); err != nil {
		return err
	}

	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)
	storage := cfg.LoadStorage()
	mt := cfg.LoadMerkele(storage)

	rootservice := cfg.LoadRootsService(client)
	claimservice := cfg.LoadClaimService(mt, rootservice, ks, acc)

	indexData := c.Args().Get(0)
	outData := c.Args().Get(1)

	claim := core.NewGenericClaim("iden3.io", "generic", []byte(indexData), []byte(outData))
	fmt.Println("clam: " + common3.BytesToHex(claim.Bytes()))

	err := claimservice.AddDirectClaim(claim)
	if err != nil {
		return err
	}
	fmt.Print("root updated: " + mt.Root().Hex())

	mp, err := mt.GenerateProof(claim.Hi())
	if err != nil {
		return err
	}
	fmt.Print("merkleproof: " + common3.BytesToHex(mp))

	return nil
}

func cmdAddClaimsFromFile(c *cli.Context) error {
	if err := cfg.MustRead(c); err != nil {
		return err
	}
	// read config
	filepath := c.Args().Get(0)

	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)
	storage := cfg.LoadStorage()
	mt := cfg.LoadMerkele(storage)

	rootservice := cfg.LoadRootsService(client)

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

		claim := core.NewGenericClaim("iden3.io", "generic", []byte(line[0]), []byte(line[1]))
		fmt.Println("clam: " + common3.BytesToHex(claim.Bytes()) + "\n")

		// add claim to merkletree, without updating the root, that will be done on the end of the loop (csv file)
		err = mt.Add(claim)
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

		claim := core.NewGenericClaim("iden3.io", "generic", []byte(line[0]), []byte(line[1]))
		fmt.Println("clam: " + common3.BytesToHex(claim.Bytes()))

		// the proofs better generate them once all claims are added
		mp, err := mt.GenerateProof(claim.Hi())
		if err != nil {
			return err
		}
		fmt.Println("merkleproof: " + common3.BytesToHex(mp) + "\n")
	}
	// update the root in the smart contract
	rootservice.SetRoot(mt.Root())
	fmt.Println("merkletree root: " + mt.Root().Hex())

	return nil
}
