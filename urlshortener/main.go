package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rprajapati0067/chapter7/urlshortener/helper"
	base62 "github.com/rprajapati0067/chapter7/urlshortener/utils"
)

type DBClient struct {
	db *sql.DB
}
type Record struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func (driver *DBClient) GenerateShortURL(w http.ResponseWriter, r *http.Request) {
	var id int
	var record Record

	postBody, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(postBody, &record)
	err = driver.db.QueryRow("INSERT INTO web_url(url) VALUES($1) RETURNING id", record.URL).Scan(&id)

	responseMap := map[string]string{"encoding_string": base62.ToBase62(id)}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(responseMap)
		w.Write(response)
	}

}
func (driver *DBClient) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	var url string
	vars := mux.Vars(r)

	id := base62.ToBase10(vars["encoded_string"])
	err := driver.db.QueryRow("SELECT url FROM web_url WHERE id = $1", id).Scan(&url)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		responseMap := map[string]interface{}{"url": url}
		response, _ := json.Marshal(responseMap)
		w.Write(response)
	}
}

func main() {

	db, err := helper.InitDB()
	if err != nil {
		log.Fatalln(err)
	}
	dbclient := &DBClient{db: db}
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := mux.NewRouter()

	r.HandleFunc("/v1/short/{encoded_string:[a-zA-Z0-9]*}", dbclient.GetOriginalURL).Methods("GET")
	r.HandleFunc("/v1/short", dbclient.GenerateShortURL).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	// x := 100

	// base62String := base62.ToBase62(x)
	// log.Println(base62String)

	// normalNumber := base62.ToBase10(base62String)
	// log.Println(normalNumber)

	// _, err := helper.InitDB()

	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println("Database table are successfully initilized")
}
