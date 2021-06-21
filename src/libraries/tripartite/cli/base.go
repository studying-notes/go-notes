//
// Created by Rustle Karl on 2021.01.01 20:16.
//
package main

import (
	"fmt"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		fmt.Println("BOOM!")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
