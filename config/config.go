package config

import (
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

var (
	// OpenRatesURL is the base URL for fiat rates based on OpenRates
	OpenRatesURL = "https://api.exchangeratesapi.io/latest?base=MXN"
	// ErrorCoinNotAvailable returns when the specified coin is not configured on the API.
	ErrorCoinNotAvailable = errors.New("coin not available")
	// ErrorNoServiceForCoin returns when is not possible to load the service for an exchange
	ErrorNoServiceForCoin = errors.New("unable to load exchange for this coin")
	// ErrorNoFallBackServiceForCoin returns when is not possible to load the fallback service for an exchange
	ErrorNoFallBackServiceForCoin = errors.New("unable to load fallback exchange for this coin")
	// ErrorNoC2CWithSameCoin returns when a C2C rate is called using the same coins
	ErrorNoC2CWithSameCoin = errors.New("cannot use the same coin on both parameters")
	// ErrorInvalidAmountOnC2C returns when try to call a C2C rate with amount, and the amount is not properly formated.
	ErrorInvalidAmountOnC2C = errors.New("invalid amount to convert from coin to coin")
	// ErrorUnknownIdForCoin returns when an id for a coin is not defined
	ErrorUnknownIdForCoin = errors.New("unknown id for coin")
	// ErrorRequestTimeout returns when a request timeout threshold is reached
	ErrorRequestTimeout = errors.New("request timeout")
	// HttpClient is a wrapper to properly timeout http.get calls
	HttpClient = &http.Client{
		Timeout: time.Second * 2,
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
