# pg_scheduleserv Release Notes

## 0.2.0 Release Notes

### New Features

* #26: Update to `vrprouting v0.3.0` and `vroom v1.11.0`.

### Fixes

* Change all fields with INTERVAL to HH:MM:SS format.

## 0.1.0 Release Notes

### New Features

* Initial release of `pg_scheduleserv`, that uses `vrprouting v0.2.0` and `vroom v1.10.0`.
* Scheduling of the tasks is done by dropping all the previous allocations.
