<div align="center">
  <img alt="pg_scheduleserv logo" src="./docs/images/logo.png" width="250px" />

# pg_scheduleserv - VRP scheduler over the web

[![Go Reference](https://pkg.go.dev/badge/github.com/Georepublic/pg_scheduleserv.svg)](https://pkg.go.dev/github.com/Georepublic/pg_scheduleserv)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Georepublic/pg_scheduleserv)](https://go.dev/doc/go1.17)
[![Go Report Card](https://goreportcard.com/badge/github.com/Georepublic/pg_scheduleserv)](https://goreportcard.com/report/github.com/Georepublic/pg_scheduleserv)
[![GitHub Release](https://img.shields.io/github/release/Georepublic/pg_scheduleserv.svg)](https://github.com/Georepublic/pg_scheduleserv/releases)
![GitHub all releases](https://img.shields.io/github/downloads/Georepublic/pg_scheduleserv/total)
![GitHub Release Date](https://img.shields.io/github/release-date/Georepublic/pg_scheduleserv)
[![License: AGPL v3](https://img.shields.io/github/license/Georepublic/pg_scheduleserv)](https://www.gnu.org/licenses/agpl-3.0)

</div>

## Introduction

A RESTful API Server for scheduling VRP tasks using [vrpRouting](https://github.com/pgRouting/vrprouting), written in [Go](https://golang.org/).

API Documentation: [docs/api.md](./docs/api.md)\
Release Notes: [NEWS.md](./NEWS.md)

## Status

| Service  | Main                                                                                                                                                                                                      | Develop                                                                                                                                                                                                         |
| -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Test     | [![Test](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/test.yml?query=branch%3Amain) | [![Test](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/test.yml/badge.svg?branch=develop)](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/test.yml?query=branch%3Adevelop) |
| Lint     | [![Lint](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/lint.yml?query=branch%3Amain) | [![Lint](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/lint.yml/badge.svg?branch=develop)](https://github.com/Georepublic/pg_scheduleserv/actions/workflows/lint.yml?query=branch%3Adevelop) |
| Coverage | [![codecov](https://img.shields.io/codecov/c/github/Georepublic/pg_scheduleserv/main?logo=codecov)](https://app.codecov.io/gh/Georepublic/pg_scheduleserv/branch/main)                                    | [![codecov](https://img.shields.io/codecov/c/github/Georepublic/pg_scheduleserv/develop?logo=codecov)](https://app.codecov.io/gh/Georepublic/pg_scheduleserv/branch/develop)                                    |

## Getting Started

### Requirements

-   Build Requirements:

    -   [Go](https://golang.org/) == 1.17

-   Usage Requirements:
    -   [vrpRouting](https://github.com/pgRouting/vrprouting) >= 0.3.0
        -   [VROOM](https://github.com/VROOM-Project/vroom) >= 1.11.0 is required to build vrpRouting.
        -   [pgRouting](https://github.com/pgRouting/pgrouting)
        -   [PostGIS](https://postgis.net/)
        -   [PostgreSQL](https://www.postgresql.org/)
        -   C and C++ compilers with C++17 standard support
        -   The Boost Graph Library (BGL) >= 1.65
        -   CMake >= 3.12

### Download and run

Builds of the latest code can be found in the [releases](https://github.com/Georepublic/pg_scheduleserv/releases).

-   Download the latest executable for Linux, say [pg_scheduleserv-0.2.0](https://github.com/Georepublic/pg_scheduleserv/releases/download/v0.2.0/pg_scheduleserv-0.2.0).
-   Change permissions to make it an executable: `chmod +x pg_scheduleserv-0.2.0`
-   Create `app.env`, and set the values to the environment variables:
    -   POSTGRES_USER=username
    -   POSTGRES_PASSWORD=password
    -   POSTGRES_HOST=localhost
    -   POSTGRES_PORT=5432
    -   POSTGRES_DB=scheduler
    -   SERVER_PORT=:9100
    -   OSRM_URL=https://router.project-osrm.org
    -   VALHALLA_URL=https://valhalla1.openstreetmap.de
-   Create the tables in the database with the help of the migrations file.
-   Run the executable to start the API server on http://localhost:9100

### Build from source

All the steps are similar as above, except that the executable will be built from source.

`pg_scheduleserv` is developed using Go 1.17. To build this project from source:

-   Ensure that the Go compiler is installed
-   Download or clone this repository.
-   Copy `app.env.example` to `app.env`, and set the values to the environment variables.
-   Apply migrations from the `migrations/000001_init.up.sql` file to the database.
-   Run `go build` command inside the directory to create the executable named `pg_scheduleserv`
-   Run the executable to start the API server on http://localhost:9100

### Usage

Any API client can be used to interact with the API server. The API can also be tested using the Swagger documentation. Sample deployed API can be found on https://api-v0.udc.pgrouting.org/

These documenations can be accessed after running the API server on the following URLs:

-   Swagger UI: http://localhost:9100
-   RapiDoc: http://localhost:9100/rapidoc
-   Redoc: http://localhost:9100/redoc

![swagger-api](https://user-images.githubusercontent.com/39548570/152192999-1f173519-61a8-4b9b-91f4-ae680f783fe1.png)

### Demo Application

A frontend demo application for the API resides in the `demo` directory. Sample deployed demo application can be found on https://demo-v0.udc.pgrouting.org/

To run the demo application:

-   Change directory to the `demo` directory: `cd demo`
-   Make sure that node is installed. Install the other requirements: `npm install`
-   Run the application: `node server.js`. The application will run at http://localhost:9101 and will use the API running at http://localhost:9100
-   To change the Server Base URL or the OSRM API URL, modify the `js/config.js` file.

![demo-application](https://user-images.githubusercontent.com/39548570/152192932-2fe42d9f-b464-42ec-9a10-47779d087c7e.png)

## LICENSE

This project is licensed under the GNU Affero General Public License v3.0. View [LICENSE](./LICENSE) for more details.
