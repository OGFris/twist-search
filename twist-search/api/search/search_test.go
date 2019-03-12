package search

import (
	"net/http"
	"testing"
)

func TestSearch(t *testing.T) {
	// NOTICE: Start the server before running this test.
	resp, err := http.Get("http://localhost:8080/api/search?q=test")
	if err != nil {
		panic(err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()
}
