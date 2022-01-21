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
    3909655254191459782::BIGINT, vehicle_id, location_id, task_id, vehicle_data, task_data,
    arrival, travel_time, setup_time, service_time, waiting_time, departure, load
  FROM vrp_vroom(
    'SELECT id, location_id, setup, service, delivery, pickup, skills, priority
     FROM jobs WHERE project_id = 3909655254191459782 AND status = ''unscheduled'' AND deleted = FALSE
     UNION
     SELECT id, location_id, setup, service, delivery, pickup, skills, 100 AS priority
     FROM jobs WHERE project_id = 3909655254191459782 AND status = ''scheduled'' AND deleted = FALSE',

    'SELECT J.id AS id, tw_open, tw_close
     FROM jobs_time_windows TW LEFT JOIN jobs J ON(TW.id = J.id)
     WHERE status = ''unscheduled'' AND project_id = 3909655254191459782
    UNION
     SELECT J.id AS id, GREATEST(tw_open, arrival - $$' || (SELECT * FROM delta) || '$$::INTERVAL) AS tw_open, LEAST(tw_close, arrival + $$' || (SELECT * FROM delta) || '$$::INTERVAL) AS tw_close
     FROM jobs_time_windows TW LEFT JOIN jobs J ON(TW.id = J.id) JOIN schedules S ON (J.id = S.task_id)
     WHERE status = ''scheduled'' AND type = ''job'' AND J.project_id = 3909655254191459782 ORDER BY id, tw_open',

    'SELECT id, p_location_id, p_setup, p_service, d_location_id, d_setup, d_service, amount, skills, priority
     FROM shipments WHERE project_id = 3909655254191459782 AND status = ''unscheduled'' AND deleted = FALSE
     UNION
     SELECT id, p_location_id, p_setup, p_service, d_location_id, d_setup, d_service, amount, skills, 100 AS priority
     FROM shipments WHERE project_id = 3909655254191459782 AND status = ''scheduled'' AND deleted = FALSE',

    'SELECT id, kind, tw_open, tw_close FROM shipments_time_windows ORDER BY id, tw_open',

    'SELECT * FROM vehicles WHERE deleted = FALSE AND project_id = 3909655254191459782',

    'SELECT * FROM breaks WHERE deleted = FALSE',
    'SELECT * FROM breaks_time_windows ORDER BY id, tw_open',
    'SELECT * FROM matrix'
  );
