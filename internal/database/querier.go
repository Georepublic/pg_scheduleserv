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
	DBCreateBreak(ctx context.Context, arg CreateBreakParams) (Break, error)
	DBListBreaks(ctx context.Context, vehicleID int64) ([]Break, error)
	DBGetBreak(ctx context.Context, id int64) (Break, error)
	DBUpdateBreak(ctx context.Context, arg UpdateBreakParams, break_id int64) (Break, error)
	DBDeleteBreak(ctx context.Context, id int64) (Break, error)

	// Break Time Window
	DBCreateBreakTimeWindow(ctx context.Context, arg CreateBreakTimeWindowParams) (BreakTimeWindow, error)
	DBListBreakTimeWindows(ctx context.Context, id int64) ([]BreakTimeWindow, error)
	DBDeleteBreakTimeWindow(ctx context.Context, id int64) (BreakTimeWindow, error)

	// Job
	DBCreateJob(ctx context.Context, arg CreateJobParams) (Job, error)
	DBListJobs(ctx context.Context, projectID int64) ([]Job, error)
	DBGetJob(ctx context.Context, id int64) (Job, error)
	DBUpdateJob(ctx context.Context, arg UpdateJobParams, job_id int64) (Job, error)
	DBDeleteJob(ctx context.Context, id int64) (Job, error)

	// Job Time Window
	DBCreateJobTimeWindow(ctx context.Context, arg CreateJobTimeWindowParams) (JobTimeWindow, error)
	DBListJobTimeWindows(ctx context.Context, id int64) ([]JobTimeWindow, error)
	DBDeleteJobTimeWindow(ctx context.Context, id int64) (JobTimeWindow, error)

	// Project
	DBCreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	DBListProjects(ctx context.Context) ([]Project, error)
	DBGetProject(ctx context.Context, id int64) (Project, error)
	DBUpdateProject(ctx context.Context, arg UpdateProjectParams, project_id int64) (Project, error)
	DBDeleteProject(ctx context.Context, id int64) (Project, error)

	// Schedule
	DBCreateSchedule(ctx context.Context, id int64) error
	DBGetSchedule(ctx context.Context, id int64) ([]util.Schedule, error)
	DBGetScheduleVehicle(ctx context.Context, id int64) ([]util.Schedule, error)
	DBDeleteSchedule(ctx context.Context, id int64) error

	// Shipment
	DBCreateShipment(ctx context.Context, arg CreateShipmentParams) (Shipment, error)
	DBListShipments(ctx context.Context, projectID int64) ([]Shipment, error)
	DBGetShipment(ctx context.Context, id int64) (Shipment, error)
	DBUpdateShipment(ctx context.Context, arg UpdateShipmentParams, shipment_id int64) (Shipment, error)
	DBDeleteShipment(ctx context.Context, id int64) (Shipment, error)

	// Shipment Time Window
	DBCreateShipmentTimeWindow(ctx context.Context, arg CreateShipmentTimeWindowParams) (ShipmentTimeWindow, error)
	DBListShipmentTimeWindows(ctx context.Context, id int64) ([]ShipmentTimeWindow, error)
	DBDeleteShipmentTimeWindow(ctx context.Context, id int64) (ShipmentTimeWindow, error)

	// Vehicle
	DBCreateVehicle(ctx context.Context, arg CreateVehicleParams) (Vehicle, error)
	DBListVehicles(ctx context.Context, projectID int64) ([]Vehicle, error)
	DBGetVehicle(ctx context.Context, id int64) (Vehicle, error)
	DBUpdateVehicle(ctx context.Context, arg UpdateVehicleParams, vehicle_id int64) (Vehicle, error)
	DBDeleteVehicle(ctx context.Context, id int64) (Vehicle, error)
}

var _ Querier = (*Queries)(nil)
