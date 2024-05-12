package db

import (
	"database/sql"
	"fmt"
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

type Conversation struct {
	ID int64
	OtherUserChatId int64
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

func (db *DB) IsUserAlreadyInQueue(uchatid int64) bool {
	var userID int64
	err := db.QueryRow("SELECT id FROM queue WHERE chatid = $1", uchatid).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to execute query: %v", err)
	}
	return err != sql.ErrNoRows
}

func (db *DB) IsUserHasConversation(uchatid int64) bool {
	var userID int64
	err := db.QueryRow("SELECT id FROM conversations WHERE user1_chatidid = $1 OR user2_chatid = $1", uchatid).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to execute query: %v", err)
	}
	return err != sql.ErrNoRows
}

func (db *DB) DeleteUserFromQueue(uchatid int64) {
	_, err := db.Exec("DELETE FROM queue WHERE chatid = $1", uchatid)
	onFail("Failed to delete user from queue: %v", err)
}

func (db *DB) DeleteUserConversation(uchatid int64) {
	_, err := db.Exec("DELETE FROM conversations WHERE user1_chatidid = $1 OR user2_chatid = $1", uchatid)
	onFail("Failed to delete user conversation: %v", err)
}

func (db *DB) GetUserConversation(uchatid int64) (conversation Conversation, err error) {
	var (
		conversationID int64
		user1ChatID    int64
		user2ChatID    int64
	)
	err = db.QueryRow("SELECT id, user1_chatidid, user2_chatid FROM conversations WHERE user1_chatidid=$1 OR user2_chatid=$1", uchatid).
		Scan(&conversationID, &user1ChatID, &user2ChatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Conversation{}, fmt.Errorf("conversation not found for user with chat id %d", uchatid)
		}
		return Conversation{}, fmt.Errorf("failed to fetch conversation: %v", err)
	}

	var otherUserChatId int64
	if user1ChatID == uchatid {
		otherUserChatId = user2ChatID
	} else {
		otherUserChatId = user1ChatID
	}

	return Conversation{
		ID: conversationID,
		OtherUserChatId: otherUserChatId,
	}, nil
}

// Связывание подходящих юзеров и удаление их из очереди
func (db *DB) BeginConversation(u1chatid int64, u2chatid int64) {
	_, err := db.Exec("INSERT INTO conversations (user1_chatidid, user2_chatid) VALUES ($1, $2)", u1chatid, u2chatid)
	if err != nil {
		log.Printf("Failed to create conversation %v", err)
		return
	}
	_, err = db.Exec("DELETE FROM queue WHERE chatid IN ($1, $2)", u1chatid, u2chatid)
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
		log.Printf(message, err)
	}
}