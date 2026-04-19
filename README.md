# migrate

Fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — Database migrations written in Go.

[![Go Reference](https://pkg.go.dev/badge/github.com/migrate/migrate.svg)](https://pkg.go.dev/github.com/migrate/migrate)
[![CI](https://github.com/migrate/migrate/actions/workflows/ci.yaml/badge.svg)](https://github.com/migrate/migrate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/migrate/migrate)](https://goreportcard.com/report/github.com/migrate/migrate)

> **Personal fork** — primarily used with PostgreSQL and local filesystem sources. Other drivers are available but untested in this fork.

## Features

- Supports multiple database drivers (PostgreSQL, MySQL, SQLite, and more)
- Supports multiple migration sources (filesystem, S3, GitHub, and more)
- CLI and Go library usage
- Atomic migrations with rollback support

## Installation

### CLI

```bash
go install github.com/migrate/migrate/v4/cmd/migrate@latest
```

### As a library

```bash
go get github.com/migrate/migrate/v4
```

## Usage

### CLI

```bash
# Apply all up migrations
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" up

# Rollback the last migration
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" down 1

# Check current migration version
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" version
```

### Library

```go
import (
    "github.com/migrate/migrate/v4"
    _ "github.com/migrate/migrate/v4/database/postgres"
    _ "github.com/migrate/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://./migrations",
        "postgres://localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Migration Files

Migration files follow the naming convention:

```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

Example:

```
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_add_email_index.up.sql
000002_add_email_index.down.sql
```

## Supported Databases

| Database   | Driver import path |
|------------|-----------------|
| PostgreSQL | `github.com/migrate/migrate/v4/database/postgres` |
| MySQL      | `github.com/migrate/migrate/v4/database/mysql` |
| SQLite     | `github.com/migrate/migrate/v4/database/sqlite3` |
| MongoDB    | `github.com/migrate/migrate/v4/database/mongodb` |

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)
