/*GRP-GNU-AGPL******************************************************************

File: schedule.go

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

package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateSchedule godoc
// @Summary Schedule the tasks
// @Description Schedule the tasks present in a project
// @Tags Schedule
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} util.Schedule
// @Failure 400 {object} util.MultiError
// @Router /projects/{project_id}/schedule [post]
func (server *Server) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	// Add the project_id path variable
	vars := mux.Vars(r)
	projectID, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	err = server.DBCreateSchedule(ctx, projectID)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}
	// Fetch the schedule
	schedule, err := server.DBGetSchedule(ctx, projectID)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, schedule)
}

// GetSchedule godoc
// @Summary Get the schedule
// @Description Get the schedule for a project
// @Tags Schedule
// @Accept application/json
// @Produce text/calendar
// @Param project_id path int true "Project ID"
// @Success 200 {object} util.Schedule
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /projects/{project_id}/schedule [get]
func (server *Server) GetSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectID, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	schedule, err := server.DBGetSchedule(ctx, projectID)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatICAL(w, http.StatusOK, schedule)
}

// DeleteSchedule godoc
// @Summary Delete the schedule
// @Description Delete the schedule for a project
// @Tags Schedule
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} util.Success
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /projects/{project_id}/schedule [delete]
func (server *Server) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	err = server.DBDeleteSchedule(ctx, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, nil)
}
