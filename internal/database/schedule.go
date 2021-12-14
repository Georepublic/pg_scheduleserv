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

const getSchedules = `
SELECT
  A.id,
  CASE
    WHEN A.type = 1 THEN 'start'
    WHEN A.type = 2 THEN 'job'
    WHEN A.type = 3 THEN 'pickup'
    WHEN A.type = 4 THEN 'delivery'
    WHEN A.type = 5 THEN 'break'
    WHEN A.type = 6 THEN 'end'
  END AS type,
  A.project_id, A.vehicle_id, A.job_id, A.shipment_id, A.break_id,
  CASE
    WHEN A.type = 1 THEN V.start_index
    WHEN A.type = 2 THEN J.location_index
    WHEN A.type = 3 THEN S.p_location_index
    WHEN A.type = 4 THEN S.d_location_index
    WHEN A.type = 6 THEN V.end_index
  END AS location_id,
  to_char(A.arrival, 'YYYY-MM-DD HH24:MI:SS') AS arrival,
  to_char(lead(A.arrival, 1, A.arrival) OVER(PARTITION BY A.vehicle_id ORDER BY A.arrival, A.type), 'YYYY-MM-DD HH24:MI:SS') AS departure,
  EXTRACT(epoch FROM A.travel_time),
  EXTRACT(epoch FROM A.service_time),
  EXTRACT(epoch FROM A.waiting_time),
  lag(A.load, 1, A.load) OVER(PARTITION BY A.vehicle_id ORDER BY A.arrival, A.type) AS start_load,
  A.load AS end_load,
  to_char(A.created_at, 'YYYY-MM-DD HH24:MI:SS'),
  to_char(A.updated_at, 'YYYY-MM-DD HH24:MI:SS')
FROM schedules A
LEFT JOIN jobs J ON job_id = J.id
LEFT JOIN shipments S ON shipment_id = S.id
LEFT JOIN vehicles V ON A.vehicle_id = V.id`

func (q *Queries) DBGetSchedule(ctx context.Context, projectID int64) ([]util.Schedule, error) {
	filter := " WHERE A.project_id = $1"
	orderBy := " ORDER BY vehicle_id, arrival, A.type"
	rows, err := q.db.Query(ctx, getSchedules+filter+orderBy, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleShipment(ctx context.Context, shipmentID int64) ([]util.Schedule, error) {
	filter := " WHERE A.shipment_id = $1"
	orderBy := " ORDER BY arrival, A.type"
	rows, err := q.db.Query(ctx, getSchedules+filter+orderBy, shipmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

func (q *Queries) DBGetScheduleVehicle(ctx context.Context, vehicleID int64) ([]util.Schedule, error) {
	filter := " WHERE A.vehicle_id = $1"
	orderBy := " ORDER BY arrival, A.type"
	rows, err := q.db.Query(ctx, getSchedules+filter+orderBy, vehicleID)
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
		var jobID *int64
		var shipmentID *int64
		var breakID *int64
		var locationID *int64
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.ProjectID,
			&i.VehicleID,
			&jobID,
			&shipmentID,
			&breakID,
			&locationID,
			&i.Arrival,
			&i.Departure,
			&i.TravelTime,
			&i.ServiceTime,
			&i.WaitingTime,
			&i.StartLoad,
			&i.EndLoad,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if jobID == nil {
			jobID = new(int64)
		}
		if shipmentID == nil {
			shipmentID = new(int64)
		}
		if breakID == nil {
			breakID = new(int64)
		}
		if locationID == nil {
			// The last location ID
			locationID = &i.LocationID
		}
		i.JobID = *jobID
		i.ShipmentID = *shipmentID
		i.BreakID = *breakID
		i.LocationID = *locationID

		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
