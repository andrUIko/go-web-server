package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/user/goforecast/models"
)

func FetchLatLong(city string) (*models.LatLong, error) {
	endpoint := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(city))
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making request to Geo API: %w", err)
	}
	defer resp.Body.Close()

	var response models.GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response from Geo API: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, errors.New("no results found")
	}

	return &response.Results[0], nil
}

func GetWeather(latLong models.LatLong) (string, error) {
	endpoint := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&hourly=temperature_2m", latLong.Latitude, latLong.Longitude)
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("error making request to Weather API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}
