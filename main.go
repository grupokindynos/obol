package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/controllers"
	"github.com/grupokindynos/obol/services"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
	err := App.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

// GetApp is used to wrap all the additions to the GIN API.
func GetApp() *gin.Engine {
	App := gin.Default()
	App.Use(cors.Default())
	ApplyRoutes(App)
	return App
}

// ApplyRoutes is used to attach all the routes to the API service.
func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		rate := limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		}
		store := memory.NewStore()
		limiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
		api.Use(limiterMiddleware)
		rateService := services.InitRateService()
		rateCtrl := controllers.RateController{RateService: rateService, RatesCache: make(map[string]controllers.CoinRate)}
		api.GET("simple/:coin", rateCtrl.GetCoinRates)
		api.GET("complex/:fromcoin/:tocoin", rateCtrl.GetCoinRateFromCoinToCoin)
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
