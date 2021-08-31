package models

type User interface {
	GetId() string
	GetName() string
	GetUserName() string
}

type UserRepository interface {
	AddUser(user User)
	AddFriend(id string, friendID string, userID string)
	CreateUser(id string, name string, username string, password string)
	RemoveUser(user User)
	FindUserById(ID string) User
	GetAllUsers() []User
	GetAllFriends(userID string) []User
	IsAlreadyFriend(userID string, friendID string) bool
}