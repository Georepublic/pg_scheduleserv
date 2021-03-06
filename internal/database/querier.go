/*GRP-GNU-AGPL******************************************************************

File: querier.go

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

//go:generate mockgen -destination=../mock/mock_querier.go -package=mock -source=querier.go

import (
	"context"

	"github.com/Georepublic/pg_scheduleserv/internal/util"
)

type Querier interface {
	// Break
	DBCreateBreakWithTw(ctx context.Context, arg CreateBreakParams) (Break, error)
	DBListBreaks(ctx context.Context, vehicleID int64) ([]Break, error)
	DBGetBreak(ctx context.Context, id int64) (Break, error)
	DBUpdateBreakWithTw(ctx context.Context, arg UpdateBreakParams, break_id int64) (Break, error)
	DBDeleteBreakWithTw(ctx context.Context, id int64) error

	// Job
	DBCreateJobWithTw(ctx context.Context, arg CreateJobParams) (Job, error)
	DBListJobs(ctx context.Context, projectID int64) ([]Job, error)
	DBGetJob(ctx context.Context, id int64) (Job, error)
	DBUpdateJobWithTw(ctx context.Context, arg UpdateJobParams, job_id int64) (Job, error)
	DBDeleteJobWithTw(ctx context.Context, id int64) error

	// Project
	DBCreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	DBListProjects(ctx context.Context) ([]Project, error)
	DBGetProject(ctx context.Context, id int64) (Project, error)
	DBUpdateProject(ctx context.Context, arg UpdateProjectParams, project_id int64) (Project, error)
	DBDeleteProject(ctx context.Context, id int64) (Project, error)

	// Schedule
	DBCreateSchedule(ctx context.Context, id int64, fresh string) error
	DBGetSchedule(ctx context.Context, id int64) (util.ScheduleData, error)
	DBGetScheduleJob(ctx context.Context, id int64) (util.ScheduleData, error)
	DBGetScheduleShipment(ctx context.Context, id int64) (util.ScheduleData, error)
	DBGetScheduleVehicle(ctx context.Context, id int64) (util.ScheduleData, error)
	DBDeleteSchedule(ctx context.Context, id int64) error

	// Shipment
	DBCreateShipmentWithTw(ctx context.Context, arg CreateShipmentParams) (Shipment, error)
	DBListShipments(ctx context.Context, projectID int64) ([]Shipment, error)
	DBGetShipment(ctx context.Context, id int64) (Shipment, error)
	DBUpdateShipmentWithTw(ctx context.Context, arg UpdateShipmentParams, shipment_id int64) (Shipment, error)
	DBDeleteShipmentWithTw(ctx context.Context, id int64) error

	// Vehicle
	DBCreateVehicle(ctx context.Context, arg CreateVehicleParams) (Vehicle, error)
	DBListVehicles(ctx context.Context, projectID int64) ([]Vehicle, error)
	DBGetVehicle(ctx context.Context, id int64) (Vehicle, error)
	DBUpdateVehicle(ctx context.Context, arg UpdateVehicleParams, vehicle_id int64) (Vehicle, error)
	DBDeleteVehicle(ctx context.Context, id int64) (Vehicle, error)

	// Locations
	DBGetProjectLocations(ctx context.Context, project_id int64) ([]int64, error)
}

var _ Querier = (*Queries)(nil)
