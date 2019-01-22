package main

import (
	"log"
	"os"

	"github.com/andrewmfiorillo/p/pgo"
)

func main() {
	app := pgo.MakeApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
