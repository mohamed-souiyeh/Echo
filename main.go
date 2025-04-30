package main

import (
	"echo/app"

	_ "embed"

	"github.com/charmbracelet/log"
)

//go:embed banner.txt
var banner string

type clientHubReq struct {
	ClientHubResChan chan clientHubRes
}

type clientHubRes struct {
}

type roomHubNotif struct {
}

// type clientRoomNotif struct {

// }

func main() {
	log.SetLevel(log.DebugLevel)

	a := app.NewApp()

	a.Start()
}
