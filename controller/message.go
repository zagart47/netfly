package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netfly/model"
)

func GetMessage(c *gin.Context) {
	var input model.Message
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m := model.Message{ToUser: input.ToUser}
	ma := model.MessageArray{}
	err := ma.GetMessageFromDb(m.ToUser)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": ma})
	}
}
