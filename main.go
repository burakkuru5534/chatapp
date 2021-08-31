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

	wsServer := NewWebsocketServer(&repository.RoomRepository{Db: db}, userRepository)
	go wsServer.Run()

	api := &API{UserRepository: userRepository}

	http.HandleFunc("/ws", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	}))

	http.HandleFunc("/api/register", api.RegisterUser)
	http.HandleFunc("/api/login", api.HandleLogin)
	http.HandleFunc("/api/friend/add", api.AddFriend)
	http.HandleFunc("/api/deletealluser", api.DeleteAllUsers)
	http.HandleFunc("/api/friends", api.ShowUserFriends)

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}