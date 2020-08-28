package main

import (
	"flag"
	"log"
)

var n string

func main() {
	flag.Parse()
	kitCmd := flag.NewFlagSet("kit", flag.ExitOnError)
	kitCmd.StringVar(&n, "n", "value", "help")

	args := flag.Args()
	switch args[0] {
	case "kit":
		_ = kitCmd.Parse(args[1:])
	}
	log.Println("n =", n)
}
