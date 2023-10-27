package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/gin-gonic/gin"
	dotenv "github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
)

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

	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
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

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run()
}
