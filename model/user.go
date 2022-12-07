package model

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html"
	"netfly/config"
	"netfly/db"
	"netfly/utils/token"
	"os"
	"strings"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
	err := UserAdd(u.Username, u.Password)
	if err != nil {
		return err
	} else {
		return nil
	}
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
	err = VerifyPassword(password, GetUserHashedPwd(name))
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
	err := config.Pool.QueryRow(context.Background(), "SELECT id, user_name, created_at, updated_at FROM netfly_users WHERE id = $1", uid).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return u, fmt.Errorf("user not found")
	}
	u.PwdCap()
	return u, nil

}

func (u *User) PwdCap() {
	u.Password = ""
}

func UserAdd(name string, password string) error {
	db.CheckConnect()
	if db.GetUserID(name) != 0 {
		return fmt.Errorf("username already registered. choose another name")
	}
	_, err := config.Pool.Query(context.Background(), "INSERT INTO netfly_users (user_name, password, created_at) VALUES ($1, $2, $3)", name, password, db.AddTimeToDb())
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
		os.Exit(1)

	}
	return nil
}

func GetUserHashedPwd(name string) string {
	db.CheckConnect()
	var pwd string
	err := config.Pool.QueryRow(context.Background(), "SELECT password FROM netfly_users WHERE user_name = $1", name).Scan(&pwd)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return pwd
}
