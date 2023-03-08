package main

import "fmt"

// This will hold all the users.
var Users []*User

type User struct {
	Username          string
	Password          string
	ConnectionAddress string // The current connection address
	CurrentRoom       *Room  // A user should be in 1 room at time
	Rooms             []string
	IsAuthenticated   bool
}

func NewUser(username string, password string, address string, isAuth bool) *User {
	new_user := User{
		Username:          username,
		Password:          password,
		ConnectionAddress: address,
		IsAuthenticated:   isAuth,
	}
	return &new_user
}

func IsUser(username string) bool {
	for _, usr := range Users {
		if usr.Username == username {
			return true
		}
	}
	return false
}

func VerifyUser(username string, password string) bool {
	is_user := IsUser(username)
	if !is_user {
		return false
	} else {
		user, _ := GetUserByUsername(username)
		return user.Password == password
	}
}

func (user *User) String() string {
	return fmt.Sprintf("Username: %s ConnectionAddress: %s", user.Username, user.ConnectionAddress)
}

func GetUserByConnectionAddress(connectionAddress string) (*User, error) {
	for _, elem := range Users {
		if elem.ConnectionAddress == connectionAddress {
			return elem, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}

func GetUserByUsername(username string) (*User, error) {
	for _, elem := range Users {
		if elem.Username == username {
			return elem, nil
		}
	}
	return nil, fmt.Errorf("User not found")
}

func (u *User) GetCurrentRoom() (*Room, error) {
	if u.CurrentRoom == nil {
		return nil, fmt.Errorf("User not in a room")
	}
	return u.CurrentRoom, nil
}
