package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nallown/animetwist-api/db"
)

// Migrate gets all animes from mysql database and add them to elasticsearch.
func Migrate(instance *gorm.DB) {
	animes := db.FindAllAnimes(instance)

	exist, err := Client.IndexExists("animes").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exist {
		result, err := Client.CreateIndex("animes").Do(context.Background())
		if err != nil {
			panic(err)
		}

		if !result.Acknowledged {
			panic(errors.New("acknowledged should be true but returned false"))
		}
	}

	for _, anime := range animes {
		_, err := Client.Index().
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
