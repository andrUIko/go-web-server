package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/goforecast/db"
)

func GetStats(c *gin.Context) {
	dbClient, ok := c.Keys["DB"].(db.DBPool)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no db client found"})
		return
	}

	cities, err := dbClient.GetLastCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "stats.html", cities)
}
