package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/user/goforecast/models"
)

type DBClient struct {
	DB *sql.DB
}

type DBPool interface {
	GetLastCities() ([]string, error)
	InsertCity(name string, latLong *models.LatLong) error
	GetLatLong(name string) (*models.LatLong, error)
}

func CreateDBClient(url string) DBPool {
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	var dbClient DBPool = &DBClient{db}
	return dbClient
}

func (c *DBClient) GetLastCities() ([]string, error) {
	var cities []string
	rows, err := c.DB.Query("SELECT name FROM cities ORDER BY id DESC LIMIT 10")
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

func (c *DBClient) InsertCity(name string, latLong *models.LatLong) error {
	_, err := c.DB.Exec("INSERT INTO cities (name, lat, long) VALUES (?, ?, ?)", name, latLong.Latitude, latLong.Longitude)
	return err
}

func (c *DBClient) GetLatLong(name string) (*models.LatLong, error) {
	var latLong = new(models.LatLong)
	rows, err := c.DB.Query("SELECT lat, long FROM cities WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	rows.Next()
	if err = rows.Scan(&latLong.Latitude, &latLong.Longitude); err == nil {
		return latLong, nil
	}
	return nil, err
}
