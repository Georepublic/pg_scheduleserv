/*GRP-GNU-AGPL******************************************************************

File: vehicle.go

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
	"time"

	"github.com/jackc/pgtype"
)

const createVehicle = `-- name: CreateVehicle :one
/*
POST /projects/{project_id}/vehicles
GET /projects/{project_id}/vehicles

GET /vehicles/{vehicle_id}
PATCH /vehicles/{vehicle_id}
DELETE /vehicles/{vehicle_id}
*/

INSERT INTO vehicles (
  start_index, end_index, capacity, skills,
  tw_open, tw_close, speed_factor, project_id, data
) VALUES (
  coord_to_id($1, $2), coord_to_id($3, $4), $5, $6, $7, $8, $9, $10, $11
)
RETURNING id, start_index, end_index, capacity, skills, tw_open, tw_close, speed_factor, project_id, data, created_at, updated_at, deleted
`

type CreateVehicleParams struct {
	Latitude    float64      `json:"latitude"`
	Longitude   float64      `json:"longitude"`
	Latitude_2  float64      `json:"latitude_2"`
	Longitude_2 float64      `json:"longitude_2"`
	Capacity    []int64      `json:"capacity"`
	Skills      []int32      `json:"skills"`
	TwOpen      time.Time    `json:"tw_open"`
	TwClose     time.Time    `json:"tw_close"`
	SpeedFactor float64      `json:"speed_factor"`
	ProjectID   int64        `json:"project_id"`
	Data        pgtype.JSONB `json:"data"`
}

func (q *Queries) CreateVehicle(ctx context.Context, arg CreateVehicleParams) (Vehicle, error) {
	row := q.db.QueryRow(ctx, createVehicle,
		arg.Latitude,
		arg.Longitude,
		arg.Latitude_2,
		arg.Longitude_2,
		arg.Capacity,
		arg.Skills,
		arg.TwOpen,
		arg.TwClose,
		arg.SpeedFactor,
		arg.ProjectID,
		arg.Data,
	)
	var i Vehicle
	err := row.Scan(
		&i.ID,
		&i.StartIndex,
		&i.EndIndex,
		&i.Capacity,
		&i.Skills,
		&i.TwOpen,
		&i.TwClose,
		&i.SpeedFactor,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}

const deleteVehicle = `-- name: DeleteVehicle :exec
UPDATE vehicles SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) DeleteVehicle(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteVehicle, id)
	return err
}

const getVehicle = `-- name: GetVehicle :one
SELECT
  id, start_index, end_index, capacity, skills,
  tw_open, tw_close, speed_factor, project_id,
  data, created_at, updated_at
FROM vehicles
WHERE id = $1 AND deleted = FALSE
LIMIT 1
`

type GetVehicleRow struct {
	ID          int64        `json:"id"`
	StartIndex  int64        `json:"start_index"`
	EndIndex    int64        `json:"end_index"`
	Capacity    []int64      `json:"capacity"`
	Skills      []int32      `json:"skills"`
	TwOpen      time.Time    `json:"tw_open"`
	TwClose     time.Time    `json:"tw_close"`
	SpeedFactor float64      `json:"speed_factor"`
	ProjectID   int64        `json:"project_id"`
	Data        pgtype.JSONB `json:"data"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (q *Queries) GetVehicle(ctx context.Context, id int64) (GetVehicleRow, error) {
	row := q.db.QueryRow(ctx, getVehicle, id)
	var i GetVehicleRow
	err := row.Scan(
		&i.ID,
		&i.StartIndex,
		&i.EndIndex,
		&i.Capacity,
		&i.Skills,
		&i.TwOpen,
		&i.TwClose,
		&i.SpeedFactor,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listVehicles = `-- name: ListVehicles :many
SELECT
  id, start_index, end_index, capacity, skills,
  tw_open, tw_close, speed_factor, project_id,
  data, created_at, updated_at
FROM vehicles
WHERE project_id = $1 AND deleted = FALSE
ORDER BY created_at
`

type ListVehiclesRow struct {
	ID          int64        `json:"id"`
	StartIndex  int64        `json:"start_index"`
	EndIndex    int64        `json:"end_index"`
	Capacity    []int64      `json:"capacity"`
	Skills      []int32      `json:"skills"`
	TwOpen      time.Time    `json:"tw_open"`
	TwClose     time.Time    `json:"tw_close"`
	SpeedFactor float64      `json:"speed_factor"`
	ProjectID   int64        `json:"project_id"`
	Data        pgtype.JSONB `json:"data"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (q *Queries) ListVehicles(ctx context.Context, projectID int64) ([]ListVehiclesRow, error) {
	rows, err := q.db.Query(ctx, listVehicles, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListVehiclesRow{}
	for rows.Next() {
		var i ListVehiclesRow
		if err := rows.Scan(
			&i.ID,
			&i.StartIndex,
			&i.EndIndex,
			&i.Capacity,
			&i.Skills,
			&i.TwOpen,
			&i.TwClose,
			&i.SpeedFactor,
			&i.ProjectID,
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

const updateVehicle = `-- name: UpdateVehicle :one
UPDATE vehicles
SET
  start_index = coord_to_id($2, $3), end_index = coord_to_id($4, $5),
  capacity = $6, skills = $7, tw_open = $8, tw_close = $9,
  speed_factor = $10, project_id = $11, data = $12
WHERE id = $1 AND deleted = FALSE
RETURNING id, start_index, end_index, capacity, skills, tw_open, tw_close, speed_factor, project_id, data, created_at, updated_at, deleted
`

type UpdateVehicleParams struct {
	ID          int64        `json:"id"`
	Latitude    float64      `json:"latitude"`
	Longitude   float64      `json:"longitude"`
	Latitude_2  float64      `json:"latitude_2"`
	Longitude_2 float64      `json:"longitude_2"`
	Capacity    []int64      `json:"capacity"`
	Skills      []int32      `json:"skills"`
	TwOpen      time.Time    `json:"tw_open"`
	TwClose     time.Time    `json:"tw_close"`
	SpeedFactor float64      `json:"speed_factor"`
	ProjectID   int64        `json:"project_id"`
	Data        pgtype.JSONB `json:"data"`
}

func (q *Queries) UpdateVehicle(ctx context.Context, arg UpdateVehicleParams) (Vehicle, error) {
	row := q.db.QueryRow(ctx, updateVehicle,
		arg.ID,
		arg.Latitude,
		arg.Longitude,
		arg.Latitude_2,
		arg.Longitude_2,
		arg.Capacity,
		arg.Skills,
		arg.TwOpen,
		arg.TwClose,
		arg.SpeedFactor,
		arg.ProjectID,
		arg.Data,
	)
	var i Vehicle
	err := row.Scan(
		&i.ID,
		&i.StartIndex,
		&i.EndIndex,
		&i.Capacity,
		&i.Skills,
		&i.TwOpen,
		&i.TwClose,
		&i.SpeedFactor,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}
