package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	*sql.DB
}

type User struct {
	ID int
	ChatId int64
	City string
}

func NewDB(dbUrl string) (*DB, error) {
	db, err := connectDB(dbUrl)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func connectDB(dbUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Связывание подходящих юзеров и удаление их из очереди
func (db *DB) BeginConversation(u1chatid int64, u2chatid int64) {
	_, err := db.Exec("INSERT INTO conversations (user1_chatidid, user2_chatid) VALUES ($1, $2)", u1chatid, u2chatid)
	if err != nil {
		log.Fatalf("Failed to create conversation %v", err)
		return
	}
	_, err = db.Exec("DELETE FROM queue WHERE id IN ($1, $2)", u1chatid, u2chatid)
	onFail("Failed to delete users from queue %v", err)
}

func (db *DB) AddUserToQueue(userChatId int64, city string) error {
	_, err := db.Exec("INSERT INTO queue (chatid, city) VALUES ($1, $2)", userChatId, city)
	onFail("Failed to add user %v", err)
	return err
}

func (db *DB) GetUsersFromQueueByCity(city string) ([]User, error) {
	rows, err := db.Query("SELECT id, chatid, city FROM queue WHERE city = $1", city)
	onFail("Failed to get queue by city %v", err)
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