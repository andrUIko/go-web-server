package db

import (
	"database/sql"
	"github.com/user/goforecast/models"
)

func GetLastCities(db *sql.DB) ([]string, error) {
	var cities []string
	rows, err := db.Query("SELECT name FROM cities ORDER BY id DESC LIMIT 10")
	if err == nil {
		var name string
		for rows.Next() {
			if err = rows.Scan(&name); err != nil {
				return nil, err
			}
			cities = append(cities, name)
		}
	}
	return cities, nil
}

func InsertCity(db *sql.DB, name string, latLong models.LatLong) error {
	_, err := db.Exec("INSERT INTO cities (name, lat, long) VALUES (?, ?, ?)", name, latLong.Latitude, latLong.Longitude)
	return err
}

func GetLatLong(db *sql.DB, name string) (*models.LatLong, error) {
	var latLong *models.LatLong = new(models.LatLong)
	rows, err := db.Query("SELECT lat, long FROM cities WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	rows.Next()
	if err = rows.Scan(&latLong.Latitude, &latLong.Longitude); err == nil {
		return latLong, nil
	}
	return nil, err
}
