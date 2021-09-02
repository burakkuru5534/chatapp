package repository

import (
	"database/sql"
	"example.com/m/models"
	"github.com/google/uuid"
	"log"

)


type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) GetId() string {
	return user.Id
}

func (user *User) GetName() string {
	return user.Name
}

func (user *User) GetUserName() string {
	return user.Username
}

func (user *User) SaveMessage()  {
	user.SaveMessage()
}



type UserRepository struct {
	Db *sql.DB
}

func (repo *UserRepository) AddUser(user models.User) {
	stmt, err := repo.Db.Prepare("INSERT INTO user(id, name, username) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(user.GetId(), user.GetName(), user.GetName())
	checkErr(err)
}

func (repo *UserRepository) AddFriend(id string, friendID string, userID string) {
	stmt, err := repo.Db.Prepare("INSERT INTO user_friend(id, friend_id, user_id) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, friendID, userID)
	checkErr(err)
}

func (repo *UserRepository) CreateUser(id string, name string, username string, password string) {

	stmt, err := repo.Db.Prepare("INSERT INTO user(id, name, username, password) values(?,?, ?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, name, username, password)
	checkErr(err)

}

func (repo *UserRepository) RemoveUser(user models.User) {
	stmt, err := repo.Db.Prepare("DELETE FROM user WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(user.GetId())
	checkErr(err)
}

func (repo *UserRepository) FindUserById(ID string) models.User {

	row := repo.Db.QueryRow("SELECT id, name FROM user where id = ? LIMIT 1", ID)

	var user User

	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user

}

func (repo *UserRepository) GetAllUsers() []models.User {

	rows, err := repo.Db.Query("SELECT id, name FROM user")

	if err != nil {
		log.Fatal(err)
	}
	var users []models.User
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users
}

func (repo *UserRepository) GetAllFriends(userID string) []models.User {

	rows, err := repo.Db.Query("SELECT id, name, username FROM user where id in (select friend_id from user_friend where user_id = ?)",userID)

	if err != nil {
		log.Fatal(err)
	}
	var users []models.User
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name, &user.Username)
		users = append(users, &user)
	}

	return users
}


func (repo *UserRepository) IsAlreadyFriend(userID string, friendID string) bool {

	rows, err := repo.Db.Query("SELECT id FROM user_friend where user_id = ? and friend_id = ?",userID, friendID)
	if err != nil {
		log.Fatal(err)
	}
	var id string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id)
		if id != "" {
			return true
		}
	}

	return false
}


func (repo *UserRepository) FindUserByUsername(username string) *User {

	row := repo.Db.QueryRow("SELECT id, name, username FROM user where username = ? LIMIT 1", username)

	var user User

	if err := row.Scan(&user.Id, &user.Name, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user
}

func (repo *UserRepository) GetUserByUserName(username string) *User {

	row := repo.Db.QueryRow("SELECT id, name, username, password FROM user where username = ? LIMIT 1", username)

	var user User

	if err := row.Scan(&user.Id, &user.Name, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}

	return &user
}

func (repo *UserRepository) SaveMessage (userID string, toID string, content string) {

	id := uuid.New().String()
	stmt, err := repo.Db.Prepare("INSERT INTO msg(id, content, user_id, to_id) values(?,?, ?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, content, userID, toID)
	checkErr(err)
	repo.Db.Close()

}

