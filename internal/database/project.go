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

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (name, data) VALUES ($1, $2)
RETURNING id, name, data, created_at, updated_at, deleted
`

type CreateProjectParams struct {
	Name *string      `json:"name" example:"sample_project" validate:"required"`
	Data *interface{} `json:"data" swaggertype:"object"`
}

func (q *Queries) DBCreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	sql, args := createResource("projects", arg)
	return_sql := " RETURNING " + util.GetOutputFields(Project{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanProjectRow(row)
}

const getProject = `-- name: GetProject :one
SELECT id, name, data, created_at, updated_at
FROM projects
WHERE id = $1 AND deleted = FALSE LIMIT 1
`

type GetProjectRow struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Data      interface{} `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (q *Queries) DBGetProject(ctx context.Context, id int64) (Project, error) {
	table_name := "projects"
	additional_query := " WHERE id = $1 AND deleted = FALSE LIMIT 1"
	sql := "SELECT " + util.GetOutputFields(Project{}) + " FROM " + table_name + additional_query
	row := q.db.QueryRow(ctx, sql, id)
	return scanProjectRow(row)
}

const listProjects = `-- name: ListProjects :many
SELECT id, name, data, created_at, updated_at
FROM projects
WHERE deleted = FALSE
ORDER BY created_at
`

type ListProjectsRow struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Data      interface{} `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

func (q *Queries) DBListProjects(ctx context.Context) ([]Project, error) {
	table_name := "projects"
	additional_query := " WHERE deleted = FALSE ORDER BY created_at"
	sql := "SELECT " + util.GetOutputFields(Project{}) + " FROM " + table_name + additional_query
	rows, err := q.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProjectRows(rows)
}

const updateProject = `-- name: UpdateProject :one
UPDATE projects
SET name = $2, data = $3
WHERE id = $1 AND deleted = FALSE
RETURNING id, name, data, created_at, updated_at, deleted
`

type UpdateProjectParams struct {
	ID   int64       `json:"id"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

func (q *Queries) DBUpdateProject(ctx context.Context, arg CreateProjectParams, project_id int64) (Project, error) {
	sql, args := updateResource("projects", arg, project_id)
	return_sql := " RETURNING " + util.GetOutputFields(Project{})
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	return scanProjectRow(row)
}

const deleteProject = `-- name: DeleteProject :exec
UPDATE projects SET deleted = TRUE
WHERE id = $1
`

func (q *Queries) DBDeleteProject(ctx context.Context, id int64) error {
	sql := "UPDATE projects SET deleted = TRUE WHERE id = $1"
	_, err := q.db.Exec(ctx, sql, id)
	return err
}

func scanProjectRow(row pgx.Row) (Project, error) {
	var i Project
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Data,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

func scanProjectRows(rows pgx.Rows) ([]Project, error) {
	items := []Project{}
	var i Project
	for rows.Next() {
		if err := rows.Scan(
			&i.ID,
			&i.Name,
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
