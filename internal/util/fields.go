/*GRP-GNU-AGPL******************************************************************

File: fields.go

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

package util

var IntervalFields = map[string]bool{
	"setup":        true,
	"service":      true,
	"p_setup":      true,
	"p_service":    true,
	"d_setup":      true,
	"d_service":    true,
	"travel_time":  true,
	"setup_time":   true,
	"service_time": true,
	"waiting_time": true,
	"max_shift":    true,
	"timeout":      true,
}

var TimestampFields = map[string]bool{
	"tw_open":    true,
	"tw_close":   true,
	"arrival":    true,
	"departure":  true,
	"created_at": true,
	"updated_at": true,
}

var AliasFields = map[string]string{
	"location":       "location_id",
	"p_location":     "p_location_id",
	"d_location":     "d_location_id",
	"start_location": "start_id",
	"end_location":   "end_id",
}
