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

func getAllCMTweets() (tweets []string) {
	var tweet string
	stmt, err := db.Prepare("SELECT tweet from masters.tweetclassification where class='CM'")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tweet)
		if err != nil {
			fmt.Println(err)
		}
		tweets = append(tweets, tweet)
	}
	return
}

func saveToUntranslatedTweets(resp Response) {
	var mitie, nltk, spacy, prose string
	if len(resp.Entities["mitie"]) > 0 {
		mitie = resp.Entities["mitie"][0]
	}

	if len(resp.Entities["nltk"]) > 0 {
		nltk = resp.Entities["nltk"][0]
	}
	if len(resp.Entities["spacy"]) > 0 {
		spacy = resp.Entities["spacy"][0]
	}
	if len(resp.Entities["prose"]) > 0 {
		prose = resp.Entities["prose"][0]
	}
	_, err := db.Exec("INSERT INTO masters.untranslatedEntities(tweet,mitie,nltk,prose,spacy) VALUES ($1,$2,$3,$4,$5);", resp.Tweet.OriginalText, mitie, nltk, prose, spacy)
	if err != nil {
		mlog.Error(fmt.Sprintf("The error is %s", err))
	}
}

func saveToTranslatedTweets(resp Response) {
	var mitie, nltk, spacy, prose string
	if len(resp.Entities["mitie"]) > 0 {
		mitie = resp.Entities["mitie"][0]
	}

	if len(resp.Entities["nltk"]) > 0 {
		nltk = resp.Entities["nltk"][0]
	}
	if len(resp.Entities["spacy"]) > 0 {
		spacy = resp.Entities["spacy"][0]
	}
	if len(resp.Entities["prose"]) > 0 {
		prose = resp.Entities["prose"][0]
	}
	_, err := db.Exec("INSERT INTO masters.translatedEntities(tweet,mitie,nltk,prose,spacy) VALUES ($1,$2,$3,$4,$5);", resp.Tweet.OriginalText, mitie, nltk, prose, spacy)
	if err != nil {
		mlog.Error(fmt.Sprintf("The error is %s", err))
	}
}
