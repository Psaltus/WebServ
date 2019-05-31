package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	app string
	url string
	// function string
	homeTempl *template.Template
	testTempl *template.Template
)

type homeDataType struct {
	Title       string
	HeaderTitle string
	Body        string
}

var homeData homeDataType
var testData homeDataType

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

	servLoggerINFO.Println("Setting up static directory")
	fs := http.FileServer(http.Dir("src"))
	http.Handle("/src/", http.StripPrefix("/src/", fs))

	// servLoggerINFO.Println("Setting up web templates...")
	// homeTempl = template.New("homeTempl")
	// testTempl = template.New("testTempl")

	// List of available URLs
	http.HandleFunc("/", homeFunc)
	http.HandleFunc("/test", testFunc)
	/*go*/ http.ListenAndServe(":8080", nil)
	//TODO: Set ListenAndServe to goroutine, insert CLI with instructional commands.

}

//FIXME: func name no longer matches use-case, need refactor
func updateURL(db *sql.DB) {
	for {

		servLoggerINFO.Println("Collecting webpages...")
		rows, err := db.Query("select title, headertitle, body from app.pageoutput")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// fmt.Println(rows)

		// i := 0

		// Setup page data TODO: fix database parsing
		for rows.Next() {
			err = rows.Scan(&homeData.Title, &homeData.HeaderTitle, &homeData.Body)
			if err != nil {
				log.Fatal(err)
			}

			//servLoggerINFO.Println(homeData)
			fmt.Println(homeData)

			rows.Next()

			err = rows.Scan(&testData.Title, &testData.HeaderTitle, &testData.Body)
			if err != nil {
				break
			}

			fmt.Println(testData)
			break
		}

		chanUpdateURL <- true
		servLoggerINFO.Println("Done pulling URLs.")
		time.Sleep(1 * time.Minute)
	}
}

func homeFunc(resp http.ResponseWriter, req *http.Request) {

	servLoggerINFO.Println("Loading webpage request.")
	//http.ServeFile(resp, req, "logs/webServ.log")
	//return
	homeTempl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Failed to parse index.html")
	}
	servLoggerINFO.Println("Executing template")

	//homeData.Body = template.HTMLEscapeString(homeData.Body)
	// homeData.Title = "Example"
	homeTempl.Execute(resp, homeData)

}

func testFunc(resp http.ResponseWriter, req *http.Request) {
	servLoggerINFO.Println("Loading test webpageURL")
	// http.ServeFile(resp, req, "logs/webServ.log")

	testTempl, err := template.ParseFiles("test.html")
	if err != nil {
		log.Fatal(err)
	}

	testTempl.Execute(resp, nil)
}
