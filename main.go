package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LatLong struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeoResponse struct {
	Results []LatLong `json:"results"`
}

// func main() {
//   latLong, err := getLatLong("Lviv")
//   if err != nil {
//     log.Fatalf("Failed to get latitude and longitude: %s", err)
//   }
//   fmt.Printf("Latitude: %f, Longitude: %f\n", latLong.Latitude, latLong.Longitude)

//   weather, err := getWeather(*latLong)

//   if err != nil {
//     log.Fatalf("Failed to get weather: %s", err)
//   }

//   fmt.Printf("Weather: %s\n", weather)
// }

func generateResponse(city string, c *gin.Context) {
	latLong, err := getLatLong(city)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	weather, err := getWeather(*latLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weather": weather})
}

func main() {
	r := gin.Default()

	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
		generateResponse(city, c)
	})

	r.GET("/weather/:city", func(c *gin.Context) {
		city := c.Param("city")
		generateResponse(city, c)
	})

	r.Run()
}
