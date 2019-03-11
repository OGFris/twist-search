package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/OGFris/twist-search/elasticsearch"
	"github.com/olivere/elastic"
	"net/http"
)

func main() {
	url := flag.String("url", "", "elasticsearch url")
	username := flag.String("username", "", "elasticsearch username")
	password := flag.String("password", "", "elasticsearch password")
	flag.Parse()

	client, err := elasticsearch.NewClient(*url, *username, *password)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		results, err := client.Search().
			Index("animes").
			Query(elastic.NewMultiMatchQuery(r.PostFormValue("query"), "title", "alt_title")).
			Do(context.Background())

		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(w).Encode(results.Hits.Hits)
		if err != nil {
			panic(err)
		}
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
