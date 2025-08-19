package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func ConnectDB() *sqlx.DB {
	databaseUrl := os.Getenv("DATABASE_URL")

	initDB, err := sqlx.Connect("postgres", databaseUrl)
	if err != nil {
		log.Fatal("Não foi possível conectar-se ao banco de dados", err)
	}

	return initDB
}