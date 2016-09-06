package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "nixie"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search PATTERN",
			Action: func(c *cli.Context) error {
				Search(c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
