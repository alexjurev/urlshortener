package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"encoding/json"
	
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY0123456789"

func shorting() string {
	b := make([]byte, 5)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isValidUrl(token string) bool  {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

type Jsonurl struct {
	URL      string `json:"url"`
  }

type Result struct {
	Link string
	Code string
	Status string
}

func codeLong(w http.ResponseWriter, r *http.Request) {
	result := Result{}
url := Jsonurl{}
w.Header().Set("Content-Type", "application/json")
			_ = json.NewDecoder(r.Body).Decode(&url)
			result.Link = url.URL
	if r.Method == "POST" {
		if !isValidUrl(result.Link) {
			fmt.Println("Что-то не так")
			result.Status = "Ссылка имеет неправильный формат!"
			result.Link = ""
		}else{
			result.Code = shorting()
			db, err := sql.Open("sqlite3", "project.db")
			if err != nil {
				panic(err)
			}
			defer db.Close()
			db.Exec("insert into links (link, short) values ($1, $2)", result.Link, result.Code)
			result.Status = "Сокращение было выполнено успешно"
			url.URL="http://localhost:8000/"+result.Code
		}
	}
	json.NewEncoder(w).Encode(url)
}

func redirectTo(w http.ResponseWriter, r *http.Request)  {
	var link string
	vars := mux.Vars(r)
	db, err := sql.Open("sqlite3", "project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows := db.QueryRow("select link from links where short=$1 limit 1", vars["key"])
	rows.Scan(&link)
	fmt.Fprintf(w, "<script>location='%s';</script>", link)
}

func findShort(w http.ResponseWriter, r *http.Request)  {
	var link string
	surl := Jsonurl{}
	url2 := Jsonurl{}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewDecoder(r.Body).Decode(&surl)
	if r.Method == "POST" {
	db, err := sql.Open("sqlite3", "project.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	sentence := surl.URL
	urlKey := string([]rune(sentence)[22:])
		rows := db.QueryRow("select link from links where short=$1 limit 1", urlKey)
	rows.Scan(&link)
	url2.URL = link
	
}
json.NewEncoder(w).Encode(url2)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/short", codeLong)
	router.HandleFunc("/long", findShort)
	router.HandleFunc("/{key}", redirectTo)
	log.Fatal(http.ListenAndServe(":8000", router))
	
}
