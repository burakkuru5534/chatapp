package main

import (
	"context"
	"example.com/m/auth"
	"example.com/m/config"
	"example.com/m/repository"
	"flag"
	"log"
	"net/http"

)

var addr = flag.String("addr", ":8080", "http server address")
var ctx = context.Background()

func main() {
	flag.Parse()

	config.CreateRedisClient()
	db := config.InitDB()
	defer db.Close()

	userRepository := &repository.UserRepository{Db: db}

	http.HandleFunc("/ws", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ServeWs( w, r)
	}))

	api := &API{UserRepository: userRepository}


	//login and register apis -- users can access this apis without token
	http.HandleFunc("/api/register", api.RegisterUser)
	http.HandleFunc("/api/login", api.HandleLogin)
	http.HandleFunc("/api/deletealluser", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.DeleteAllUsers(w,r)
	}))
	//user apis -- users can access this apis with token
	http.HandleFunc("/api/friend/add", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.AddFriend(w,r)
	}))

	http.HandleFunc("/api/friends", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.ShowUserFriends(w,r)
	}))
	http.HandleFunc("/api/user/sended/messages", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.ShowUsersPostedMessagesToFriend(w,r)
	}))
	http.HandleFunc("/api/user/received/messages", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		api.ShowUsersReceivedMessagesFromFriend(w,r)
	}))
	//end of the user apis -- users can access this apis with token


	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}