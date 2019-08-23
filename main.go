package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	_ = godotenv.Load()
}
func main() {
	port := os.Getenv("PORT")
	App := GetApp()
	err := App.Run(":" + port)
	if err != nil {
		log.Panic(err)
	}
}

func GetApp() *gin.Engine {
	App := gin.Default()
	App.Use(cors.Default())
	api.ApplyRoutes(App)
	return App
}




