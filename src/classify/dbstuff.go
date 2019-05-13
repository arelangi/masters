package main

import (
	"fmt"

	"gitlab.logzero.in/arelangi/mlog"
)

type Record struct {
	Id    int64  `json:"id"`
	Tweet string `json:"tweet"`
	Class string `json:"class"`
}

func getCurrentTweet() (record Record) {
	var id int64
	var tweet string
	stmt, err := db.Prepare("SELECT id, tweet from masters.tweetclassification where id = (select max(id) from masters.tweetclassification where class is not null and id < 50000 limit 1)")
	checkErr(err)

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &tweet)
		checkErr(err)
		record.Id = id
		record.Tweet = tweet
	}
	return
}

func getCurrentTweetL2() (record Record) {
	var id int64
	var tweet string
	stmt, err := db.Prepare("SELECT id, tweet from masters.tweetclassification where id = (select max(id)-1 from masters.tweetclassification where class is not null limit 1)")
	checkErr(err)

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &tweet)
		checkErr(err)
		record.Id = id
		record.Tweet = tweet
	}
	return
}

func getPrevTweet(id int64) (record Record) {
	var tweet string
	stmt, err := db.Prepare("SELECT id, tweet from masters.tweetclassification where id = $1-1;")
	checkErr(err)

	rows, err := stmt.Query(id)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &tweet)
		checkErr(err)
		record.Id = id
		record.Tweet = tweet
	}
	return
}

func getNextTweet(id int64) (record Record) {
	var tweet string
	stmt, err := db.Prepare("SELECT id, tweet from masters.tweetclassification where id = $1+1;")
	checkErr(err)

	rows, err := stmt.Query(id)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &tweet)
		checkErr(err)
		record.Id = id
		record.Tweet = tweet
	}
	return
}

func getNextTweetL2(id int64) (record Record) {
	var tweet string
	stmt, err := db.Prepare("SELECT id, tweet from masters.tweetclassification where id = (select max(id) from masters.tweetclassification where class is null limit 1)")
	checkErr(err)

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &tweet)
		checkErr(err)
		record.Id = id
		record.Tweet = tweet
	}
	return
}

func updateTweet(record Record) bool {
	_, err := db.Exec("UPDATE masters.tweetclassification set class=$1 where id=$2", record.Class, record.Id)
	if err != nil {
		mlog.Error(fmt.Sprintf("The error is %s", err))
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		mlog.Error(err.Error())
	}
}
