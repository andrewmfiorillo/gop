package main

import (
	"log"
	"os"

	gop "github.com/andrewmfiorillo/gop"
)

func main() {
	app := gop.MakeApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
