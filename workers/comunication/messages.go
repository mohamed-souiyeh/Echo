package comunication

type ClientID = int32

type ClientMessage struct {
	ClientId       ClientID
	ClientUsername string
	Msg            string
}

type RoomMessage struct {
	SenderID       ClientID
	SenderUsername string
	Msg            string
}
