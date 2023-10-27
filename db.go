package main

import (
	sql "database/sql"
)

func insertCity(db *sql.DB, name string, latLong LatLong) error {
	_, err := db.Exec("INSERT INTO cities (name, lat, long) VALUES (?, ?, ?)", name, latLong.Latitude, latLong.Longitude)
	return err
}

func getLatLong(db *sql.DB, name string) (*LatLong, error) {
	var latLong *LatLong = new(LatLong)

	rows, err := db.Query("SELECT lat, long FROM cities WHERE name = ?", name)
	if err == nil {
		for rows.Next() {
			if err = rows.Scan(&latLong.Latitude, &latLong.Longitude); err == nil {
				return latLong, nil
			}
		}
	}

	latLong, err = fetchLatLong(name)
	if err != nil {
		return nil, err
	}

	err = insertCity(db, name, *latLong)
	if err != nil {
		return nil, err
	}

	return latLong, nil
}

func getLastCities(db *sql.DB) ([]string, error) {
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
