package workers

import (
	"context"
	"echo/workers/comunication"
	"errors"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/log"
)

type ComunicationHub struct {
	ClientReqChan <-chan ClientHubReq

	Rooms map[RoomID]chan<- HubRoomReq
}

func NewComunicationHub(clientReqChan <-chan ClientHubReq) *ComunicationHub {
	return &ComunicationHub{
		ClientReqChan: clientReqChan,
		Rooms:         make(map[RoomID]chan<- HubRoomReq),
	}
}

// TODO make the hub more resilient (by handling errors well like rooms not found in the map or errors in sending or recieving from a channel) then implement the actual request response structure
func (hub *ComunicationHub) Run(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Hub PANICKED: %v", r)
			log.Errorf("Stack trace:\n%s", string(debug.Stack()))
		}
		log.Info("Hub defer function finished.")
	}()

	log.Info("Central Communication Hub started.")

	for {
		select {
		case request, ok := <-hub.ClientReqChan:
			if ok {
				hub.handleClientCommand(ctx, request)
			}
		}
	}
}

// TODO handle the client commands
func (hub *ComunicationHub) handleClientCommand(ctx context.Context, req ClientHubReq) {
	roomChan, ok := hub.Rooms[req.RoomId]

	if !ok {
		err := hub.createRoom(ctx, req.RoomId)
		if err != nil {
			// handle failed room creation gracefully
			req.ResponseChanne <- ClientHubRes{
				Code: RoomCreationFailed,
				Msg:  "oops failed to create the room",
				err:  errors.New("room creation failed"),
			}
			return
		}

		roomChan = hub.Rooms[req.RoomId]
	}

	roomReq := hub.craftRoomRequest(req)

	select {
	case roomChan <- roomReq:
		roomRes, ok := <-roomReq.ResponseChan
		if !ok {
			// handle the error
			req.ResponseChanne <- ClientHubRes{
				Code: ChannelFailure,
				Msg:  "oops the room failed to respond",
				err:  errors.New("room failed to respond"),
			}
			return
		}
		clientRes := hub.craftClientResponse(roomRes)

		req.ResponseChanne <- clientRes
	case <-time.After(1 * time.Second): // Timeout for sending
		log.Warnf("Hub: Timeout sending command %q to broadcaster for room %d", req.ReqType.string(), req.RoomId)
		req.ResponseChanne <- ClientHubRes{
			Code: RequestTimeout,
			Msg:  "oops it took too long",
			err:  errors.New("timeout"),
		}
	case <-ctx.Done():
		log.Warnf("Hub: Shutdown requested while sending command %q to broadcaster for room %d", req.ReqType.string(), req.RoomId)
		// Send error back to client if possible
		req.ResponseChanne <- ClientHubRes{
			Code: ServerShuttingDown,
			Msg:  "oops The server looks like is shutting down",
			err:  errors.New("the server is shutting down"),
		}
	}
}

func (hub *ComunicationHub) craftClientResponse(res HubRoomRes) ClientHubRes {
	return ClientHubRes{
		Code:     res.Code,
		RoomChan: res.roomChan,
		Msg:      res.Msg,
		err:      res.err,
	}
}

func (hub *ComunicationHub) craftRoomRequest(req ClientHubReq) HubRoomReq {
	return HubRoomReq{
		ClientId:     req.ClientId,
		Type:         req.ReqType,
		clientChan:   req.ClientChan,
		Msg:          "hello room",
		ResponseChan: make(chan HubRoomRes, 1),
	}
}

func (hub *ComunicationHub) createRoom(ctx context.Context, ID RoomID) error {
	roomReqChan := make(chan HubRoomReq, 1024)
	hub.Rooms[ID] = roomReqChan

	go func() {
		msgsChannel := make(chan comunication.ClientMessage, 4096)
		clients := make(map[ClientID]chan<- comunication.RoomMessage)
		for {
			select {
			case request, ok := <-roomReqChan:
				if ok && request.Type == JOIN {
					clients[request.ClientId] = request.clientChan
					response := HubRoomRes{
						Code:     AllGood,
						roomChan: msgsChannel,
						Msg:      "welcome to the room",
						err:      nil,
					}
					request.ResponseChan <- response
					} else if ok && request.Type == LEAVE {
						delete(clients, request.ClientId)
						response := HubRoomRes{
							Code:     AllGood,
							roomChan: nil,
							Msg:      "see you soon ^^",
							err:      nil,
						}
						request.ResponseChan <- response
				}
			case msg, ok := <-msgsChannel:
				if ok {
					// log.Debugf("got msgs from client with id '%d': %q", msg.ClientId, msg.Msg)
					for _, clientChan := range clients {
						clientChan <- comunication.RoomMessage{
							SenderID:       msg.ClientId,
							SenderUsername: msg.ClientUsername,
							Msg:            msg.Msg,
						}
					}
				}
			}
		}
	}()

	return nil
}
