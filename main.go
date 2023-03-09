package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
}

func main() {
	PORT := os.Getenv("PORT")
	address := fmt.Sprintf("0.0.0.0:%v", PORT)
	server := NewServer(address)

	// Listening to incomming messages with a coroutine and sending messages
	// to the resepctive room using a custom goroutine.
	go func() {
		for msg := range server.MessageChannel {
			for _, connection := range msg.Room.Connections {
				if msg.From.ConnectionAddress != connection.RemoteAddr().String() {
					message := fmt.Sprintf("[%v]  %v ::: [%v] %v\n\n", msg.Room.Name, msg.Timestamp.Format("2006/01/02 15:04:05"), msg.From.Username, msg.Payload)
					connection.Write([]byte(message))
				} else {
					// There is no normal reason to do that, beside
					// keeping regular spaces in the terminals. :)
					connection.Write([]byte("\n"))
				}
			}
		}
	}()

	log.Fatal(server.Start())
}
