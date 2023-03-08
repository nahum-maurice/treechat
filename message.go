package main

import "time"

type Message struct {
	ID        string
	From      string
	Payload   string
	Timestamp time.Time
	Room      *Room // where the message should be sent to
}

func NewMessage(username string, payload string, r *Room) Message {
	new_message := Message{
		From:      username,
		Payload:   payload,
		Timestamp: time.Now(),
		Room:      r,
	}
	return new_message
}
