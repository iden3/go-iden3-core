package commands

import (
	"fmt"

	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/urfave/cli"
)

var DbCommands = []cli.Command{{
	Name:  "db",
	Usage: "operate with database",
	Subcommands: []cli.Command{
		{
			Name:   "info",
			Usage:  "show database information",
			Action: cmdDbInfo,
		},
	},
}}

func cmdDbInfo(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}
	storage := cfg.LoadStorage()
	fmt.Println(storage.Info())
	return nil
}
