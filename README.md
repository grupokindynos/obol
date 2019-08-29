# Obol

[![CircleCI](https://circleci.com/gh/grupokindynos/obol.svg?style=svg)](https://circleci.com/gh/grupokindynos/obol)
[![codecov](https://codecov.io/gh/grupokindynos/obol/branch/master/graph/badge.svg?token=doQVTdPAe5)](https://codecov.io/gh/grupokindynos/obol)
[![Go Report](https://goreportcard.com/badge/github.com/grupokindynos/obol)](https://goreportcard.com/report/github.com/grupokindynos/obol)
[![GoDocs](https://godoc.org/github.com/grupokindynos/obol?status.svg)](http://godoc.org/github.com/grupokindynos/obol)

> The obol was a form of ancient Greek currency and weight

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/grupokindynos/obol/blob/master/)

Obol is a microservice API for multiple cryptocurrency rates.

## Building

To run Obol from the source code, first you need to install golang, follow this guide:

```
https://golang.org/doc/install
```

To run Obol simply clone de repository:

```
git clone https://github.com/grupokindynos/obol && cd obol
```

Install dependencies

```
go mod tidy
```

Build it or Run it:

```
go build && ./obol
```

```
go run main.go
```

Make sure the port is configured under en enviroment variable `PORT=8080`

## API Reference

### Get Rates:

Retrieves the rate to many currencies for one coin.

**GET method:**

```
https://obol-rates.herokuapp.com/simple/:coin
```

This will retrieve an array with rates result based on exchange real time price:

```
{
    "data": {
        "AUD": 0.934,
        "BGN": 1.1131,
        "BRL": 2.6261,
        "BTC": 0.00006171,
        "CAD": 0.8359,
        "CHF": 0.6194,
        "CNY": 4.5256,
        "CZK": 14.6914,
        "DKK": 4.2444,
        "GBP": 0.5146,
        "HKD": 4.9576,
        "HRK": 4.2117,
        "HUF": 187.158,
        "IDR": 9004.686,
        "ILS": 2.2236,
        "INR": 45.2016,
        "ISK": 78.9382,
        "JPY": 66.8328,
        "KRW": 766.2924,
        "MXN": 12.6183,
        "MYR": 2.6561,
        "NOK": 5.6867,
        "NZD": 0.9911,
        "PHP": 33.0635,
        "PLN": 2.4868,
        "RON": 2.6931,
        "RUB": 42.0271,
        "SEK": 6.0925,
        "SGD": 0.877,
        "THB": 19.3133,
        "TRY": 3.6822,
        "USD": 0.6319,
        "ZAR": 9.6663
    },
    "status": 1
}
```

### Get Rates from one coin to another:

Retrieves the rate from one coin to another.

**GET method:**

```
https://obol-rates.herokuapp.com/complex/:fromcoin/:tocoin
```

This will get the current amount you need to convert from one coin to another:

```
{
    "data": 144.88721438,
    "status": 1
}
```

### Get Rates from amount

Retrieves the rate from one coin to another, provided the amount of coins you want to convert.

**GET method:**

```
https://obol-rates.herokuapp.com/complex/:fromcoin/:tocoin?amount=100
```

This will get the current amount you need to convert from one coin to another, provided the amount of coins you want to convert:

```
{
    "data": 140.88721438,
    "status": 1
}
```

## Testing

Simply run:

```
go test ./...
```

## Contributing

Pull Requests accepted.

To add a new coin, please add a new coin configuration under `models/coin-factory/coins.go`

To add a new exchange please add it to `services/exchanges/` under a new folder with same functions and test cases.
