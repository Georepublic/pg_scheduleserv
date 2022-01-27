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
	"fmt"

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
)

type CreateBreakParams struct {
	VehicleID   *int64       `json:"vehicle_id,string" example:"1234567812345678" validate:"required" swaggerignore:"true"`
	Service     *string      `json:"service" validate:"omitempty" example:"00:02:00"`
	Data        *interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	TimeWindows *[][]string  `json:"time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
}

type UpdateBreakParams struct {
	VehicleID   *int64       `json:"vehicle_id,string" example:"1234567812345678" swaggerignore:"true"`
	Service     *string      `json:"service" validate:"omitempty" example:"00:02:00"`
	Data        *interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	TimeWindows *[][]string  `json:"time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
}

func (q *Queries) DBCreateBreak(ctx context.Context, arg CreateBreakParams) (int64, error) {
	tableName := "breaks"
	sql, args := createResource(tableName, arg)
	return_sql := " RETURNING id"
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanID(row)
}

func (q *Queries) DBGetBreak(ctx context.Context, id int64) (Break, error) {
	tableName := "breaks"
	joinSelectQuery := fmt.Sprintf(
		", array_agg(ARRAY[%s, %s])",
		util.GetFormattedTimestamp("tw_open"),
		util.GetFormattedTimestamp("tw_close"),
	)
	joinTableQuery := fmt.Sprintf(" LEFT JOIN breaks_time_windows TW on(%s.id = TW.id) ", tableName)
	additionalQuery := fmt.Sprintf(" WHERE %s.id = $1 AND deleted = FALSE GROUP BY %s.id LIMIT 1", tableName, tableName)
	sql := "SELECT " + util.GetOutputFields(Break{}, tableName) + joinSelectQuery + " FROM " + tableName + joinTableQuery + additionalQuery
	row := q.db.QueryRow(ctx, sql, id)
	return scanBreakRow(row)
}

func (q *Queries) DBListBreaks(ctx context.Context, vehicleID int64) ([]Break, error) {
	_, err := q.DBGetProject(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	tableName := "breaks"
	joinSelectQuery := fmt.Sprintf(
		", array_agg(ARRAY[%s, %s])",
		util.GetFormattedTimestamp("tw_open"),
		util.GetFormattedTimestamp("tw_close"),
	)
	joinTableQuery := fmt.Sprintf(" LEFT JOIN breaks_time_windows TW on(%s.id = TW.id) ", tableName)
	additionalQuery := fmt.Sprintf(" WHERE vehicle_id = $1 AND deleted = FALSE GROUP BY %s.id ORDER BY %s.created_at", tableName, tableName)
	sql := "SELECT " + util.GetOutputFields(Break{}, tableName) + joinSelectQuery + " FROM " + tableName + joinTableQuery + additionalQuery
	rows, err := q.db.Query(ctx, sql, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBreakRows(rows)
}

func (q *Queries) DBUpdateBreak(ctx context.Context, arg UpdateBreakParams, break_id int64) error {
	tableName := "breaks"
	sql, args := updateResource(tableName, arg, break_id)
	_, err := q.db.Exec(ctx, sql, args...)
	return err
}

func (q *Queries) DBDeleteBreak(ctx context.Context, id int64) error {
	tableName := "breaks"
	sql := "UPDATE " + tableName + " SET deleted = TRUE WHERE id = $1"
	_, err := q.db.Exec(ctx, sql, id)
	return err
}

func (q *Queries) DBCreateBreakWithTw(ctx context.Context, arg CreateBreakParams) (Break, error) {
	id, err := q.execCreateTx(ctx, func(q *Queries) (int64, error) {
		id, err := q.DBCreateBreak(ctx, arg)
		if err != nil {
			return 0, err
		}

		// create time windows from arg and pass to DBCreateBreakTimeWindows
		timeWindows := []TimeWindowParams{}
		if arg.TimeWindows != nil {
			for _, tw := range *arg.TimeWindows {
				// append time window to timeWindows
				timeWindows = append(timeWindows, TimeWindowParams{
					TwOpen:  tw[0],
					TwClose: tw[1],
				})
			}
		}
		err = q.DBCreateBreakTimeWindows(ctx, id, timeWindows)
		return id, err
	})
	if err != nil {
		return Break{}, err
	}
	return q.DBGetBreak(ctx, id)
}

func (q *Queries) DBUpdateBreakWithTw(ctx context.Context, arg UpdateBreakParams, break_id int64) (Break, error) {
	err := q.execUpdateTx(ctx, func(q *Queries) error {
		if err := q.DBUpdateBreak(ctx, arg, break_id); err != nil {
			return err
		}
		// delete all time windows
		if err := q.DBDeleteBreakTimeWindows(ctx, break_id); err != nil {
			return err
		}

		// create time windows from arg and pass to DBCreateBreakTimeWindows
		timeWindows := []TimeWindowParams{}
		if arg.TimeWindows != nil {
			for _, tw := range *arg.TimeWindows {
				// append time window to timeWindows
				timeWindows = append(timeWindows, TimeWindowParams{
					TwOpen:  tw[0],
					TwClose: tw[1],
				})
			}
		}
		err := q.DBCreateBreakTimeWindows(ctx, break_id, timeWindows)
		return err
	})
	if err != nil {
		return Break{}, err
	}
	return q.DBGetBreak(ctx, break_id)
}

func (q *Queries) DBDeleteBreakWithTw(ctx context.Context, break_id int64) error {
	err := q.execUpdateTx(ctx, func(q *Queries) error {
		// delete all time windows
		if err := q.DBDeleteBreakTimeWindows(ctx, break_id); err != nil {
			return err
		}
		err := q.DBDeleteBreak(ctx, break_id)
		return err
	})
	return err
}

func scanBreakRow(row pgx.Row) (Break, error) {
	var i Break
	var timeWindows [][]*string
	err := row.Scan(
		&i.ID,
		&i.VehicleID,
		&i.Service,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&timeWindows,
	)
	i.TimeWindows = util.GetTimeWindows(timeWindows)
	err = util.HandleDBError(err)
	return i, err
}

func scanBreakRows(rows pgx.Rows) ([]Break, error) {
	items := []Break{}
	var i Break
	var timeWindows [][]*string
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&i.VehicleID,
			&i.Service,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&timeWindows,
		); err != nil {
			return nil, err
		}
		i.TimeWindows = util.GetTimeWindows(timeWindows)
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
