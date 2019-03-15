package main

import (
	"flag"
	"github.com/OGFris/twist-search/api/search"
	"github.com/OGFris/twist-search/elasticsearch"
	"github.com/gorilla/mux"
	"github.com/nallown/animetwist-api/db"
	"net/http"
)

var (
	Username string
	Password string
	Url      string
)

func init() {
	flag.StringVar(&Url, "url", "", "elasticsearch url")
	flag.StringVar(&Username, "username", "", "elasticsearch username")
	flag.StringVar(&Password, "password", "", "elasticsearch password")
	flag.Parse()
}

func main() {
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}

	client, err := elasticsearch.NewClient(Url, Username, Password)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		search.Search(w, r, client)
	})

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
