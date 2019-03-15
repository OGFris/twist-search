package search

import (
	"bitbucket.org/ascendile/doctors-addresses/utils"
	"encoding/json"
	"github.com/OGFris/twist-search/elasticsearch"
	"github.com/nallown/animetwist-api/db"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	// NOTICE: Start the server before running this test.
	var (
		anime        db.Anime
		outputAnimes []elasticsearch.ElasticAnime
	)

	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}

	db.Instance.Find(&anime, &db.Anime{Hidden: 0})

	urlReplace := strings.NewReplacer(" ", "%20")

	resp, err := http.Get("http://localhost:8080/api/search?q=" + urlReplace.Replace(anime.Title))
	if err != nil {
		panic(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	utils.Assert(t, resp.StatusCode, http.StatusOK)

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &outputAnimes)
	if err != nil {
		panic(err)
	}

	utils.Assert(t, len(outputAnimes), 1)
	utils.Assert(t, outputAnimes[0].Title, anime.Title)
	utils.Assert(t, outputAnimes[0].AltTitle, anime.AltTitle)
	utils.Assert(t, outputAnimes[0].HbID, anime.HbID)
	utils.Assert(t, outputAnimes[0].MalID, anime.MalID)
	utils.Assert(t, outputAnimes[0].Ongoing, anime.Ongoing)
	utils.Assert(t, outputAnimes[0].Season, anime.Season)
}
