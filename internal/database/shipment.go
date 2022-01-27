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
	"fmt"

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
)

type CreateShipmentParams struct {
	PLocation    *util.LocationParams `json:"p_location" validate:"required"`
	PSetup       *string              `json:"p_setup"    validate:"omitempty" example:"00:00:00"`
	PService     *string              `json:"p_service"  validate:"omitempty" example:"00:02:00"`
	DLocation    *util.LocationParams `json:"d_location" validate:"required"`
	DSetup       *string              `json:"d_setup"    validate:"omitempty" example:"00:00:00"`
	DService     *string              `json:"d_service"  validate:"omitempty" example:"00:02:00"`
	Amount       *[]int64             `json:"amount"     validate:"omitempty,dive,min=0" example:"5,15"`
	Skills       *[]int32             `json:"skills"     validate:"omitempty,dive,min=0" example:"1,5"`
	PTimeWindows *[][]string          `json:"p_time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
	DTimeWindows *[][]string          `json:"d_time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
	Priority     *int32               `json:"priority"   validate:"omitempty,min=0,max=100" example:"10"`
	ProjectID    *int64               `json:"project_id,string" validate:"required" swaggerignore:"true"`
	Data         *interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

type UpdateShipmentParams struct {
	PLocation    *util.LocationParams `json:"p_location"`
	PSetup       *string              `json:"p_setup"    validate:"omitempty" example:"00:00:00"`
	PService     *string              `json:"p_service"  validate:"omitempty" example:"00:02:00"`
	DLocation    *util.LocationParams `json:"d_location"`
	DSetup       *string              `json:"d_setup"    validate:"omitempty" example:"00:00:00"`
	DService     *string              `json:"d_service"  validate:"omitempty" example:"00:02:00"`
	Amount       *[]int64             `json:"amount"     validate:"omitempty,dive,min=0" example:"5,15"`
	Skills       *[]int32             `json:"skills"     validate:"omitempty,dive,min=0" example:"1,5"`
	PTimeWindows *[][]string          `json:"p_time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
	DTimeWindows *[][]string          `json:"d_time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
	Priority     *int32               `json:"priority"   validate:"omitempty,min=0,max=100" example:"10"`
	ProjectID    *int64               `json:"project_id,string" swaggerignore:"true"`
	Data         *interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

func (q *Queries) DBCreateShipment(ctx context.Context, arg CreateShipmentParams) (int64, error) {
	tableName := "shipments"
	sql, args := createResource(tableName, arg)
	return_sql := " RETURNING id"
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanID(row)
}

func (q *Queries) DBGetShipment(ctx context.Context, id int64) (Shipment, error) {
	tableName := "shipments"
	joinSelectQuery := fmt.Sprintf(
		", array_agg(kind), array_agg(ARRAY[%s, %s])",
		util.GetFormattedTimestamp("tw_open"),
		util.GetFormattedTimestamp("tw_close"),
	)
	joinTableQuery := fmt.Sprintf(" LEFT JOIN shipments_time_windows TW on(%s.id = TW.id) ", tableName)
	additionalQuery := fmt.Sprintf(" WHERE %s.id = $1 AND deleted = FALSE GROUP BY %s.id LIMIT 1", tableName, tableName)
	sql := "SELECT " + util.GetOutputFields(Shipment{}, tableName) + joinSelectQuery + " FROM " + tableName + joinTableQuery + additionalQuery
	row := q.db.QueryRow(ctx, sql, id)
	return scanShipmentRow(row)
}

func (q *Queries) DBListShipments(ctx context.Context, projectID int64) ([]Shipment, error) {
	_, err := q.DBGetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	tableName := "shipments"
	joinSelectQuery := fmt.Sprintf(
		", array_agg(kind), array_agg(ARRAY[%s, %s])",
		util.GetFormattedTimestamp("tw_open"),
		util.GetFormattedTimestamp("tw_close"),
	)
	joinTableQuery := fmt.Sprintf(" LEFT JOIN shipments_time_windows TW on(%s.id = TW.id) ", tableName)
	additionalQuery := fmt.Sprintf(" WHERE project_id = $1 AND deleted = FALSE GROUP BY %s.id ORDER BY %s.created_at", tableName, tableName)
	sql := "SELECT " + util.GetOutputFields(Shipment{}, tableName) + joinSelectQuery + " FROM " + tableName + joinTableQuery + additionalQuery
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanShipmentRows(rows)
}

func (q *Queries) DBUpdateShipment(ctx context.Context, arg UpdateShipmentParams, shipment_id int64) error {
	tableName := "shipments"
	sql, args := updateResource(tableName, arg, shipment_id)
	_, err := q.db.Exec(ctx, sql, args...)
	err = util.HandleDBError(err)
	return err
}

func (q *Queries) DBDeleteShipment(ctx context.Context, id int64) error {
	tableName := "shipments"
	sql := "UPDATE " + tableName + " SET deleted = TRUE WHERE id = $1"
	_, err := q.db.Exec(ctx, sql, id)
	return err
}

func (q *Queries) DBCreateShipmentWithTw(ctx context.Context, arg CreateShipmentParams) (Shipment, error) {
	id, err := q.execCreateTx(ctx, func(q *Queries) (int64, error) {
		id, err := q.DBCreateShipment(ctx, arg)
		if err != nil {
			return 0, err
		}

		// create time windows from arg and pass to DBCreateShipmentTimeWindows
		timeWindows := []ShipmentTimeWindowParams{}
		if arg.PTimeWindows != nil {
			for _, tw := range *arg.PTimeWindows {
				// append time window to timeWindows
				timeWindows = append(timeWindows, ShipmentTimeWindowParams{
					Kind:    tw[0],
					TwOpen:  tw[1],
					TwClose: tw[2],
				})
			}
		}
		if arg.DTimeWindows != nil {
			for _, tw := range *arg.DTimeWindows {
				timeWindows = append(timeWindows, ShipmentTimeWindowParams{
					Kind:    tw[0],
					TwOpen:  tw[1],
					TwClose: tw[2],
				})
			}
		}
		err = q.DBCreateShipmentTimeWindows(ctx, id, timeWindows)
		return id, err
	})
	if err != nil {
		return Shipment{}, err
	}
	return q.DBGetShipment(ctx, id)
}

func (q *Queries) DBUpdateShipmentWithTw(ctx context.Context, arg UpdateShipmentParams, shipment_id int64) (Shipment, error) {
	err := q.execUpdateTx(ctx, func(q *Queries) error {
		if err := q.DBUpdateShipment(ctx, arg, shipment_id); err != nil {
			return err
		}
		// delete all time windows
		if err := q.DBDeleteShipmentTimeWindows(ctx, shipment_id); err != nil {
			return err
		}

		// create time windows from arg and pass to DBCreateShipmentTimeWindows
		timeWindows := []ShipmentTimeWindowParams{}
		if arg.PTimeWindows != nil {
			for _, tw := range *arg.PTimeWindows {
				// append time window to timeWindows
				timeWindows = append(timeWindows, ShipmentTimeWindowParams{
					Kind:    "p",
					TwOpen:  tw[0],
					TwClose: tw[1],
				})
			}
		}
		if arg.DTimeWindows != nil {
			for _, tw := range *arg.DTimeWindows {
				timeWindows = append(timeWindows, ShipmentTimeWindowParams{
					Kind:    "d",
					TwOpen:  tw[0],
					TwClose: tw[1],
				})
			}
		}
		err := q.DBCreateShipmentTimeWindows(ctx, shipment_id, timeWindows)
		return err
	})
	if err != nil {
		return Shipment{}, err
	}
	return q.DBGetShipment(ctx, shipment_id)
}

func (q *Queries) DBDeleteShipmentWithTw(ctx context.Context, shipment_id int64) error {
	err := q.execUpdateTx(ctx, func(q *Queries) error {
		// delete all time windows
		if err := q.DBDeleteShipmentTimeWindows(ctx, shipment_id); err != nil {
			return err
		}
		err := q.DBDeleteShipment(ctx, shipment_id)
		return err
	})
	return err
}

func scanShipmentRow(row pgx.Row) (Shipment, error) {
	var i Shipment
	var p_location_id, d_location_id int64
	var kind []*string
	var timeWindows [][]*string
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
		&kind,
		&timeWindows,
	)

	i.PTimeWindows, i.DTimeWindows = util.GetShipmentTimeWindows(kind, timeWindows)

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
	var kind []*string
	var timeWindows [][]*string
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
			&kind,
			&timeWindows,
		); err != nil {
			return nil, err
		}

		i.PTimeWindows, i.DTimeWindows = util.GetShipmentTimeWindows(kind, timeWindows)

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
