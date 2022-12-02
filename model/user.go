package model

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html"
	"netfly/config"
	"netfly/db"
	"netfly/utils/token"
	"strings"
)

type User struct {
	ID       int64  `json:"id"`
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
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(name string, password string) (string, error) {
	var err error

	err = db.CheckUser(name)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}
	err = VerifyPassword(password, db.GetUserHashedPwd(name))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", fmt.Errorf("wrong password")
	}
	token, err := token.GenerateToken(db.GetUserID(name))
	if err != nil {
		return "", err
	}

	return token, nil

}

func GetUserByID(uid uint) (User, error) {
	var u User
	db.CheckConnect()
	err := config.Pool.QueryRow(context.Background(), "SELECT user_name FROM netfly_users WHERE id = $1", uid).Scan(&u.Username)
	if err != nil {
		return u, fmt.Errorf("user not found")
	}
	u.PrepareGive()
	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}
