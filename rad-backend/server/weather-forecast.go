package server

import (
	"fmt"
	"log"
	"math"
	"time"
)

//Round ..
func Round(x, unit float64) float64 {
	if x < 0 {
		unit = unit * (-1)
	}
	return math.Ceil(x/unit) * unit
}

var apiKey = "b9b851872af95daec834a07df35632d5"

type forecast struct {
	cacher     Cacher
	auctioneer Auctioneer
}

//NewForecast return new dispatch service
func NewForecast(cacher Cacher, auctioneer Auctioneer) ForecastService {
	return &forecast{cacher: cacher, auctioneer: auctioneer}
}

func (f *forecast) Forecast(lat float64, lon float64, inputTime time.Time) (ForecastResponse, error) {
	///->get from Redis

	// lat = Round(lat, 0.05)
	// lon = Round(lon, 0.05)

	//ngubah waktu dari WIB jadi UTC
	utc := inputTime.UTC()
	jam := utc.Hour() / 3
	jam = jam * 3
	waktu := fmt.Sprint(utc.Year(), int(utc.Month()), utc.Day(), jam, utc.Minute())

	// fmt.Println(inputTime, waktu)
	input, _ := time.Parse("2006 01 02 15 04", waktu)
	///fmt.Println("ad request", lat, lon, input.Round(time.Duration(time.Hour)))

	a := f.auctioneer.AddCurrent(lat, lon, input)
	fmt.Println("Current Weather Data :", a)

	// f.auctioneer.AddForecast(lat, lon, input)

	res, err := f.cacher.GetForecast5(lat, lon, input)
	fmt.Println("Weather Data from Cache :", res)
	if err != nil {
		log.Println(err)
	}

	// if res.RetrieveTime == (time.Time{}) || res.RetrieveTime.Sub(time.Now()) > time.Duration(time.Hour*6) {
	// 	if input.Sub(time.Now()) < time.Duration(time.Hour*3) {

	// 		res = f.auctioneer.AddCurrent(lat, lon, input)

	// 	}
	// 	//------------>redis not found
	// 	f.auctioneer.AddForecast(lat, lon, input)
	// 	res, err = f.cacher.GetForecast5(lat, lon, input)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

	return res, nil
}
