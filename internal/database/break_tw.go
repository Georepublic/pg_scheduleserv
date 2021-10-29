/*GRP-GNU-AGPL******************************************************************

File: break_tw.go

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
	"github.com/sirupsen/logrus"
)

const createBreakTimeWindow = `-- name: CreateBreakTimeWindow :one
INSERT INTO breaks_time_windows (id, tw_open, tw_close)
VALUES ($1, $2, $3)
RETURNING id, tw_open, tw_close, created_at, updated_at
`

type CreateBreakTimeWindowParams struct {
	ID      *int64  `json:"id,string" example:"1234567890123456789" validate:"required" swaggerignore:"true"`
	TwOpen  *string `json:"tw_open" validate:"required,datetime=2006-01-02 15:04:05"`
	TwClose *string `json:"tw_close" validate:"required,datetime=2006-01-02 15:04:05"`
}

func (q *Queries) DBCreateBreakTimeWindow(ctx context.Context, arg CreateBreakTimeWindowParams) (BreakTimeWindow, error) {
	sql, args := createResource("breaks_time_windows", arg)
	logrus.Debug(sql)
	logrus.Debug(args)
	return_sql := " RETURNING " + util.GetOutputFields(BreakTimeWindow{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanBreakTimeWindowRow(row)
}

const listBreakTimeWindow = `-- name: ListBreakTimeWindow :many
SELECT id, tw_open, tw_close, created_at, updated_at
FROM breaks_time_windows
WHERE id = $1
ORDER BY created_at
`

func (q *Queries) DBListBreakTimeWindow(ctx context.Context, id int64) ([]BreakTimeWindow, error) {
	table_name := "breaks_time_windows"
	additional_query := " WHERE id = $1 ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(BreakTimeWindow{}) + " FROM " + table_name + additional_query
	rows, err := q.db.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBreakTimeWindowRows(rows)
}

const deleteBreakTimeWindow = `-- name: DeleteBreakTimeWindow :exec
DELETE FROM breaks_time_windows
WHERE id = $1 AND tw_open = $2 AND tw_close = $3
`

type DeleteBreakTimeWindowParams struct {
	ID      int64  `json:"id"`
	TwOpen  string `json:"tw_open"`
	TwClose string `json:"tw_close"`
}

func (q *Queries) DBDeleteBreakTimeWindow(ctx context.Context, arg DeleteBreakTimeWindowParams) error {
	sql := "DELETE FROM breaks_time_windows WHERE id = $1 AND tw_open = $2 AND tw_close = $3"
	_, err := q.db.Exec(ctx, sql, arg.ID, arg.TwOpen, arg.TwClose)
	return err
}

func scanBreakTimeWindowRow(row pgx.Row) (BreakTimeWindow, error) {
	var i BreakTimeWindow
	err := row.Scan(
		&i.ID,
		&i.TwOpen,
		&i.TwClose,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func scanBreakTimeWindowRows(rows pgx.Rows) ([]BreakTimeWindow, error) {
	items := []BreakTimeWindow{}
	var i BreakTimeWindow
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
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
