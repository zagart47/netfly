package model

import (
	"context"
	"fmt"
	"log"
	"netfly/config"
	"netfly/db"
)

type Message struct {
	ID       int64  `json:"ID"`
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Text     string `json:"text"`
	SendTime string `json:"sendTime"`
}

type MessageArray []Message

func (ma *MessageArray) GetMessageFromDb(ToUser string) error {
	CheckMessageTable()
	messages, err := config.Pool.Query(context.Background(), "SELECT * FROM netfly_messages WHERE to_user = $1", ToUser)
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
			ID:       values[0].(int64),
			FromUser: values[1].(string),
			ToUser:   values[2].(string),
			Text:     values[3].(string),
			SendTime: values[4].(string),
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
		queryAdd := fmt.Sprint("CREATE TABLE netfly_messages(id bigserial primary key, from_user text, to_user text, message_text text, send_time text); ")
		queryOwner := fmt.Sprint("ALTER TABLE netfly_messages OWNER TO postgres;")
		config.Pool.QueryRow(context.Background(), queryAdd)
		config.Pool.QueryRow(context.Background(), queryOwner)
	}
}
