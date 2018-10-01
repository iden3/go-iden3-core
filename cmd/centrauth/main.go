package main

import (
	"os"

	"github.com/iden3/go-iden3/cmd/centrauth/commands"
	log "github.com/sirupsen/logrus"
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
	}

	app.Flags = flags

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
