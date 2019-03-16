package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nallown/animetwist-api/db"
	"github.com/olivere/elastic"
)

type ElasticAnime struct {
	Title    string `json:"title"`
	AltTitle string `json:"alt_title"`
	Season   int    `json:"season"`
	Ongoing  int    `json:"ongoing"`
	Slug     string `json:"slug"`
	SlugID   uint   `json:"slug_id"`
	HbID     int    `json:"hb_id"`
	MalID    int    `json:"mal_id"`
}

// Migrate gets all animes from mysql database and add them to elasticsearch.
func Migrate(instance *gorm.DB, client *elastic.Client) {
	animes := db.FindAllAnimes(instance)

	exist, err := client.IndexExists("animes").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exist {
		body := `
		{
			"settings" : {
				"analysis" : {
					"analyzer" : {
						"default" : {
							"tokenizer" : "standard",
								"filter" : ["asciifolding", "lowercase"]
						}
					}
				}
			}
		}`

		result, err := client.CreateIndex("animes").BodyString(body).Do(context.Background())
		if err != nil {
			panic(err)
		}

		if !result.Acknowledged {
			panic(errors.New("acknowledged should be true but returned false"))
		}
	} else {
		resp, err := client.DeleteIndex("animes").Do(context.Background())
		if err != nil {
			panic(err)
		}

		if !resp.Acknowledged {
			panic(errors.New("acknowledged should be true but returned false"))
		}

		body := `
		{
			"settings" : {
				"analysis" : {
					"analyzer" : {
						"default" : {
							"tokenizer" : "standard",
								"filter" : ["asciifolding", "lowercase"]
						}
					}
				}
			}
		}`

		result, err := client.CreateIndex("animes").BodyString(body).Do(context.Background())
		if err != nil {
			panic(err)
		}

		if !result.Acknowledged {
			panic(errors.New("acknowledged should be true but returned false"))
		}

	}

	bulk := client.Bulk().Index("animes").Type("_doc")
	for _, anime := range animes {
		bulk.Add(elastic.NewBulkIndexRequest().Index("animes").Type("_doc").Id(fmt.Sprint(anime.ID)).
			Doc(&ElasticAnime{
				Title:    anime.Title,
				AltTitle: anime.AltTitle,
				Season:   anime.Season,
				Ongoing:  anime.Ongoing,
				Slug:     anime.Slug.Slug,
				SlugID:   anime.Slug.ID,
				HbID:     anime.HbID,
				MalID:    anime.MalID,
			}))

		if err != nil {
			panic(err)
		}
	}

	response, err := bulk.Refresh("true").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if response.Errors {
		panic("bulk response returned an error or more")
	}
}
