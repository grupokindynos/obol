package main

import (
	"github.com/eabz/cache"
	"github.com/eabz/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/controllers"
	"github.com/grupokindynos/obol/services"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	limit "github.com/yangxikun/gin-limit-by-key"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"time"
)

func init() {
	_ = godotenv.Load()
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	App := GetApp()
	_ = App.Run(":" + port)
}

func GetApp() *gin.Engine {
	App := gin.Default()
	App.Use(cors.Default())
	ApplyRoutes(App)
	return App
}

func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		api.Use(limit.NewRateLimiter(func(c *gin.Context) string {
			return c.ClientIP()
		}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
			return rate.NewLimiter(rate.Every(100*time.Hour), 10), time.Hour
		}, func(c *gin.Context) {
			c.AbortWithStatus(429)
		}))
		store := persistence.NewInMemoryStore(time.Second)
		rateService := services.InitRateService()
		rateCtrl := controllers.RateController{RateService: rateService}
		api.GET("simple/:coin", cache.CachePage(store, time.Minute*5, rateCtrl.GetCoinRates))
		api.GET("complex/:fromcoin/:tocoin", cache.CachePage(store, time.Minute*5, rateCtrl.GetCoinRateFromCoinToCoin))
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
