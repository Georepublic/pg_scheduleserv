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
	"time"
)

const createShipmentTimeWindow = `-- name: CreateShipmentTimeWindow :one
/*
POST /shipments/{shipment_id}/time_windows/
GET /shipments/{shipment_id}/time_windows/
DELETE /shipments/{shipment_id}/time_windows/{p|d}/YYYYMMDDThhmmssZ/YYYYMMDDThhmmssZ
*/

INSERT INTO shipments_time_windows (id, kind, tw_open, tw_close)
VALUES ($1, $2, $3, $4)
RETURNING id, kind, tw_open, tw_close, created_at, updated_at
`

type CreateShipmentTimeWindowParams struct {
	ID      int64     `json:"id"`
	Kind    string    `json:"kind"`
	TwOpen  time.Time `json:"tw_open"`
	TwClose time.Time `json:"tw_close"`
}

func (q *Queries) CreateShipmentTimeWindow(ctx context.Context, arg CreateShipmentTimeWindowParams) (ShipmentsTimeWindow, error) {
	row := q.db.QueryRow(ctx, createShipmentTimeWindow,
		arg.ID,
		arg.Kind,
		arg.TwOpen,
		arg.TwClose,
	)
	var i ShipmentsTimeWindow
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
	ID      int64     `json:"id"`
	Kind    string    `json:"kind"`
	TwOpen  time.Time `json:"tw_open"`
	TwClose time.Time `json:"tw_close"`
}

func (q *Queries) DeleteShipmentTimeWindow(ctx context.Context, arg DeleteShipmentTimeWindowParams) error {
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

func (q *Queries) ListShipmentTimeWindows(ctx context.Context, id int64) ([]ShipmentsTimeWindow, error) {
	rows, err := q.db.Query(ctx, listShipmentTimeWindows, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ShipmentsTimeWindow{}
	for rows.Next() {
		var i ShipmentsTimeWindow
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
