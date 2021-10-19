#!/bin/bash

<<//
/*GRP-GNU-AGPL******************************************************************

File: create_migrations.sh - script to create a new database migration file

Copyright (C) 2021  Team Georepublic <info@georepublic.de>

Developer(s):
Copyright (C) 2021  Ashish Kumar <ashishkr23438@gmail.com>

-----

This file is part of pg_scheduleserv.

pg_scheduleserv is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

pg_scheduleserv is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with pg_scheduleserv.  If not, see <https://www.gnu.org/licenses/>.

******************************************************************GRP-GNU-AGPL*/
//

if [[ -z "$1" ]]; then
    echo "Please provide the migration file name."
    exit 1
fi

command -v migrate >/dev/null 2>&1 || {
    echo >&2 "Migrate command not found. Please install golang-migrate";
    exit 1;
}

migrate create -ext sql -dir migrations -seq $1

TEMPLATE='BEGIN;



END;'

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

for m in $(ls migrations | tail -n 2); do
    echo "$TEMPLATE" >> "migrations/$m";
    "$SCRIPT_DIR/../license/add_license.sh" "migrations/$m"
done
