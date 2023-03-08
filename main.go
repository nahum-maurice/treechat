package main

import (
	"fmt"
	"log"
)

func main() {
	server := NewServer(":3000")

	// Listening to incomming messages with a coroutine and sending messages
	// to the resepctive room using a custom goroutine.
	go func() {
		for msg := range server.MessageChannel {
			for _, connection := range msg.Room.Connections {
				message := fmt.Sprintf("\n[%v] ::: [%v] %v.\n\n", msg.Room.Name, msg.From, msg.Payload)
				connection.Write([]byte(message))
			}
		}
	}()

	log.Fatal(server.Start())
}
