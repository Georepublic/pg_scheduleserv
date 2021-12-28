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
)

func (q *Queries) DBCreateSchedule(ctx context.Context, projectID int64) error {
	query := fmt.Sprintf("SELECT create_schedule(%d)", projectID)
	_, err := q.db.Exec(ctx, query)
	return err
}

func (q *Queries) DBGetSchedule(ctx context.Context, projectID int64) ([]util.Schedule, error) {
	filter := " WHERE project_id = $1"
	orderBy := " ORDER BY arrival, type, vehicle_id"
	sql := "SELECT " + util.GetOutputFields(util.Schedule{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleJob(ctx context.Context, shipmentID int64) ([]util.Schedule, error) {
	filter := " WHERE task_id = $1 AND type = 'job'"
	orderBy := " ORDER BY arrival, type, vehicle_id"
	sql := "SELECT " + util.GetOutputFields(util.Schedule{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleShipment(ctx context.Context, shipmentID int64) ([]util.Schedule, error) {
	filter := " WHERE task_id = $1  AND (type = 'pickup' OR type = 'delivery')"
	orderBy := " ORDER BY arrival, type, vehicle_id"
	sql := "SELECT " + util.GetOutputFields(util.Schedule{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleVehicle(ctx context.Context, vehicleID int64) ([]util.Schedule, error) {
	filter := " WHERE vehicle_id = $1"
	orderBy := " ORDER BY arrival, type, vehicle_id"
	sql := "SELECT " + util.GetOutputFields(util.Schedule{}) + " FROM schedules" + filter + orderBy
	rows, err := q.db.Query(ctx, sql, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

const deleteSchedule = `DELETE FROM schedules WHERE project_id = $1`

func (q *Queries) DBDeleteSchedule(ctx context.Context, projectID int64) error {
	_, err := q.db.Exec(ctx, deleteSchedule, projectID)
	return err
}

func scanScheduleRows(rows pgx.Rows) ([]util.Schedule, error) {
	var i util.Schedule
	items := []util.Schedule{}
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
			return nil, err
		}

		latitude, longitude := util.GetCoordinates(locationID)
		i.Location = util.LocationParams{
			Latitude:  &latitude,
			Longitude: &longitude,
		}

		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
