# Sakila Film Service

## Introduction

The Sakila Film Service serves film data from the [Sakila Sample Database](https://dev.mysql.com/doc/sakila/en/).

This project has been built as an example of how to build a microservice around a given database.

## Features

- REST API
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
| name           | description                                     | type    | optional | default      |
|----------------|-------------------------------------------------|---------|----------|--------------|
| PORT           | The server port                                 | string  | yes      | 3000         |
| LOGGER         | The logger type (TEST, DEVELOPMENT, PRODUCTION) | string  | yes      | DEVELOPMENT  |
| MYSQL_USER     | The database user                               | string  | no       |              |
| MYSQL_PASSWORD | The database password                           | string  | no       |              |
| MYSQL_HOST     | The database host                               | string  | no       |              |
| MYSQL_PORT     | The database port                               | string  | no       |              |
| MYSQL_NAME     | The database name                               | string  | no       |              |
| REDIS_URL      | The cache URL                                   | string  | no       |              |
| REDIS_PASSWORD | The cache password                              | string  | no       |              |

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

## Coming Soon

- [ ] CI
