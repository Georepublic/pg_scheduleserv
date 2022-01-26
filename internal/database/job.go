/*GRP-GNU-AGPL******************************************************************

File: job.go

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

type CreateJobParams struct {
	Location    *util.LocationParams `json:"location" validate:"required"`
	Setup       *string              `json:"setup"    validate:"omitempty" example:"00:00:00"`
	Service     *string              `json:"service"  validate:"omitempty" example:"00:02:00"`
	Delivery    *[]int64             `json:"delivery" validate:"omitempty,dive,min=0" example:"10,20"`
	Pickup      *[]int64             `json:"pickup"   validate:"omitempty,dive,min=0" example:"5,15"`
	Skills      *[]int32             `json:"skills"   validate:"omitempty,dive,min=0" example:"1,5"`
	Priority    *int32               `json:"priority" validate:"omitempty,min=0,max=100" example:"10"`
	TimeWindows *[][]string          `json:"time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
	ProjectID   *int64               `json:"project_id,string" validate:"required" swaggerignore:"true"`
	Data        *interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

type UpdateJobParams struct {
	Location    *util.LocationParams `json:"location"`
	Setup       *string              `json:"setup"    validate:"omitempty" example:"00:00:00"`
	Service     *string              `json:"service"  validate:"omitempty" example:"00:02:00"`
	Delivery    *[]int64             `json:"delivery" validate:"omitempty,dive,min=0" example:"10,20"`
	Pickup      *[]int64             `json:"pickup"   validate:"omitempty,dive,min=0" example:"5,15"`
	Skills      *[]int32             `json:"skills"   validate:"omitempty,dive,min=0" example:"1,5"`
	Priority    *int32               `json:"priority" validate:"omitempty,min=0,max=100" example:"10"`
	TimeWindows *[][]string          `json:"time_windows" validate:"omitempty,dive,min=2,max=2,dive,datetime=2006-01-02T15:04:05"`
	ProjectID   *int64               `json:"project_id,string" swaggerignore:"true"`
	Data        *interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

func (q *Queries) DBCreateJob(ctx context.Context, arg CreateJobParams) (int64, error) {
	tableName := "jobs"
	sql, args := createResource(tableName, arg)
	return_sql := " RETURNING id"
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanID(row)
}

func (q *Queries) DBGetJob(ctx context.Context, id int64) (Job, error) {
	tableName := "jobs"
	joinSelectQuery := fmt.Sprintf(
		", array_agg(ARRAY[%s, %s])",
		util.GetFormattedTimestamp("tw_open"),
		util.GetFormattedTimestamp("tw_close"),
	)
	joinTableQuery := fmt.Sprintf(" LEFT JOIN jobs_time_windows TW on(%s.id = TW.id) ", tableName)
	additionalQuery := fmt.Sprintf(" WHERE %s.id = $1 AND deleted = FALSE GROUP BY %s.id LIMIT 1", tableName, tableName)
	sql := "SELECT " + util.GetOutputFields(Job{}, tableName) + joinSelectQuery + " FROM " + tableName + joinTableQuery + additionalQuery
	row := q.db.QueryRow(ctx, sql, id)
	return scanJobRow(row)
}

func (q *Queries) DBListJobs(ctx context.Context, projectID int64) ([]Job, error) {
	_, err := q.DBGetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	tableName := "jobs"
	joinSelectQuery := fmt.Sprintf(
		", array_agg(ARRAY[%s, %s])",
		util.GetFormattedTimestamp("tw_open"),
		util.GetFormattedTimestamp("tw_close"),
	)
	joinTableQuery := fmt.Sprintf(" LEFT JOIN jobs_time_windows TW on(%s.id = TW.id) ", tableName)
	additionalQuery := fmt.Sprintf(" WHERE project_id = $1 AND deleted = FALSE GROUP BY %s.id ORDER BY %s.created_at", tableName, tableName)
	sql := "SELECT " + util.GetOutputFields(Job{}, tableName) + joinSelectQuery + " FROM " + tableName + joinTableQuery + additionalQuery
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobRows(rows)
}

func (q *Queries) DBUpdateJob(ctx context.Context, arg UpdateJobParams, job_id int64) error {
	tableName := "jobs"
	sql, args := updateResource(tableName, arg, job_id)
	_, err := q.db.Exec(ctx, sql, args...)
	return err
}

func (q *Queries) DBDeleteJob(ctx context.Context, id int64) (Job, error) {
	tableName := "jobs"
	sql := "UPDATE " + tableName + " SET deleted = TRUE WHERE id = $1"
	return_sql := " RETURNING " + util.GetOutputFields(Job{}, tableName)
	row := q.db.QueryRow(ctx, sql+return_sql, id)
	return scanJobRow(row)
}

func (q *Queries) DBCreateJobWithTw(ctx context.Context, arg CreateJobParams) (Job, error) {
	id, err := q.execCreateTx(ctx, func(q *Queries) (int64, error) {
		id, err := q.DBCreateJob(ctx, arg)
		if err != nil {
			return 0, err
		}

		// create time windows from arg and pass to DBCreateJobTimeWindows
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
		err = q.DBCreateJobTimeWindows(ctx, id, timeWindows)
		return id, err
	})
	if err != nil {
		return Job{}, err
	}
	return q.DBGetJob(ctx, id)
}

func (q *Queries) DBUpdateJobWithTw(ctx context.Context, arg UpdateJobParams, job_id int64) (Job, error) {
	err := q.execUpdateTx(ctx, func(q *Queries) error {
		if err := q.DBUpdateJob(ctx, arg, job_id); err != nil {
			return err
		}
		// delete all time windows
		if err := q.DBDeleteJobTimeWindows(ctx, job_id); err != nil {
			return err
		}

		// create time windows from arg and pass to DBCreateJobTimeWindows
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
		err := q.DBCreateJobTimeWindows(ctx, job_id, timeWindows)
		return err
	})
	if err != nil {
		return Job{}, err
	}
	return q.DBGetJob(ctx, job_id)
}

func scanID(row pgx.Row) (int64, error) {
	var id int64
	err := row.Scan(&id)
	return id, err
}

func scanJobRow(row pgx.Row) (Job, error) {
	var i Job
	var location_id int64
	var timeWindows [][]*string
	err := row.Scan(
		&i.ID,
		&location_id,
		&i.Setup,
		&i.Service,
		&i.Delivery,
		&i.Pickup,
		&i.Skills,
		&i.Priority,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&timeWindows,
	)

	i.TimeWindows = util.GetTimeWindows(timeWindows)
	latitude, longitude := util.GetCoordinates(location_id)
	i.Location = util.LocationParams{
		Latitude:  &latitude,
		Longitude: &longitude,
	}
	err = util.HandleDBError(err)
	return i, err
}

func scanJobRows(rows pgx.Rows) ([]Job, error) {
	var i Job
	items := []Job{}
	var location_id int64
	var timeWindows [][]*string
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&location_id,
			&i.Setup,
			&i.Service,
			&i.Delivery,
			&i.Pickup,
			&i.Skills,
			&i.Priority,
			&i.ProjectID,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
			&timeWindows,
		); err != nil {
			return nil, err
		}
		i.TimeWindows = util.GetTimeWindows(timeWindows)
		latitude, longitude := util.GetCoordinates(location_id)
		i.Location = util.LocationParams{
			Latitude:  &latitude,
			Longitude: &longitude,
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
