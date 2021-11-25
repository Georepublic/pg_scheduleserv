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
	"github.com/Georepublic/pg_scheduleserv/internal/util"
)

type Break struct {
	ID        int64       `json:"id,string" example:"1234567812345678"`
	VehicleID int64       `json:"vehicle_id,string" example:"1234567812345678"`
	Service   int64       `json:"service" example:"120"`
	Data      interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt string      `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string      `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type BreakTimeWindow struct {
	ID        int64  `json:"id,string" example:"1234567812345678"`
	TwOpen    string `json:"tw_open" example:"2021-12-31 23:00:00"`
	TwClose   string `json:"tw_close" example:"2021-12-31 23:59:00"`
	CreatedAt string `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type Job struct {
	ID        int64               `json:"id,string" example:"1234567812345678"`
	Location  util.LocationParams `json:"location"`
	Service   int64               `json:"service" example:"120"`
	Delivery  []int64             `json:"delivery" example:"10,20"`
	Pickup    []int64             `json:"pickup" example:"5,15"`
	Skills    []int32             `json:"skills" example:"1,5"`
	Priority  int32               `json:"priority" example:"10"`
	ProjectID int64               `json:"project_id,string" example:"1234567812345678"`
	Data      interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt string              `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string              `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type JobTimeWindow struct {
	ID        int64  `json:"id,string" example:"1234567812345678"`
	TwOpen    string `json:"tw_open" example:"2021-12-31 23:00:00"`
	TwClose   string `json:"tw_close" example:"2021-12-31 23:59:00"`
	CreatedAt string `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type Project struct {
	ID        int64       `json:"id,string" example:"1234567812345678"`
	Name      string      `json:"name"`
	Data      interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt string      `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string      `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type Shipment struct {
	ID        int64               `json:"id,string" example:"1234567812345678"`
	PLocation util.LocationParams `json:"p_location"`
	PService  int64               `json:"p_service"`
	DLocation util.LocationParams `json:"d_location"`
	DService  int64               `json:"d_service"`
	Amount    []int64             `json:"amount"`
	Skills    []int32             `json:"skills" example:"1,5"`
	Priority  int32               `json:"priority" example:"10"`
	ProjectID int64               `json:"project_id,string" example:"1234567812345678"`
	Data      interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt string              `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string              `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type ShipmentTimeWindow struct {
	ID        int64  `json:"id,string" example:"1234567812345678"`
	Kind      string `json:"kind"`
	TwOpen    string `json:"tw_open" example:"2021-12-31 23:00:00"`
	TwClose   string `json:"tw_close" example:"2021-12-31 23:59:00"`
	CreatedAt string `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt string `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type Vehicle struct {
	ID            int64               `json:"id,string" example:"1234567812345678"`
	StartLocation util.LocationParams `json:"start_location"`
	EndLocation   util.LocationParams `json:"end_location"`
	Capacity      []int64             `json:"capacity"`
	Skills        []int32             `json:"skills" example:"1,5"`
	TwOpen        string              `json:"tw_open" example:"2021-12-31 23:00:00"`
	TwClose       string              `json:"tw_close" example:"2021-12-31 23:59:00"`
	SpeedFactor   float64             `json:"speed_factor"`
	ProjectID     int64               `json:"project_id,string" example:"1234567812345678"`
	Data          interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt     string              `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt     string              `json:"updated_at" example:"2021-12-01 13:00:00"`
}
