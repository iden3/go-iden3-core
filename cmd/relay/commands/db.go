package commands

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/db"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli"
)

var DbCommands = []cli.Command{{
	Name:  "db",
	Usage: "operate with database",
	Subcommands: []cli.Command{
		{
			Name:   "rawdump",
			Usage:  "dump database raw key values",
			Action: cmdDbRawDump,
		},
		{
			Name:   "ipfsexport",
			Usage:  "export database values to ipfs",
			Action: cmdDbIPFSexport,
		},
	},
}}

func cmdDbRawDump(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}
	storage := cfg.LoadStorage()
	ldb := (storage.(*db.LevelDbStorage)).LevelDB()
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Println(hex.EncodeToString(iter.Key()), " ", hex.EncodeToString(iter.Value()))
	}
	iter.Release()
	return nil
}

func cmdDbIPFSexport(c *cli.Context) error {
	if err := cfg.MustRead(c); err != nil {
		return err
	}
	storage := cfg.LoadStorage()
	ldb := (storage.(*db.LevelDbStorage)).LevelDB()
	iter := ldb.NewIterator(nil, nil)
	for iter.Next() {
		sh := shell.NewShell("localhost:5001") // ipfs daemon IP:Port
		cid, err := sh.Add(bytes.NewReader(iter.Value()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
		fmt.Println("value of key "+common3.BytesToHex(iter.Key())+" added, ipfs hash: ", cid)
	}
	iter.Release()
	return nil
}
