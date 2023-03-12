package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/nahum-maurice/treechat/utils"
)

type Command struct {
	Text      string
	Formatter *utils.Formatter
	Logger    *utils.Logger
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
	return &Command{
		Text: text,
		Formatter: utils.NewFormatter("System"),
	}
}

// TODO : Change return type from any to the actual constrained
// return types
func (c *Command) Handle(conn net.Conn) {
	args := strings.Split(c.Text, " ")
	key := strings.TrimSpace(args[0])[1:]

	switch key {
	case CSignUp:
		if len(args) <= 2 {
			msg := "Please provide your username and password."
			conn.Write([]byte(c.Formatter.MessageCLI(msg, "System", "")))
			return
		}
		username, password := strings.TrimSpace(args[1]), strings.TrimSpace(args[2])
		HandleSignUp(c.Formatter, username, password, conn)
	case CLogin:
		if len(args) <= 2 {
			msg := "Please provide your username and password."
			conn.Write([]byte(c.Formatter.MessageCLI(msg, "System", "")))
			return
		}
		username, password := strings.TrimSpace(args[1]), strings.TrimSpace(args[2])
		HandleLogin(c.Formatter, username, password, conn)
	case CJoin:
		if len(args) <= 1 {
			msg := "Please provide a room name."
			conn.Write([]byte(c.Formatter.MessageCLI(msg, "System", "")))
			return
		}
		room_name := strings.TrimSpace(args[1])
		HandleJoin(c.Formatter, room_name, conn)
	case CMembers:
		HandleMembers(c.Formatter, conn)
	case CNewRoom:
		room_name := strings.TrimSpace(args[1])
		HandleNewRoom(c.Formatter, room_name, conn)
	case COnline:
		HandleOnline(c.Formatter, conn)
	case CQuit:
		HandleQuit(c.Formatter, conn)
	case CRooms:
		HandleRooms(c.Formatter, conn)
	default:
		conn.Write([]byte(c.Formatter.MessageCLI("Unknown command.", "System", "")))
	}
}

func HandleLogin(f *utils.Formatter, username string, password string, conn net.Conn) {
	is_user := IsUser(username)
	if !is_user {
		msg := "Sorry, there is no such user. To create a new account, please use '/signup <username> <password>'."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}
	is_verified := VerifyUser(username, password)

	if !is_verified {
		msg := "Sorry, username/password combination don't match."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}
	the_user := NewUser(username, password, conn.RemoteAddr().String(), true)

	Users = append(Users, the_user)
	msg := "Welcome back " + username + "!"
	conn.Write([]byte(f.MessageCLI(msg, "System", "")))

}

func HandleSignUp(f *utils.Formatter, username string, password string, conn net.Conn) {
	is_user := IsUser(username)
	if is_user {
		msg := "Sorry, this username is already taken. Please, try with another one."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}
	the_user := NewUser(username, password, conn.RemoteAddr().String(), true)

	Users = append(Users, the_user)
	msg := "Welcome " + username + "!"
	conn.Write([]byte(f.MessageCLI(msg, "System", "")))
}

func HandleRooms(f *utils.Formatter, conn net.Conn) {
	if (len(Rooms)) == 0 {
		msg := "There are no rooms. To create a new room, please type '/newroom <room_name>'."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}
	roomsString := f.MessagePrimaryCLI("The available rooms are:", "System", "")
	for _, room := range Rooms {
		roomsString +=  f.MessageSecondaryCLI(" # " + room.Name, "System", "")
	}
	conn.Write([]byte(roomsString + "\n"))
}

func HandleNewRoom(f *utils.Formatter, room_name string, conn net.Conn) {
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		msg := "Sorry, You are not authenticated. Please, login or signup."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}
	if user.IsAuthenticated {
		// TODO: Check if the room names doesn't already exist
		a_room, _ := GetRoomByName(room_name)
		if a_room != nil {
			msg := "There is already a room with that name."
			conn.Write([]byte(f.MessageCLI(msg, "System", "")))
			return
		}

		newRoom := NewRoom(room_name, user.Username)
		Rooms = append(Rooms, newRoom)
		msg := "Success, now you can join by typing '/join " + newRoom.Name + "'"
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
	}
}

func HandleOnline(f *utils.Formatter, conn net.Conn) {
	// This function responds "You don't belong to any room" if
	// the user didn't join any room. Otherwise, it responds
	// with the actual members that are present in the current
	// room the user is.
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		// log.Fatal(err)
		msg := "Sorry, You are not authenticated. Please, login or signup."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}

	if user.CurrentRoom == nil {
		msg := "You don't belong to any room."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}

	response := f.MessagePrimaryCLI("The online users are:", "System", "")
	for _, room := range Rooms {
		if room.Name == user.CurrentRoom.Name {
			for _, usr := range room.Online {
				response += f.MessageSecondaryCLI(" @ " + usr, "System", "")
			}
		}
	}
	conn.Write([]byte(response + "\n"))
}

func HandleMembers(f *utils.Formatter, conn net.Conn) {
	// This function responds "You are not a member of any room" if
	// the user didn't join any room. Otherwise, it responds
	// with the actual members of a room
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		msg := "Sorry, You are not authenticated. Please, login or signup."
		conn.Write([]byte(f.MessageCLI("@ " + msg, "System", "")))
		return
	}

	if user.CurrentRoom == nil {
		msg := "You are not a member of any room."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
		return
	}

	response := f.MessagePrimaryCLI(
		fmt.Sprintf("The members of the room '%v' are:", user.CurrentRoom.Name),
		"System",
		"",
	)
	for _, usr := range user.CurrentRoom.Members {
		response += f.MessageSecondaryCLI(" @ " + usr, "System", "")
	}
	conn.Write([]byte(response + "\n"))
}

func HandleJoin(f *utils.Formatter, room_name string, conn net.Conn) {
	user, err := GetUserByConnectionAddress(conn.RemoteAddr().String())
	if err != nil {
		msg := "Sorry, You are not authenticated. Please, login or signup."
		conn.Write([]byte(f.MessageCLI(msg, "System", "")))
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
					message := fmt.Sprintf("%s joined the room '%v'.\n", user.Username, elem.Name)
					connection.Write([]byte(f.MessageCLI(message, "System", "")))
					break
				}
			}

			if !contains(elem.Members, user.Username) {
				elem.Members = append(elem.Members, user.Username)
			}

			user.CurrentRoom = elem
			msg := "You joined the room '" + elem.Name + "'."
			conn.Write([]byte(f.MessageCLI(msg, "System", "")))

			break
		}
	}
	if count == 0 {
		conn.Write([]byte(f.MessageCLI("There is no such room.", "System", "")))
		return
	}
}

func HandleQuit(f *utils.Formatter, conn net.Conn) {
	// Remove the user from the Online channel
	usr, _ := GetUserByConnectionAddress(conn.RemoteAddr().String())
	for _, elem := range Users {
		if elem == usr {
			// Remove the user from the Online field of its last room
			usr.CurrentRoom.RemoveOnline(usr.Username)
			usr.CurrentRoom.RemoveConnection(conn)
			// notify every one from the room it was
			for _, connection := range usr.CurrentRoom.Connections {
				message := fmt.Sprintf("%s left the room.\n", usr.Username)
				connection.Write([]byte(f.MessageCLI(message, "System", "")))
			}
			// We don't need this (the next)
			// Remove the user from the list of users.
			// Users = append(Users[:i], Users[i+1:]...)
		}
	}
	conn.Write([]byte(f.MessageCLI("Bye!", "System", "")))
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
