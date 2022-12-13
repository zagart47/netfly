package model

import (
	"context"
	"fmt"
	"log"
	"netfly/config"
	"netfly/db"
)

type Message struct {
	ID         int64  `json:"ID"`
	FromUserID int64  `json:"fromUserID"`
	ToUserID   int64  `json:"toUserID" binding:"required"`
	Text       string `json:"text" binding:"required"`
	SendTime   string `json:"sendTime"`
}

type MessageArray []Message

func (ma *MessageArray) GetMessageFromDb(ToUserID int64) error {
	CheckMessageTable()
	messages, err := config.Pool.Query(context.Background(), "SELECT * FROM netfly_messages WHERE to_user_id = $1", ToUserID)
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

func CheckMessageTable() {
	if config.Pool.Ping(context.Background()) != nil {
		db.ConnectDb()
	}
	var tableStatus bool
	err := config.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'netfly_messages');").Scan(&tableStatus)
	if err != nil {
		log.Fatal(err)
	}
	if tableStatus != true {
		queryAdd := fmt.Sprint("CREATE TABLE netfly_messages(id bigserial primary key, from_user_id bigint, FOREIGN KEY (from_user_id) REFERENCES netfly_users(id), to_user_id bigint, FOREIGN KEY (to_user_id) REFERENCES netfly_users(id), message_text text, send_time text);")
		queryOwner := fmt.Sprint("ALTER TABLE netfly_messages OWNER TO postgres;")
		config.Pool.QueryRow(context.Background(), queryAdd)
		config.Pool.QueryRow(context.Background(), queryOwner)
	}
}
