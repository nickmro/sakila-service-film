# Sakila Film Service
[![Build Status](https://travis-ci.com/nickmro/sakila-service-film.svg?branch=master)](https://travis-ci.com/nickmro/sakila-service-film)

## Introduction

The Sakila Film Service serves film data from the [Sakila Sample Database](https://dev.mysql.com/doc/sakila/en/).

This project has been built as an example of how to build a microservice around a given database.

## Features

- GraphQL API

## Installation

Install Go: https://golang.org/doc/install

Install golangci-lint: https://github.com/golangci/golangci-lint

Install MySQL:
```bash
brew install mysql
```

Install Redis:
```bash
brew install redis
```

Download and install the Sakila database following the instructions here: https://dev.mysql.com/doc/sakila/en/sakila-installation.html

Create `.env` file:
```bash
cp .env.template .env
```

Fill `.env` with the applicable environment variables.

## Run

```bash
go build -o ./bin/serve -i ./cmd/serve
./bin/serve
```

## Docker

```bash
docker build -t sakila/service-film:1.0 .
docker run --name sakila-service-film --publish 3000:3000 sakila/service-film:1.0
```

## Documentation

The REST API documentation is provided by `swagger.yml`. It can be loaded into an API client such as [Postman](https://www.postman.com/) or [Insomnia](https://insomnia.rest/).

The GraphQL schema and documentation may also be viewed by loading the schema into one of the API clients above.

## Environment Variables
#### Runtime
| name                   | description                                     | type    | optional | default      |
|------------------------|-------------------------------------------------|---------|---------|--------------|
| PORT                   | The server port                                 | string  | yes      | 3000         |
| LOGGER                 | The logger type (TEST, DEVELOPMENT, PRODUCTION) | string  | yes      | DEVELOPMENT  |
| MYSQL_USER             | The database user                               | string  | yes      |              |
| MYSQL_PASSWORD         | The database password                           | string  | yes      |              |
| MYSQL_HOST             | The database host                               | string  | no       |              |
| MYSQL_PORT             | The database port                               | string  | no       |              |
| MYSQL_NAME             | The database name                               | string  | no       |              |
| REDIS_HOST             | The cache host                                  | string  | no       |              |
| REDIS_PORT             | The cache port                                  | string  | no       |              |
| REDIS_PASSWORD         | The cache password                              | string  | yes      |              |
| REDIS_KEY_PREFIX       | The cache key prefix                            | string  | no       |              |

## Test

Tests are written using [ginkgo](https://onsi.github.io/ginkgo/) and [gomega](http://onsi.github.io/gomega/).

To run tests:
```bash
go test ./...
```

To use ginkgo, first install:
```
go get -u github.com/onsi/ginkgo/ginkgo
```

Then run:
```
ginkgo -r
```

## Lint

To run the linter:
```
golangci-lint run
```
