package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
)

var (
	app      string
	url      string
	// function string
)

var chanUpdateURL = make(chan bool)

var servLoggerINFO *log.Logger

func main() {
	f, err := os.OpenFile("logs/webServ.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	servLoggerINFO = log.New(f, "INFO ", log.LstdFlags)

	servLoggerINFO.Println("Starting database driver...")
	db, err := sql.Open("mysql", "app:123456@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if db.Ping() != nil {
		log.Fatal("Unable to open database.")
	}

	go updateURL(db)
	servLoggerINFO.Println("Waiting up to 5 minutes to pull URLs...")

	//Wait 5 minutes for pull, or until signal received from goroutine updateURL.
	select {
	case <-chanUpdateURL:
		break
	case <-time.After(5 * time.Minute):
		servLoggerINFO.Println("Unable to complete pulling URL list")
		break
	}

	servLoggerINFO.Println("Setting up web templates...")
	homeTempl := template.New("homeTempl")


	// List of available URLs
	http.HandleFunc("/", homeFunc)
	http.HandleFunc("/test", testFunc)
	go servLoggerINFO.Println(http.ListenAndServe(":8080", nil))

}

func updateURL(db *sql.DB) {
	for {

		servLoggerINFO.Println("Collecting webpages...")
		rows, err := db.Query("select appname, pageurl from app.url")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// i := 0

		for rows.Next() {
			err := rows.Scan(&app, &url)

			if err != nil {
				log.Fatal(err)
			}

			servLoggerINFO.Println(app, url)

		}

		chanUpdateURL <- true
		servLoggerINFO.Println("Done pulling URLs.")
		time.Sleep(15 * time.Minute)
	}
}

func homeFunc(resp http.ResponseWriter, req *http.Request) {
	
	servLoggerINFO.Println("Loading webpage request.")
	http.ServeFile(resp, req, "logs/webServ.log")

}

func testFunc(resp http.ResponseWriter, req *http.Request)  {
	servLoggerINFO.Println("Loading test webpageURL")
	http.ServeFile(resp, req, "logs/webServ.log")
}