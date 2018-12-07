package main

import (
	"strconv"
	"time"

	svc "rad-backend/server"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(Cors())
	router.POST("/radApp", getData)

	router.Run(":8080")
}

func getData(c *gin.Context) {
	chHost := "127.0.0.1:6379"

	cacher := svc.NewRedisCache(chHost)

	auctioneer := svc.NewAuctioneer(cacher)
	service := svc.NewForecast(cacher, auctioneer)
	// t := time.Date(2018, 11, 07, 10, 10, 10, 10, time.UTC)
	// fmt.Println("t :", t)
	// -6.246110, 106.801120
	latitude, _ := strconv.ParseFloat((c.PostForm("latitude")), 64)
	longitude, _ := strconv.ParseFloat((c.PostForm("longitude")), 64)
	if data, err := service.Forecast(latitude, longitude, time.Now()); err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, data)
	}

}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
