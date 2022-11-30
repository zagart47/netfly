package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var _ = godotenv.Load()
var DbHost = os.Getenv("DBHOST")
var Pool, err = pgxpool.New(context.Background(), DbHost)
var TokenLifespan, _ = strconv.Atoi(os.Getenv("TOKENLIFESPAN"))

func Configure() {
	if err != nil {
		log.Fatalln(err)
	}
}
