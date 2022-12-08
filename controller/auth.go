package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"netfly/model"
	"netfly/utils/token"
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
		log.Fatal(err)
	}
	err = u.SaveToDb()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "user registration was successful"})
	}

}

func Login(c *gin.Context) {
	var err error
	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u model.User
	token, err := u.LoginCheck(input.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password incorrect"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CurrentUser(c *gin.Context) {
	u := model.User{}
	var err error
	u.ID, err = token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = u.GetUserByID()
	if err != nil {
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}
