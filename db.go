package main

import (
	"database/sql"
	"math/rand"
	"net/http"
	"time"
)
import _ "modernc.org/sqlite"

var db *sql.DB

func CreateDB() {
	var err error
	db, err = sql.Open("sqlite", "./URLS.db")
	if err != nil {
		panic(err)
	}
}

func CreateTables() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS urls (id INTEGER PRIMARY KEY, url TEXT, short TEXT)")
	if err != nil {
		panic(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateLetters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateShort(url string) string {
	short := GenerateLetters(5)
	_, err := db.Exec("INSERT INTO urls (url, short) VALUES (?, ?)", url, short)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/"+short, HandleCode)
	return short
}

func GetURL(short string) string {
	var url string
	err := db.QueryRow("SELECT url FROM urls WHERE short=?", short).Scan(&url)
	if err != nil {
		panic(err)
	}
	return url
}

func GetShort(url string) string {
	var short string
	err := db.QueryRow("SELECT short FROM urls WHERE url=?", url).Scan(&short)
	if err != nil {
		return ""
	}
	return short
}

func Startup() {
	rows, err := db.Query("SELECT * FROM urls")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var url string
		var short string
		err = rows.Scan(&id, &url, &short)
		if err != nil {
			panic(err)
		}
		http.HandleFunc("/"+short, HandleCode)
	}
}
