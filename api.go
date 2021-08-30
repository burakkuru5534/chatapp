package main

import (
	"encoding/json"
	"example.com/m/auth"
	"example.com/m/cmn"
	"example.com/m/config"
	"example.com/m/repository"
	"example.com/m/response"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/go-chi/chi"
)
type API struct {
	UserRepository *repository.UserRepository
	Auth *auth.Auth
}
type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type queueStatement struct {
	sql          string
	params       []interface{}
	isNamedExec  bool
	namedExecArg interface{}
	// db *DbHandle
}
type DbHandle struct {
	*sqlx.DB
	Logger     *Logger
	LimitOfset string
	DriverName string

	writeQueue chan queueStatement
}
type Logger struct {
	zerolog.Logger
}

func (api *API) RegisterUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var isUserNameExist bool
	isUserNameExist = true
	// Find the user in the database by username
	dbUser := api.UserRepository.FindUserByUsername(username)
	if dbUser == nil {
		isUserNameExist = false
	}

	id := uuid.New().String()
	if !isUserNameExist {
		hashedPassword, err := cmn.HashPasswd(password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}
		api.UserRepository.CreateUser(id, username,username,hashedPassword)

		api.ResponseSuccess(w, id)
		return
	}else {
		http.Error(w, "username allready taken by other users", http.StatusUnauthorized)
		return
	}

}
func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {

	var user LoginUser

	// Try to decode the JSON request to a LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the user in the database by username
	dbUser := api.UserRepository.FindUserByUsername(user.Username)
	if dbUser == nil {
		returnErrorResponse(w)
		return
	}

	// Check if the passwords match
	ok, err := auth.ComparePassword(user.Password, dbUser.Password)

	if !ok || err != nil {
		returnErrorResponse(w)
		return
	}

	// Create a JWT
	token, err := auth.CreateJWTToken(dbUser)

	if err != nil {
		returnErrorResponse(w)
		return
	}

	w.Write([]byte(token))

}
//func (api *API) LoginUser(w http.ResponseWriter, r *http.Request) {
//	var (
//		upassFromDb *string
//
//		userID        *int64
//		userName      *string
//		userFirstName *string
//		userLastName  *string
//		userEmail     *string
//	)
//
//	err := r.ParseForm()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	username := r.FormValue("username")
//	password := r.FormValue("password")
//
//	qs := "select id, username, first_name, last_name, email, userpass from usr where username = $1 or email = $1 and is_active"
//	err = api.Db.QueryRow(qs, username).Scan(&userID, &userName, &userFirstName, &userLastName, &userEmail, &upassFromDb)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//
//		} else {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//			return
//		}
//		return
//	}
//
//	if cmn.PtrToString(upassFromDb) != password {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	tc := auth.TokenClaims{
//		UserID:       *userID,
//		UserName:     cmn.PtrToString(userName),
//		UserEmail:    cmn.PtrToString(userEmail),
//		UserFullName: fmt.Sprintf("%s %s", cmn.PtrToString(userFirstName), cmn.PtrToString(userLastName)),
//		UserType:     "user",
//	}
//
//	at := auth.AuthResponse{
//		UserID:       *userID,
//		UserName:     cmn.PtrToString(userName),
//		UserEmail:    cmn.PtrToString(userEmail),
//		UserFullName: fmt.Sprintf("%s %s", cmn.PtrToString(userFirstName), cmn.PtrToString(userLastName)),
//		AccessToken:  api.Auth.NewTokenString(&tc),
//	}
//
//	api.ResponseSuccess(w, at)
//}
func (api *API) DeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	db := config.InitDB()
	defer db.Close()

	_, err:= db.Exec("delete from user")
	if err != nil {
		returnErrorResponse(w)
		return
	}


}
func (api *API) ShowUserFriends(w http.ResponseWriter, r *http.Request) {
	userID := StrToInt64(chi.URLParam(r, "id"))

	db := config.InitDB()
	defer db.Close()

	_, err := db.Exec("select * from user where id in (select friend_id from user_friend where user_id = $1)",userID)
	if err != nil {
		returnErrorResponse(w)
		return
	}


}
func (api *API) AddFriend(w http.ResponseWriter, r *http.Request) {

	userID := StrToInt64(chi.URLParam(r, "id"))

	friendUserName := r.FormValue("fUserName")
	var isExist bool
	isExist = false
	db := config.InitDB()
	defer db.Close()

	row := db.QueryRow("select id from user where username = $1  ", friendUserName)
	if row != nil {
		isExist = true
	}
	var friendID int64
	err := row.Scan(friendID)
	if err != nil {
		returnErrorResponse(w)
		return
	}

	if isExist {

		sq := fmt.Sprintf(`insert into user_friend (user_id,friend_id) values (%d,%d)`,userID,friendID)

		_, err = db.Exec(sq)
		if err != nil {
			returnErrorResponse(w)
			return
		}

	}
}

func (api *API) ResponseSuccess(w http.ResponseWriter, data interface{}) {
	ar, err := response.NewResponse(true, "", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = ar.Send(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func returnErrorResponse(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\": \"error\"}"))
}
func StrToInt64(aval string) int64 {
	aval = strings.Trim(strings.TrimSpace(aval), "\n")
	i, err := strconv.ParseInt(aval, 10, 64)
	if err != nil {
		return 0
	}
	return i
}