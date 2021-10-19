/*GRP-GNU-AGPL******************************************************************

File: models.go

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
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
)

type Break struct {
	ID        int64        `json:"id"`
	VehicleID int64        `json:"vehicle_id"`
	Service   int64        `json:"service"`
	Data      pgtype.JSONB `json:"data"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Deleted   bool         `json:"deleted"`
}

type BreaksTimeWindow struct {
	ID        int64     `json:"id"`
	TwOpen    time.Time `json:"tw_open"`
	TwClose   time.Time `json:"tw_close"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Job struct {
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
	Deleted       bool         `json:"deleted"`
}

type JobsTimeWindow struct {
	ID        int64     `json:"id"`
	TwOpen    time.Time `json:"tw_open"`
	TwClose   time.Time `json:"tw_close"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Location struct {
	ID        int64           `json:"id"`
	Location  interface{}     `json:"location"`
	Latitude  sql.NullFloat64 `json:"latitude"`
	Longitude sql.NullFloat64 `json:"longitude"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type Matrix struct {
	StartVid  int64     `json:"start_vid"`
	EndVid    int64     `json:"end_vid"`
	AggCost   int32     `json:"agg_cost"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	Data      pgtype.JSONB `json:"data"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Deleted   bool         `json:"deleted"`
}

type ProjectLocation struct {
	ProjectID  int64     `json:"project_id"`
	LocationID int64     `json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Shipment struct {
	ID             int64        `json:"id"`
	PLocationIndex int64        `json:"p_location_index"`
	PService       int64        `json:"p_service"`
	DLocationIndex int64        `json:"d_location_index"`
	DService       int64        `json:"d_service"`
	Amount         []int64      `json:"amount"`
	Skills         []int32      `json:"skills"`
	Priority       int32        `json:"priority"`
	ProjectID      int64        `json:"project_id"`
	Data           pgtype.JSONB `json:"data"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Deleted        bool         `json:"deleted"`
}

type ShipmentsTimeWindow struct {
	ID        int64     `json:"id"`
	Kind      string    `json:"kind"`
	TwOpen    time.Time `json:"tw_open"`
	TwClose   time.Time `json:"tw_close"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Vehicle struct {
	ID          int64        `json:"id"`
	StartIndex  int64        `json:"start_index"`
	EndIndex    int64        `json:"end_index"`
	Capacity    []int64      `json:"capacity"`
	Skills      []int32      `json:"skills"`
	TwOpen      time.Time    `json:"tw_open"`
	TwClose     time.Time    `json:"tw_close"`
	SpeedFactor float64      `json:"speed_factor"`
	ProjectID   int64        `json:"project_id"`
	Data        pgtype.JSONB `json:"data"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Deleted     bool         `json:"deleted"`
}
