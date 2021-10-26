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

	"github.com/Georepublic/pg_scheduleserv/internal/database"
	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// CreateBreakTimeWindows godoc
// @Summary Create a new break time window
// @Description Create a new break time window with the input payload
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Param BreakTimeWindow body database.CreateBreakTimeWindowParams true "Create break time window"
// @Success 200 {object} database.BreakTimeWindow
// @Router /breaks/{break_id}/time_windows [post]
func (server *Server) createBreakTimeWindow(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	logrus.Debug(r.Body)
	json.NewDecoder(r.Body).Decode(&userInput)

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
	json.Unmarshal(userInputString, &v_break)

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
