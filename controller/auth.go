package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"netfly/model"
)

func Register(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := model.User{}

	u.Username = input.Username
	u.Password = input.Password

	err := u.CryptPwd()
	if err != nil {
		log.Fatalln(err)
	}
	err = u.SaveToDb()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "validated!"})
}

func Login(c *gin.Context) {
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	u := model.User{}

	u.Username = input.Username
	u.Password = input.Password

	token, err := model.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
