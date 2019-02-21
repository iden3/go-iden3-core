package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/iden3/go-iden3/cmd/claimserver/commands"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "claimserveri3"
	app.Version = "0.1.0-alpha"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config"},
	}

	app.Commands = []cli.Command{}
	app.Commands = append(app.Commands, commands.ServerCommands...)
	app.Commands = append(app.Commands, commands.DbCommands...)
	app.Commands = append(app.Commands, commands.ClaimCommands...)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
