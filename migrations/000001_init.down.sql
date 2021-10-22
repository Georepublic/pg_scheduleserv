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
  SELECT string_agg('DROP TRIGGER IF EXISTS tgr_updated_at_field
    ON ' || quote_ident(T) || ';', E'\n')
  FROM unnest('{locations, projects, project_locations, jobs,
    jobs_time_windows, shipments, shipments_time_windows, vehicles,
    breaks, breaks_time_windows, matrix}'::text[]) T
  );
END
$$;
DROP FUNCTION IF EXISTS tgr_updated_at_field_func;

DROP TRIGGER IF EXISTS tgr_matrix_insert ON matrix;
DROP FUNCTION IF EXISTS tgr_matrix_insert_func;

DROP TRIGGER IF EXISTS tgr_project_locations_insert ON project_locations;
DROP FUNCTION IF EXISTS tgr_project_locations_insert_func;

DROP TRIGGER IF EXISTS tgr_vehicles_insert ON vehicles;
DROP FUNCTION IF EXISTS tgr_vehicles_insert_func;

DROP TRIGGER IF EXISTS tgr_shipments_insert ON shipments;
DROP FUNCTION IF EXISTS tgr_shipments_insert_func;

DROP TRIGGER IF EXISTS tgr_jobs_insert ON jobs;
DROP FUNCTION IF EXISTS tgr_jobs_insert_func;

DROP TABLE IF EXISTS matrix;
DROP TABLE IF EXISTS breaks_time_windows;
DROP TABLE IF EXISTS breaks;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS shipments_time_windows;
DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS jobs_time_windows;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS project_locations;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS locations;

DROP FUNCTION IF EXISTS random_bigint;
DROP FUNCTION IF EXISTS id_to_geom;
DROP FUNCTION IF EXISTS geom_to_id;
DROP FUNCTION IF EXISTS coord_to_id;

DROP EXTENSION IF EXISTS vrprouting;
DROP EXTENSION IF EXISTS pgrouting;
DROP EXTENSION IF EXISTS postgis;

END;
