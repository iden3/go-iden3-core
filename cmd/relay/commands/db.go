package commands

import (
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/urfave/cli"
)

var DbCommands = []cli.Command{{
	Name:  "db",
	Usage: "operate with database",
	Subcommands: []cli.Command{
		{
			Name:   "rawdump",
			Usage:  "dump database raw key values",
			Action: genericserver.CmdDbRawDump,
		},
		{
			Name:   "rawimport",
			Usage:  "import database raw from dumped key values",
			Action: genericserver.CmdDbRawImport,
		},
		{
			Name:   "ipfsexport",
			Usage:  "export database values to ipfs",
			Action: genericserver.CmdDbIPFSexport,
		},
	},
}}
