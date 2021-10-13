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
