package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/db"
)

type Config struct {
	DB *sql.DB
}

func (app *Config) Start() {
	fmt.Println("We have started the server successfully.")
}

func main() {
	err := godotenv.Load("../../../.env")

	if err != nil {
		log.Panicf("Error loading .env file: %s", err)
	}

	db_conn, err := db.InitDB()

	if err != nil {
		log.Panicf("Error in initialising db: %s", err)
	}

	err = db.RunMigrations(db_conn)

	if err != nil {
		log.Panicf("Error in running migrations: %s", err)
	}

	app := Config{DB: db_conn}

	app.Start()
}
