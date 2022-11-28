package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"netfly/model"
)

func Register(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "validated!"})
}
