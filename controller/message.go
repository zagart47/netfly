package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netfly/model"
)

func GetMessage(c *gin.Context) {
	ma := model.MessageArray{}
	err := ma.GetMessageFromDb(CurrentUser(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": ma})
	}
}
