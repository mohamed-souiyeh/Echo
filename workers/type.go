package workers

import (
	"echo/workers/comunication"
	"fmt"
)

type RequestType int

const (
	JOIN RequestType = iota
	LEAVE
)

func (rt RequestType) string() string {
	switch rt {
	case JOIN:
		return "JOIN"
	case LEAVE:
		return "LEAVE"
	default:
		return "Unknown"
	}
}

type ClientHubReq struct {
	ClientId       ClientID
	ReqType        RequestType
	RoomId         RoomID
	ClientChan     chan<- comunication.RoomMessage
	Msg            string
	ResponseChanne chan ClientHubRes
}

type ClientHubRes struct {
	Code     StatusCode
	RoomChan chan<- comunication.ClientMessage
	Msg      string
	err      error
}

type RoomID = int32
type ClientID = int32

type HubRoomReq struct {
	ClientId     ClientID
	Type         RequestType
	clientChan   chan<- comunication.RoomMessage
	Msg          string
	ResponseChan chan HubRoomRes
}

type HubRoomRes struct {
	Code     StatusCode
	roomChan chan<- comunication.ClientMessage
	Msg      string
	err      error
}

// StatusCode represents a response status code.
type StatusCode int

// String returns the text representation of the HTTPStatusCode.
func (s StatusCode) String() string {
	switch s {
	case RequestTimeout:
		return "Request Timeout"
	case ChannelFailure:
		return "Channel Failure"
	case RoomCreationFailed:
		return "Room Creation Failed"
	default:
		return fmt.Sprintf("Unknown Status Code (%d)", s)
	}
}

// Constants for HTTP Status Codes
const (
	AllGood StatusCode = iota
	RequestTimeout
	ChannelFailure
	RoomCreationFailed
	ServerShuttingDown
)
