package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./chatdb.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `	
	CREATE TABLE IF NOT EXISTS room (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		private TINYINT NULL
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	sqlStmt = ` 
    CREATE TABLE IF NOT EXISTS user (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        username VARCHAR(255) NULL,
        password VARCHARR(255) NULL
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	sqlStmt = ` 
    CREATE TABLE IF NOT EXISTS user_friend (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        user_id  VARCHAR(255) not null,
        friend_id VARCHAR(255) not null,

         constraint fk_user_friend_user foreign key (user_id)
        references user (id) on update cascade on delete no action,
        
        constraint fk_user_friend_friend foreign key (friend_id)
        references user (id) on update cascade on delete no action
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	sqlStmt = ` 
    CREATE TABLE IF NOT EXISTS msg (
        id VARCHAR(255) NOT NULL PRIMARY KEY,
        content VARCHAR(255) not null,
        user_id VARCHAR(255) not null,
        to_id VARCHAR(255) not null,

         constraint fk_msg_user foreign key (user_id)
        references user (id) on update cascade on delete no action,
        
        constraint fk_msg_to foreign key (to_id)
        references room (id) on update cascade on delete no action
    );
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
	}

	return db
}
