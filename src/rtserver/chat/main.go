package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8090", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/home" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func serveUser(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/user" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "user.html")
}

type predictReq struct {
	Text string `json:"text"`
}

func predictIndividual(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	log.Println(r.URL)
	if r.URL.Path != "/predict" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg predictReq
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("The inpur we received is")
	fmt.Println(msg)

	msg1 := RTTweet{OriginalText: msg.Text, NormalizedText: normalize(msg.Text)}
	resp, _ := predictions(msg1)
	jsonResp, _ := json.Marshal(resp)

	w.Write(jsonResp)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	go streamProcessing()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/home", serveHome)
	http.HandleFunc("/user", serveUser)
	http.HandleFunc("/predict", predictIndividual)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
