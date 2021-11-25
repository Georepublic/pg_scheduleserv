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

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
)

type CreateJobParams struct {
	Location  *util.LocationParams `json:"location" validate:"required"`
	Service   *int64               `json:"service"`
	Delivery  *[]int64             `json:"delivery"`
	Pickup    *[]int64             `json:"pickup"`
	Skills    *[]int32             `json:"skills"`
	Priority  *int32               `json:"priority"`
	ProjectID *int64               `json:"project_id,string" validate:"required" swaggerignore:"true"`
	Data      *interface{}         `json:"data" swaggertype:"object"`
}

func (q *Queries) DBCreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	sql, args := createResource("jobs", arg)
	return_sql := " RETURNING " + util.GetOutputFields(Job{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanJobRow(row)
}

type GetJobRow struct {
	ID            int64       `json:"id"`
	LocationIndex int64       `json:"location_index"`
	Service       int64       `json:"service"`
	Delivery      []int64     `json:"delivery"`
	Pickup        []int64     `json:"pickup"`
	Skills        []int32     `json:"skills"`
	Priority      int32       `json:"priority"`
	ProjectID     int64       `json:"project_id"`
	Data          interface{} `json:"data"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

func (q *Queries) DBGetJob(ctx context.Context, id int64) (Job, error) {
	table_name := "jobs"
	additional_query := " WHERE id = $1 AND deleted = FALSE LIMIT 1"
	sql := "SELECT " + util.GetOutputFields(Job{}) + " FROM " + table_name + additional_query
	row := q.db.QueryRow(ctx, sql, id)
	return scanJobRow(row)
}

type ListJobsRow struct {
	ID            int64       `json:"id"`
	LocationIndex int64       `json:"location_index"`
	Service       int64       `json:"service"`
	Delivery      []int64     `json:"delivery"`
	Pickup        []int64     `json:"pickup"`
	Skills        []int32     `json:"skills"`
	Priority      int32       `json:"priority"`
	ProjectID     int64       `json:"project_id"`
	Data          interface{} `json:"data"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

func (q *Queries) DBListJobs(ctx context.Context, projectID int64) ([]Job, error) {
	table_name := "jobs"
	additional_query := " WHERE project_id = $1 AND deleted = FALSE ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(Job{}) + " FROM " + table_name + additional_query
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanJobRows(rows)
}

type UpdateJobParams struct {
	Location  *util.LocationParams `json:"location"`
	Service   *int64               `json:"service"`
	Delivery  *[]int64             `json:"delivery"`
	Pickup    *[]int64             `json:"pickup"`
	Skills    *[]int32             `json:"skills"`
	Priority  *int32               `json:"priority"`
	ProjectID *int64               `json:"project_id,string"`
	Data      *interface{}         `json:"data" swaggertype:"object"`
}

func (q *Queries) DBUpdateJob(ctx context.Context, arg UpdateJobParams, job_id int64) (Job, error) {
	sql, args := updateResource("jobs", arg, job_id)
	return_sql := " RETURNING " + util.GetOutputFields(Job{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanJobRow(row)
}

func (q *Queries) DBDeleteJob(ctx context.Context, id int64) (Job, error) {
	sql := "UPDATE jobs SET deleted = TRUE WHERE id = $1"
	return_sql := " RETURNING " + util.GetOutputFields(Job{})
	row := q.db.QueryRow(ctx, sql+return_sql, id)
	return scanJobRow(row)
}

func scanJobRow(row pgx.Row) (Job, error) {
	var i Job
	var location_index int64
	err := row.Scan(
		&i.ID,
		&location_index,
		&i.Service,
		&i.Delivery,
		&i.Pickup,
		&i.Skills,
		&i.Priority,
		&i.ProjectID,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	latitude, longitude := util.GetCoordinates(location_index)
	i.Location = util.LocationParams{
		Latitude:  &latitude,
		Longitude: &longitude,
	}
	return i, err
}

func scanJobRows(rows pgx.Rows) ([]Job, error) {
	var i Job
	items := []Job{}
	var location_index int64
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&location_index,
			&i.Service,
			&i.Delivery,
			&i.Pickup,
			&i.Skills,
			&i.Priority,
			&i.ProjectID,
			&i.Data,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		latitude, longitude := util.GetCoordinates(location_index)
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
