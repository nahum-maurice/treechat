package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Command struct {
	Text string
}

const (
	CAuth    string = "auth"
	CJoin    string = "join"
	CNewRoom string = "newroom"
	COnline  string = "online"
	CQuit    string = "quit"
	CRooms   string = "rooms"
)

func NewCommand(text string) *Command {
	return &Command{Text: text}
}

// TODO : Change return type from any to the actual constrained
// return types
func (c *Command) Handle(conn net.Conn) {
	args := strings.Split(c.Text, " ")
	key := strings.TrimSpace(args[0])[1:]

	switch key {
	case CAuth:
		if len(args) <= 2 {
			conn.Write([]byte("\n[System] ::: Please provide your username and password.\n\n"))
		} else {
			username, password := strings.TrimSpace(args[1]), strings.TrimSpace(args[2])
			HandleAuth(username, password, conn)
		}
	case CJoin:
		room_name := strings.TrimSpace(args[1])
		HandleJoin(room_name, conn)
	case CNewRoom:
		room_name := strings.TrimSpace(args[1])
		HandleNewRoom(room_name, conn)
	case COnline:
		HandleOnline(conn)
	case CQuit:
		HandleQuit(conn)
	case CRooms:
		HandleRooms(conn)
	default:
		conn.Write([]byte("Unknown command.\n"))
	}
}

func HandleAuth(username string, password string, conn net.Conn) {
	// TODO
	// Check whether the user exist in our database. If not,
	// create a new user with these credentials.
	//
	// Hash the password before storing it

	the_user := NewUser(username, password, conn.RemoteAddr().String(), true)

	// Add the newly created user to the list of the other users
	// Hint: Since all the authenticated users will be in that
	//       slice, we can just look inside of it to check if the
	//       user is authenticated further.
	Users = append(Users, the_user)
	conn.Write([]byte("\n[System] ::: Welcome " + username + "! You are authenticated!\n\n"))
}

func HandleRooms(conn net.Conn) {
	if (len(Rooms)) == 0 {
		conn.Write([]byte("\n[System] ::: There are no rooms. To create a new room, please type '/newroom <room_name>'.\n\n"))
	} else {
		rooms_string := "\n[System] ::: The available rooms are: \n"
		for _, room := range Rooms {
			rooms_string += "........       -" + room.Name + ".\n"
		}
		conn.Write([]byte(rooms_string + "\n"))
	}
}

func HandleNewRoom(room_name string, conn net.Conn) {
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	if user.IsAuthenticated {
		// TODO: Check if the room names doesn't already exist

		new_room := NewRoom(room_name, user.Username)
		Rooms = append(Rooms, new_room)
		conn.Write([]byte("\n[System] ::: Success, now you can join by typing '/join " + new_room.Name + "'.\n\n"))
	}
}

func HandleOnline(conn net.Conn) {
	// This function responds "You don't belong to any room" if
	// the user didn't join any room. Otherwise, it responds
	// with the actual members that are present in the current
	// room the user is.
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	if user.CurrentRoom == nil {
		conn.Write([]byte("\nYou don't belong to any room.\n\n"))
	}

	response := "\n[System] ::: The online users are:\n"
	for _, room := range Rooms {
		if room.Name == user.CurrentRoom.Name {
			for _, usr := range room.Online {
				response += "........     -" + usr + ".\n"
			}
		}
	}
	conn.Write([]byte(response + "\n"))
}

func HandleJoin(room_name string, conn net.Conn) {
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for _, elem := range Rooms {
		// Check if the room name already exists. If yes, add the
		// user to the list of members if he's not already there
		// and add the user to the list of online people
		if elem.Name == room_name {
			count += 1
			elem.Online = append(elem.Online, user.Username)
			elem.Connections = append(elem.Connections, conn)

			if !contains(elem.Members, user.Username) {
				elem.Members = append(elem.Members, user.Username)
			}

			user.CurrentRoom = elem
			conn.Write([]byte("\n[System] ::: You joined the room '" + elem.Name + "'.\n\n"))

			break
		}
	}
	if count == 0 {
		conn.Write([]byte("\n[System] ::: There is no such room."))
	}
}

func HandleQuit(conn net.Conn) {
	// Remove the user from the Online channel
	usr, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	for i, elem := range Users {
		if elem == usr {
			// Remove the user from the Online field of its last room
			for _, a_room := range Rooms {
				a_room.RemoveOnline(usr.Username)
				a_room.RemoveConnection(conn)
			}
			// Remove the user from the list of users.
			Users = append(Users[:i], Users[i+1:]...)
			fmt.Printf("[System] ::: User %v left.\n", usr.Username)
		}
		// TODO: take necessary actions to clean after him
	}
	conn.Write([]byte("Bye!\n\n"))
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
