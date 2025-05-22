package messages

import (
	"echo/workers"
	"echo/workers/comunication"
)

type ChatMsg struct {
	SenderID       int32
	SenderUsername string
	Msg            string
}

type JoinedChatRoomsMsg struct {
	Msg      string
	RoomId   workers.RoomID
	RoomChan chan<- comunication.ClientMessage
}

type AllowMsgSend struct {
	
}
