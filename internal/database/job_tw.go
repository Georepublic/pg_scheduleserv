/*GRP-GNU-AGPL******************************************************************

File: job_tw.go

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
)

const createJobTimeWindow = `-- name: CreateJobTimeWindow :one
/*
POST /jobs/{job_id}/time_windows/
GET /jobs/{job_id}/time_windows/
DELETE /jobs/{job_id}/time_windows/YYYYMMDDThhmmssZ/YYYYMMDDThhmmssZ
*/

INSERT INTO jobs_time_windows (id, tw_open, tw_close)
VALUES ($1, $2, $3)
RETURNING id, tw_open, tw_close, created_at, updated_at
`

type CreateJobTimeWindowParams struct {
	ID      int64     `json:"id"`
	TwOpen  time.Time `json:"tw_open"`
	TwClose time.Time `json:"tw_close"`
}

func (q *Queries) CreateJobTimeWindow(ctx context.Context, arg CreateJobTimeWindowParams) (JobsTimeWindow, error) {
	row := q.db.QueryRow(ctx, createJobTimeWindow, arg.ID, arg.TwOpen, arg.TwClose)
	var i JobsTimeWindow
	err := row.Scan(
		&i.ID,
		&i.TwOpen,
		&i.TwClose,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteJobTimeWindow = `-- name: DeleteJobTimeWindow :exec
DELETE FROM jobs_time_windows
WHERE id = $1 AND tw_open = $2 AND tw_close = $3
`

type DeleteJobTimeWindowParams struct {
	ID      int64     `json:"id"`
	TwOpen  time.Time `json:"tw_open"`
	TwClose time.Time `json:"tw_close"`
}

func (q *Queries) DeleteJobTimeWindow(ctx context.Context, arg DeleteJobTimeWindowParams) error {
	_, err := q.db.Exec(ctx, deleteJobTimeWindow, arg.ID, arg.TwOpen, arg.TwClose)
	return err
}

const listJobTimeWindows = `-- name: ListJobTimeWindows :many
SELECT id, tw_open, tw_close, created_at, updated_at
FROM jobs_time_windows
WHERE id = $1
ORDER BY created_at
`

func (q *Queries) ListJobTimeWindows(ctx context.Context, id int64) ([]JobsTimeWindow, error) {
	rows, err := q.db.Query(ctx, listJobTimeWindows, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []JobsTimeWindow{}
	for rows.Next() {
		var i JobsTimeWindow
		if err := rows.Scan(
			&i.ID,
			&i.TwOpen,
			&i.TwClose,
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
