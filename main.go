package main

import (
	"echo/app"

	_ "embed"

	"github.com/charmbracelet/log"
)

//go:embed banner.txt
var banner string


func main() {
	log.SetLevel(log.DebugLevel)

	a := app.NewApp()

	a.Start()
}
