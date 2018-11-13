package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3/cmd/backupserver/commands"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "backupserveri3"
	app.Version = "0.1.0-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{},
			Usage:   "start the server",
			Action:  commands.Start,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
