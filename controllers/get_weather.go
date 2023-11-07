package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/goforecast/db"
	"github.com/user/goforecast/domain"
	"github.com/user/goforecast/services"
)

func GetWeather(c *gin.Context) {
	city := c.Query("city")
	DB, ok := c.Keys["DB"].(db.DBPool)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no db client found"})
	}

	latLong, err := DB.GetLatLong(city)

	if err != nil {

		if latLong, err = services.FetchLatLong(city); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = DB.InsertCity(city, latLong)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	weather, err := services.GetWeather(*latLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	weatherDisplay, err := domain.ExtractWeatherData(city, weather)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "weather.html", weatherDisplay)
}
