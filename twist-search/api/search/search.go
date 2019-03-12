package search

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"net/http"
)

// Search is the api route function for /api/search which search for animes using the given query on elasticsearch.
func Search(w http.ResponseWriter, r *http.Request) {
	client := r.Context().Value("elastic_client").(*elastic.Client)

	results, err := client.Search().
		Index("animes").
		Query(elastic.NewMultiMatchQuery(r.URL.Query().Get("q"), "title", "alt_title")).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(results.Hits.Hits)
	if err != nil {
		panic(err)
	}
}
