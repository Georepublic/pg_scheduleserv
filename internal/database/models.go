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

	"github.com/Georepublic/pg_scheduleserv/internal/util"
)

type Break struct {
	ID        int64       `json:"id,string"`
	VehicleID int64       `json:"vehicle_id,string"`
	Service   int64       `json:"service"`
	Data      interface{} `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type BreakTimeWindow struct {
	ID        int64  `json:"id,string"`
	TwOpen    string `json:"tw_open"`
	TwClose   string `json:"tw_close"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Job struct {
	ID        int64               `json:"id,string"`
	Location  util.LocationParams `json:"location"`
	Service   int64               `json:"service"`
	Delivery  []int64             `json:"delivery"`
	Pickup    []int64             `json:"pickup"`
	Skills    []int32             `json:"skills"`
	Priority  int32               `json:"priority"`
	ProjectID int64               `json:"project_id,string"`
	Data      interface{}         `json:"data"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

type JobTimeWindow struct {
	ID        int64  `json:"id,string"`
	TwOpen    string `json:"tw_open"`
	TwClose   string `json:"tw_close"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Location struct {
	ID        int64           `json:"id,string"`
	Location  interface{}     `json:"location"`
	Latitude  sql.NullFloat64 `json:"latitude"`
	Longitude sql.NullFloat64 `json:"longitude"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

type Matrix struct {
	StartVid  int64  `json:"start_vid,string"`
	EndVid    int64  `json:"end_vid,string"`
	AggCost   int32  `json:"agg_cost"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Project struct {
	ID        int64       `json:"id,string"`
	Name      string      `json:"name"`
	Data      interface{} `json:"data"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

type ProjectLocation struct {
	ProjectID  int64  `json:"project_id,string"`
	LocationID int64  `json:"location_id,string"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type Shipment struct {
	ID        int64               `json:"id,string"`
	PLocation util.LocationParams `json:"p_location"`
	PService  int64               `json:"p_service"`
	DLocation util.LocationParams `json:"d_location"`
	DService  int64               `json:"d_service"`
	Amount    []int64             `json:"amount"`
	Skills    []int32             `json:"skills"`
	Priority  int32               `json:"priority"`
	ProjectID int64               `json:"project_id,string"`
	Data      interface{}         `json:"data"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

type ShipmentTimeWindow struct {
	ID        int64  `json:"id,string"`
	Kind      string `json:"kind"`
	TwOpen    string `json:"tw_open"`
	TwClose   string `json:"tw_close"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Vehicle struct {
	ID            int64               `json:"id,string"`
	StartLocation util.LocationParams `json:"start_location"`
	EndLocation   util.LocationParams `json:"end_location"`
	Capacity      []int64             `json:"capacity"`
	Skills        []int32             `json:"skills"`
	TwOpen        string              `json:"tw_open"`
	TwClose       string              `json:"tw_close"`
	SpeedFactor   float64             `json:"speed_factor"`
	ProjectID     int64               `json:"project_id,string"`
	Data          interface{}         `json:"data"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}
