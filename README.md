# Obol

[![CircleCI](https://circleci.com/gh/grupokindynos/obol.svg?style=svg)](https://circleci.com/gh/grupokindynos/obol)
[![codecov](https://codecov.io/gh/grupokindynos/obol/branch/master/graph/badge.svg?token=doQVTdPAe5)](https://codecov.io/gh/grupokindynos/obol)
[![Go Report](https://goreportcard.com/badge/github.com/grupokindynos/obol)](https://goreportcard.com/report/github.com/grupokindynos/obol)
[![GoDocs](https://godoc.org/github.com/grupokindynos/obol?status.svg)](http://godoc.org/github.com/grupokindynos/obol)

> The obol was a form of ancient Greek currency and weight

Obol is a microservice API for multiple cryptocurrency rates.

## Deploy

#### Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/grupokindynos/obol/blob/master/)

#### Docker

To deploy to docker, simply pull the image
```
docker pull kindynos/obol:latest
```
Run the docker image
```
docker run -p 8080:8080 kindynos/plutus:latest 
```

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
go mod download
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

Documentation: [API Reference](https://documenter.getpostman.com/view/6539894/SVfNx9qJ?version=latest)

## Testing

Simply run:

```
go test ./...
```

## Contributing

To contribute to this repository, please fork it, create a new branch and submit a pull request.

To add a new coin, please add a new coin configuration under `models/coin-factory/coins.go`

To add a new exchange please add it to `services/exchanges/` under a new folder with same functions and test cases.
