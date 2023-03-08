package main

import "time"

type Message struct {
	From      *User
	Payload   string
	Timestamp time.Time
	Room      *Room // where the message should be sent to
}

func NewMessage(user *User, payload string, r *Room) Message {
	newMessage := Message{
		From:      user,
		Payload:   payload,
		Timestamp: time.Now(),
		Room:      r,
	}
	return newMessage
}
