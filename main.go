package meteociel

import (
	"errors"
	"io"
	"log"
	"net/http"
)

// MeteocielClient This struct represents a Meteociel client.
type MeteocielClient struct {
	HttpClient http.Client
}

// CreateClient This method creates a new Meteociel client.
func CreateClient() *MeteocielClient {
	return &MeteocielClient{}
}

func (client MeteocielClient) makeRequest(endpoint string) (reader io.Reader, err error) {
	log.Println("Making request to: ", endpoint)
	response, err := client.HttpClient.Get(endpoint)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return nil, errors.New("bad http status code: " + response.Status)
	}

	reader = response.Body
	return
}
