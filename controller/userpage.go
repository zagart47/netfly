package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netfly/model"
)

func Profile(c *gin.Context) {
	u := model.User{}
	u.ID = CurrentUser(c)
	err := u.GetUserByID()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": u})
	}
}
