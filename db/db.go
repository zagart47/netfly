package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"netfly/config"
	"os"
)

func ConnectDb() {
	var err error
	config.Pool, err = pgxpool.New(context.Background(), config.DbHost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer config.Pool.Close()
}

func UserAdd(name string, password string) error {
	CheckConnect()
	_, err := config.Pool.Query(context.Background(), "INSERT INTO netfly_users (user_name, password) VALUES ($1, $2)", name, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return err
		os.Exit(1)

	}
	return nil
}

func CheckConnect() {
	if config.Pool.Ping(context.Background()) != nil {
		ConnectDb()
	}
	CheckTable()
}

func CheckTable() {
	var tableStatus bool
	err := config.Pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'netfly_users');").Scan(&tableStatus)
	if err != nil {
		log.Fatal(err)
	}
	if tableStatus != true {
		queryAdd := fmt.Sprint("CREATE TABLE netfly_users(id bigserial primary key, user_name text, password text); ")
		queryOwner := fmt.Sprint("ALTER TABLE netfly_users OWNER TO postgres;")
		config.Pool.QueryRow(context.Background(), queryAdd)
		config.Pool.QueryRow(context.Background(), queryOwner)
	}
}

func CheckUser(name string) error {
	CheckConnect()
	rows, err1 := config.Pool.Query(context.Background(), "SELECT user_name FROM netfly_users WHERE user_name=$1", name)
	if err1 != nil {
		log.Fatal(fmt.Errorf("ошибка запроса в БД"))
	}
	for rows.Next() {
		values, err2 := rows.Values()
		if err2 != nil {
			return fmt.Errorf("ошибка парсинга БД")
		}
		if name == values[0].(string) {
			return nil
		}
	}
	return fmt.Errorf("не найдено совпадений в имени")
}

func GetUserID(name string) uint {
	CheckConnect()
	var ID uint
	err := config.Pool.QueryRow(context.Background(), "SELECT id FROM netfly_users WHERE user_name = $1", name).Scan(&ID)
	if err != nil {
		return 0
	}
	return ID
}

func GetUserHashedPwd(name string) string {
	CheckConnect()
	var pwd string
	err := config.Pool.QueryRow(context.Background(), "SELECT password FROM netfly_users WHERE user_name = $1", name).Scan(&pwd)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return pwd
}
