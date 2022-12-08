package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"netfly/config"
	"os"
	"time"
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
		queryAdd := fmt.Sprint("CREATE TABLE netfly_users(id bigserial primary key, user_name text, password text, created_at text, updated_at text); ")
		queryOwner := fmt.Sprint("ALTER TABLE netfly_users OWNER TO postgres;")
		config.Pool.QueryRow(context.Background(), queryAdd)
		config.Pool.QueryRow(context.Background(), queryOwner)
	}
}

func AddTimeToDb() string {
	dateTimeToDb := time.Now()
	return fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d", dateTimeToDb.Day(), dateTimeToDb.Month(), dateTimeToDb.Year(), dateTimeToDb.Hour(), dateTimeToDb.Minute(), dateTimeToDb.Second())
}
