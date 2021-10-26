/*GRP-GNU-AGPL******************************************************************

File: shipment.go

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

const createShipment = `-- name: CreateShipment :one
INSERT INTO shipments (
  p_location_index, p_service, d_location_index, d_service,
  amount, skills, priority, project_id, data
) VALUES (
  coord_to_id($1, $2), $3, coord_to_id($4, $5), $6, $7, $8, $9, $10, $11
)
RETURNING id, p_location_index, p_service, d_location_index, d_service, amount, skills, priority, project_id, data, created_at, updated_at, deleted
`

type CreateShipmentParams struct {
	PLocation *util.LocationParams `json:"p_location" validate:"required"`
	PService  *int64               `json:"p_service"`
	DLocation *util.LocationParams `json:"d_location" validate:"required"`
	DService  *int64               `json:"d_service"`
	Amount    *[]int64             `json:"amount"`
	Skills    *[]int32             `json:"skills"`
	Priority  *int32               `json:"priority"`
	ProjectID *int64               `json:"project_id,string" validate:"required" swaggerignore:"true"`
	Data      *interface{}         `json:"data" swaggertype:"object"`
}

func (q *Queries) DBCreateShipment(ctx context.Context, arg CreateShipmentParams) (Shipment, error) {
	sql, args := createResource("shipments", arg)
	logrus.Debug(sql)
	logrus.Debug(args)
	var i Shipment
	return_sql := util.GetReturnSql(i)
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	var p_location_index, d_location_index int64
	err := row.Scan(
		&i.ID,
		&p_location_index,
		&i.PService,
		&d_location_index,
		&i.DService,
		&i.Amount,
		&i.Skills,
		&i.Priority,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	p_latitude, p_longitude := util.GetCoordinates(p_location_index)
	d_latitude, d_longitude := util.GetCoordinates(d_location_index)
	i.PLocation = util.LocationParams{
		Latitude:  &p_latitude,
		Longitude: &p_longitude,
	}
	i.DLocation = util.LocationParams{
		Latitude:  &d_latitude,
		Longitude: &d_longitude,
	}
	return i, err
}

const deleteShipment = `-- name: DeleteShipment :exec
UPDATE shipments SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) DBDeleteShipment(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteShipment, id)
	return err
}

const getShipment = `-- name: GetShipment :one
SELECT
  id, p_location_index, p_service, d_location_index, d_service,
  amount, skills, priority, project_id, data, created_at, updated_at
FROM shipments
WHERE id = $1 AND deleted = FALSE
LIMIT 1
`

type GetShipmentRow struct {
	ID             int64       `json:"id"`
	PLocationIndex int64       `json:"p_location_index"`
	PService       int64       `json:"p_service"`
	DLocationIndex int64       `json:"d_location_index"`
	DService       int64       `json:"d_service"`
	Amount         []int64     `json:"amount"`
	Skills         []int32     `json:"skills"`
	Priority       int32       `json:"priority"`
	ProjectID      int64       `json:"project_id"`
	Data           interface{} `json:"data"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
}

func (q *Queries) DBGetShipment(ctx context.Context, id int64) (GetShipmentRow, error) {
	row := q.db.QueryRow(ctx, getShipment, id)
	var i GetShipmentRow
	err := row.Scan(
		&i.ID,
		&i.PLocationIndex,
		&i.PService,
		&i.DLocationIndex,
		&i.DService,
		&i.Amount,
		&i.Skills,
		&i.Priority,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listShipments = `-- name: ListShipments :many
SELECT
  id, p_location_index, p_service, d_location_index, d_service,
  amount, skills, priority, project_id, data, created_at, updated_at
FROM shipments
WHERE project_id = $1 AND deleted = FALSE
ORDER BY created_at
`

type ListShipmentsRow struct {
	ID             int64       `json:"id"`
	PLocationIndex int64       `json:"p_location_index"`
	PService       int64       `json:"p_service"`
	DLocationIndex int64       `json:"d_location_index"`
	DService       int64       `json:"d_service"`
	Amount         []int64     `json:"amount"`
	Skills         []int32     `json:"skills"`
	Priority       int32       `json:"priority"`
	ProjectID      int64       `json:"project_id"`
	Data           interface{} `json:"data"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
}

func (q *Queries) DBListShipments(ctx context.Context, projectID int64) ([]ListShipmentsRow, error) {
	rows, err := q.db.Query(ctx, listShipments, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListShipmentsRow{}
	for rows.Next() {
		var i ListShipmentsRow
		if err := rows.Scan(
			&i.ID,
			&i.PLocationIndex,
			&i.PService,
			&i.DLocationIndex,
			&i.DService,
			&i.Amount,
			&i.Skills,
			&i.Priority,
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

const updateShipment = `-- name: UpdateShipment :one
UPDATE shipments
SET
  p_location_index = coord_to_id($2, $3), p_service = $4,
  d_location_index = coord_to_id($5, $6), d_service = $7,
  amount = $8, skills = $9, priority = $10,
  project_id = $11, data = $12
WHERE id = $1 AND deleted = FALSE
RETURNING id, p_location_index, p_service, d_location_index, d_service, amount, skills, priority, project_id, data, created_at, updated_at, deleted
`

type UpdateShipmentParams struct {
	ID          int64       `json:"id"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	PService    int64       `json:"p_service"`
	Latitude_2  float64     `json:"latitude_2"`
	Longitude_2 float64     `json:"longitude_2"`
	DService    int64       `json:"d_service"`
	Amount      []int64     `json:"amount"`
	Skills      []int32     `json:"skills"`
	Priority    int32       `json:"priority"`
	ProjectID   int64       `json:"project_id"`
	Data        interface{} `json:"data"`
}

func (q *Queries) DBUpdateShipment(ctx context.Context, arg UpdateShipmentParams) (Shipment, error) {
	row := q.db.QueryRow(ctx, updateShipment,
		arg.ID,
		arg.Latitude,
		arg.Longitude,
		arg.PService,
		arg.Latitude_2,
		arg.Longitude_2,
		arg.DService,
		arg.Amount,
		arg.Skills,
		arg.Priority,
		arg.ProjectID,
		arg.Data,
	)
	var i Shipment
	var p_location_index, d_location_index int64
	err := row.Scan(
		&i.ID,
		&p_location_index,
		&i.PService,
		&d_location_index,
		&i.DService,
		&i.Amount,
		&i.Skills,
		&i.Priority,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}
