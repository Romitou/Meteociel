package meteociel

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

type ForecastType string

const (
	ForecastGFS       ForecastType = "https://www.meteociel.fr/previsions/{stationId}/{station}.htm"
	ForecastWRF       ForecastType = "https://www.meteociel.fr/previsions-wrf/{stationId}/{station}.htm"
	ForecastWRF1H     ForecastType = "https://www.meteociel.fr/previsions-wrf-1h/{stationId}/{station}.htm"
	ForecastAROME     ForecastType = "https://www.meteociel.fr/previsions-arome/{stationId}/{station}.htm"
	ForecastAROME1H   ForecastType = "https://www.meteociel.fr/previsions-arome-1h/{stationId}/{station}.htm"
	ForecastARPEGE1H  ForecastType = "https://www.meteociel.fr/previsions-arpege-1h/{stationId}/{station}.htm"
	ForecastICONEU    ForecastType = "https://www.meteociel.fr/previsions-iconeu/{stationId}/{station}.htm"
	ForecastICOND2    ForecastType = "https://www.meteociel.fr/previsions-icond2/{stationId}/{station}.htm"
	ForecastTrends10J ForecastType = "https://www.meteociel.fr/tendances/{stationId}/{station}.htm"
)

type WeatherType struct {
	Name string
}

var WeatherTypes = map[string]WeatherType{
	"soleil": {
		Name: "Sunny",
	},
	"voile": {
		Name: "Little cloudy",
	},
	"peu_nuageux": {
		Name: "Partly cloudy",
	},
	"mitige": {
		Name: "Mixed",
	},
	"nuageux": {
		Name: "Cloudy",
	},
	"brouillard": {
		Name: "Foggy",
	},
	"pluie": {
		Name: "Rainy",
	},
	"grele": {
		Name: "Hail",
	},
	"neige": {
		Name: "Snowy",
	},
	"averse_pluiefaible": {
		Name: "Light rain shower",
	},
	"averse_pluie": {
		Name: "Rain shower",
	},
	"averse_neige": {
		Name: "Snow shower",
	},
	"averse_orage": {
		Name: "Thunderstorm",
	},
	"averse_pluieneige": {
		Name: "Rain and snow shower",
	},
	"pluie_neige": {
		Name: "Rain and snow",
	},
	"oragefaible": {
		Name: "Major thunderstorm",
	},
}

var MeteocielDays = []string{"Lun", "Mar", "Mer", "Jeu", "Ven", "Sam", "Dim"}

const MeteocielHourFormat = "15:04"

type MeteocielForecast struct {
	Time          time.Time
	Temperature   int8
	WindDirection int16

	WindSpeed int8
	WindGust  int8

	Rainfall float32
	Humidity int8

	Pressure int16
	Weather  WeatherType
}

// GetForecast This method returns the forecast for a given station. Depending on the type of forecast you have chosen,
// the data may be different and some elements may not be present.
func (client MeteocielClient) GetForecast(forecast ForecastType, station MeteocielStation) (forecasts []MeteocielForecast, err error) {
	endpoint := strings.Replace(string(forecast), "{stationId}", station.ID, 1)
	endpoint = strings.Replace(endpoint, "{station}", station.Name, 1)

	reader, err := client.makeRequest(endpoint)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return
	}

	selection := doc.Find("td.Style1 center table tbody tr td table tbody").First().Find("tr[bgcolor='#CCFFFF'], tr[bgcolor='#DDEEFF']")

	var currentDate time.Time

	selection.Each(func(i int, tr *goquery.Selection) {
		dateSpan := tr.Find("td[rowspan]").First()
		if dateSpan.Size() == 1 {
			dayText := dateSpan.Text()[3:]

			var dayNumber int
			dayNumber, err = strconv.Atoi(dayText)
			if err != nil {
				return
			}

			currentDate = time.Date(time.Now().Year(), time.Now().Month(), dayNumber, 0, 0, 0, 0, time.Now().Location())
		}

		date := tr.Find("td").First().Text()

		var forecastDate time.Time
		forecastDate, err = time.Parse(MeteocielHourFormat, date)
		if err != nil {
			return
		}

		currentDate = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), forecastDate.Hour(), forecastDate.Minute(), 0, 0, time.Now().Location())

		forecast := MeteocielForecast{
			Time: currentDate,
		}

		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			switch i {
			case 1:
				withoutCelsius := strings.Replace(td.Text(), " \xb0C", "", 1)
				var converted int64
				converted, err = strconv.ParseInt(withoutCelsius, 10, 8)
				if err != nil {
					return
				}
				forecast.Temperature = int8(converted)
			case 3:
				imageTitle, exists := td.Find("img").First().Attr("title")
				if !exists {
					log.Println("No image title found")
					return
				}
				windDirection := strings.Split(imageTitle, " : ")[1]
				withoutDegree := strings.Replace(windDirection, " \xb0", "", 1)
				var converted int64
				converted, err = strconv.ParseInt(withoutDegree, 10, 16)
				if err != nil {
					return
				}
				forecast.WindDirection = int16(converted)
			case 4:
				var converted int64
				converted, err = strconv.ParseInt(td.Text(), 10, 8)
				if err != nil {
					return
				}
				forecast.WindSpeed = int8(converted)
			case 5:
				var converted int64
				converted, err = strconv.ParseInt(td.Text(), 10, 8)
				if err != nil {
					return
				}
				forecast.WindGust = int8(converted)
			case 6:
				if td.Text() == "--" {
					return
				}

				withoutMm := strings.Replace(td.Text(), " mm", "", 1)

				var converted float64
				converted, err = strconv.ParseFloat(withoutMm, 32)
				if err != nil {
					return
				}
				forecast.Rainfall = float32(converted)
			case 7:
				withoutPercent := strings.Replace(td.Text(), " %", "", 1)
				var converted int64
				converted, err = strconv.ParseInt(withoutPercent, 10, 8)
				if err != nil {
					return
				}
				forecast.Humidity = int8(converted)
			case 8:
				withoutHpa := strings.Replace(td.Text(), " hPa", "", 1)
				var converted int64
				converted, err = strconv.ParseInt(withoutHpa, 10, 16)
				if err != nil {
					return
				}
				forecast.Pressure = int16(converted)
			case 9:
				imageSrc, exists := td.Find("img").First().Attr("src")
				if !exists {
					log.Println("No image src found")
					return
				}
				srcSplit := strings.Split(imageSrc, "/")
				imageName := srcSplit[len(srcSplit)-1]
				imageNameWithoutExtension := strings.Replace(imageName, ".gif", "", 1)
				forecast.Weather = WeatherTypes[imageNameWithoutExtension]
			}
		})

		forecasts = append(forecasts, forecast)
	})

	return
}
