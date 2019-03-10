package elasticsearch

import (
	"github.com/nallown/animetwist-api/db"
	"testing"
)

func TestMigrate(t *testing.T) {
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}

	Migrate(db.Instance)
}
