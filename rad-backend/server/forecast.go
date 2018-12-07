package server

import (
	"fmt"
	"log"
	"time"

	owm "github.com/briandowns/openweathermap"
)

type auctioneer struct {
	// readWriter ReadWriter
	cacher Cacher
}

//NewAuctioneer returns new Auctioneer
func NewAuctioneer(cacher Cacher) Auctioneer {
	return &auctioneer{cacher: cacher}
}

func (a *auctioneer) AddForecast(lat float64, lon float64, inputTime time.Time) {
	w, err := owm.NewForecast("5", "C", "EN", apiKey)
	if err != nil {
		log.Println(err)
	}

	w.DailyByCoordinates(&owm.Coordinates{Latitude: lat, Longitude: lon}, 8)

	// fmt.Printf("%+v \n ", w.ForecastWeatherJson)
	data := w.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	fmt.Println(data)

	a.cacher.AddForecast5(data, lat, lon, inputTime)
	// a.readWriter.WriteForecast(forecast)
}

func (a *auctioneer) AddCurrent(lat float64, lon float64, inputTime time.Time) ForecastResponse {
	c, err := owm.NewCurrent("C", "EN", apiKey)

	//	fmt.Println("current summoned")
	if err != nil {
		log.Println(err)
	}
	c.CurrentByCoordinates(&owm.Coordinates{Latitude: lat, Longitude: lon})
	res := a.cacher.AddCurrentWeather(c, lat, lon, inputTime)
	return res
}
