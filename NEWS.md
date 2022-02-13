# pg_scheduleserv Release Notes

## v0.2.0 Release Notes

To see all issues & pull requests closed by this release see the [Git closed milestone for v0.2.0](https://github.com/Georepublic/pg_scheduleserv/issues?q=milestone%3Av0.2.0+) on Github.

### New Features

- #26: Update to `vrprouting v0.3.0` and `vroom v1.11.0`.
- #21: Extend the schedule output to provide metadata and the complete schedule, using optional "overview" query parameter.
  - Also extended all API responses to provide message and status code.
- #12: Create a demo application
- #15, #22: Schedule new tasks of a project without changing the schedule of existing tasks, Schedule request with dropping all previous allocations.
  - Fresh scheduling (`fresh=true` query parameter in Schedule POST API endpoint):
    - Schedule request with dropping all previous allocations.
  - Normal scheduling - DEFAULT (`fresh=false` query parameter in Schedule POST API endpoint):
    - Schedule request such that previous tasks are not unscheduled, and are shifted by the "max_shift" interval.
    - Add "max_shift" field in the projects, with default value of `'00:30:00'::INTERVAL` (30 mins).
- #33: Add "duration_calc" field in projects to choose whether to use OSRM API (2a39915), Valhalla API (23983a4), or Euclidean distance with speed = 9.0 m/sec (3423468) for computing durations.
- d282638: Add "exploration_level" (default = 5) and "timeout" (default = "00:10:00"::INTERVAL) field to projects.
- b313ca2: Enable Cross-Origin Resource Sharing (CORS)- currently, allow all.
- ff48656: Change pgx connection to connection pool for concurrent requests.
- e0e9354, 56a2400, 60490c8: Included time_window as job, shipment(p&d), and break field, instead of separate API endpoint.
- 4b1750f: Remove matrix and project_locations table, so that the durations are calculated only when the schedule is requested.

### Fixes

- 13b4805: Change all fields with INTERVAL to HH:MM:SS format.
- 1ba6d35: Change datetime output from `yyyy-mm-dd hh:tt:ss` to `yyyy-mm-ddThh:tt:ss`.
- b3b5817: Returning 404 NotFound Error instead of empty list when the project/task id is wrong.
- ff35d18: Error in displaying schedule with multiple vehicles.
- 7daa973: Change ical schedule id for start/end entry, to make it a unique entry.
- 3e9bc59: Fix Break List API error, displaying custom error message for DB Errors in Job and Shipment.
- 697c6ec: Change data field in shipment to p_data and d_data.

## v0.1.0 Release Notes

### New Features

- Initial release of `pg_scheduleserv`, that uses `vrprouting v0.2.0` and `vroom v1.10.0`.
- Scheduling of the tasks is done by dropping all the previous allocations.
