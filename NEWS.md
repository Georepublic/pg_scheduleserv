# pg_scheduleserv Release Notes

## 0.2.0 Release Notes

### New Features

* #26: Update to `vrprouting v0.3.0` and `vroom v1.11.0`.
* #21: Extend the schedule output to provide metadata and the complete schedule, using optional "overview" query parameter.
  * Also extended all API responses to provide message and status code.

### Fixes

* Change all fields with INTERVAL to HH:MM:SS format.
* Returning 404 NotFound Error instead of empty list when the project/task id is wrong.

## 0.1.0 Release Notes

### New Features

* Initial release of `pg_scheduleserv`, that uses `vrprouting v0.2.0` and `vroom v1.10.0`.
* Scheduling of the tasks is done by dropping all the previous allocations.
