# Meteociel

This repository and code was created as a result of a personal need to have access to high quality weather data and specifically forecasts for my projects. Having searched for a long time for a French weather API, I finally decided to create a library to get weather information from a well known French site, [meteociel.fr](https://www.meteociel.fr/). Note that the data is not retrieved via a standard API - not existing -, but via a scraping process to retrieve data from HTML pages.

## Examples

### Get weather forecast for a city

```go
package main

import (
	"fmt"
	"log"
	"strconv"
	
	"github.com/romitou/meteociel"
)

func main() {
	client := meteociel.CreateClient()
	station, err := client.GetStationForCity("75000")
	if err != nil {
		log.Fatal(err)
	}

	forecasts, err := client.GetForecast(meteociel.ForecastGFS, station)
	if err != nil {
		log.Fatal(err)
	}

	for _, forecast := range forecasts {
		log.Println("Here is the forecast for " + forecast.Time.String() + ":")
		log.Println("  - Weather: " + forecast.Weather.Name)
		log.Println("  - Temperature: " + strconv.Itoa(int(forecast.Temperature)) + "°C")
		log.Println("  - Wind direction: " + strconv.Itoa(int(forecast.WindDirection)) + "°")
		log.Println("  - Wind speed: " + strconv.Itoa(int(forecast.WindSpeed)) + " km/h")
		log.Println("  - Wind gust: " + strconv.Itoa(int(forecast.WindGust)) + " km/h")
		log.Println("  - Rainfall: " + fmt.Sprint(forecast.Rainfall) + " mm")
		log.Println("  - Humidity: " + strconv.Itoa(int(forecast.Humidity)) + "%")
		log.Println("  - Pressure: " + strconv.Itoa(int(forecast.Pressure)) + " hPa")

	}

}
```