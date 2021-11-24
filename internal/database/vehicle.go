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

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
)

type CreateVehicleParams struct {
	StartLocation *util.LocationParams `json:"start_location" validate:"required"`
	EndLocation   *util.LocationParams `json:"end_location" validate:"required"`
	Capacity      *[]int64             `json:"capacity"`
	Skills        *[]int32             `json:"skills"`
	TwOpen        *string              `json:"tw_open" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	TwClose       *string              `json:"tw_close" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	SpeedFactor   *float64             `json:"speed_factor"`
	ProjectID     *int64               `json:"project_id,string" validate:"required" swaggerignore:"true"`
	Data          *interface{}         `json:"data" swaggertype:"object"`
}

func (q *Queries) DBCreateVehicle(ctx context.Context, arg CreateVehicleParams) (Vehicle, error) {
	sql, args := createResource("vehicles", arg)
	return_sql := " RETURNING " + util.GetOutputFields(Vehicle{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanVehicleRow(row)
}

type GetVehicleRow struct {
	ID          int64       `json:"id"`
	StartIndex  int64       `json:"start_index"`
	EndIndex    int64       `json:"end_index"`
	Capacity    []int64     `json:"capacity"`
	Skills      []int32     `json:"skills"`
	TwOpen      string      `json:"tw_open"`
	TwClose     string      `json:"tw_close"`
	SpeedFactor float64     `json:"speed_factor"`
	ProjectID   int64       `json:"project_id"`
	Data        interface{} `json:"data"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

func (q *Queries) DBGetVehicle(ctx context.Context, id int64) (Vehicle, error) {
	table_name := "vehicles"
	additional_query := " WHERE id = $1 AND deleted = FALSE LIMIT 1"
	sql := "SELECT " + util.GetOutputFields(Vehicle{}) + " FROM " + table_name + additional_query
	row := q.db.QueryRow(ctx, sql, id)
	return scanVehicleRow(row)
}

type ListVehiclesRow struct {
	ID          int64       `json:"id"`
	StartIndex  int64       `json:"start_index"`
	EndIndex    int64       `json:"end_index"`
	Capacity    []int64     `json:"capacity"`
	Skills      []int32     `json:"skills"`
	TwOpen      string      `json:"tw_open"`
	TwClose     string      `json:"tw_close"`
	SpeedFactor float64     `json:"speed_factor"`
	ProjectID   int64       `json:"project_id"`
	Data        interface{} `json:"data"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

func (q *Queries) DBListVehicles(ctx context.Context, projectID int64) ([]Vehicle, error) {
	table_name := "vehicles"
	additional_query := " WHERE project_id = $1 AND deleted = FALSE ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(Vehicle{}) + " FROM " + table_name + additional_query
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVehicleRows(rows)
}

type UpdateVehicleParams struct {
	StartLocation *util.LocationParams `json:"start_location"`
	EndLocation   *util.LocationParams `json:"end_location"`
	Capacity      *[]int64             `json:"capacity"`
	Skills        *[]int32             `json:"skills"`
	TwOpen        *string              `json:"tw_open" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	TwClose       *string              `json:"tw_close" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	SpeedFactor   *float64             `json:"speed_factor"`
	ProjectID     *int64               `json:"project_id,string"`
	Data          *interface{}         `json:"data" swaggertype:"object"`
}

func (q *Queries) DBUpdateVehicle(ctx context.Context, arg UpdateVehicleParams, vehicle_id int64) (Vehicle, error) {
	sql, args := updateResource("vehicles", arg, vehicle_id)
	return_sql := " RETURNING " + util.GetOutputFields(Vehicle{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanVehicleRow(row)
}

func (q *Queries) DBDeleteVehicle(ctx context.Context, id int64) (Vehicle, error) {
	sql := "UPDATE vehicles SET deleted = TRUE WHERE id = $1"
	return_sql := " RETURNING " + util.GetOutputFields(Vehicle{})
	row := q.db.QueryRow(ctx, sql+return_sql, id)
	return scanVehicleRow(row)
}

func scanVehicleRow(row pgx.Row) (Vehicle, error) {
	var i Vehicle
	var start_index, end_index int64
	err := row.Scan(
		&i.ID,
		&start_index,
		&end_index,
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
	start_latitude, start_longitude := util.GetCoordinates(start_index)
	end_latitude, end_longitude := util.GetCoordinates(end_index)
	i.StartLocation = util.LocationParams{
		Latitude:  &start_latitude,
		Longitude: &start_longitude,
	}
	i.EndLocation = util.LocationParams{
		Latitude:  &end_latitude,
		Longitude: &end_longitude,
	}
	return i, err
}

func scanVehicleRows(rows pgx.Rows) ([]Vehicle, error) {
	var i Vehicle
	items := []Vehicle{}
	var start_index, end_index int64
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&start_index,
			&end_index,
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
		start_latitude, start_longitude := util.GetCoordinates(start_index)
		end_latitude, end_longitude := util.GetCoordinates(end_index)
		i.StartLocation = util.LocationParams{
			Latitude:  &start_latitude,
			Longitude: &start_longitude,
		}
		i.EndLocation = util.LocationParams{
			Latitude:  &end_latitude,
			Longitude: &end_longitude,
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
