package search

import (
	"encoding/json"
	"github.com/OGFris/twist-search/elasticsearch"
	"github.com/nallown/animetwist-api/db"
	"github.com/olivere/elastic"
	"net/http"
	"strings"
)

// Search is the api route function for /api/search which search for animes using the given query on elasticsearch.
func Search(w http.ResponseWriter, r *http.Request, client *elastic.Client) {
	q := r.URL.Query().Get("q")

	if client != nil {
		// By default, if the client is available, it uses elasticsearch.
		results, err := client.Search().
			Index("animes").
			Query(elastic.NewMultiMatchQuery(q, "title", "alt_title")).
			Do(r.Context())

		if err != nil {
			panic(err)
		}

		var sources []elasticsearch.ElasticAnime
		for _, source := range results.Hits.Hits {
			var anime elasticsearch.ElasticAnime
			bytes, err := source.Source.MarshalJSON()
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bytes, &anime)
			if err != nil {
				panic(err)
			}

			sources = append(sources, anime)
		}

		err = json.NewEncoder(w).Encode(sources)
		if err != nil {
			panic(err)
		}
	} else {
		// If the client isn't passed, it'd use the database instead.
		var (
			results []elasticsearch.ElasticAnime
			animes  []db.Anime
		)

		db.Instance.Preload("Slug").Find(&animes, &db.Anime{Hidden: 0})

		for _, anime := range animes {
			if strings.Contains(anime.Title, q) || strings.Contains(anime.AltTitle, q) {
				results = append(results, elasticsearch.ElasticAnime{
					Title:    anime.Title,
					AltTitle: anime.AltTitle,
					Season:   anime.Season,
					Ongoing:  anime.Ongoing,
					Slug:     anime.Slug.Slug,
					SlugID:   anime.Slug.ID,
					HbID:     anime.HbID,
					MalID:    anime.MalID,
				})
			}
		}
	}
}
