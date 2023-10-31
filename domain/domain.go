package domain

import (
	"encoding/json"
	"fmt"
	"time"

	database "github.com/user/goforecast/db"
	"github.com/user/goforecast/models"
	"github.com/user/goforecast/services"
)

func GetLatLong(dbClient database.DBPool, name string) (*models.LatLong, error) {
	if latLong, err := dbClient.GetLatLong(name); err == nil {
		return latLong, nil
	}

	latLong, err := services.FetchLatLong(name)
	if err != nil {
		return nil, err
	}

	err = dbClient.InsertCity(name, latLong)
	if err != nil {
		return nil, err
	}

	return latLong, nil
}

func ExtractWeatherData(city string, rawWeather string) (models.WeatherDisplay, error) {
	var weatherResponse models.WeatherResponse
	if err := json.Unmarshal([]byte(rawWeather), &weatherResponse); err != nil {
		return models.WeatherDisplay{}, err
	}

	var forecasts []models.Forecast
	for i, t := range weatherResponse.Hourly.Time {
		data, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			return models.WeatherDisplay{}, err
		}
		forecast := models.Forecast{
			Date:        data.Format("Mon 15:04"),
			Temperature: fmt.Sprintf("%.1f°C", weatherResponse.Hourly.Temperature2m[i]),
		}
		forecasts = append(forecasts, forecast)
	}
	return models.WeatherDisplay{
		City:      city,
		Forecasts: forecasts,
	}, nil
}
