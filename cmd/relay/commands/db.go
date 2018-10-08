package commands

import (
	"encoding/hex"
	"fmt"

	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/db"
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
