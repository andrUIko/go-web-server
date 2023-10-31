package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	dotenv "github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
	database "github.com/user/goforecast/db"
	"github.com/user/goforecast/domain"
	"github.com/user/goforecast/services"
)

func main() {
	if err := dotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUrl := "libsql://" + os.Getenv("DATABASE_NAME") + ".turso.io?authToken=" + os.Getenv("DATABASE_ACCESS_TOKEN")
	dbClient := database.CreateDBClient(dbUrl)

	r := gin.Default()

	r.LoadHTMLGlob("views/*")

	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
		latLong, err := domain.GetLatLong(dbClient, city)
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
	})

	r.GET("/stats", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTH_LOGIN"): os.Getenv("AUTH_PASSWORD"),
	}), func(c *gin.Context) {
		cities, err := dbClient.GetLastCities()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "stats.html", cities)
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run()
}
