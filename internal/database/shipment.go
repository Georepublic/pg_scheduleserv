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
	"github.com/jackc/pgx/v4"
)

type CreateShipmentParams struct {
	PLocation *util.LocationParams `json:"p_location" validate:"required"`
	PSetup    *string              `json:"p_setup"    validate:"omitempty" example:"00:00:00"`
	PService  *string              `json:"p_service"  validate:"omitempty" example:"00:02:00"`
	DLocation *util.LocationParams `json:"d_location" validate:"required"`
	DSetup    *string              `json:"d_setup"    validate:"omitempty" example:"00:00:00"`
	DService  *string              `json:"d_service"  validate:"omitempty" example:"00:02:00"`
	Amount    *[]int64             `json:"amount"     validate:"omitempty,dive,min=0" example:"5,15"`
	Skills    *[]int32             `json:"skills"     validate:"omitempty,dive,min=0" example:"1,5"`
	Priority  *int32               `json:"priority"   validate:"omitempty,min=0,max=100" example:"10"`
	ProjectID *int64               `json:"project_id,string" validate:"required" swaggerignore:"true"`
	Data      *interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

type UpdateShipmentParams struct {
	PLocation *util.LocationParams `json:"p_location"`
	PSetup    *string              `json:"p_setup"    validate:"omitempty" example:"00:00:00"`
	PService  *string              `json:"p_service"  validate:"omitempty" example:"00:02:00"`
	DLocation *util.LocationParams `json:"d_location"`
	DSetup    *string              `json:"d_setup"    validate:"omitempty" example:"00:00:00"`
	DService  *string              `json:"d_service"  validate:"omitempty" example:"00:02:00"`
	Amount    *[]int64             `json:"amount"     validate:"omitempty,dive,min=0" example:"5,15"`
	Skills    *[]int32             `json:"skills"     validate:"omitempty,dive,min=0" example:"1,5"`
	Priority  *int32               `json:"priority"   validate:"omitempty,min=0,max=100" example:"10"`
	ProjectID *int64               `json:"project_id,string" swaggerignore:"true"`
	Data      *interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

func (q *Queries) DBCreateShipment(ctx context.Context, arg CreateShipmentParams) (Shipment, error) {
	tableName := "shipments"
	sql, args := createResource(tableName, arg)
	return_sql := " RETURNING " + util.GetOutputFields(Shipment{}, tableName)
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanShipmentRow(row)
}

func (q *Queries) DBGetShipment(ctx context.Context, id int64) (Shipment, error) {
	tableName := "shipments"
	additionalQuery := " WHERE id = $1 AND deleted = FALSE LIMIT 1"
	sql := "SELECT " + util.GetOutputFields(Shipment{}, tableName) + " FROM " + tableName + additionalQuery
	row := q.db.QueryRow(ctx, sql, id)
	return scanShipmentRow(row)
}

func (q *Queries) DBListShipments(ctx context.Context, projectID int64) ([]Shipment, error) {
	_, err := q.DBGetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	tableName := "shipments"
	additionalQuery := " WHERE project_id = $1 AND deleted = FALSE ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(Shipment{}, tableName) + " FROM " + tableName + additionalQuery
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShipmentRows(rows)
}

func (q *Queries) DBUpdateShipment(ctx context.Context, arg UpdateShipmentParams, shipment_id int64) (Shipment, error) {
	tableName := "shipments"
	sql, args := updateResource(tableName, arg, shipment_id)
	return_sql := " RETURNING " + util.GetOutputFields(Shipment{}, tableName)
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanShipmentRow(row)
}

func (q *Queries) DBDeleteShipment(ctx context.Context, id int64) (Shipment, error) {
	tableName := "shipments"
	sql := "UPDATE " + tableName + " SET deleted = TRUE WHERE id = $1"
	return_sql := " RETURNING " + util.GetOutputFields(Shipment{}, tableName)
	row := q.db.QueryRow(ctx, sql+return_sql, id)
	return scanShipmentRow(row)
}

func scanShipmentRow(row pgx.Row) (Shipment, error) {
	var i Shipment
	var p_location_id, d_location_id int64
	err := row.Scan(
		&i.ID,
		&p_location_id,
		&i.PSetup,
		&i.PService,
		&d_location_id,
		&i.DSetup,
		&i.DService,
		&i.Amount,
		&i.Skills,
		&i.Priority,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	p_latitude, p_longitude := util.GetCoordinates(p_location_id)
	d_latitude, d_longitude := util.GetCoordinates(d_location_id)
	i.PLocation = util.LocationParams{
		Latitude:  &p_latitude,
		Longitude: &p_longitude,
	}
	i.DLocation = util.LocationParams{
		Latitude:  &d_latitude,
		Longitude: &d_longitude,
	}
	err = util.HandleDBError(err)
	return i, err
}

func scanShipmentRows(rows pgx.Rows) ([]Shipment, error) {
	items := []Shipment{}
	var i Shipment
	var p_location_id, d_location_id int64
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&p_location_id,
			&i.PSetup,
			&i.PService,
			&d_location_id,
			&i.DSetup,
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
		p_latitude, p_longitude := util.GetCoordinates(p_location_id)
		d_latitude, d_longitude := util.GetCoordinates(d_location_id)
		i.PLocation = util.LocationParams{
			Latitude:  &p_latitude,
			Longitude: &p_longitude,
		}
		i.DLocation = util.LocationParams{
			Latitude:  &d_latitude,
			Longitude: &d_longitude,
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
