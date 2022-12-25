package model

import (
	"context"
	"fmt"
	"netfly/config"
	"netfly/db"
	"strconv"
)

type Message struct {
	ID         int64  `json:"ID"`
	FromUserID int64  `json:"fromUserID"`
	ToUserID   int64  `json:"toUserID" binding:"required"`
	Text       string `json:"text" binding:"required"`
	SendTime   string `json:"sendTime"`
}

type MessageArray []Message

func (ma *MessageArray) ReadAllMessagesFromDb(ToUserID int64) error {
	queryUserMessageTable := fmt.Sprintf("SELECT * FROM messages.id%s", strconv.FormatInt(ToUserID, 10))
	messages, err := config.Pool.Query(context.Background(), queryUserMessageTable)
	if err != nil {
		return err
	}
	for messages.Next() {
		values, err := messages.Values()
		if err != nil {
			return err
		}
		if values == nil {
			return fmt.Errorf("user have no messages")
		}
		*ma = append(*ma, Message{
			ID:         values[0].(int64),
			FromUserID: values[1].(int64),
			ToUserID:   values[2].(int64),
			Text:       values[3].(string),
			SendTime:   values[4].(string),
		})
	}
	return nil

}

func (u *User) CreateUserMessageDb() {
	if config.Pool.Ping(context.Background()) != nil {
		db.ConnectDb()
	}
	/*_, err := config.Pool.Exec(context.Background(), "CREATE SCHEMA IF NOT EXISTS messages")
	if err != nil {
		fmt.Println(err.Error())
	}*/
	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS messages.id%s (
		id bigserial primary key,
		from_user_id bigint,
		FOREIGN KEY (from_user_id) REFERENCES netfly_users(id), 
		to_user_id bigint, 
		FOREIGN KEY (to_user_id) REFERENCES netfly_users(id), 
		message_text text, send_time text);`, strconv.FormatInt(u.ID, 10))
	_, err := config.Pool.Exec(context.Background(), createTableQuery)
	if err != nil {
		fmt.Println(err.Error())
	}
	ownerQuery := fmt.Sprintf("ALTER TABLE messages.id%s OWNER TO postgres;", strconv.FormatInt(u.ID, 10))
	_, err = config.Pool.Exec(context.Background(), ownerQuery)
	if err != nil {
		return
	}
}
