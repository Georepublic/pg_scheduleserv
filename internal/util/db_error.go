/*GRP-GNU-AGPL******************************************************************

File: db_error.go

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

package util

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

func HandleDBError(err error) error {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "breaks_time_windows_id_fkey":
				err = fmt.Errorf("Break with the given 'break_id' does not exist")
			case "project_locations_project_id_fkey":
				err = fmt.Errorf("Project with the given 'project_id' does not exist")
			case "jobs_time_windows_id_fkey":
				err = fmt.Errorf("Job with the given 'job_id' does not exist")
			case "shipments_time_windows_id_fkey":
				err = fmt.Errorf("Shipment with the given 'shipment_id' does not exist")
			case "breaks_vehicle_id_fkey":
				err = fmt.Errorf("Vehicle with the given 'vehicle_id' does not exist")
			case "jobs_check":
				err = fmt.Errorf("Field 'pickup' and 'delivery' must have same length")
			case "breaks_time_windows_check":
				err = fmt.Errorf("Field 'tw_open' must be less than or equal to field 'tw_close'")
			case "jobs_time_windows_check":
				err = fmt.Errorf("Field 'tw_open' must be less than or equal to field 'tw_close'")
			case "vehicles_check":
				err = fmt.Errorf("Field 'tw_open' must be less than or equal to field 'tw_close'")
			case "shipments_time_windows_check":
				err = fmt.Errorf("Field 'tw_open' must be less than or equal to field 'tw_close'")
			case "jobs_time_windows_pkey":
				err = fmt.Errorf("Jobs time window with given values already exist")
			case "shipments_time_windows_pkey":
				err = fmt.Errorf("Shipments time window with given values already exist")
			case "breaks_time_windows_pkey":
				err = fmt.Errorf("Breaks time window with given values already exist")
			}
		}
	}
	return err
}
