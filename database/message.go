package database

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        int64
	ChatId    int64
	UserId    int64
	UserName  string
	Message   string
	CreatedAt time.Time
}

func CreateMessage(db *sql.DB, message Message) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO message (chat_id, user_id, user_name, message) 
		VALUES (?, ?, ?, ?)`,
		message.ChatId, message.UserId, message.UserName, message.Message,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
