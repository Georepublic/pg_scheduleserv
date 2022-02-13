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


CREATE OR REPLACE FUNCTION get_coordinates_from_id(id BIGINT)
RETURNS FLOAT[] AS $BODY$
DECLARE
  latitude FLOAT;
  longitude FLOAT;
BEGIN
  latitude = (id / 100000000) / 10000.0;
  IF latitude >= 1000 THEN
    latitude = -(latitude - 1000);
  END IF;
  longitude = (id - id / 100000000 * 100000000) / 10000.0;
  IF longitude >= 1000 THEN
    longitude = -(longitude - 1000);
  END IF;
  RETURN ARRAY[latitude, longitude];
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


DO $$ BEGIN
  CREATE TYPE distance_calc_type AS ENUM ('euclidean', 'valhalla', 'osrm');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

-- PROJECTS TABLE start
CREATE TABLE IF NOT EXISTS projects (
  id                BIGINT                DEFAULT random_bigint() PRIMARY KEY,
  name              VARCHAR               NOT NULL,
  distance_calc     DISTANCE_CALC_TYPE    NOT NULL DEFAULT 'euclidean',
  exploration_level INTEGER               NOT NULL DEFAULT 5,
  timeout           INTERVAL              NOT NULL DEFAULT '00:10:00'::INTERVAL,
  max_shift         INTERVAL              NOT NULL DEFAULT '00:30:00'::INTERVAL,

  data        JSONB     NOT NULL DEFAULT '{}'::JSONB,
  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(exploration_level >= 0 AND exploration_level <= 5)
);
-- PROJECTS TABLE end


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
  status          TEXT      NOT NULL DEFAULT 'unscheduled'::TEXT,
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
  status            TEXT      NOT NULL DEFAULT 'unscheduled'::TEXT,
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


-- Create schedule for a project (such that any previous scheduled tasks are not likely to be unscheduled)
CREATE OR REPLACE FUNCTION create_schedule(
  project_id_param BIGINT,
  start_ids BIGINT[],
  end_ids BIGINT[],
  durations BIGINT[]
)
RETURNS void
AS $BODY$

  CREATE TABLE schedules_copy AS TABLE schedules;

  -- DELETE the schedules without changing the status field of jobs/shipments. Status field will be set by insert trigger later.
  ALTER TABLE schedules DISABLE TRIGGER tgr_schedule_delete;
  DELETE FROM schedules WHERE project_id = project_id_param;
  ALTER TABLE schedules ENABLE TRIGGER tgr_schedule_delete;

  WITH delta AS (SELECT max_shift FROM projects WHERE id = project_id_param)
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
    project_id_param::BIGINT, vehicle_id, location_id, task_id, vehicle_data, task_data,
    arrival, travel_time, setup_time, service_time, waiting_time, departure, load
  FROM vrp_vroom(
    -- jobs (Unscheduled jobs + Scheduled jobs with 100 priority)
    'SELECT id, location_id, setup, service, delivery, pickup, skills, priority
     FROM jobs WHERE project_id = ' || project_id_param || ' AND status = ''unscheduled'' AND deleted = FALSE
     UNION
     SELECT id, location_id, setup, service, delivery, pickup, skills, 100 AS priority
     FROM jobs WHERE project_id = ' || project_id_param || ' AND status = ''scheduled'' AND deleted = FALSE',

    -- jobs_time_windows (For unscheduled, select original time windows. For scheduled, alter the time window with a delta interval from the arrival time)
    'SELECT J.id AS id, tw_open, tw_close
     FROM jobs_time_windows TW LEFT JOIN jobs J ON(TW.id = J.id)
     WHERE status = ''unscheduled'' AND project_id = ' || project_id_param || '
    UNION
     SELECT
      J.id AS id,
      GREATEST(tw_open, arrival - $$' || (SELECT * FROM delta) || '$$::INTERVAL) AS tw_open,
      LEAST(tw_close, arrival + $$' || (SELECT * FROM delta) || '$$::INTERVAL) AS tw_close
     FROM jobs_time_windows TW RIGHT JOIN jobs J ON(TW.id = J.id) JOIN schedules_copy S ON (J.id = S.task_id)
     WHERE
      GREATEST(tw_open, arrival - $$' || (SELECT * FROM delta) || '$$::INTERVAL) <= LEAST(tw_close, arrival + $$' || (SELECT * FROM delta) || '$$::INTERVAL)
      AND status = ''scheduled'' AND type = ''job'' AND J.project_id = ' || project_id_param || ' ORDER BY id, tw_open',

    -- shipments (Unscheduled shipments + Scheduled shipments with 100 priority)
    'SELECT id, p_location_id, p_setup, p_service, d_location_id, d_setup, d_service, amount, skills, priority
     FROM shipments WHERE project_id = ' || project_id_param || ' AND status = ''unscheduled'' AND deleted = FALSE
     UNION
     SELECT id, p_location_id, p_setup, p_service, d_location_id, d_setup, d_service, amount, skills, 100 AS priority
     FROM shipments WHERE project_id = ' || project_id_param || ' AND status = ''scheduled'' AND deleted = FALSE',

    -- shipments_time_windows
    -- For unscheduled, select original time windows.
    -- For scheduled, alter the time window with a delta interval from the arrival time
    -- TODO: When time windows are "edited" such that the delta range falls outside new time windows, then the time window is ignored because the <= condition fails
    'SELECT S.id AS id, kind, tw_open, tw_close
     FROM shipments_time_windows TW LEFT JOIN shipments S ON(TW.id = S.id)
     WHERE status = ''unscheduled'' AND project_id = ' || project_id_param || '
    UNION
     SELECT
      S.id AS id,
      kind,
      GREATEST(tw_open, arrival - $$' || (SELECT * FROM delta) || '$$::INTERVAL) AS tw_open,
      LEAST(tw_close, arrival + $$' || (SELECT * FROM delta) || '$$::INTERVAL) AS tw_close
     FROM shipments_time_windows TW RIGHT JOIN shipments S ON(TW.id = S.id) JOIN schedules_copy S2 ON (S.id = S2.task_id)
     WHERE
      GREATEST(tw_open, arrival - $$' || (SELECT * FROM delta) || '$$::INTERVAL) <= LEAST(tw_close, arrival + $$' || (SELECT * FROM delta) || '$$::INTERVAL)
      AND status = ''scheduled'' AND ((type = ''pickup'' AND kind = ''p'') OR (type = ''delivery'' AND kind = ''d''))
      AND S.project_id = ' || project_id_param || ' ORDER BY id, tw_open',

    -- vehicles
    'SELECT * FROM vehicles WHERE deleted = FALSE AND project_id = ' || project_id_param || '',

    -- breaks
    'SELECT * FROM breaks WHERE deleted = FALSE',
    'SELECT * FROM breaks_time_windows ORDER BY id, tw_open',

    -- matrix
    'SELECT unnest(ARRAY[' || array_to_string(start_ids, ',') || ']::BIGINT[]) AS start_id,
     unnest(ARRAY[' || array_to_string(end_ids, ',') || ']::BIGINT[]) AS end_id,
     make_interval(secs => unnest(ARRAY[' || array_to_string(durations, ',') || ']::BIGINT[])) AS duration',

    exploration_level => (SELECT exploration_level FROM projects WHERE id = project_id_param),
    timeout => (SELECT timeout FROM projects WHERE id = project_id_param)
  );
  DROP TABLE schedules_copy;
$BODY$ LANGUAGE sql VOLATILE;


-- Create schedule for a project (fresh scheduling, deleting any previous schedule)
CREATE OR REPLACE FUNCTION create_fresh_schedule(
  project_id_param BIGINT,
  start_ids BIGINT[],
  end_ids BIGINT[],
  durations BIGINT[]
)
RETURNS void
AS $BODY$
  DELETE FROM schedules WHERE project_id = project_id_param;
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
    project_id_param::BIGINT, vehicle_id, location_id, task_id, vehicle_data, task_data,
    arrival, travel_time, setup_time, service_time, waiting_time, departure, load
  FROM vrp_vroom(
    'SELECT * FROM jobs WHERE deleted = FALSE AND project_id = ' || project_id_param,
    'SELECT * FROM jobs_time_windows ORDER BY id, tw_open',
    'SELECT * FROM shipments WHERE deleted = FALSE AND project_id = ' || project_id_param,
    'SELECT * FROM shipments_time_windows ORDER BY id, tw_open',
    'SELECT * FROM vehicles WHERE deleted = FALSE AND project_id = ' || project_id_param,
    'SELECT * FROM breaks WHERE deleted = FALSE',
    'SELECT * FROM breaks_time_windows ORDER BY id, tw_open',
    'SELECT unnest(ARRAY[' || array_to_string(start_ids, ',') || ']::BIGINT[]) AS start_id,
     unnest(ARRAY[' || array_to_string(end_ids, ',') || ']::BIGINT[]) AS end_id,
     make_interval(secs => unnest(ARRAY[' || array_to_string(durations, ',') || ']::BIGINT[])) AS duration',
    exploration_level => (SELECT exploration_level FROM projects WHERE id = project_id_param),
    timeout => (SELECT timeout FROM projects WHERE id = project_id_param)
  );
$BODY$ LANGUAGE sql VOLATILE;


-- AFTER INSERT Trigger for schedule, update the status field in jobs or shipments for all the rows
CREATE OR REPLACE FUNCTION tgr_schedule_insert_func()
RETURNS TRIGGER
AS $trig$
DECLARE
  project_id_param BIGINT;
BEGIN
  SELECT DISTINCT project_id FROM new_table INTO project_id_param;

  -- Update schedule jobs status
  UPDATE jobs SET status = 'scheduled'::TEXT
    WHERE project_id = project_id_param AND id IN (
      SELECT task_id FROM schedules
      WHERE project_id = project_id_param AND type = 'job'::STEP_TYPE AND vehicle_id > 0
    );

  -- Update unschedule jobs status
  UPDATE jobs SET status = 'unscheduled'::TEXT
    WHERE project_id = project_id_param AND id IN (
      SELECT task_id FROM schedules
      WHERE project_id = project_id_param AND type = 'job'::STEP_TYPE AND vehicle_id = -1
    );

  -- Update schedule shipments status
  UPDATE shipments SET status = 'scheduled'::TEXT
    WHERE project_id = project_id_param AND id IN (
      SELECT task_id FROM schedules
      WHERE project_id = project_id_param AND type = 'pickup'::STEP_TYPE AND vehicle_id > 0
    );  -- Pickup and delivery always occur with same id

  -- Update unschedule shipments status
  UPDATE shipments SET status = 'unscheduled'::TEXT
    WHERE project_id = project_id_param AND id IN (
      SELECT task_id FROM schedules
      WHERE project_id = project_id_param AND type = 'pickup'::STEP_TYPE AND vehicle_id = -1
    );

  RETURN NULL;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_schedule_insert
AFTER INSERT ON schedules
REFERENCING NEW TABLE AS new_table
FOR EACH STATEMENT EXECUTE FUNCTION tgr_schedule_insert_func();


-- AFTER DELETE Trigger for schedule, update the status field in jobs or shipments for the deleted rows
CREATE OR REPLACE FUNCTION tgr_schedule_delete_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  -- Update jobs status as unscheduled
  UPDATE jobs SET status = 'unscheduled'::TEXT
    WHERE id IN (
      SELECT task_id FROM old_table
      WHERE type = 'job'::STEP_TYPE
    );

  -- Pickup and delivery always occur with the same id
  -- Update shipments status as unscheduled
  UPDATE shipments SET status = 'unscheduled'::TEXT
    WHERE id IN (
      SELECT task_id FROM old_table
      WHERE type = 'pickup'::STEP_TYPE
    );

  RETURN NULL;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_schedule_delete
AFTER DELETE ON schedules
REFERENCING OLD TABLE AS old_table
FOR EACH STATEMENT EXECUTE FUNCTION tgr_schedule_delete_func();


-------------------------------------------------------------------------------
-- TRIGGERS
-------------------------------------------------------------------------------

-- BEFORE INSERT OR UPDATE Trigger for jobs, inserts rows into locations
CREATE OR REPLACE FUNCTION tgr_jobs_insert_update_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO locations (id)
  SELECT NEW.location_id
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_jobs_insert_update
BEFORE INSERT OR UPDATE ON jobs
FOR EACH ROW EXECUTE PROCEDURE tgr_jobs_insert_update_func();


-- BEFORE INSERT OR UPDATE Trigger for shipments, inserts rows into locations
CREATE OR REPLACE FUNCTION tgr_shipments_insert_update_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO locations (id)
  SELECT NEW.p_location_id
  UNION
  SELECT NEW.d_location_id
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_shipments_insert_update
BEFORE INSERT OR UPDATE ON shipments
FOR EACH ROW EXECUTE PROCEDURE tgr_shipments_insert_update_func();


-- BEFORE INSERT OR UPDATE Trigger for vehicles, inserts rows into locations
CREATE OR REPLACE FUNCTION tgr_vehicles_insert_update_func()
RETURNS TRIGGER
AS $trig$
BEGIN
  INSERT INTO locations (id)
  SELECT NEW.start_id
  UNION
  SELECT NEW.end_id
  ON CONFLICT DO NOTHING;

  RETURN NEW;
END;
$trig$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_vehicles_insert_update
BEFORE INSERT OR UPDATE ON vehicles
FOR EACH ROW EXECUTE PROCEDURE tgr_vehicles_insert_update_func();




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
  FROM unnest('{locations, projects, jobs,
    jobs_time_windows, shipments, shipments_time_windows, vehicles,
    breaks, breaks_time_windows, schedules}'::text[]) T
  );
END
$$;

END;
