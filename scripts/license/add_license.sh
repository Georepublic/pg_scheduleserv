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

year="2021"
author="Ashish Kumar"
author_mail="ashishkr23438@gmail.com"

license_file="GNU_AGPL.txt"

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
license=$(cat $SCRIPT_DIR/$license_file)
license=${license//\[year\]/"$year"}
license=${license//\[name of author\]/"$author"}
license=${license//\[mail of author\]/"$author_mail"}

function add_license {
    echo "Adding license to $1"
    filename=$(basename $1)
    final_license=${license//\[filename with brief idea of what it does\]/"$filename"}
    if [[ $(sed -n '1{/^#!/p};q' $1) ]]; then
        cat <(head -1 $1) <(echo -e "\n<<//\n$final_license\n//\n") <(tail -n+2 $1) > "$1.new" && mv "$1.new" "$1"
    else
        cat <(echo -e "$final_license\n") "$1" > "$1.new" && mv "$1.new" "$1"
    fi
}

if [[ -z "$1" ]]; then
    for file in $(find . -name '*.go' -or -name '*.sql' -or -name '*.sh');
    do
        if ! grep -q "GRP-GNU-AGPL" "$file"
        then
            add_license "$file"
        fi
    done
else
    if ! grep -q "GRP-GNU-AGPL" "$1"
    then
        add_license "$1"
    fi
fi
