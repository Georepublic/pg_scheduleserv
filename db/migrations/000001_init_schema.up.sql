CREATE EXTENSION postgis;
CREATE EXTENSION pgrouting;
CREATE EXTENSION vrprouting;


-------------------------------------------------------------------------------
-- FUNCTIONS
-------------------------------------------------------------------------------


-- Generate id of location from the latitude and longitude of the coordinates,
-- considering the values upto 4 decimal places.
-- id = 16 digits (or less), where the last 8 digits denote the longitude, and
-- the rest 8 or less digits denote the latitude.
CREATE FUNCTION coord_to_id(latitude FLOAT, longitude FLOAT)
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
CREATE FUNCTION geom_to_id(location geometry)
RETURNS BIGINT
AS $BODY$
  SELECT
    coord_to_id(ST_Y(location), ST_X(location));
$BODY$ LANGUAGE SQL IMMUTABLE STRICT;


-- Generate geometry from the id (upto 4 decimal places of latitude and longitude)
CREATE FUNCTION id_to_geom(id BIGINT)
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
CREATE FUNCTION random_bigint()
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
CREATE TABLE locations (
  id          BIGINT PRIMARY KEY,
  location    geometry(Point, 4326) GENERATED ALWAYS AS (id_to_geom(id)) STORED,
  latitude    FLOAT GENERATED ALWAYS AS (ST_Y(id_to_geom(id))) STORED,
  longitude   FLOAT GENERATED ALWAYS AS (ST_X(id_to_geom(id))) STORED,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp
);
-- LOCATIONS TABLE end


-- PROJECTS TABLE start
CREATE TABLE projects (
  id          BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  name        VARCHAR   NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0)
);
-- PROJECTS TABLE end


-- PROJECT LOCATIONS TABLE start
-- Aggregates all the locations in a project, eases inserting rows in the matrix
CREATE TABLE project_locations (
  project_id  BIGINT    NOT NULL REFERENCES projects(id),
  location_id BIGINT    NOT NULL REFERENCES locations(id),

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp
);
-- PROJECT LOCATIONS TABLE end


-- JOBS TABLE start
CREATE TABLE jobs (
  id              BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  location_index  BIGINT    NOT NULL REFERENCES locations(id),
  service         INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  delivery        BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  pickup          BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  skills          INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
  priority        INTEGER   NOT NULL DEFAULT 0,

  created_at      TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at      TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted         BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(service >= '00:00:00'::INTERVAL),
  CHECK(0 <= ALL(delivery)),
  CHECK(0 <= ALL(pickup)),
  CHECK(0 <= ALL(skills)),
  CHECK(priority >= 0 AND priority <= 100)
);
-- JOBS TABLE end


-- JOBS TIME WINDOWS TABLE start
CREATE TABLE jobs_time_windows (
  id          BIGINT    NOT NULL REFERENCES jobs(id),
  tw_open     TIMESTAMP NOT NULL,
  tw_close    TIMESTAMP NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(tw_open <= tw_close)
);
-- JOBS TIME WINDOWS TABLE end


-- SHIPMENTS TABLE start
CREATE TABLE shipments (
  id                BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  p_location_index  BIGINT    NOT NULL REFERENCES locations(id),
  p_service         INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  d_location_index  BIGINT    NOT NULL REFERENCES locations(id),
  d_service         INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,
  amount            BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  skills            INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
  priority          INTEGER   NOT NULL DEFAULT 0,

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
CREATE TABLE shipments_time_windows (
  id          BIGINT    NOT NULL REFERENCES shipments(id),
  kind        CHAR(1)   NOT NULL,
  tw_open     TIMESTAMP NOT NULL,
  tw_close    TIMESTAMP NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(kind = 'p' OR kind = 'd'),
  CHECK(tw_open <= tw_close)
);
-- SHIPMENTS TIME WINDOWS TABLE end


-- VEHICLES TABLE start
CREATE TABLE vehicles (
  id            BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  start_index   BIGINT    NOT NULL REFERENCES locations(id),
  end_index     BIGINT    NOT NULL REFERENCES locations(id),
  capacity      BIGINT[]  NOT NULL DEFAULT ARRAY[]::BIGINT[],
  skills        INTEGER[] NOT NULL DEFAULT ARRAY[]::INTEGER[],
  tw_open       TIMESTAMP NOT NULL DEFAULT (to_timestamp(0) at time zone 'UTC'),
  tw_close      TIMESTAMP NOT NULL DEFAULT (to_timestamp(2147483647) at time zone 'UTC'),
  speed_factor  FLOAT     NOT NULL DEFAULT 1.0,

  created_at    TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at    TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted       BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(0 <= ALL(capacity)),
  CHECK(0 <= ALL(skills)),
  CHECK(tw_open <= tw_close),
  CHECK(speed_factor > 0.0)
);
-- VEHICLES TABLE end


-- BREAKS TABLE start
CREATE TABLE breaks (
  id          BIGINT    DEFAULT random_bigint() PRIMARY KEY,
  vehicle_id  BIGINT    NOT NULL REFERENCES vehicles(id),
  service     INTERVAL  NOT NULL DEFAULT '00:00:00'::INTERVAL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(id >= 0),
  CHECK(service >= '00:00:00'::INTERVAL)
);
-- BREAKS TABLE end


-- BREAKS TIME WINDOWS TABLE start
CREATE TABLE breaks_time_windows (
  id          BIGINT    NOT NULL REFERENCES breaks(id),
  tw_open     TIMESTAMP NOT NULL,
  tw_close    TIMESTAMP NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  deleted     BOOLEAN   NOT NULL DEFAULT FALSE,

  CHECK(tw_open <= tw_close)
);
-- BREAKS TIME WINDOWS TABLE end


-- MATRIX TABLE start
CREATE TABLE matrix (
  start_vid   BIGINT    NOT NULL REFERENCES locations(id),
  end_vid     BIGINT    NOT NULL REFERENCES locations(id),
  agg_cost    INTEGER   NOT NULL,

  created_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,
  updated_at  TIMESTAMP NOT NULL DEFAULT current_timestamp,

  PRIMARY KEY (start_vid, end_vid),

  CHECK(agg_cost >= 0)
);
-- MATRIX TABLE end
