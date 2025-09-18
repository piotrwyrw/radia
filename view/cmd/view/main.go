package main

import (
	"github.com/piotrwyrw/otherproj/internal/ui"
	"github.com/piotrwyrw/radia/radia/radia"
)

func main() {
	radia.Initialize()
	ui.CreateUI()
}
