/*GRP-GNU-AGPL******************************************************************

File: shipment_tw.go

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
	"fmt"
)

type ShipmentTimeWindowParams struct {
	Kind    string `json:"kind" validate:"required,oneof=p d" example:"p"`
	TwOpen  string `json:"tw_open" validate:"required,datetime=2006-01-02T15:04:05" example:"2021-12-31T23:00:00"`
	TwClose string `json:"tw_close" validate:"required,datetime=2006-01-02T15:04:05" example:"2021-12-31T23:59:00"`
}

func (q *Queries) DBCreateShipmentTimeWindows(ctx context.Context, id int64, arg []ShipmentTimeWindowParams) error {
	if len(arg) == 0 {
		return nil
	}

	// create an sql query to insert multiple rows
	sql := "INSERT INTO shipments_time_windows (id, kind, tw_open, tw_close) VALUES "
	for i := range arg {
		if i > 0 {
			sql += ","
		}
		sql += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
	}

	// create a slice of arguments for the query
	args := make([]interface{}, len(arg)*4)
	for i, v := range arg {
		args[i*4] = id
		args[i*4+1] = v.Kind
		args[i*4+2] = v.TwOpen
		args[i*4+3] = v.TwClose
	}

	// execute the query
	_, err := q.db.Exec(ctx, sql, args...)
	return err
}

func (q *Queries) DBDeleteShipmentTimeWindows(ctx context.Context, id int64) error {
	tableName := "shipments_time_windows"
	sql := "DELETE FROM " + tableName + " WHERE id = $1"
	_, err := q.db.Exec(ctx, sql, id)
	return err
}
