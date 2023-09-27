package main

import (
	"fmt"
	"net/http"
	"strings"
)

func ParseURL(url string) string {
	_url := strings.Split(url, "://")[0]
	_parsed := strings.Split(url, "://")[1]
	if _url == "http" {
		return _parsed
	}
	if _url == "https" {
		return _parsed
	} else {
		return "https://" + _parsed
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		url := r.FormValue("url")
		short := GetShort(url)
		if short == "" {
			short = GenerateShort(url)
		}
		parsed := ParseURL(url)
		_, err := http.Get(parsed)
		if err != nil {
			w.Write([]byte("Invalid URL"))
			return
		}
		w.Write([]byte("http://45.142.114.111/" + short))
		return
	} else if r.Method == "GET" {
		short := r.FormValue("short")
		if short != "" {
			url := GetURL(short)
			if url != "" {
				http.Redirect(w, r, url, http.StatusFound)
				return
			}
		}
		http.ServeFile(w, r, "index.html")
		return
	}

}

func HandleCode(w http.ResponseWriter, r *http.Request) {
	url := GetURL(r.URL.Path[1:])
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	CreateDB()
	CreateTables()
	Startup()
	fmt.Println("Server started")
	http.HandleFunc("/", handleConnection)
	err := http.ListenAndServe("0.0.0.0:80", nil)
	if err != nil {
		panic(err)
	}
}
