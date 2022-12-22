package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"netfly/config"
	"netfly/db"
	"netfly/model"
	"os"
)

func ReadMessage(c *gin.Context) {
	ma := model.MessageArray{}
	err := ma.GetMessageFromDb(CurrentUser(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": ma})
	}
}

func SendMessage(c *gin.Context) {
	type Input struct {
		Recipient string `json:"recipient"`
		Text      string `json:"text"`
	}

	var User model.User

	type MessageToDb model.Message
	input := Input{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
	}
	u := User
	u.Username = input.Recipient
	err := u.GetUserFromDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		c.Abort()
	}
	m := MessageToDb{
		FromUserID: CurrentUser(c),
		ToUserID:   u.ID,
		Text:       input.Text,
	}
	if m.FromUserID == m.ToUserID {
		c.JSON(http.StatusOK, gin.H{"error": "you cannot send a message to yourself"})
		c.Abort()
	} else {
		insertQueryFrom := fmt.Sprintf("INSERT INTO messages.id%d (from_user_id, to_user_id, message_text, send_time) VALUES (%d, %d, '%s', '%s');", m.FromUserID, m.FromUserID, m.ToUserID, m.Text, db.AddTimeToDb())
		_, err := config.Pool.Query(context.Background(), insertQueryFrom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		insertQueryTo := fmt.Sprintf("INSERT INTO messages.id%d (from_user_id, to_user_id, message_text, send_time) VALUES (%d, %d, '%s', '%s');", m.ToUserID, m.FromUserID, m.ToUserID, m.Text, db.AddTimeToDb())
		_, err = config.Pool.Query(context.Background(), insertQueryTo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
	}

}
