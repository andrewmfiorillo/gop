package main

import (
	"log"
	"os"

	pgo "github.com/andrewmfiorillo/p"
)

func main() {
	app := pgo.MakeApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
