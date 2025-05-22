package commands

import (
	"context"
	"echo/tui/messages"
	"echo/workers"
	"echo/workers/comunication"

	tea "github.com/charmbracelet/bubbletea"
)

func ChatMsgsListenerCmd(ctx context.Context, listeningChan <-chan comunication.RoomMessage, clientId int32) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-ctx.Done():
			// simply closing this listner as part of the ressources cleaning in a client
			return nil
		case msg, ok := <-listeningChan:
			if ok {
				return messages.ChatMsg{
					SenderID:       msg.SenderID,
					SenderUsername: msg.SenderUsername,
					Msg:            msg.Msg,
				}
			}
		}
		return nil
	}
}

func JoinLobbyChatRoomCmd(hubReqChan chan<- workers.ClientHubReq, req workers.ClientHubReq) tea.Cmd {
	return func() tea.Msg {
		hubReqChan <- req
		res := <-req.ResponseChanne

		if res.Code == workers.AllGood {
			return messages.JoinedChatRoomsMsg{
				Msg:      res.Msg,
				RoomId:   req.RoomId,
				RoomChan: res.RoomChan,
			}
		}
		//TODO handle the error returned by the cetral hub appropriatly
		return nil
	}
}

func LeaveLobbyChatRoomCmd(hubReqChan chan<- workers.ClientHubReq, req workers.ClientHubReq) tea.Cmd {
	return func() tea.Msg {
		hubReqChan <- req
		res := <-req.ResponseChanne
		
		if res.Code == workers.AllGood {
			return messages.JoinedChatRoomsMsg{
				Msg:      res.Msg,
				RoomId:   req.RoomId,
				RoomChan: res.RoomChan,
			}
		}
		//TODO handle the error returned by the cetral hub appropriatly
		return nil
	}
}
