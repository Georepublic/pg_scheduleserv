/*GRP-GNU-AGPL******************************************************************

File: shipment_tw.go

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
	"github.com/sirupsen/logrus"
)

const createShipmentTimeWindow = `-- name: CreateShipmentTimeWindow :one
INSERT INTO shipments_time_windows (id, kind, tw_open, tw_close)
VALUES ($1, $2, $3, $4)
RETURNING id, kind, tw_open, tw_close, created_at, updated_at
`

type CreateShipmentTimeWindowParams struct {
	ID      *int64  `json:"id,string" example:"1234567890123456789" validate:"required" swaggerignore:"true"`
	Kind    *string `json:"kind" validate:"required"`
	TwOpen  *string `json:"tw_open" validate:"required,datetime=2006-01-02 15:04:05"`
	TwClose *string `json:"tw_close" validate:"required,datetime=2006-01-02 15:04:05"`
}

func (q *Queries) DBCreateShipmentTimeWindow(ctx context.Context, arg CreateShipmentTimeWindowParams) (ShipmentTimeWindow, error) {
	sql, args := createResource("shipments_time_windows", arg)
	logrus.Debug(sql)
	logrus.Debug(args)
	var i ShipmentTimeWindow
	return_sql := util.GetReturnSql(i)
	row := q.db.QueryRow(ctx, sql+return_sql, args...)
	err := row.Scan(
		&i.ID,
		&i.Kind,
		&i.TwOpen,
		&i.TwClose,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteShipmentTimeWindow = `-- name: DeleteShipmentTimeWindow :exec
DELETE FROM shipments_time_windows
WHERE id = $1 AND kind = $2 AND tw_open = $3 AND tw_close = $4
`

type DeleteShipmentTimeWindowParams struct {
	ID      int64  `json:"id"`
	Kind    string `json:"kind"`
	TwOpen  string `json:"tw_open"`
	TwClose string `json:"tw_close"`
}

func (q *Queries) DBDeleteShipmentTimeWindow(ctx context.Context, arg DeleteShipmentTimeWindowParams) error {
	_, err := q.db.Exec(ctx, deleteShipmentTimeWindow,
		arg.ID,
		arg.Kind,
		arg.TwOpen,
		arg.TwClose,
	)
	return err
}

const listShipmentTimeWindows = `-- name: ListShipmentTimeWindows :many
SELECT id, kind, tw_open, tw_close, created_at, updated_at
FROM shipments_time_windows
WHERE id = $1
ORDER BY created_at
`

func (q *Queries) DBListShipmentTimeWindows(ctx context.Context, id int64) ([]ShipmentTimeWindow, error) {
	rows, err := q.db.Query(ctx, listShipmentTimeWindows, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ShipmentTimeWindow{}
	for rows.Next() {
		var i ShipmentTimeWindow
		if err := rows.Scan(
			&i.ID,
			&i.Kind,
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
