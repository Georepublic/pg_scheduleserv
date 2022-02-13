/*GRP-GNU-AGPL******************************************************************

File: location.go

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
	"context"

	"github.com/jackc/pgx/v4"
)

func (q *Queries) DBGetProjectLocations(ctx context.Context, project_id int64) ([]int64, error) {
	sql := `
	SELECT unnest(ARRAY[location_id]) AS location_id FROM jobs WHERE project_id = $1 AND deleted=FALSE UNION
	SELECT unnest(ARRAY[p_location_id, d_location_id]) FROM shipments WHERE project_id = $1 AND deleted=FALSE UNION
	SELECT unnest(ARRAY[start_id, end_id]) FROM vehicles WHERE project_id = $1 AND deleted=FALSE`

	rows, err := q.db.Query(ctx, sql, project_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProjectLocationRows(rows)
}

func scanProjectLocationRows(rows pgx.Rows) ([]int64, error) {
	var locations []int64
	for rows.Next() {
		var location_id int64
		err := rows.Scan(&location_id)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location_id)
	}
	return locations, nil
}
