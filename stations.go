package meteociel

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type StationEndpoint string

const (
	StationSearch StationEndpoint = "https://www.meteociel.fr/prevville.php?action=getville&villeid=&ville={city}&envoyer=OK"
)

// MeteocielStation This struct represents a Meteociel station.
// You can use it to get forecasts for a given station.
type MeteocielStation struct {
	ID   string
	Name string
}

// GetStationForCity This method returns the Meteociel station for a given city.
// You MUST provide the exact ZIP code or the name of the city, otherwise the method will return an error.
func (client MeteocielClient) GetStationForCity(exactName string) (station MeteocielStation, err error) {
	endpoint := strings.Replace(string(StationSearch), "{city}", exactName, 1)
	reader, err := client.makeRequest(endpoint)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}

	selection := doc.Find("body table tbody tr td table tbody tr td p script")
	if selection.Length() != 1 {
		err = errors.New("no station found")
		return
	}

	split := strings.Split(selection.Text(), "/")
	stationID := split[2]
	stationName := strings.Replace(split[3], ".htm';", "", 1)

	return MeteocielStation{
		ID:   stationID,
		Name: stationName,
	}, nil
}
