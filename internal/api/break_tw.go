/*GRP-GNU-AGPL******************************************************************

File: break_tw.go

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
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Georepublic/pg_scheduleserv/internal/database"
	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// CreateBreakTimeWindow godoc
// @Summary Create a new break time window
// @Description Create a new break time window with the input payload
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Param BreakTimeWindow body database.CreateBreakTimeWindowParams true "Create break time window"
// @Success 200 {object} database.BreakTimeWindow
// @Failure 400 {object} util.MultiError
// @Router /breaks/{break_id}/time_windows [post]
func (server *Server) CreateBreakTimeWindow(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	// Add the break_id path variable
	vars := mux.Vars(r)
	userInput["id"] = vars["break_id"]

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateBreakTimeWindowParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	v_break := database.CreateBreakTimeWindowParams{}
	if err = json.Unmarshal(userInputString, &v_break); err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("%v", userInput)
	logrus.Debugf("%+v", v_break)

	// Validate the struct
	if err := server.validate.Struct(v_break); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_break, err := server.DBCreateBreakTimeWindow(ctx, v_break)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_break)
}

// ListBreakTimeWindows godoc
// @Summary List break time windows for a break
// @Description Get a list of break time windows for a break with break_id
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Success 200 {object} database.BreakTimeWindow
// @Failure 400 {object} util.MultiError
// @Router /breaks/{break_id}/time_windows [get]
func (server *Server) ListBreakTimeWindows(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	break_id, err := strconv.ParseInt(vars["break_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_vehicle, err := server.DBListBreakTimeWindows(ctx, break_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_vehicle)
}

// DeleteBreakTimeWindows godoc
// @Summary Delete break time windows
// @Description Delete break time windows for a break with break_id
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Param tw_open path string true "Break opening Time Window"
// @Param tw_close path string true "Break closing Time Window"
// @Success 200 {object} util.Success
// @Failure 400 {object} util.MultiError
// @Router /breaks/{break_id}/time_windows/{tw_open}/{tw_close} [delete]
func (server *Server) DeleteBreakTimeWindow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	break_id, err := strconv.ParseInt(vars["break_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	userInput := map[string]interface{}{
		"id":       break_id,
		"tw_open":  vars["tw_open"],
		"tw_close": vars["tw_close"],
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateBreakTimeWindowParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	break_tw := database.CreateBreakTimeWindowParams{}
	if err = json.Unmarshal(userInputString, &break_tw); err != nil {
		logrus.Error(err)
	}

	// Validate the struct
	if err := server.validate.Struct(break_tw); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	_, err = server.DBDeleteBreakTimeWindow(ctx, break_tw)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, nil)
}
