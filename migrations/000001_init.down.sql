/*GRP-GNU-AGPL******************************************************************

File: 000001_init.down.sql - Initial migrations down file

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

BEGIN;

DO
$$
BEGIN
  EXECUTE (
  SELECT string_agg('DROP TRIGGER tgr_updated_at_field
    ON ' || quote_ident(T) || ';', E'\n')
  FROM unnest('{locations, projects, project_locations, jobs,
    jobs_time_windows, shipments, shipments_time_windows, vehicles,
    breaks, breaks_time_windows, matrix}'::text[]) T
  );
END
$$;
DROP FUNCTION tgr_updated_at_field_func;

DROP TRIGGER tgr_matrix_insert ON matrix;
DROP FUNCTION tgr_matrix_insert_func;

DROP TRIGGER tgr_project_locations_insert ON project_locations;
DROP FUNCTION tgr_project_locations_insert_func;

DROP TRIGGER tgr_vehicles_insert ON vehicles;
DROP FUNCTION tgr_vehicles_insert_func;

DROP TRIGGER tgr_shipments_insert ON shipments;
DROP FUNCTION tgr_shipments_insert_func;

DROP TRIGGER tgr_jobs_insert ON jobs;
DROP FUNCTION tgr_jobs_insert_func;

DROP TABLE matrix;
DROP TABLE breaks_time_windows;
DROP TABLE breaks;
DROP TABLE vehicles;
DROP TABLE shipments_time_windows;
DROP TABLE shipments;
DROP TABLE jobs_time_windows;
DROP TABLE jobs;
DROP TABLE project_locations;
DROP TABLE projects;
DROP TABLE locations;

DROP FUNCTION random_bigint;
DROP FUNCTION id_to_geom;
DROP FUNCTION geom_to_id;
DROP FUNCTION coord_to_id;

DROP EXTENSION vrprouting;
DROP EXTENSION pgrouting;
DROP EXTENSION postgis;

END;
