package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	router = mux.NewRouter()
}

func main() {

	router.HandleFunc("/", helloWorldHandler)
	router.HandleFunc("/current", getCurrent)
	router.HandleFunc("/currentL2", getCurrent)
	router.HandleFunc("/next", getNext)
	router.HandleFunc("/nextl2", getNextL2)
	router.HandleFunc("/prev", getPrev)
	router.HandleFunc("/classify", classify)

	log.Fatal(http.ListenAndServe(":8881", router))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, nice meeting you. \n What are you doing here?"))
}

func getCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	resp := getCurrentTweet()
	respString, _ := json.Marshal(resp)
	w.Write(respString)
}

func getCurrentL2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	resp := getCurrentTweetL2()
	respString, _ := json.Marshal(resp)
	w.Write(respString)
}

func getNext(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var idVal int64
	id := r.PostFormValue("id")
	if idVal, err = strconv.ParseInt(id, 10, 64); err != nil {
		checkErr(err)
	}
	resp := getNextTweet(idVal)
	respString, _ := json.Marshal(resp)
	w.Write(respString)
}

func getNextL2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var idVal int64
	id := r.PostFormValue("id")
	if idVal, err = strconv.ParseInt(id, 10, 64); err != nil {
		checkErr(err)
	}
	resp := getNextTweetL2(idVal)
	respString, _ := json.Marshal(resp)
	w.Write(respString)
}

func getPrev(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var idVal int64
	id := r.PostFormValue("id")
	if idVal, err = strconv.ParseInt(id, 10, 64); err != nil {
		checkErr(err)
	}
	resp := getPrevTweet(idVal)
	respString, _ := json.Marshal(resp)
	w.Write(respString)
}

func classify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var err error
	var idVal int64
	id := r.PostFormValue("id")
	if idVal, err = strconv.ParseInt(id, 10, 64); err != nil {
		checkErr(err)
	}
	class := r.PostFormValue("class")
	record := Record{Id: idVal, Class: class}
	if updateTweet(record) {
		//Let's return the next stuff
		resp := getNextTweet(idVal)
		respString, _ := json.Marshal(resp)
		w.Write(respString)
	}
}
