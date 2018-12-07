package server

import (
	"time"

	owm "github.com/briandowns/openweathermap"
)

const (
	//ServiceID ..
	ServiceID = "forecast.bluebird.id"
)

//ForecastRequest ..
type ForecastRequest struct {
	Latitude  float64
	Longitude float64
	InputTime time.Time
}

//ForecastResponse ..
type ForecastResponse struct {
	Name    string
	Main    Main
	Weather Weather
	Clouds  Clouds
	Wind    Wind
	// RetrieveTime time.Time
}

// Wind struct contains the speed and degree of the wind.
type Wind struct {
	Speed float64
	Deg   float64
}

// Weather struct holds high-level, basic info on the returned
// data.
type Weather struct {
	ID          int64
	Main        string
	Description string
	Icon        string
}

// Main struct contains the temperates, humidity, pressure for the request.
type Main struct {
	Temp      float64
	TempMin   float64
	TempMax   float64
	Pressure  float64
	SeaLevel  float64
	GrndLevel float64
	Humidity  int64
}

//Rain ..
type Rain struct {
	Rain1h  string
	Rain3h  string
	Rain24h string
}

//Snow ..
type Snow struct {
	Snow1h  string
	Snow3h  string
	Snow24h string
}

// Clouds struct holds data regarding cloud cover.
type Clouds struct {
	All int64
}

//Cacher is data cacher
type Cacher interface {
	// GetForecast5(float64, float64, time.Time) (Forecast, error)
	AddForecast5(forecast *owm.Forecast5WeatherData, lat float64, lon float64, input time.Time)
	AddCurrentWeather(current *owm.CurrentWeatherData, lat float64, lon float64, input time.Time) ForecastResponse
	GetForecast5(lat float64, lon float64, time time.Time) (ForecastResponse, error)
}

//Auctioneer ..
type Auctioneer interface {
	AddForecast(float64, float64, time.Time)
	AddCurrent(float64, float64, time.Time) ForecastResponse
}

//ForecastService ..
type ForecastService interface {
	Forecast(float64, float64, time.Time) (ForecastResponse, error)
}

//BigQueryWriter ..
type BigQueryWriter interface {
	AddLocation() error
}

//ReadWriter is data reader/writer
// type ReadWriter interface {
// 	WriteForecast(owm.Forecast5WeatherData) error
// }

//Json Data

//Locations ..
type Locations struct {
	Locations []Location `json:"locations"`
}

//Location ..
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

//WeatherLoaders ..
type WeatherLoaders struct {
	Dt          int
	DtTime      time.Time
	CityID      string
	CityName    string
	Lat         string
	Lon         string
	Temp        float64
	TempMin     float64
	TempMax     float64
	Pressure    float64
	SeaLevel    float64
	GrndLevel   float64
	Humidity    int64
	Speed       float64
	Deg         float64
	Rain1h      string
	Rain3h      string
	Rain24h     string
	Snow1h      string
	Snow3h      string
	Snow24h     string
	Clouds      string
	ID          int64
	Main        string
	Description string
	Icon        string
}

// {"Serang": {
//     "id": "file",
//     "value": "File",
//     "popup": {
//       "menuitem": [
//         {"value": "New", "onclick": "CreateNewDoc()"},
//         {"value": "Open", "onclick": "OpenDoc()"},
//         {"value": "Close", "onclick": "CloseDoc()"}
//       ]
//     }
//   }}
