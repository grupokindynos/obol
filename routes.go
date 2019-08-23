package main

import (
	"github.com/eabz/cache"
	"github.com/eabz/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/services"
	"net/http"
	"time"
)

func ApplyRoutes(r *gin.Engine) {
	api := r.Group("/")
	{
		store := persistence.NewInMemoryStore(time.Second)
		rateService := services.InitRateService()
		rateCtrl := RateController{RateService: rateService}
		api.GET("simple/:coin", cache.CachePage(store, time.Minute*5, rateCtrl.GetCoinRates))
		api.GET("complex/:fromcoin/:tocoin", cache.CachePage(store, time.Minute*5, rateCtrl.GetCoinRateFromCoinToCoin))
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
