package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	app      string
	url      string
	function string
)

var chanUpdateURL = make(chan bool)

func main() {
	fmt.Println("Starting database driver...")
	db, err := sql.Open("mysql", "app:123456@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if db.Ping() != nil {
		log.Fatal("Unable to open database.")
	}

	go updateURL(db)

	fmt.Println("done")
	http.HandleFunc(url, decodeURL)
	go fmt.Println(http.ListenAndServe(":8080", nil))

}

func updateURL(db *sql.DB) {
	for {
		chanUpdateURL <- true

		fmt.Println("Collecting webpages...")
		rows, err := db.Query("select appname, pageurl, urlfunction from app.url")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		for rows.Next() {
			err := rows.Scan(&app, &url, &function)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(app, url)
		}
	}
}

func decodeURL(resp http.ResponseWriter, req *http.Request) {

}
