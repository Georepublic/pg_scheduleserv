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
	"github.com/sirupsen/logrus"
)

const createBreak = `-- name: CreateBreak :one
INSERT INTO breaks (vehicle_id, service, data) VALUES ($1, $2, $3)
RETURNING id, vehicle_id, service, data, created_at, updated_at, deleted
`

type CreateBreakParams struct {
	VehicleID *int64       `json:"vehicle_id,string" example:"1234567890123456789" validate:"required" swaggerignore:"true"`
	Service   *int64       `json:"service"`
	Data      *interface{} `json:"data" swaggertype:"object"`
}

func (q *Queries) DBCreateBreak(ctx context.Context, arg CreateBreakParams) (Break, error) {
	sql, args := createResource("breaks", arg)
	logrus.Debug(sql)
	logrus.Debug(args)
	var i Break
	return_sql := util.GetReturnSql(i)
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	err := row.Scan(
		&i.ID,
		&i.VehicleID,
		&i.Service,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}

const deleteBreak = `-- name: DeleteBreak :exec
UPDATE breaks SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) DBDeleteBreak(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteBreak, id)
	return err
}

const getBreak = `-- name: GetBreak :one
SELECT
  id, vehicle_id, service, data, created_at, updated_at
FROM breaks
WHERE id = $1 AND deleted = FALSE
LIMIT 1
`

type GetBreakRow struct {
	ID        int64       `json:"id"`
	VehicleID int64       `json:"vehicle_id"`
	Service   int64       `json:"service"`
	Data      interface{} `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (q *Queries) DBGetBreak(ctx context.Context, id int64) (GetBreakRow, error) {
	row := q.db.QueryRow(ctx, getBreak, id)
	var i GetBreakRow
	err := row.Scan(
		&i.ID,
		&i.VehicleID,
		&i.Service,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listBreaks = `-- name: ListBreaks :many
SELECT
  id, vehicle_id, service, data, created_at, updated_at
FROM breaks
WHERE vehicle_id = $1 AND deleted = FALSE
ORDER BY created_at
`

type ListBreaksRow struct {
	ID        int64       `json:"id"`
	VehicleID int64       `json:"vehicle_id"`
	Service   int64       `json:"service"`
	Data      interface{} `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (q *Queries) DBListBreaks(ctx context.Context, vehicleID int64) ([]ListBreaksRow, error) {
	rows, err := q.db.Query(ctx, listBreaks, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListBreaksRow{}
	for rows.Next() {
		var i ListBreaksRow
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

const updateBreak = `-- name: UpdateBreak :one
UPDATE breaks
SET vehicle_id = $2, service = $3, data = $4
WHERE id = $1 AND deleted = FALSE
RETURNING id, vehicle_id, service, data, created_at, updated_at, deleted
`

type UpdateBreakParams struct {
	ID        int64       `json:"id"`
	VehicleID int64       `json:"vehicle_id"`
	Service   int64       `json:"service"`
	Data      interface{} `json:"data"`
}

func (q *Queries) DBUpdateBreak(ctx context.Context, arg UpdateBreakParams) (Break, error) {
	row := q.db.QueryRow(ctx, updateBreak,
		arg.ID,
		arg.VehicleID,
		arg.Service,
		arg.Data,
	)
	var i Break
	err := row.Scan(
		&i.ID,
		&i.VehicleID,
		&i.Service,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}
