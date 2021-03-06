package main

import (
	"github.com/grupokindynos/obol/services/exchanges/birake"
	"github.com/grupokindynos/obol/services/exchanges/bithumb"
	"github.com/grupokindynos/obol/services/exchanges/pancake"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/controllers"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/services"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/bitrue"
	"github.com/grupokindynos/obol/services/exchanges/bittrex"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/graviex"
	"github.com/grupokindynos/obol/services/exchanges/hitbtc"
	"github.com/grupokindynos/obol/services/exchanges/kucoin"
	"github.com/grupokindynos/obol/services/exchanges/lukki"
	"github.com/grupokindynos/obol/services/exchanges/southxhcange"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
	rateService := &services.RateSevice{
		FiatRates: &models.FiatRates{
			Rates:       nil,
			LastUpdated: time.Time{},
		},
		BittrexService:      bittrex.InitService(),
		BinanceService:      binance.InitService(),
		Crex24Service:       crex24.InitService(),
		StexService:         stex.InitService(),
		SouthXChangeService: southxhcange.InitService(),
		KuCoinService:       kucoin.InitService(),
		GraviexService:      graviex.InitService(),
		BitrueService:       bitrue.InitService(),
		HitBTCService:       hitbtc.InitService(),
		LukkiService:        lukki.InitService(),
		BithumbService:      bithumb.InitService(),
		BirakeService:       birake.InitService(),
		PancakeService: 	 pancake.InitService(),
	}
	err := rateService.LoadFiatRates()
	if err != nil {
		panic(err)
	}
	api := r.Group("/")
	{
		rate := limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		}
		store := memory.NewStore()
		limiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
		api.Use(limiterMiddleware)
		rateCtrl := controllers.RateController{RateService: rateService, RatesCache: make(map[string]controllers.CoinRate)}
		api.GET("simple/:coin", rateCtrl.GetCoinRates)
		api.GET("complex/:fromcoin/:tocoin", rateCtrl.GetCoinRateFromCoinToCoin)
		api.GET("liquidity/:coin", rateCtrl.GetCoinLiquidity)
		api.GET("rate/margin/:fromCoin/:toCoin", rateCtrl.GetCoinToCoinRateWithExchangeMargin)

	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	apiv2 := r.Group("/v2/")
	{
		rate := limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		}
		store := memory.NewStore()
		limiterMiddleware := mgin.NewMiddleware(limiter.New(store, rate))
		apiv2.Use(limiterMiddleware)
		rateCtrl := controllers.RateController{RateService: rateService, RatesCache: make(map[string]controllers.CoinRate)}
		apiv2.GET("simple/:coin", rateCtrl.GetCoinRatesV2)
		apiv2.GET("complexfiat/:fromcoin/:tocoin", rateCtrl.GetCoinToFIATRate)
		apiv2.GET("node/:coin", rateCtrl.GetNodeProvider)

	}

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
