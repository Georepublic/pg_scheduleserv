# pg_scheduleserv

[![Test](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/test.yml)
[![Lint](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/lint.yml)
[![codecov](https://img.shields.io/codecov/c/github/Georepublic/pg_scheduleserv/main?logo=codecov)](https://codecov.io/gh/Georepublic/pg_scheduleserv)
[![License: AGPL v3](https://img.shields.io/github/license/Georepublic/pg_scheduleserv)](https://www.gnu.org/licenses/agpl-3.0)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Georepublic/pg_scheduleserv)](https://go.dev/doc/go1.17)
[![Go Reference](https://pkg.go.dev/badge/github.com/Georepublic/pg_scheduleserv.svg)](https://pkg.go.dev/github.com/Georepublic/pg_scheduleserv)
[![Go Report Card](https://goreportcard.com/badge/github.com/Georepublic/pg_scheduleserv)](https://goreportcard.com/report/github.com/Georepublic/pg_scheduleserv)
[![GitHub Release](https://img.shields.io/github/release/Georepublic/pg_scheduleserv.svg)](https://github.com/Georepublic/pg_scheduleserv/releases)
![GitHub all releases](https://img.shields.io/github/downloads/Georepublic/pg_scheduleserv/total)
![GitHub Release Date](https://img.shields.io/github/release-date/Georepublic/pg_scheduleserv)

A RESTful API Server for scheduling VRP tasks using [vrpRouting](https://github.com/pgRouting/vrprouting), written in [Go](https://golang.org/).

API Documentation: [docs/api.md](./docs/api.md)  
Release Notes: [NEWS.md](./NEWS.md)

## Getting Started

### Requirements

- [vrpRouting](https://github.com/pgRouting/vrprouting) >= 0.2.0
  - [VROOM](https://github.com/VROOM-Project/vroom) >= 1.10.0 is required to build.
  - [pgRouting](https://github.com/pgRouting/pgrouting)
  - [PostGIS](https://postgis.net/)
  - [PostgreSQL](https://www.postgresql.org/)
  - C and C++ compilers with C++17 standard support
  - The Boost Graph Library (BGL) >= 1.65
  - CMake >= 3.12
- [Go](https://golang.org/) == 1.17

### Build from source

`pg_scheduleserv` is developed using Go 1.17. It may work with earlier versions. To build this project:
- Ensure that the Go compiler is installed
- Download or clone this repository.
- Copy `app.env.example` to `app.env`, and set the value of `DATABASE_URL`.
- Apply migrations from the `migrations/000001_init.up.sql` file to the database.
- Run `go build` command inside the directory to create the executable named `pg_scheduleserv`
- Run the executable to start the API server on http://localhost:9100

```
git clone https://github.com/Georepublic/pg_scheduleserv
cd pg_scheduleserv
cp app.env.example app.env
vi app.env  # Set the variables
psql -U username -d database -f migrations/000001_init.up.sql
go build
./pg_scheduleserv
```

### Usage

The swagger documentation is available on http://localhost:9100, which can also be used to test the API. Or any other API client can be independently used to interact with the API server.

## LICENSE

This project is licensed under the GNU Affero General Public License v3.0. View [LICENSE](./LICENSE) for more details.
