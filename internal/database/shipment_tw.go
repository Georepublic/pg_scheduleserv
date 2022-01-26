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

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
)

type CreateShipmentTimeWindowParams struct {
	ID      *int64  `json:"id,string" example:"1234567812345678" validate:"required" swaggerignore:"true"`
	Kind    *string `json:"kind" validate:"required,oneof=p d" example:"p"`
	TwOpen  *string `json:"tw_open" validate:"required,datetime=2006-01-02T15:04:05" example:"2021-12-31T23:00:00"`
	TwClose *string `json:"tw_close" validate:"required,datetime=2006-01-02T15:04:05" example:"2021-12-31T23:59:00"`
}

func (q *Queries) DBCreateShipmentTimeWindow(ctx context.Context, arg CreateShipmentTimeWindowParams) (ShipmentTimeWindow, error) {
	tableName := "shipments_time_windows"
	sql, args := createResource(tableName, arg)
	return_sql := " RETURNING " + util.GetOutputFields(ShipmentTimeWindow{}, tableName)
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanShipmentTimeWindowRow(row)
}

func (q *Queries) DBListShipmentTimeWindows(ctx context.Context, id int64) ([]ShipmentTimeWindow, error) {
	_, err := q.DBGetShipment(ctx, id)
	if err != nil {
		return nil, err
	}
	tableName := "shipments_time_windows"
	additionalQuery := " WHERE id = $1 ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(ShipmentTimeWindow{}, tableName) + " FROM " + tableName + additionalQuery
	rows, err := q.db.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShipmentTimeWindowRows(rows)
}

func (q *Queries) DBDeleteShipmentTimeWindow(ctx context.Context, id int64) (ShipmentTimeWindow, error) {
	tableName := "shipments_time_windows"
	sql := "DELETE FROM " + tableName + " WHERE id = $1"
	return_sql := " RETURNING " + util.GetOutputFields(ShipmentTimeWindow{}, tableName)
	row := q.db.QueryRow(ctx, sql+return_sql, id)
	return scanShipmentTimeWindowRow(row)
}

func scanShipmentTimeWindowRow(row pgx.Row) (ShipmentTimeWindow, error) {
	var i ShipmentTimeWindow
	err := row.Scan(
		&i.ID,
		&i.Kind,
		&i.TwOpen,
		&i.TwClose,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	err = util.HandleDBError(err)
	return i, err
}

func scanShipmentTimeWindowRows(rows pgx.Rows) ([]ShipmentTimeWindow, error) {
	items := []ShipmentTimeWindow{}
	var i ShipmentTimeWindow
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&i.Kind,
			&i.TwOpen,
			&i.TwClose,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
