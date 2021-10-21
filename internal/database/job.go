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
	"time"

	"github.com/jackc/pgtype"
)

const createJob = `-- name: CreateJob :one
/*
POST /projects/{project_id}/jobs
GET /projects/{project_id}/jobs

GET /jobs/{job_id}
PATCH /jobs/{job_id}
DELETE /jobs/{job_id}
*/

INSERT INTO jobs (
  location_index, service, delivery, pickup, skills, priority, project_id, data
) VALUES (
  coord_to_id($1, $2), $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, location_index, service, delivery, pickup, skills, priority, project_id, data, created_at, updated_at, deleted
`

type CreateJobParams struct {
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Service   int64        `json:"service"`
	Delivery  []int64      `json:"delivery"`
	Pickup    []int64      `json:"pickup"`
	Skills    []int32      `json:"skills"`
	Priority  int32        `json:"priority"`
	ProjectID int64        `json:"project_id"`
	Data      pgtype.JSONB `json:"data"`
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRow(ctx, createJob,
		arg.Latitude,
		arg.Longitude,
		arg.Service,
		arg.Delivery,
		arg.Pickup,
		arg.Skills,
		arg.Priority,
		arg.ProjectID,
		arg.Data,
	)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.LocationIndex,
		&i.Service,
		&i.Delivery,
		&i.Pickup,
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

const deleteJob = `-- name: DeleteJob :exec
UPDATE jobs SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) DeleteJob(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteJob, id)
	return err
}

const getJob = `-- name: GetJob :one
SELECT
  id, location_index, service, delivery, pickup,
  skills, priority, project_id, data, created_at, updated_at
FROM jobs
WHERE id = $1 AND deleted = FALSE
LIMIT 1
`

type GetJobRow struct {
	ID            int64        `json:"id"`
	LocationIndex int64        `json:"location_index"`
	Service       int64        `json:"service"`
	Delivery      []int64      `json:"delivery"`
	Pickup        []int64      `json:"pickup"`
	Skills        []int32      `json:"skills"`
	Priority      int32        `json:"priority"`
	ProjectID     int64        `json:"project_id"`
	Data          pgtype.JSONB `json:"data"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (q *Queries) GetJob(ctx context.Context, id int64) (GetJobRow, error) {
	row := q.db.QueryRow(ctx, getJob, id)
	var i GetJobRow
	err := row.Scan(
		&i.ID,
		&i.LocationIndex,
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
	return i, err
}

const listJobs = `-- name: ListJobs :many
SELECT
  id, location_index, service, delivery, pickup,
  skills, priority, project_id, data, created_at, updated_at
FROM jobs
WHERE project_id = $1 AND deleted = FALSE
ORDER BY created_at
`

type ListJobsRow struct {
	ID            int64        `json:"id"`
	LocationIndex int64        `json:"location_index"`
	Service       int64        `json:"service"`
	Delivery      []int64      `json:"delivery"`
	Pickup        []int64      `json:"pickup"`
	Skills        []int32      `json:"skills"`
	Priority      int32        `json:"priority"`
	ProjectID     int64        `json:"project_id"`
	Data          pgtype.JSONB `json:"data"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (q *Queries) ListJobs(ctx context.Context, projectID int64) ([]ListJobsRow, error) {
	rows, err := q.db.Query(ctx, listJobs, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListJobsRow{}
	for rows.Next() {
		var i ListJobsRow
		if err := rows.Scan(
			&i.ID,
			&i.LocationIndex,
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
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateJob = `-- name: UpdateJob :one
UPDATE jobs
SET
  location_index = coord_to_id($2, $3), service = $4, delivery = $5,
  pickup = $6, skills = $7, priority = $8, project_id = $9, data = $10
WHERE id = $1 AND deleted = FALSE
RETURNING id, location_index, service, delivery, pickup, skills, priority, project_id, data, created_at, updated_at, deleted
`

type UpdateJobParams struct {
	ID        int64        `json:"id"`
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
	Service   int64        `json:"service"`
	Delivery  []int64      `json:"delivery"`
	Pickup    []int64      `json:"pickup"`
	Skills    []int32      `json:"skills"`
	Priority  int32        `json:"priority"`
	ProjectID int64        `json:"project_id"`
	Data      pgtype.JSONB `json:"data"`
}

func (q *Queries) UpdateJob(ctx context.Context, arg UpdateJobParams) (Job, error) {
	row := q.db.QueryRow(ctx, updateJob,
		arg.ID,
		arg.Latitude,
		arg.Longitude,
		arg.Service,
		arg.Delivery,
		arg.Pickup,
		arg.Skills,
		arg.Priority,
		arg.ProjectID,
		arg.Data,
	)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.LocationIndex,
		&i.Service,
		&i.Delivery,
		&i.Pickup,
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