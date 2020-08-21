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

## Documentation

The graphQL documentation is available at endpoint `/graphql`.

## Environment Variables
#### Runtime
| name           | description                                     | type    | optional | default      |
|----------------|-------------------------------------------------|---------|----------|--------------|
| PORT           | The server port                                 | string  | yes      | 3000         |
| LOGGER         | The logger type (TEST, DEVELOPMENT, PRODUCTION) | string  | yes      | DEVELOPMENT  |
| MYSQL_USER     | The database user                               | string  | no       |              |
| MYSQL_PASSWORD | The database password                           | string  | no       |              |
| MYSQL_HOST     | The database host                               | string  | no       |              |
| MYSQL_PORT     | The database port                               | string  | no       | 3306         |
| MYSQL_NAME     | The database name                               | string  | no       | sakila       |
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

## Coming Soon

[ ] Health checker
[ ] Rest API documentation
[ ] Linter
