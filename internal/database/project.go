/*GRP-GNU-AGPL******************************************************************

File: project.go

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

const createProject = `-- name: CreateProject :one
/*
POST /projects
GET /projects/{project_id}
PATCH /projects/{project_id}
DELETE /projects/{project_id}
*/

INSERT INTO projects (name, data) VALUES ($1, $2)
RETURNING id, name, data, created_at, updated_at, deleted
`

type CreateProjectParams struct {
	Name string       `json:"name"`
	Data pgtype.JSONB `json:"data"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, createProject, arg.Name, arg.Data)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}

const deleteProject = `-- name: DeleteProject :exec
UPDATE projects SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) DeleteProject(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteProject, id)
	return err
}

const getProject = `-- name: GetProject :one
SELECT id, name, data, created_at, updated_at
FROM projects
WHERE id = $1 AND deleted = FALSE LIMIT 1
`

type GetProjectRow struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Data      pgtype.JSONB `json:"data"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (q *Queries) GetProject(ctx context.Context, id int64) (GetProjectRow, error) {
	row := q.db.QueryRow(ctx, getProject, id)
	var i GetProjectRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateProject = `-- name: UpdateProject :one
UPDATE projects
SET name = $2, data = $3
WHERE id = $1 AND deleted = FALSE
RETURNING id, name, data, created_at, updated_at, deleted
`

type UpdateProjectParams struct {
	ID   int64        `json:"id"`
	Name string       `json:"name"`
	Data pgtype.JSONB `json:"data"`
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, updateProject, arg.ID, arg.Name, arg.Data)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Deleted,
	)
	return i, err
}
