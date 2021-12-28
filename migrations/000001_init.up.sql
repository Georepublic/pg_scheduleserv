/*GRP-GNU-AGPL******************************************************************

File: 000001_init.up.sql - Initial migrations up file

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

CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgrouting;
CREATE EXTENSION IF NOT EXISTS vrprouting WITH VERSION '0.3.0';


-------------------------------------------------------------------------------
-- FUNCTIONS
-------------------------------------------------------------------------------


-- Generate id of location from the latitude and longitude of the coordinates,
-- considering the values upto 4 decimal places.
-- id = 16 digits (or less), where the last 8 digits denote the longitude, and
-- the rest 8 or less digits denote the latitude.
CREATE OR REPLACE FUNCTION coord_to_id(latitude FLOAT, longitude FLOAT)
RETURNS BIGINT
AS $BODY$
DECLARE
  -- First digit out of 8 denote whether the value is positive (0) or negative (1).
  lat_prefix CHAR(1) := '0';
  lon_prefix CHAR(1) := '0';
BEGIN
  IF latitude < 0 THEN
    lat_prefix := '1';
  END IF;
  IF longitude < 0 THEN
    lon_prefix := '1';
  END IF;
  RETURN
    CONCAT(
      lat_prefix,
      LPAD(ROUND(10000 * ABS(latitude))::TEXT, 7, '0'),
      lon_prefix,
      LPAD(ROUND(10000 * ABS(longitude))::TEXT, 7, '0')
    )::BIGINT;
END;
$BODY$ LANGUAGE plpgsql IMMUTABLE STRICT;


-- Generate id by concatenating the latitude and longitude upto 4 decimal places
CREATE OR REPLACE FUNCTION geom_to_id(location geometry)
RETURNS BIGINT
AS $BODY$
  SELECT
    coord_to_id(ST_Y(location), ST_X(location));
$BODY$ LANGUAGE SQL IMMUTABLE STRICT;


-- Generate geometry from the id (upto 4 decimal places of latitude and longitude)
CREATE OR REPLACE FUNCTION id_to_geom(id BIGINT)
RETURNS geometry
AS $BODY$
DECLARE
  latitude FLOAT;
  longitude FLOAT;
BEGIN
  latitude := (id/100000000)::FLOAT/10000::FLOAT;
  IF latitude >= 1000 THEN
    latitude = -(latitude - 1000);
  END IF;
  longitude := (id - id/100000000*100000000)::FLOAT/10000::FLOAT;
  IF longitude >= 1000 THEN
    longitude = -(longitude - 1000);
  END IF;
  RETURN
    ST_SetSRID(
      ST_Point(longitude, latitude), 4326
    );
END;
$BODY$ LANGUAGE plpgsql IMMUTABLE STRICT;


-- Generate a random bigint by using the first half of 128-bit hex uuid
CREATE OR REPLACE FUNCTION random_bigint()
RETURNS BIGINT
AS $BODY$
  SELECT
    ABS(
      ('x' || translate(gen_random_uuid()::TEXT, '-', ''))::BIT(64)::BIGINT
    );
$BODY$ LANGUAGE SQL VOLATILE STRICT;



-------------------------------------------------------------------------------
-- TABLES
-------------------------------------------------------------------------------

-- LOCATIONS TABLE start
CREATE TABLE IF NOT EXISTS locations (
  id          BIGINT PRIMARY KEY,
  location    geometry(Point, 4326) GENERATED ALWAYS AS (id_to_geom(id)) STORED,
  latitude    FLOAT GENERATED ALWAYS AS (ST_Y(id_to_geom(id))) STORED,
  longitude   FLOAT GENERATED ALWAYS AS (ST_X(id_to_geom(id))) STORED,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,

  CHECK(latitude >= -90 AND latitude <= 90),
  CHECK(longitude >= -180 AND longitude <= 180)
);
-- LOCATIONS TABLE end


-- PROJECTS TABLE start
CREATE TABLE IF NOT EXISTS projects (
  id                BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  name              VARCHAR   NOT NULL,
  exploration_level INTEGER   NOT NULL DEFAULT 5,
  timeout           INTERVAL  NOT NULL DEFAULT '-00:00:01'::INTERVAL,

  data        JSONB     NOT NULL DEFAULT '{}'::JSONB,
  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(exploration_level >= 0 AND exploration_level <= 5)
);
-- PROJECTS TABLE end


-- PROJECT LOCATIONS TABLE start
-- Aggregates all the locations in a project, eases inserting rows in the matrix
CREATE TABLE IF NOT EXISTS project_locations (
  project_id  BIGINT    NOT NULL REFERENCES projects(id),
  location_id BIGINT    NOT NULL REFERENCES locations(id),

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp
);
-- PROJECT LOCATIONS TABLE end


-- JOBS TABLE start
CREATE TABLE IF NOT EXISTS jobs (
  id              BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  location_id     BIGINT    NOT NULL REFERENCES locations(id),
  setup           INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  service         INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  delivery        BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  pickup          BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  skills          INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
  priority        INTEGER   NOT NULL DEFAULT 0,

  project_id      BIGINT    NOT NULL REFERENCES projects(id),

  data            JSONB     NOT NULL DEFAULT '{}'::JSONB,
  created_at      TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at      TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted         BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(service >= '00:00:00'::INTERVAL),
  CHECK(0 <= ALL(delivery)),
  CHECK(0 <= ALL(pickup)),
  CHECK(0 <= ALL(skills)),
  CHECK(priority >= 0 AND priority <= 100),
  CHECK(array_length(delivery, 1) = array_length(pickup, 1))
);
-- JOBS TABLE end


-- JOBS TIME WINDOWS TABLE start
CREATE TABLE IF NOT EXISTS jobs_time_windows (
  id          BIGINT    NOT NULL REFERENCES jobs(id),
  tw_open     TIMESTAMP NOT NULL,
  tw_close    TIMESTAMP NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,

  PRIMARY KEY(id, tw_open, tw_close),

  CHECK(tw_open <= tw_close)
);
-- JOBS TIME WINDOWS TABLE end


-- SHIPMENTS TABLE start
CREATE TABLE IF NOT EXISTS shipments (
  id                BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  p_location_id     BIGINT    NOT NULL REFERENCES locations(id),
  p_setup           INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  p_service         INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  d_location_id     BIGINT    NOT NULL REFERENCES locations(id),
  d_setup           INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  d_service         INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  amount            BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  skills            INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
  priority          INTEGER   NOT NULL DEFAULT 0,

  project_id        BIGINT    NOT NULL REFERENCES projects(id),

  data              JSONB     NOT NULL DEFAULT '{}'::JSONB,
  created_at        TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at        TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted           BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(p_service >= '00:00:00'::INTERVAL),
  CHECK(d_service >= '00:00:00'::INTERVAL),
  CHECK(0 <= ALL(amount)),
  CHECK(0 <= ALL(skills)),
  CHECK(priority >= 0 AND priority <= 100)
);
-- SHIPMENTS TABLE end


-- SHIPMENTS TIME WINDOWS TABLE start
CREATE TABLE IF NOT EXISTS shipments_time_windows (
  id          BIGINT    NOT NULL REFERENCES shipments(id),
  kind        CHAR(1)   NOT NULL,
  tw_open     TIMESTAMP NOT NULL,
  tw_close    TIMESTAMP NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,

  PRIMARY KEY(id, kind, tw_open, tw_close),

  CHECK(kind = 'p' OR kind = 'd'),
  CHECK(tw_open <= tw_close)
);
-- SHIPMENTS TIME WINDOWS TABLE end


-- VEHICLES TABLE start
CREATE TABLE IF NOT EXISTS vehicles (
  id            BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  start_id      BIGINT    NOT NULL REFERENCES locations(id),
  end_id        BIGINT    NOT NULL REFERENCES locations(id),
  capacity      BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  skills        INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
  tw_open       TIMESTAMP NOT NULL DEFAULT (to_timestamp(0) at time zone 'UTC'),
  tw_close      TIMESTAMP NOT NULL DEFAULT (to_timestamp(2147483647) at time zone 'UTC'),
  speed_factor  FLOAT     NOT NULL DEFAULT 1.0,
  max_tasks     INTEGER   NOT NULL DEFAULT 2147483647,

  project_id    BIGINT    NOT NULL REFERENCES projects(id),

  data          JSONB     NOT NULL DEFAULT '{}'::JSONB,
  created_at    TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at    TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted       BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(0 <= ALL(capacity)),
  CHECK(0 <= ALL(skills)),
  CHECK(tw_open <= tw_close),
  CHECK(speed_factor > 0.0),
  CHECK(max_tasks >= 0)
);
-- VEHICLES TABLE end


-- BREAKS TABLE start
CREATE TABLE IF NOT EXISTS breaks (
  id          BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  vehicle_id  BIGINT    NOT NULL REFERENCES vehicles(id),
  service     INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,

  data        JSONB     NOT NULL DEFAULT '{}'::JSONB,
  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(service >= '00:00:00'::INTERVAL)
);
-- BREAKS TABLE end


-- BREAKS TIME WINDOWS TABLE start
CREATE TABLE IF NOT EXISTS breaks_time_windows (
  id          BIGINT    NOT NULL REFERENCES breaks(id),
  tw_open     TIMESTAMP NOT NULL,
  tw_close    TIMESTAMP NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,

  PRIMARY KEY(id, tw_open, tw_close),

  CHECK(tw_open <= tw_close)
);
-- BREAKS TIME WINDOWS TABLE end


-- MATRIX TABLE start
CREATE TABLE IF NOT EXISTS matrix (
  start_id    BIGINT    NOT NULL REFERENCES locations(id),
  end_id      BIGINT    NOT NULL REFERENCES locations(id),
  duration    INTERVAL  NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,

  PRIMARY KEY (start_id, end_id),

  CHECK(duration >= '00:00:00'::INTERVAL)
);
-- MATRIX TABLE end


DO $$ BEGIN
  CREATE TYPE step_type AS ENUM ('summary', 'start', 'job', 'pickup', 'delivery', 'break', 'end');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;


-- SCHEDULES TABLE start
CREATE TABLE IF NOT EXISTS schedules (
  type          STEP_TYPE   NOT NULL,
  project_id    BIGINT      NOT NULL REFERENCES projects(id),

  vehicle_id    BIGINT      NOT NULL,
  task_id       BIGINT      NOT NULL,
  location_id   BIGINT      NOT NULL,

  arrival       TIMESTAMP   NOT NULL,
  departure     TIMESTAMP   NOT NULL,
  travel_time   INTERVAL    NOT NULL,
  setup_time    INTERVAL    NOT NULL,
  service_time  INTERVAL    NOT NULL,
  waiting_time  INTERVAL    NOT NULL,
  load          BIGINT[]    NOT NULL,

  vehicle_data  JSONB       NOT NULL,
  task_data     JSONB       NOT NULL,

  created_at    TIMESTAMP   NOT NULL DEFAULT current_timestamp,
  updated_at    TIMESTAMP   NOT NULL DEFAULT current_timestamp,

  PRIMARY KEY (task_id, type, vehicle_id, project_id),
  CHECK(travel_time >= '00:00:00'::INTERVAL),
  CHECK(setup_time >= '00:00:00'::INTERVAL),
  CHECK(service_time >= '00:00:00'::INTERVAL),
  CHECK(waiting_time >= '00:00:00'::INTERVAL),
  CHECK(0 <= ALL(load))
);
-- SCHEDULE TABLE end

-- Create schedule for a project
CREATE OR REPLACE FUNCTION create_schedule(BIGINT)
RETURNS void
AS $BODY$
  DELETE FROM schedules WHERE project_id = $1;
  INSERT INTO schedules
    (type, project_id, vehicle_id, location_id, task_id, vehicle_data, task_data,
    arrival, travel_time, setup_time, service_time, waiting_time, departure, load)
  SELECT
    CASE
      WHEN step_type = 0 THEN 'summary'::STEP_TYPE
      WHEN step_type = 1 THEN 'start'::STEP_TYPE
      WHEN step_type = 2 THEN 'job'::STEP_TYPE
      WHEN step_type = 3 THEN 'pickup'::STEP_TYPE
      WHEN step_type = 4 THEN 'delivery'::STEP_TYPE
      WHEN step_type = 5 THEN 'break'::STEP_TYPE
      WHEN step_type = 6 THEN 'end'::STEP_TYPE
    END,
    $1::BIGINT, vehicle_id, location_id, task_id, vehicle_data, task_data,
    arrival, travel_time, setup_time, service_time, waiting_time, departure, load
  FROM vrp_vroom(
    'SELECT * FROM jobs WHERE deleted = FALSE AND project_id = ' || $1,
    'SELECT * FROM jobs_time_windows ORDER BY id, tw_open',
    'SELECT * FROM shipments WHERE deleted = FALSE AND project_id = ' || $1,
    'SELECT * FROM shipments_time_windows ORDER BY id, tw_open',
    'SELECT * FROM vehicles WHERE deleted = FALSE AND project_id = ' || $1,
    'SELECT * FROM breaks WHERE deleted = FALSE',
    'SELECT * FROM breaks_time_windows ORDER BY id, tw_open',
    'SELECT * FROM matrix',
    exploration_level => (SELECT exploration_level FROM projects WHERE id = $1),
    timeout => (SELECT timeout FROM projects WHERE id = $1)
  );
$BODY$ LANGUAGE sql VOLATILE STRICT;



-------------------------------------------------------------------------------
-- TRIGGERS
-------------------------------------------------------------------------------

-- BEFORE INSERT OR UPDATE Trigger for jobs, inserts rows into locations and project_locations
CREATE OR REPLACE FUNCTION tgr_jobs_insert_update_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO locations (id)
  SELECT NEW.location_id
  ON CONFLICT DO NOTHING;

  INSERT INTO project_locations (project_id, location_id)
  SELECT NEW.project_id, NEW.location_id
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_jobs_insert_update
BEFORE INSERT OR UPDATE ON jobs
FOR EACH ROW EXECUTE PROCEDURE tgr_jobs_insert_update_func();


-- BEFORE INSERT OR UPDATE Trigger for shipments, inserts rows into locations and project_locations
CREATE OR REPLACE FUNCTION tgr_shipments_insert_update_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO locations (id)
  SELECT NEW.p_location_id
  UNION
  SELECT NEW.d_location_id
  ON CONFLICT DO NOTHING;

  INSERT INTO project_locations (project_id, location_id)
  SELECT NEW.project_id, NEW.p_location_id
  UNION
  SELECT NEW.project_id, NEW.d_location_id
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_shipments_insert_update
BEFORE INSERT OR UPDATE ON shipments
FOR EACH ROW EXECUTE PROCEDURE tgr_shipments_insert_update_func();


-- BEFORE INSERT OR UPDATE Trigger for vehicles, inserts rows into locations and project_locations
CREATE OR REPLACE FUNCTION tgr_vehicles_insert_update_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO locations (id)
  SELECT NEW.start_id
  UNION
  SELECT NEW.end_id
  ON CONFLICT DO NOTHING;

  INSERT INTO project_locations (project_id, location_id)
  SELECT NEW.project_id, NEW.start_id
  UNION
  SELECT NEW.project_id, NEW.end_id
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_vehicles_insert_update
BEFORE INSERT OR UPDATE ON vehicles
FOR EACH ROW EXECUTE PROCEDURE tgr_vehicles_insert_update_func();


-- AFTER INSERT Trigger for project locations, inserts rows into matrix
CREATE OR REPLACE FUNCTION tgr_project_locations_insert_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO matrix(start_id, end_id, duration)
  SELECT
    NEW.location_id,
    PL.location_id,
    make_interval(secs => ROUND(ST_distance(
      id_to_geom(NEW.location_id)::geography,
      id_to_geom(PL.location_id)::geography)
    ))
    FROM project_locations AS PL
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_project_locations_insert
AFTER INSERT ON project_locations
FOR EACH ROW EXECUTE PROCEDURE tgr_project_locations_insert_func();


-- AFTER INSERT Trigger for matrix, inserts reverse direction cost
CREATE OR REPLACE FUNCTION tgr_matrix_insert_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO matrix(start_id, end_id, duration)
  SELECT NEW.end_id, NEW.start_id, NEW.duration
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_matrix_insert
AFTER INSERT ON matrix
FOR EACH ROW EXECUTE PROCEDURE tgr_matrix_insert_func();


-- BEFORE UPDATE TRIGGER for all tables, auto-update updated_at field
CREATE OR REPLACE FUNCTION tgr_updated_at_field_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  NEW.updated_at = current_timestamp;
  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

DO
$$
BEGIN
  EXECUTE (
  SELECT string_agg('CREATE TRIGGER tgr_updated_at_field
    BEFORE UPDATE ON ' || quote_ident(T) || '
    FOR EACH ROW EXECUTE PROCEDURE tgr_updated_at_field_func();', E'\n')
  FROM unnest('{locations, projects, project_locations, jobs,
    jobs_time_windows, shipments, shipments_time_windows, vehicles,
    breaks, breaks_time_windows, matrix, schedules}'::text[]) T
  );
END
$$;

END;
