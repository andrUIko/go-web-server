package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	dotenv "github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
	"github.com/user/goforecast/controllers"
	database "github.com/user/goforecast/db"
)

func main() {
	if err := dotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var (
		DATABASE_NAME         = os.Getenv("DATABASE_NAME")
		DATABASE_ACCESS_TOKEN = os.Getenv("DATABASE_ACCESS_TOKEN")
		AUTH_LOGIN            = os.Getenv("AUTH_LOGIN")
		AUTH_PASSWORD         = os.Getenv("AUTH_PASSWORD")
	)

	dbUrl := "libsql://" + DATABASE_NAME + ".turso.io?authToken=" + DATABASE_ACCESS_TOKEN
	dbClient := database.CreateDBClient(dbUrl)

	r := gin.Default()
	r.Static("/assets", "./assets")
	r.Use(func(ctx *gin.Context) {
		ctx.Keys = make(map[string]interface{})
		ctx.Keys["DB"] = dbClient
	})

	r.LoadHTMLGlob("views/*")

	r.GET("/weather", controllers.GetWeather)

	r.GET("/stats", gin.BasicAuth(gin.Accounts{
		AUTH_LOGIN: AUTH_PASSWORD,
	}), controllers.GetStats)

	r.GET("/", controllers.GetIndex)

	r.Run()
}
