/*GRP-GNU-AGPL******************************************************************

File: schedule.go

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
	"fmt"

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func (q *Queries) DBCreateSchedule(ctx context.Context, projectID int64) error {
	query := fmt.Sprintf("SELECT create_schedule(%d)", projectID)
	_, err := q.db.Exec(ctx, query)
	return err
}

func (q *Queries) DBGetSchedule(ctx context.Context, projectID int64) (util.ScheduleData, error) {
	_, err := q.DBGetProject(ctx, projectID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	filter := " WHERE project_id = $1"
	orderBy := " ORDER BY vehicle_id, arrival, type"
	sql := "SELECT " + util.GetOutputFields(util.ScheduleDB{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleJob(ctx context.Context, jobID int64) (util.ScheduleData, error) {
	_, err := q.DBGetJob(ctx, jobID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	filter := " WHERE task_id = $1 AND type = 'job'"
	orderBy := " ORDER BY vehicle_id, arrival, type"
	sql := "SELECT " + util.GetOutputFields(util.ScheduleDB{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, jobID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleShipment(ctx context.Context, shipmentID int64) (util.ScheduleData, error) {
	_, err := q.DBGetShipment(ctx, shipmentID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	filter := " WHERE task_id = $1  AND (type = 'pickup' OR type = 'delivery')"
	orderBy := " ORDER BY vehicle_id, arrival, type"
	sql := "SELECT " + util.GetOutputFields(util.ScheduleDB{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, shipmentID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleVehicle(ctx context.Context, vehicleID int64) (util.ScheduleData, error) {
	_, err := q.DBGetVehicle(ctx, vehicleID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	filter := " WHERE vehicle_id = $1"
	orderBy := " ORDER BY vehicle_id, arrival, type"
	sql := "SELECT " + util.GetOutputFields(util.ScheduleDB{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, vehicleID)
	if err != nil {
		return util.ScheduleData{}, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

const deleteSchedule = `DELETE FROM schedules WHERE project_id = $1`

func (q *Queries) DBDeleteSchedule(ctx context.Context, projectID int64) error {
	_, err := q.DBGetProject(ctx, projectID)
	if err != nil {
		return err
	}
	_, err = q.db.Exec(ctx, deleteSchedule, projectID)
	return err
}

func scanScheduleRows(rows pgx.Rows) (util.ScheduleData, error) {
	var projectID int64
	schedule := []util.ScheduleResponse{}

	summary := []util.ScheduleSummary{}
	var totalSummary map[string]string = map[string]string{
		"total_travel":  "00:00:00",
		"total_setup":   "00:00:00",
		"total_service": "00:00:00",
		"total_waiting": "00:00:00",
	}
	unassigned := []util.ScheduleUnassigned{}

	var route []util.ScheduleRoute

	var i, prevI util.ScheduleDB
	fullSummaryFound := false

	for rows.Next() {
		var locationID int64
		if err := rows.Scan(
			&i.Type,
			&i.ProjectID,
			&i.VehicleID,
			&i.TaskID,
			&locationID,
			&i.Arrival,
			&i.Departure,
			&i.TravelTime,
			&i.SetupTime,
			&i.ServiceTime,
			&i.WaitingTime,
			&i.Load,
			&i.VehicleData,
			&i.TaskData,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return util.ScheduleData{}, err
		}

		latitude, longitude := util.GetCoordinates(locationID)
		i.Location = util.LocationParams{
			Latitude:  &latitude,
			Longitude: &longitude,
		}

		if i.VehicleID > 0 && i.Type != "summary" {
			// Complete schedule of tasks
			currentRoute := util.ScheduleRoute{
				Type:        i.Type,
				TaskID:      i.TaskID,
				Location:    i.Location,
				Arrival:     i.Arrival,
				Departure:   i.Departure,
				TravelTime:  i.TravelTime,
				SetupTime:   i.SetupTime,
				ServiceTime: i.ServiceTime,
				WaitingTime: i.WaitingTime,
				Load:        i.Load,
				TaskData:    i.TaskData,
				CreatedAt:   i.CreatedAt,
				UpdatedAt:   i.UpdatedAt,
			}
			if i.VehicleID == prevI.VehicleID {
				route = append(route, currentRoute)
			} else {
				if route != nil {
					schedule = append(schedule, util.ScheduleResponse{
						VehicleID:   prevI.VehicleID,
						VehicleData: prevI.VehicleData,
						Route:       route,
					})
					route = nil
				}
				route = append(route, currentRoute)
			}
			prevI = i
		} else if i.VehicleID > 0 {
			// Schedule summary for a vehicle
			summary = append(summary, util.ScheduleSummary{
				VehicleID:   i.VehicleID,
				TravelTime:  i.TravelTime,
				SetupTime:   i.SetupTime,
				ServiceTime: i.ServiceTime,
				WaitingTime: i.WaitingTime,
				VehicleData: i.VehicleData,
			})
		} else if i.VehicleID == 0 {
			fullSummaryFound = true
			// Schedule summary for the complete problem
			totalSummary = map[string]string{
				"total_travel":  i.TravelTime,
				"total_setup":   i.SetupTime,
				"total_service": i.ServiceTime,
				"total_waiting": i.WaitingTime,
			}
		} else if i.VehicleID == -1 {
			// Unassigned tasks
			unassigned = append(unassigned, util.ScheduleUnassigned{
				Type:     i.Type,
				TaskID:   i.TaskID,
				Location: i.Location,
				TaskData: i.TaskData,
			})
		} else {
			logrus.Error("Got Invalid Schedule Response")
		}
		projectID = i.ProjectID
	}
	if err := rows.Err(); err != nil {
		return util.ScheduleData{}, err
	}

	if route != nil {
		schedule = append(schedule, util.ScheduleResponse{
			VehicleID:   prevI.VehicleID,
			VehicleData: prevI.VehicleData,
			Route:       route,
		})
	}

	if !fullSummaryFound && len(summary) != 0 {
		if len(summary) >= 2 {
			logrus.Error("More than 1 vehicles found, but total summary not found")
		}
		totalSummary = map[string]string{
			"total_travel":  summary[0].TravelTime,
			"total_setup":   summary[0].SetupTime,
			"total_service": summary[0].ServiceTime,
			"total_waiting": summary[0].WaitingTime,
		}
	}

	items := util.ScheduleData{
		Schedule: schedule,
		Metadata: util.MetadataResponse{
			Summary:      summary,
			Unassigned:   unassigned,
			TotalTravel:  totalSummary["total_travel"],
			TotalSetup:   totalSummary["total_setup"],
			TotalService: totalSummary["total_service"],
			TotalWaiting: totalSummary["total_waiting"],
		},
		ProjectID: projectID,
	}
	return items, nil
}
