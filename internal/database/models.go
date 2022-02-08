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
	ID          int64       `json:"id,string" example:"1234567812345678"`
	VehicleID   int64       `json:"vehicle_id,string" example:"1234567812345678"`
	Service     string      `json:"service" example:"00:02:00"`
	Data        interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt   string      `json:"created_at" example:"2021-12-01T13:00:00"`
	UpdatedAt   string      `json:"updated_at" example:"2021-12-01T13:00:00"`
	TimeWindows [][]string  `json:"time_windows"`
}

// add description to all fields and name, so that it can be displayed in the frontend
// only one data field is present in shipment, instead of p_ and d_

// change timezone example (add T)

type Job struct {
	ID          int64               `json:"id,string" example:"1234567812345678"`
	Location    util.LocationParams `json:"location"`
	Setup       string              `json:"setup" example:"00:00:00"`
	Service     string              `json:"service" example:"00:02:00"`
	Delivery    []int64             `json:"delivery" example:"10,20"`
	Pickup      []int64             `json:"pickup" example:"5,15"`
	Skills      []int32             `json:"skills" example:"1,5"`
	Priority    int32               `json:"priority" example:"10"`
	ProjectID   int64               `json:"project_id,string" example:"1234567812345678"`
	Data        interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt   string              `json:"created_at" example:"2021-12-01T13:00:00"`
	UpdatedAt   string              `json:"updated_at" example:"2021-12-01T13:00:00"`
	TimeWindows [][]string          `json:"time_windows"`
}

type Project struct {
	ID           int64       `json:"id,string" example:"1234567812345678"`
	Name         string      `json:"name" example:"Sample Project"`
	DistanceCalc string      `json:"distance_calc" example:"euclidean"`
	Data         interface{} `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt    string      `json:"created_at" example:"2021-12-01T13:00:00"`
	UpdatedAt    string      `json:"updated_at" example:"2021-12-01T13:00:00"`
}

type Shipment struct {
	ID           int64               `json:"id,string" example:"1234567812345678"`
	PLocation    util.LocationParams `json:"p_location" `
	PSetup       string              `json:"p_setup" example:"00:00:00"`
	PService     string              `json:"p_service" example:"00:02:00"`
	DLocation    util.LocationParams `json:"d_location"`
	DSetup       string              `json:"d_setup" example:"00:00:00"`
	DService     string              `json:"d_service" example:"00:02:00"`
	Amount       []int64             `json:"amount" example:"5,15"`
	Skills       []int32             `json:"skills" example:"1,5"`
	Priority     int32               `json:"priority" example:"10"`
	ProjectID    int64               `json:"project_id,string" example:"1234567812345678"`
	Data         interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt    string              `json:"created_at" example:"2021-12-01T13:00:00"`
	UpdatedAt    string              `json:"updated_at" example:"2021-12-01T13:00:00"`
	PTimeWindows [][]string          `json:"p_time_windows"`
	DTimeWindows [][]string          `json:"d_time_windows"`
}

type Vehicle struct {
	ID            int64               `json:"id,string" example:"1234567812345678"`
	StartLocation util.LocationParams `json:"start_location"`
	EndLocation   util.LocationParams `json:"end_location"`
	Capacity      []int64             `json:"capacity" example:"50,25"`
	Skills        []int32             `json:"skills" example:"1,5"`
	TwOpen        string              `json:"tw_open" example:"2021-12-31T23:00:00"`
	TwClose       string              `json:"tw_close" example:"2021-12-31T23:59:00"`
	SpeedFactor   float64             `json:"speed_factor" example:"1.0"`
	MaxTasks      int32               `json:"max_tasks" example:"20"`
	ProjectID     int64               `json:"project_id,string" example:"1234567812345678"`
	Data          interface{}         `json:"data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt     string              `json:"created_at" example:"2021-12-01T13:00:00"`
	UpdatedAt     string              `json:"updated_at" example:"2021-12-01T13:00:00"`
}
