package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	Address        string
	Listener       net.Listener
	QuitChannel    chan struct{}
	Rooms          []*Room
	MessageChannel chan Message
}

func NewServer(address string) *Server {
	return &Server{
		Address:        address,
		QuitChannel:    make(chan struct{}),
		MessageChannel: make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Treechat server up and running on address: %v\n\n", s.Address)

	// When terminating, please, close the listener
	defer listener.Close()
	s.Listener = listener

	// Start listening to incoming connections
	go s.acceptLoop()

	// Wait until the quit channel is invoked
	<-s.QuitChannel
	close(s.MessageChannel)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
			// We don't return cause we don't want to close the loop
			// since other incomming requests may pass in
			continue
		}

		fmt.Println("New connection to the server:", conn.RemoteAddr())

		// First display for the user upon connection.
		conn.Write([]byte(
			"\n" +
				"[System] ::: Welcome folk! You are connected to Treechat.\n" +
				"........     Here are the commands you can use:\n" +
				"........       /signup  <username> <password>   --> Sign up to Treechat.\n" +
				"........       /login <username> <password>     --> Login to Treechat.\n" +
				"........       /rooms                           --> Show the list of all rooms.\n" +
				"........       /join <room>                     --> Join a room.\n" +
				"........       /newroom <room>                  --> Create a new room.\n" +
				"........       /online                          --> Show the people that are online in a room.\n" +
				"........       /quit                            --> Quit the server.\n\n"))

		// Start the reading loop. This loop is responsible to listen
		// to incomming messages, make sure it recognizes and respond
		// to commands and publish messages to the latest room the
		// user added itself.
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()

	// Read the message sent by the user. The buffer has a size of
	// 256 bytes. This limit is imposed by the server in order to
	// prevent undesirable behaviours.
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Lost connection: %v\n", conn.RemoteAddr().String())
			fmt.Println(err)
			return
		}

		msg := string(buf[:n])
		msg = strings.Trim(msg, "\r\n")

		// Conduct a test to see whether the msg is a command. If it
		// is a command, it should be handled as such. The way it's
		// handled is defined by the commands packages.
		if msg[0] == '/' {
			command := NewCommand(msg)
			command.Handle(conn)
		} else {
			curr_user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
			// If the user is not authenticated, we should not allow them to send
			// messages to any room. This could be changed later, but for now, that's
			// the best possible behaviour I think for the sake of simplicity.
			if err != nil {
				conn.Write([]byte("\n[System] ::: You are not authenticated. Please log in by typing '/auth <username> <password>'.\n\n"))
			} else {
				// If the user didn't join any room, we should neither allow them to send
				// messages to anywhere. This could change later but for the sake of
				// simplicity, let's keep it like that now
				destination_room, err := curr_user.GetCurrentRoom()
				if err != nil {
					msg := fmt.Sprintf(
						"\n" +
							"[System] ::: Please, join a room. To see all rooms, type '/rooms'.\n" +
							"........     To create a new room, type '/newroom <name>'.\n" +
							"........     To join an existing room, type '/join <room>'.\n\n")
					conn.Write([]byte(msg))
				} else {
					// When a message is sent, we then need to send the message to the
					// room where the sender is currently present and therefore broadcast
					// the message to the other members of the room.
					formated_message := NewMessage(curr_user, msg, destination_room)
					s.MessageChannel <- formated_message
				}
			}
		}
	}
}
