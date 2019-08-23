package config

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	OpenRatesURL = "https://api.exchangeratesapi.io/latest?base=MXN"
	ErrorUnableToParseStringToFloat = errors.New("unable to convert string to float")
	ErrorCoinNotAvailable           = errors.New("coin not available")
	HttpClient = &http.Client{
		Timeout: time.Second * 10,
	}
)

func GlobalResponse(result interface{}, err error, c *gin.Context) *gin.Context {
	if err != nil {
		c.JSON(500, gin.H{"message": "Error", "error": err.Error(), "status": -1})
	} else {
		c.JSON(200, gin.H{"data": result, "status": 1})
	}
	return c
}
