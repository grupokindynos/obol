package config

import (
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

var (
	OpenRatesURL            = "https://api.exchangeratesapi.io/latest?base=MXN"
	ErrorCoinNotAvailable   = errors.New("coin not available")
	ErrorNoServiceForCoin   = errors.New("unable to load exchange for this coin")
	ErrorNoC2CWithBTC       = errors.New("coin to coin function doesn't work using BTC")
	ErrorNoC2CWithSameCoin  = errors.New("cannot use the same coin on both parameters")
	ErrorInvalidAmountOnC2C = errors.New("invalid amount to convert from coin to coin")
	ErrorUnknownIdForCoin   = errors.New("unknown id for coin")
	ErrorRequestTimeout     = errors.New("request timeout")
	HttpClient              = &http.Client{
		Timeout: time.Second * 1,
	}
)

// GlobalResponse is used to wrap all the API responses under the same model.
// Automatically detect if there is an error and return status and code according
func GlobalResponse(result interface{}, err error, c *gin.Context) *gin.Context {
	if err != nil {
		c.JSON(500, gin.H{"message": "Error", "error": err.Error(), "status": -1})
		return c
	}
	// If is a float, truncate it to sats
	value, isfloat := result.(float64)
	if isfloat {
		value := math.Floor(value*1e8) / 1e8
		c.JSON(200, gin.H{"data": value, "status": 1})
		return c
	}
	c.JSON(200, gin.H{"data": result, "status": 1})
	return c
}
