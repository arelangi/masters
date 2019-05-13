package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gitlab.logzero.in/arelangi/mlog"
)

var (
	_dbUser     = "aditya"
	_dbPassword = "aditya"
	_dbName     = "aditya"
	_dbHost     = "localhost"
	db          *sql.DB
)

func init() {
	var err error
	user := os.Getenv("DBUSER")
	if user != "" {
		_dbUser = user
	}

	password := os.Getenv("DBPASSWORD")
	if password != "" {
		_dbPassword = password
	}

	name := os.Getenv("DBNAME")
	if name != "" {
		_dbName = name
	}

	host := os.Getenv("DBHOST")
	if host != "" {
		_dbHost = host
	}

	logLevel := os.Getenv("logLevel")
	if logLevel != "" {
		mlog.SetLogLevel(logLevel)
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		_dbUser, _dbPassword, _dbName, _dbHost)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		mlog.Error("Failed to connect to the database", mlog.Items{"error": err})
	}

}
