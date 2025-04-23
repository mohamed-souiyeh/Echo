package main

import (
	"echo/app"

	_ "embed"
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
	a := app.NewApp()

	a.Start()
}
