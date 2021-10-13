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
