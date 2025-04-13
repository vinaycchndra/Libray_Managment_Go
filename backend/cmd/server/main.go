package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/handlers"
	"github.com/vinaycchndra/Libray_Managment_Go/backend/backend/internal/api/routes"
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

	// Closing the db connection.
	defer app.DB.Close()

	app.Start()

	router := gin.Default()

	// CORS Middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true

	// Allowing router to use the cors middleware.
	router.Use(cors.New(config))

	apiRoutes := router.Group("/api")

	{
		routes.SetupGenericRoutes(apiRoutes, handlers.NewGenericHandler())
	}
	router.Run()
}
