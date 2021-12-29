/*GRP-GNU-AGPL******************************************************************

File: break.go

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

type CreateBreakParams struct {
	VehicleID *int64       `json:"vehicle_id,string" example:"1234567812345678" validate:"required" swaggerignore:"true"`
	Service   *string      `json:"service" validate:"omitempty" example:"00:02:00"`
	Data      *interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

type UpdateBreakParams struct {
	VehicleID *int64       `json:"vehicle_id,string" example:"1234567812345678" swaggerignore:"true"`
	Service   *string      `json:"service" validate:"omitempty" example:"00:02:00"`
	Data      *interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

func (q *Queries) DBCreateBreak(ctx context.Context, arg CreateBreakParams) (Break, error) {
	sql, args := createResource("breaks", arg)
	return_sql := " RETURNING " + util.GetOutputFields(Break{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanBreakRow(row)
}

func (q *Queries) DBGetBreak(ctx context.Context, id int64) (Break, error) {
	table_name := "breaks"
	additional_query := " WHERE id = $1 AND deleted = FALSE LIMIT 1"
	sql := "SELECT " + util.GetOutputFields(Break{}) + " FROM " + table_name + additional_query
	row := q.db.QueryRow(ctx, sql, id)
	return scanBreakRow(row)
}

func (q *Queries) DBListBreaks(ctx context.Context, vehicleID int64) ([]Break, error) {
	table_name := "breaks"
	additional_query := " WHERE vehicle_id = $1 AND deleted = FALSE ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(Break{}) + " FROM " + table_name + additional_query
	rows, err := q.db.Query(ctx, sql, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBreakRows(rows)
}

func (q *Queries) DBUpdateBreak(ctx context.Context, arg UpdateBreakParams, break_id int64) (Break, error) {
	sql, args := updateResource("breaks", arg, break_id)
	return_sql := " RETURNING " + util.GetOutputFields(Break{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanBreakRow(row)
}

func (q *Queries) DBDeleteBreak(ctx context.Context, id int64) (Break, error) {
	sql := "UPDATE breaks SET deleted = TRUE WHERE id = $1"
	return_sql := " RETURNING " + util.GetOutputFields(Break{})
	row := q.db.QueryRow(ctx, sql+return_sql, id)
	return scanBreakRow(row)
}

func scanBreakRow(row pgx.Row) (Break, error) {
	var i Break
	err := row.Scan(
		&i.ID,
		&i.VehicleID,
		&i.Service,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	err = util.HandleDBError(err)
	return i, err
}

func scanBreakRows(rows pgx.Rows) ([]Break, error) {
	items := []Break{}
	var i Break
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&i.VehicleID,
			&i.Service,
			&i.Data,
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
