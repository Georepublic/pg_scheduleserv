# pg_scheduleserv

_...This project is a Work in Progress..._

A RESTful API Server for scheduling VRP tasks using [vrpRouting](https://github.com/pgRouting/vrprouting), written in [Go](https://golang.org/).

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

The swagger documentation is available on http://localhost:9100 _(work in progress)_, which can also be used to test the API. Or any other API client can be independently used to interact with the API server.

## LICENSE

* This project is licensed under the GNU Affero General Public License v3.0. View [LICENSE](./LICENSE) for more details.
