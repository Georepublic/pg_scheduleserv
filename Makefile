# GRP-GNU-AGPL******************************************************************

# File: Makefile

# Copyright (C) 2021  Team Georepublic <info@georepublic.de>

# Developer(s):
# Copyright (C) 2021  Ashish Kumar <ashishkr23438@gmail.com>

# -----

# This file is part of pg_scheduleserv.

# pg_scheduleserv is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published
# by the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# pg_scheduleserv is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.

# You should have received a copy of the GNU Affero General Public License
# along with pg_scheduleserv.  If not, see <https://www.gnu.org/licenses/>.

# ******************************************************************GRP-GNU-AGPL

app.env:
	cp app.env.example app.env

include app.env

migrateup:
	migrate -path migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB) -verbose up
.PHONY: migrateup

migratedown:
	migrate -path migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB) -verbose down
.PHONY: migratedown

doc:
	swag init
	swagger generate markdown -f docs/swagger.json --output=docs/api.md
.PHONY: doc

check-mod:
	go mod tidy
	git diff --exit-code go.mod
.PHONY: check-mod

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run
.PHONY: lint

test:
	go install gotest.tools/gotestsum@latest
	gotestsum --format=testname -- -count=1 -timeout=20m -coverpkg=./... -coverprofile=coverage.out ./...
.PHONY: test

test-coverage:
	go tool cover -func=./coverage.out
.PHONY: test-coverage

html-coverage:
	go tool cover -html=coverage.out -o coverage.html
.PHONY: html-coverage
