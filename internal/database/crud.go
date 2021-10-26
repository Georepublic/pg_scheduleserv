/*GRP-GNU-AGPL******************************************************************

File: crud.go

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

package database

import (
	"fmt"

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/sirupsen/logrus"
)

// A utility function for creating any resource
// Takes the resource name and resource object as parameter,
// and returns the sql query with the input arguments
// (only those arguments which are not nil)
func createResource(resource string, resourceStruct interface{}) (sql string, args []interface{}) {
	partialSQL := util.GetPartialSQL(resourceStruct)
	logrus.Debugf("Resource struct: %+v", resourceStruct)
	logrus.Debugf("Partial SQL: %+v", partialSQL)
	sqlFields := ""
	values := ""
	for i, field := range partialSQL.Fields {
		val := fmt.Sprintf("$%d", i+1)

		// Convert any interval field to its type
		if _, intervalFieldFound := util.IntervalFields[field]; intervalFieldFound {
			val = val + "* '1 sec'::INTERVAL"
		}

		if i == 0 {
			sqlFields += field
			values += val
		} else {
			sqlFields += ", " + field
			values += ", " + val
		}
	}
	sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", resource, sqlFields, values)
	args = partialSQL.Args
	return
}
