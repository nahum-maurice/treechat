package main

import (
	"fmt"
	"net"
	"strings"
)

type Command struct {
	Text string
}

const (
	CLogin   string = "login"
	CSignUp  string = "signup"
	CJoin    string = "join"
	CMembers string = "members"
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
	case CSignUp:
		if len(args) <= 2 {
			conn.Write([]byte("\n[System] ::: Please provide your username and password.\n\n"))
			return
		}
		username, password := strings.TrimSpace(args[1]), strings.TrimSpace(args[2])
		HandleSignUp(username, password, conn)
	case CLogin:
		if len(args) <= 2 {
			conn.Write([]byte("\n[System] ::: Please provide your username and password.\n\n"))
			return
		}
		username, password := strings.TrimSpace(args[1]), strings.TrimSpace(args[2])
		HandleLogin(username, password, conn)
	case CJoin:
		if len(args) <= 1 {
			conn.Write([]byte("Please provide a room name.\n\n"))
			return
		}
		room_name := strings.TrimSpace(args[1])
		HandleJoin(room_name, conn)
	case CMembers:
		HandleMembers(conn)
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
		conn.Write([]byte("\n[System] ::: Unknown command.\n\n"))
	}
}

func HandleLogin(username string, password string, conn net.Conn) {
	is_user := IsUser(username)
	if !is_user {
		conn.Write([]byte("\n[System] ::: Sorry, there is no such user. To create a new account, please use '/signup <username> <password>'.\n\n"))
		return
	}
	is_verified := VerifyUser(username, password)

	if !is_verified {
		conn.Write([]byte("\n[System] ::: Sorry, username/password combination don't match.\n\n"))
		return
	}
	the_user := NewUser(username, password, conn.RemoteAddr().String(), true)

	Users = append(Users, the_user)
	conn.Write([]byte("\n[System] ::: Welcome back " + username + "!\n\n"))

}

func HandleSignUp(username string, password string, conn net.Conn) {
	is_user := IsUser(username)
	if is_user {
		conn.Write([]byte("\n[System] ::: Sorry, this username is already taken. Please, try with another one.\n\n"))
		return
	}
	the_user := NewUser(username, password, conn.RemoteAddr().String(), true)

	Users = append(Users, the_user)
	conn.Write([]byte("\n[System] ::: Welcome " + username + "!\n\n"))
}

func HandleRooms(conn net.Conn) {
	if (len(Rooms)) == 0 {
		conn.Write([]byte("\n[System] ::: There are no rooms. To create a new room, please type '/newroom <room_name>'.\n\n"))
		return
	}
	rooms_string := "\n[System] ::: The available rooms are: \n"
	for _, room := range Rooms {
		rooms_string += "........       - " + room.Name + "\n"
	}
	conn.Write([]byte(rooms_string + "\n"))
}

func HandleNewRoom(room_name string, conn net.Conn) {
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		conn.Write([]byte("\n[System] ::: Sorry, You are not authenticated. Please, login or signup.\n\n"))
		return
	}
	if user.IsAuthenticated {
		// TODO: Check if the room names doesn't already exist
		a_room, _ := GetRoomByName(room_name)
		if a_room != nil {
			conn.Write([]byte("\n[System] ::: There is already a room with that name.\n\n"))
			return
		}

		newRoom := NewRoom(room_name, user.Username)
		Rooms = append(Rooms, newRoom)
		conn.Write([]byte("\n[System] ::: Success, now you can join by typing '/join " + newRoom.Name + "'.\n\n"))
	}
}

func HandleOnline(conn net.Conn) {
	// This function responds "You don't belong to any room" if
	// the user didn't join any room. Otherwise, it responds
	// with the actual members that are present in the current
	// room the user is.
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		// log.Fatal(err)
		conn.Write([]byte("\n[System] ::: Sorry, You are not authenticated. Please, login or signup.\n\n"))
		return
	}

	if user.CurrentRoom == nil {
		conn.Write([]byte("\n[System] ::: You don't belong to any room.\n\n"))
		return
	}

	response := "\n[System] ::: The online users are:\n"
	for _, room := range Rooms {
		if room.Name == user.CurrentRoom.Name {
			for _, usr := range room.Online {
				response += "........     - " + usr + "\n"
			}
		}
	}
	conn.Write([]byte(response + "\n"))
}

func HandleMembers(conn net.Conn) {
	// This function responds "You are not a member of any room" if 
	// the user didn't join any room. Otherwise, it responds
	// with the actual members of a room
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		conn.Write([]byte("\n[System] ::: Sorry, You are not authenticated. Please, login or signup.\n\n"))
		return
	}

	if user.CurrentRoom == nil {
		conn.Write([]byte("\n[System] ::: You are not a member of any room.\n\n"))
		return
	}

	response := fmt.Sprintf("\n[System] ::: The members of the room '%v' are:\n", user.CurrentRoom.Name)
	for _, usr := range user.CurrentRoom.Members {
		response += "........     - " + usr + "\n"
	}
	conn.Write([]byte(response + "\n"))	
}

func HandleJoin(room_name string, conn net.Conn) {
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		conn.Write([]byte("\n[System] ::: Sorry, You are not authenticated. Please, login or signup.\n\n"))
		return
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

			// notify everyone in the room
			for _, connection := range elem.Connections {
				if user.ConnectionAddress != connection.RemoteAddr().String() {
					message := fmt.Sprintf("\n[System] ::: %s joined the room '%v'.\n\n\n", user.Username, elem.Name)
					connection.Write([]byte(message))
					break
				}
			}

			if !contains(elem.Members, user.Username) {
				elem.Members = append(elem.Members, user.Username)
			}

			user.CurrentRoom = elem
			conn.Write([]byte("\n[System] ::: You joined the room '" + elem.Name + "'.\n\n"))

			break
		}
	}
	if count == 0 {
		conn.Write([]byte("\n[System] ::: There is no such room.\n\n"))
		return
	}
}

func HandleQuit(conn net.Conn) {
	// Remove the user from the Online channel
	usr, _ := GetUserByConnectionAddress(conn.RemoteAddr().String())
	for _, elem := range Users {
		if elem == usr {
			// Remove the user from the Online field of its last room
			usr.CurrentRoom.RemoveOnline(usr.Username)
			usr.CurrentRoom.RemoveConnection(conn)
			// notify every one from the room it was
			for _, connection := range usr.CurrentRoom.Connections {
				message := fmt.Sprintf("\n[System] ::: %s left the room.\n\n\n", usr.Username)
				connection.Write([]byte(message))
			}
			// We don't need this (the next)
			// Remove the user from the list of users.
			// Users = append(Users[:i], Users[i+1:]...)
		}
	}
	conn.Write([]byte("\n[System] ::: Bye!\n\n"))
	conn.Close()
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
