package db

import (
	"database/sql"
	"log"
)

type DB struct {
	*sql.DB
}

type User struct {
	ID int
	ChatId int
	City string
}

func NewDB() (*DB, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://anon_chat_tg@localhost/anon_chat_tg?sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) AddUserToQueue(userChatId int, city string) error {
	_, err := db.Exec("INSERT INTO users (chatid, city) VALUES ($1, $2)")
	onFail("Failed to add user %v", err)
	return err
}

func (db *DB) GetUsersFromQueueByCity(city string) ([]User, error) {
	rows, err := db.Query("SELECT id, chatid, city FROM users WHERE city = $1", city)
	onFail("Failed to get users by city %v", err)
	defer rows.Close()
	
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.ChatId, &user.City)
		onFail("Failed to scan user row %v", err)
		users = append(users, user)
	}
	return users, nil
}

func onFail(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}