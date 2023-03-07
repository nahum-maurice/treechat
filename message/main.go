package message

import (
	"time"

	"github.com/nahum-maurice/treechat/room"
)

type Message struct {
	ID        string
	From      string
	Payload   string
	Timestamp time.Time
	Room     *room.Room // where the message should be sent to
}

func NewMessage(username string, payload string, room *room.Room) Message {
	new_message := Message{
		From:      username,
		Payload:   payload,
		Timestamp: time.Now(),
		Room: room,
	}
	return new_message
}
