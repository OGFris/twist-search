package elasticsearch

import (
	"github.com/olivere/elastic"
	"net/http"
)

var Client *elastic.Client

type Transport struct {
	Username string
	Password string
}

// RoundTrip implementation used to login.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultClient.Do(req)
}

// NewClient creates a new client to the variable Client.
func NewClient(url, username, password string) error {
	var err error

	Client, err = elastic.NewClient(
		elastic.SetHttpClient(&http.Client{Transport: &Transport{
			Username: username,
			Password: password,
		}}),
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHealthcheckTimeoutStartup(0),
	)

	return err
}
