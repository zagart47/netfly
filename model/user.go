package model

import (
	"golang.org/x/crypto/bcrypt"
	"html"
	"netfly/db"
	"netfly/utils/token"
	"strings"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) CryptPwd() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil

}

func (u *User) SaveToDb() error {
	err := db.UserAdd(u.Username, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))
}

func LoginCheck(name string, password string) (string, error) {
	var err error
	u := User{}

	err = db.CheckUser(name)
	if err != nil {
		return "", err
	}
	err = VerifyPassword(password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(db.CheckUserID(name))
	if err != nil {
		return "", err
	}

	return token, nil

}
