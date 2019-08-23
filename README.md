# Obol

[![CircleCI](https://circleci.com/gh/grupokindynos/obol.svg?style=svg)](https://circleci.com/gh/grupokindynos/obol) 
[![codecov](https://codecov.io/gh/grupokindynos/obol/branch/master/graph/badge.svg?token=doQVTdPAe5)](https://codecov.io/gh/grupokindynos/obol) 
[![Go Report](https://goreportcard.com/badge/github.com/grupokindynos/obol)](https://goreportcard.com/report/github.com/grupokindynos/obol) 
[![GoDocs](https://godoc.org/github.com/grupokindynos/obol?status.svg)](http://godoc.org/github.com/grupokindynos/obol)


> The obol was a form of ancient Greek currency and weight

Obol is a microservice API for multiple cryptocurrency rates.

## Building

To run Obol simply clone de repository:

```
git clone https://github.com/grupokindynos/obol 
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
@TODO

## Testing

Simply run:
```
go test ./...
```

## Contributing

Pull Requests accepted.

To add a new coin, please add a new coin configuration under `models/coin-factory/coins.go`

To add a new exchange please add it to `services/exchanges/` under a new folder with same functions and test cases.