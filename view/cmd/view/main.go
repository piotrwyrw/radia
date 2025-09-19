package main

import (
	"log"

	"github.com/piotrwyrw/otherproj/internal/ui"
	"github.com/piotrwyrw/radia/radia/radia"
)

func main() {
	radia.Initialize()
	err := ui.CreateUI()
	if err != nil {
		log.Fatal(err)
	}
}
