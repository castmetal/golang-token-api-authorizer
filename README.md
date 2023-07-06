# Golang api boilerplate

API Boilerplate using GIN Framework, Swagger, Redis and PostgreSQL

## Installation

This API boilerplate requires [Golang](https://go.dev/dl/) v17+ to run.

Install the dependencies and devDependencies and start the server.

```sh
go mod tidy
go install
```

### Run Migrations

```sh
make migrate-db
```

### Sqlc

You need sqlc to generate more queries that you'll need.

```sh
make generate-sql
```

For dev development...

```sh
go run main.go
```

## Plugins

SwaggerGo

| Plugin | README                                     |
| ------ | ------------------------------------------ |
| Swaggo | [https://github.com/swaggo/swag][pldb]     |
| Sqlc   | [https://github.com/kyleconroy/sqlc][plgh] |

## Docker

Todo

## Tests

Todo

## License

MIT
