package elasticsearch

import (
	"github.com/nallown/animetwist-api/db"
	"os"
	"testing"
)

func TestMigrate(t *testing.T) {
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}

	Migrate(db.Instance, os.Getenv("URL"), os.Getenv("CLIENT_USERNAME"), os.Getenv("CLIENT_PASSWORD"))
}
