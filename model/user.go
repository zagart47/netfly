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
	err := u.UserAdd()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (u *User) VerifyPassword(inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
}

func (u *User) LoginCheck(inputPassword string) (string, error) {
	err := u.GetUserFromDb()
	if err != nil {
		return "", err
	}
	err = u.VerifyPassword(inputPassword)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) GetUserByID() error {
	db.CheckConnect()
	err := config.Pool.QueryRow(context.Background(), "SELECT id, user_name, created_at, updated_at FROM netfly_users WHERE id = $1", u.ID).Scan(&u.ID, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}
	u.PwdCap()
	return nil

}

func (u *User) PwdCap() {
	u.Password = ""
}

func (u *User) UserAdd() error {
	u.GetUserFromDb()
	if u.ID != 0 {
		return fmt.Errorf("username already registered. choose another name")
	}
	_, err := config.Pool.Query(context.Background(), "INSERT INTO netfly_users (user_name, password, created_at, updated_at) VALUES ($1, $2, $3, $3)", u.Username, u.Password, db.AddTimeToDb())
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
		os.Exit(1)

	}
	u.GetUserFromDb()
	return nil
}

func (u *User) GetUserFromDb() error {
	db.CheckConnect()
	err := config.Pool.QueryRow(context.Background(), "SELECT id, user_name, password, created_at, updated_at FROM netfly_users WHERE user_name=$1", u.Username).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("user not found in the db")
	}
	return nil
}
