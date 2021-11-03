package database

//go:generate mockgen -destination=../mock/mock_querier.go -package=mock -source=querier.go

import (
	"context"
)

type Querier interface {
	// Break
	DBCreateBreak(ctx context.Context, arg CreateBreakParams) (Break, error)
	DBGetBreak(ctx context.Context, id int64) (Break, error)
	DBListBreaks(ctx context.Context, vehicleID int64) ([]Break, error)
	DBUpdateBreak(ctx context.Context, arg UpdateBreakParams, break_id int64) (Break, error)
	DBDeleteBreak(ctx context.Context, id int64) error

	// Break Time Window
	DBCreateBreakTimeWindow(ctx context.Context, arg CreateBreakTimeWindowParams) (BreakTimeWindow, error)
	DBListBreakTimeWindows(ctx context.Context, id int64) ([]BreakTimeWindow, error)
	DBDeleteBreakTimeWindow(ctx context.Context, arg CreateBreakTimeWindowParams) error

	// Job
	DBCreateJob(ctx context.Context, arg CreateJobParams) (Job, error)
	DBDeleteJob(ctx context.Context, id int64) error
	DBGetJob(ctx context.Context, id int64) (Job, error)
	DBUpdateJob(ctx context.Context, arg UpdateJobParams, job_id int64) (Job, error)
	DBListJobs(ctx context.Context, projectID int64) ([]Job, error)

	// Job Time Window
	DBCreateJobTimeWindow(ctx context.Context, arg CreateJobTimeWindowParams) (JobTimeWindow, error)
	DBListJobTimeWindows(ctx context.Context, id int64) ([]JobTimeWindow, error)
	DBDeleteJobTimeWindow(ctx context.Context, arg CreateJobTimeWindowParams) error

	// Project
	DBCreateProject(ctx context.Context, arg CreateProjectParams) (Project, error)
	DBGetProject(ctx context.Context, id int64) (Project, error)
	DBListProjects(ctx context.Context) ([]Project, error)
	DBUpdateProject(ctx context.Context, arg UpdateProjectParams, project_id int64) (Project, error)
	DBDeleteProject(ctx context.Context, id int64) (Project, error)

	// Shipment
	DBCreateShipment(ctx context.Context, arg CreateShipmentParams) (Shipment, error)
	DBGetShipment(ctx context.Context, id int64) (Shipment, error)
	DBListShipments(ctx context.Context, projectID int64) ([]Shipment, error)
	DBUpdateShipment(ctx context.Context, arg UpdateShipmentParams, shipment_id int64) (Shipment, error)
	DBDeleteShipment(ctx context.Context, id int64) error

	// Shipment Time Window
	DBCreateShipmentTimeWindow(ctx context.Context, arg CreateShipmentTimeWindowParams) (ShipmentTimeWindow, error)
	DBDeleteShipmentTimeWindow(ctx context.Context, arg CreateShipmentTimeWindowParams) error
	DBListShipmentTimeWindows(ctx context.Context, id int64) ([]ShipmentTimeWindow, error)

	// Vehicle
	DBCreateVehicle(ctx context.Context, arg CreateVehicleParams) (Vehicle, error)
	DBGetVehicle(ctx context.Context, id int64) (Vehicle, error)
	DBListVehicles(ctx context.Context, projectID int64) ([]Vehicle, error)
	DBUpdateVehicle(ctx context.Context, arg UpdateVehicleParams, vehicle_id int64) (Vehicle, error)
	DBDeleteVehicle(ctx context.Context, id int64) error
}

var _ Querier = (*Queries)(nil)
