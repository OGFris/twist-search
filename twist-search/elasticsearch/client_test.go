package elasticsearch

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Couldn't load the .env file")
		panic(err)
	}

	err = NewClient(os.Getenv("URL"), os.Getenv("CLIENT_USERNAME"), os.Getenv("CLIENT_PASSWORD"))
	if err != nil {
		panic(err)
	}

	_, _, err = Client.Ping(os.Getenv("URL")).Do(context.Background())
	if err != nil {
		panic(err)
	}
}
