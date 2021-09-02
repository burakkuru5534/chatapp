package main

import (
	"encoding/json"
	"example.com/m/auth"
	"example.com/m/cmn"
	"example.com/m/config"
	"example.com/m/models"
	"example.com/m/repository"
	"example.com/m/response"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
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
type AnonUser struct {
	Id string `json:"id"`
	Name string `json:"name"`
}
type contextKey string
const UserContextKey = contextKey("user")

func (api *API) RegisterUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var isUserNameExist bool
	isUserNameExist = true
	// Find the user in the database by username
	dbUser, err := api.UserRepository.FindUserByUsername(username,"register")
	if err != nil {
		api.ResponseFail(w,"wrong username")
		return
	}
	if dbUser == nil {
		isUserNameExist = false
	}

	id := uuid.New().String()
	if !isUserNameExist {
		hashedPassword, err := cmn.HashPasswd(password)
		if err != nil {
			api.ResponseFail(w,"hash password error")
			return

		}
		api.UserRepository.CreateUser(id, username,username,hashedPassword)

		api.ResponseSuccess(w, id)
		return
	}else {
		api.ResponseFail(w,"username already exist error")
		return
	}

}
func (api *API) HandleLogin(w http.ResponseWriter, r *http.Request) {

	var user LoginUser

	// Try to decode the JSON request to a LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		api.ResponseFail(w,"request body decode error")
		return
	}

	// Find the user in the database by username
	dbUser := api.UserRepository.GetUserByUserName(user.Username)
	if dbUser == nil {
		api.ResponseFail(w,"find user by username error")
		return
	}

	// Check if the passwords match
	ok:= cmn.CheckPass(dbUser.Password,user.Password)
	if !ok {
		api.ResponseFail(w,"password not correct")
		return
	}

	// Create a JWT
	token, err := auth.CreateJWTToken(dbUser)
	if err != nil {
		api.ResponseFail(w,"token create error")
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
		api.ResponseFail(w,"delete all user db exec error")
		return
	}

	api.ResponseSuccess(w, "All Users deleted with Success")


}
func (api *API) ShowUserFriends(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")

	db := config.InitDB()
	defer db.Close()

	userFriends := api.UserRepository.GetAllFriends(userID)

	api.ResponseSuccess(w, userFriends)


}
func (api *API) AddFriend(w http.ResponseWriter, r *http.Request) {

	token, _ := r.URL.Query()["bearer"]
	user, err := ValidateToken(token[0])
	if err != nil {
		api.ResponseFail(w,"wrong token")
		return
	}

	userID := user.GetId()
	friendUserName := r.FormValue("fUserName")
	friendUser, err := api.UserRepository.FindUserByUsername(friendUserName,"add-friend")
	if err != nil {
		api.ResponseFail(w,"wrong friend name")
		return
	}
	friendUserID := friendUser.GetId()
	isAlreadyFriend := api.UserRepository.IsAlreadyFriend(userID,friendUserID)
	if  friendUserID != "" && !isAlreadyFriend{
		id := uuid.New().String()
		api.UserRepository.AddFriend(id,friendUserID,userID)
		api.ResponseSuccess(w,"Friend Added with Success")
		return
	}else {
		api.ResponseFail(w, "friend already exist")
		return
	}

}
func (api *API) ShowUsersPostedMessagesToFriend(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")
	toID := r.FormValue("toID")

	db := config.InitDB()
	defer db.Close()

	userSentMessagesToSomeOne := api.UserRepository.GetUserSentMessagesByToID(userID,toID)

	api.ResponseSuccess(w, userSentMessagesToSomeOne)
}

func (api *API) ShowUsersReceivedMessagesFromFriend(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")
	toID := r.FormValue("toID")

	db := config.InitDB()
	defer db.Close()

	userSentMessagesFromSomeOne := api.UserRepository.GetUserReceivedMessagesByToID(userID,toID)

	api.ResponseSuccess(w, userSentMessagesFromSomeOne)
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
func (api *API) ResponseFail(w http.ResponseWriter, msg string) {
	ar, err := response.NewResponse(false, msg, nil)
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

func printContextInternals(ctx interface{}, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				printContextInternals(reflectValue.Interface(), true)
			} else {
				fmt.Printf("field name: %+v\n", reflectField.Name)
				fmt.Printf("value: %+v\n", reflectValue.Interface())
			}
		}
	} else {
		fmt.Printf("context is empty (int)\n")
	}
}

func ValidateToken(tokenString string) (models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSecret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

type Claims struct {
	ID string `json:"id"`
	Name string `json:"name"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

const hmacSecret = "SecretValueReplaceThis"

func (c *Claims) GetId() string {
	return c.ID
}

func (c *Claims) GetName() string {
	return c.Name
}

func (c *Claims) GetUserName() string {
	return c.UserName
}
