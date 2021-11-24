#!/bin/bash

# GRP-GNU-AGPL******************************************************************

# File: check_license.sh - script to check file license

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

EXCLUDE_LIST="internal/mock"

mylicensecheck() {
    licensecheck -r --copyright -l 30 --tail 0 -i "$EXCLUDE_LIST" docs e2e_test internal migrations scripts Makefile main.go
}

missing=$(! { mylicensecheck;}  | grep "No copyright\|UNKNOWN" | grep -v "generated file")

if [[ $missing ]]; then
  echo " ****************************************************"
  echo " *** Found source files without valid license headers"
  echo " ****************************************************"
  echo "$missing"
  exit 1
fi
