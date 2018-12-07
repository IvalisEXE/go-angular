package server

import (
	"strconv"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/go-redis/redis"
)

const (
	keyForecast      = "Weather_"
	keyTemp          = "Temperature:"
	keyWeath         = "Weather:"
	keyClouds        = "Clouds:"
	keyWind          = "Wind:"
	fieldName        = "Name:"
	fieldTemp        = "Temp"
	fieldTempMin     = "TempMin"
	fieldTempMax     = "TempMax"
	fieldPressure    = "Pressure"
	fiieldSeaLevel   = "SeaLevel"
	fieldGrndLevel   = "GrndLevel"
	fieldHumidity    = "Humidity"
	fieldTempKf      = "TempKf"
	fieldWeatherID   = "WeatherID"
	fieldCondition   = "Condition"
	fieldDescription = "Description"
	fieldIcon        = "Icon"
	fieldAll         = "All"
	fieldSpeed       = "Speed"
	fieldDeg         = "Deg"
	fieldDtTxt       = "DtTxt"
	fieldRetriveTime = "Retrieve_Time"
)

type redisCache struct {
	client *redis.Client
}

//NewRedisCache returns new redis cacher
func NewRedisCache(url string) Cacher {
	return &redisCache{client: redis.NewClient(&redis.Options{Addr: url})}
}

func (r *redisCache) AddCurrentWeather(current *owm.CurrentWeatherData, lat float64, lon float64, insertTime time.Time) ForecastResponse {

	//fmt.Println("waktu sekarang", insertTime)
	pipe := r.client.Pipeline()
	mapCurrent := make(map[string]interface{})
	round := time.Duration(time.Hour * 3)
	waktubaru := insertTime.Round(round)
	//fmt.Println("waktubaru", waktubaru, round)
	datetime := waktubaru.Format("20060102_15")
	pipe.RPush(keyForecast + datetime + "_" + strconv.FormatFloat(lat, 'G', -1, 64) +
		"_" + strconv.FormatFloat(lon, 'G', -1, 64))

	mapCurrent[fieldTemp] = strconv.FormatFloat((current.Main.Temp), 'G', -1, 64)
	mapCurrent[fieldTempMin] = strconv.FormatFloat((current.Main.TempMin), 'G', -1, 64)
	mapCurrent[fieldTempMax] = strconv.FormatFloat((current.Main.TempMax), 'G', -1, 64)
	mapCurrent[fieldPressure] = strconv.FormatFloat((current.Main.Pressure), 'G', -1, 64)
	mapCurrent[fiieldSeaLevel] = strconv.FormatFloat((current.Main.SeaLevel), 'G', -1, 64)
	mapCurrent[fieldGrndLevel] = strconv.FormatFloat((current.Main.GrndLevel), 'G', -1, 64)
	mapCurrent[fieldHumidity] = strconv.FormatInt(int64(current.Main.Humidity), 10)

	for _, j := range current.Weather {
		mapCurrent[fieldWeatherID] = strconv.FormatInt(int64(j.ID), 10)
		mapCurrent[fieldCondition] = j.Main
		mapCurrent[fieldDescription] = j.Description
		mapCurrent[fieldIcon] = j.Icon
	}

	mapCurrent[fieldAll] = strconv.FormatInt(int64(current.Clouds.All), 10)

	mapCurrent[fieldSpeed] = strconv.FormatFloat((current.Wind.Speed), 'G', -1, 64)
	mapCurrent[fieldDeg] = strconv.FormatFloat((current.Wind.Deg), 'G', -1, 64)
	// mapCurrent[fieldRetriveTime] = datetime
	mapCurrent[fieldName] = current.Name
	pipe.HMSet(keyForecast+datetime+"_"+strconv.FormatFloat(lat, 'G', -1, 64)+
		"_"+strconv.FormatFloat(lon, 'G', -1, 64), mapCurrent)
	pipe.Exec()

	var res ForecastResponse
	res.Main.Temp = current.Main.Temp
	res.Main.TempMin = current.Main.TempMin
	res.Main.TempMax = current.Main.TempMax
	res.Main.Pressure = current.Main.Pressure
	res.Main.GrndLevel = current.Main.GrndLevel
	res.Main.Humidity = int64(current.Main.Humidity)

	for _, j := range current.Weather {
		res.Weather.ID = int64(j.ID)
		res.Weather.Main = j.Main
		res.Weather.Description = j.Description
		res.Weather.Icon = j.Icon
	}
	res.Clouds.All = int64(current.Clouds.All)

	res.Wind.Speed = current.Wind.Speed
	res.Wind.Deg = current.Wind.Deg
	// res.RetrieveTime = insertTime
	res.Name = current.Name

	return res
}

func (r *redisCache) AddForecast5(forecast *owm.Forecast5WeatherData, lat float64, lon float64, insertTime time.Time) {
	// fmt.Println("MASUK ADDFORECAST5 REDIS.GO")
	// fmt.Println(forecast)

	insertTime.UTC()
	pipe := r.client.Pipeline()

	mapForecast := make(map[string]interface{})
	for _, i := range forecast.List {
		datetime := i.DtTxt.Format("20060102_15")
		pipe.RPush(keyForecast + datetime + "_" + strconv.FormatFloat(lat, 'G', -1, 64) +
			"_" + strconv.FormatFloat(lon, 'G', -1, 64))

		mapForecast[fieldTemp] = strconv.FormatFloat((i.Main.Temp), 'G', -1, 64)
		mapForecast[fieldTempMin] = strconv.FormatFloat((i.Main.TempMin), 'G', -1, 64)
		mapForecast[fieldTempMax] = strconv.FormatFloat((i.Main.TempMax), 'G', -1, 64)
		mapForecast[fieldPressure] = strconv.FormatFloat((i.Main.Pressure), 'G', -1, 64)
		mapForecast[fiieldSeaLevel] = strconv.FormatFloat((i.Main.SeaLevel), 'G', -1, 64)
		mapForecast[fieldGrndLevel] = strconv.FormatFloat((i.Main.GrndLevel), 'G', -1, 64)
		mapForecast[fieldHumidity] = strconv.FormatInt(int64(i.Main.Humidity), 10)

		for _, j := range i.Weather {
			mapForecast[fieldWeatherID] = strconv.FormatInt(int64(j.ID), 10)
			mapForecast[fieldCondition] = j.Main
			mapForecast[fieldDescription] = j.Description
			mapForecast[fieldIcon] = j.Icon
		}

		mapForecast[fieldAll] = strconv.FormatInt(int64(i.Clouds.All), 10)

		mapForecast[fieldSpeed] = strconv.FormatFloat((i.Wind.Speed), 'G', -1, 64)
		mapForecast[fieldDeg] = strconv.FormatFloat((i.Wind.Deg), 'G', -1, 64)
		mapForecast[fieldDtTxt] = i.DtTxt.String()
		// mapForecast[fieldRetriveTime] = insertTime.String()

		pipe.HMSet(keyForecast+datetime+"_"+strconv.FormatFloat(lat, 'G', -1, 64)+
			"_"+strconv.FormatFloat(lon, 'G', -1, 64), mapForecast)
		pipe.Exec()
	}
}

func (r *redisCache) GetForecast5(lat float64, lon float64, keytime time.Time) (ForecastResponse, error) {
	var res ForecastResponse
	// layout := "20060102_15"
	// keytime.UTC()

	round := time.Duration(time.Hour * 3)
	keytime = keytime.Round(round)
	inputTime := keytime.Format("20060102_15")
	//fmt.Println("keytime", inputTime)
	mapForecast, err := r.client.HGetAll(keyForecast + inputTime + "_" + strconv.FormatFloat(lat, 'G', -1, 64) +
		"_" + strconv.FormatFloat(lon, 'G', -1, 64)).Result()
	// fmt.Println("MF :", keyForecast+inputTime+"_"+strconv.FormatFloat(lat, 'G', -1, 64)+
	// 	"_"+strconv.FormatFloat(lon, 'G', -1, 64))
	if err == nil {
		res.Main.Temp, _ = strconv.ParseFloat(mapForecast[fieldTemp], 64)
		res.Main.TempMin, _ = strconv.ParseFloat(mapForecast[fieldTempMin], 64)
		res.Main.TempMax, _ = strconv.ParseFloat(mapForecast[fieldTempMax], 64)
		res.Main.Pressure, _ = strconv.ParseFloat(mapForecast[fieldPressure], 64)
		res.Main.GrndLevel, _ = strconv.ParseFloat(mapForecast[fieldGrndLevel], 64)
		res.Main.Humidity, _ = strconv.ParseInt(mapForecast[fieldHumidity], 10, 64)

		res.Weather.ID, _ = strconv.ParseInt(mapForecast[fieldWeatherID], 10, 64)
		res.Weather.Main, _ = mapForecast[fieldCondition]
		res.Weather.Description, _ = mapForecast[fieldDescription]
		res.Weather.Icon, _ = mapForecast[fieldIcon]

		res.Clouds.All, _ = strconv.ParseInt(mapForecast[fieldAll], 10, 64)

		res.Wind.Speed, _ = strconv.ParseFloat(mapForecast[fieldSpeed], 64)
		res.Wind.Deg, _ = strconv.ParseFloat(mapForecast[fieldDeg], 64)
		res.Name, _ = mapForecast[fieldName]
		// res.RetrieveTime, _ = time.Parse(layout, mapForecast[fieldRetriveTime])

	}

	return res, err
}
