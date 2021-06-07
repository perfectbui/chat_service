package main

import "github.com/perfectbui/chat/job"

type WsServer struct {
	rooms map[*Room]bool
}

func newWsServer() *WsServer {
	return &WsServer{
		rooms: make(map[*Room]bool),
	}
}

func (WsServer *WsServer) createRoom(name string) *Room {
	room := newRoom(name)
	WsServer.rooms[room] = true
	go job.CreateRoomJob(name)
	go room.runRoom()
	return room
}

func (wsServer *WsServer) findRoomByName(name string) *Room {
	for room := range wsServer.rooms {
		if room.getName() == name {
			return room
		}
	}
	return nil
}
