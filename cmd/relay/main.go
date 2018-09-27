package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3/cmd/relay/commands"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()

	flags := []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{},
			Usage:   "start the server",
			Action:  commands.Start,
		},
		{
			Name:    "contract",
			Aliases: []string{},
			Usage:   "operate with contracts",
			Subcommands: []cli.Command{
				{
					Name:   "info",
					Usage:  "show information about contracts",
					Action: commands.ContractInfo,
				},
				{
					Name:   "deploy",
					Usage:  "deploy contract",
					Action: commands.ContractDeploy,
				},
			},
		},
	}

	app.Flags = flags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
