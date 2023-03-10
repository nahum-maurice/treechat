package main

import (
	"fmt"
	"net"

	"github.com/nahum-maurice/treechat/utils"
)

// Will contain the existing rooms.
// TODO: Make a init function that will be getting all the
// existing rooms and set them to memory upon start.
var Rooms []*Room

type Room struct {
	Name        string
	creator     string
	Members     []string
	Online      []string
	Connections []net.Conn
	QuitChannel chan struct{}
	Logger      *utils.Logger
}

func NewRoom(name string, creator string) *Room {
	new := Room{
		Name:        name,
		creator:     creator,
		Members:     []string{creator},
		QuitChannel: make(chan struct{}),
		Logger:      utils.NewLogger("System"),
	}
	newLog := fmt.Sprintf("New room created: %s", new.Name)
	new.Logger.Info(newLog)
	return &new
}

func GetRoomByName(name string) (*Room, error) {
	for _, room := range Rooms {
		if room.Name == name {
			return room, nil
		}
	}
	return nil, fmt.Errorf("Room not found")
}

func (r *Room) String() string {
	return fmt.Sprintf("Room: %s", r.Name)
}

func (r *Room) AddMember(member string) {
	r.Members = append(r.Members, member)
}

func (r *Room) AddOnline(member string) {
	r.Online = append(r.Online, member)
}

func (r *Room) RemoveOnline(member string) {
	for i, v := range r.Online {
		if v == member {
			r.Online = append(r.Online[:i], r.Online[i+1:]...)
		}
	}
}

func (r *Room) RemoveConnection(member net.Conn) {
	for i, v := range r.Connections {
		if v == member {
			r.Connections = append(r.Connections[:i], r.Connections[i+1:]...)
		}
	}
}
