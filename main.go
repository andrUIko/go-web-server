package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sql "database/sql"

	gin "github.com/gin-gonic/gin"
	dotenv "github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
)

type LatLong struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeoResponse struct {
	Results []LatLong `json:"results"`
}

type Forecast struct {
	Date        string
	Temperature string
}

type WeatherDisplay struct {
	City      string
	Forecasts []Forecast
}

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Hourly    struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

func extractWeatherData(city string, rawWeather string) (WeatherDisplay, error) {
	var weatherResponse WeatherResponse
	if err := json.Unmarshal([]byte(rawWeather), &weatherResponse); err != nil {
		return WeatherDisplay{}, err
	}

	var forecasts []Forecast
	for i, t := range weatherResponse.Hourly.Time {
		data, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			return WeatherDisplay{}, err
		}
		forecast := Forecast{
			Date:        data.Format("Mon 15:04"),
			Temperature: fmt.Sprintf("%.1f°C", weatherResponse.Hourly.Temperature2m[i]),
		}
		forecasts = append(forecasts, forecast)
	}
	return WeatherDisplay{
		City:      city,
		Forecasts: forecasts,
	}, nil
}

func generateResponse(city string, c *gin.Context, db *sql.DB) {
	latLong, err := getLatLong(db, city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	weather, err := getWeather(*latLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	weatherDisplay, err := extractWeatherData(city, weather)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "weather.html", weatherDisplay)
}

func main() {

	if err := dotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	accessToken, dbName := os.Getenv("DATABASE_ACCESS_TOKEN"), os.Getenv("DATABASE_NAME")
	dbUrl := "libsql://" + dbName + ".turso.io?authToken=" + accessToken
	db, err := sql.Open("libsql", dbUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}

	r := gin.Default()
	r.LoadHTMLGlob("views/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
		generateResponse(city, c, db)
	})

	r.GET("/weather/:city", func(c *gin.Context) {
		city := c.Param("city")
		generateResponse(city, c, db)
	})

	r.GET("/stats", gin.BasicAuth(gin.Accounts{
		os.Getenv("AUTH_LOGIN"): os.Getenv("AUTH_PASSWORD"),
	}), func(c *gin.Context) {
		cities, err := getLastCities(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "stats.html", cities)
	})

	r.Run()
}
