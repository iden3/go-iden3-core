package commands

import (
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/urfave/cli"
)

var ClaimCommands = []cli.Command{
	{
		Name:  "claim",
		Usage: "claim add",
		Subcommands: []cli.Command{{
			Name:   "add",
			Usage:  "claim add",
			Action: genericserver.CmdAddClaim,
		}},
	},
	{
		Name:  "claims",
		Usage: "claims import from file",
		Subcommands: []cli.Command{{
			Name:   "fromfile",
			Usage:  "import claims from file",
			Action: genericserver.CmdAddClaimsFromFile,
		}},
	},
}
