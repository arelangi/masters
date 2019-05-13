package main

import (
	"database/sql"

	"gitlab.logzero.in/arelangi/mlog"
)

type Record struct {
	Id    int64  `json:"id"`
	Tweet string `json:"tweet"`
	Class string `json:"class"`
}

func getCurrentTweet(id int64) (record Record) {
	var tweet string
	var class *sql.NullString
	stmt, err := db.Prepare("SELECT id, tweet,class from masters.tweetclassification where id = $1")
	checkErr(err)

	rows, err := stmt.Query(id)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &tweet, &class)
		checkErr(err)
		record.Id = id
		record.Tweet = tweet
		if class != nil && class.Valid {
			record.Class = class.String
		}
	}
	return
}

func saveNormalizedTweet(record Record) bool {
	_, err := db.Exec("INSERT INTO masters.normalizeAll(id,tweet,class) VALUES ($1,$2,$3);", record.Id, record.Tweet, record.Class)
	if err != nil {
		//		mlog.Error(fmt.Sprintf("The error is %s", err))
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		mlog.Error(err.Error())
	}
}
