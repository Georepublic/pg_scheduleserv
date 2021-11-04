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

const createSchedule = `
INSERT INTO schedules
  (type, project_id, vehicle_id, job_id, shipment_id, break_id,
   arrival, travel_time, service_time, waiting_time, load)
SELECT
  step_type, $1::BIGINT, vehicle_id,
  CASE WHEN step_type = 2 THEN task_id ELSE NULL END,
  CASE WHEN step_type = 3 OR step_type = 4 THEN task_id ELSE NULL END,
  CASE WHEN step_type = 5 THEN task_id ELSE NULL END,
  arrival, travel_time, service_time, waiting_time, load
FROM vrp_vroom(
  'SELECT * FROM jobs WHERE project_id = ' || $1,
  'SELECT * FROM jobs_time_windows',
  'SELECT * FROM shipments WHERE project_id = ' || $1,
  'SELECT * FROM shipments_time_windows',
  'SELECT * FROM vehicles WHERE project_id = ' || $1,
  'SELECT * FROM breaks',
  'SELECT * FROM breaks_time_windows',
  'SELECT * FROM matrix'
);
`

func (q *Queries) DBCreateSchedule(ctx context.Context, projectID int64) error {
	_, err := q.db.Exec(ctx, createSchedule, fmt.Sprintf("%d", projectID))
	return err
}

const getSchedules = `
SELECT
  id, type, project_id, vehicle_id, job_id, shipment_id, break_id,
  to_char(arrival, 'YYYY-MM-DD HH24:MI:SS') AS arrival,
  to_char(lead(arrival, 1, arrival) OVER(PARTITION BY vehicle_id ORDER BY arrival), 'YYYY-MM-DD HH24:MI:SS') AS departure,
  EXTRACT(epoch FROM travel_time),
  EXTRACT(epoch FROM service_time),
  EXTRACT(epoch FROM waiting_time),
  lag(load, 1, load) OVER(PARTITION BY vehicle_id ORDER BY arrival) AS start_load,
  load AS end_load,
  to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
  to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
FROM schedules
WHERE project_id = $1
ORDER BY vehicle_id, arrival
`

func (q *Queries) DBGetSchedule(ctx context.Context, projectID int64) ([]util.Schedule, error) {
	rows, err := q.db.Query(ctx, getSchedules, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanScheduleRows(rows)
}

const updateSchedule = `-- name: UpdateSchedule :one
UPDATE schedules
SET
  location_index = coord_to_id($2, $3), service = $4, delivery = $5,
  pickup = $6, skills = $7, priority = $8, project_id = $9, data = $10
WHERE id = $1 AND deleted = FALSE
RETURNING id, location_index, service, delivery, pickup, skills, priority, project_id, data, created_at, updated_at, deleted
`

const deleteSchedule = `
DELETE FROM schedules
WHERE project_id = $1
`

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
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.ProjectID,
			&i.VehicleID,
			&jobID,
			&shipmentID,
			&breakID,
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
		i.JobID = *jobID
		i.ShipmentID = *shipmentID
		i.BreakID = *breakID

		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
