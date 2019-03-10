package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nallown/animetwist-api/db"
)

// Migrate gets all animes from mysql database and add them to elasticsearch.
func Migrate(instance *gorm.DB, url, username, password string) {
	animes := db.FindAllAnimes(instance)
	client, err := NewClient(url, username, password)

	exist, err := client.IndexExists("animes").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exist {
		result, err := client.CreateIndex("animes").Do(context.Background())
		if err != nil {
			panic(err)
		}

		if !result.Acknowledged {
			panic(errors.New("acknowledged should be true but returned false"))
		}
	}

	for _, anime := range animes {
		_, err := client.Index().
			Index("animes").
			Type("_doc").
			Id(fmt.Sprint(anime.ID)).
			BodyJson(&anime).
			Refresh("true").
			Do(context.Background())

		if err != nil {
			panic(err)
		}
	}
}
