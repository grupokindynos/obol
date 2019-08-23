package config

import "errors"

var (
	ErrorUnableToParseStringToFloat = errors.New("unable to convert string to float")
	ErrorCoinNotAvailable           = errors.New("coin not available")
)
