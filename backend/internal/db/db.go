package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

var counts int64

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter, sqldblogger.WithSQLQueryAsMessage(true))
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB() (*sql.DB, error) {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		fmt.Printf("%s is the dsn string \n", dsn)
		if err != nil {
			log.Println("Postgress not yet ready...")
			counts++
		} else {
			log.Println("Connected to postgress...")
			return connection, nil
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("Backing off for 2 seconds....")
		time.Sleep(2 * time.Second)
		continue
	}

}
