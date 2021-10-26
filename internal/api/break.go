/*GRP-GNU-AGPL******************************************************************

File: break.go

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

// CreateBreak godoc
// @Summary Create a new break
// @Description Create a new break with the input payload
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param vehicle_id path int true "Vehicle ID"
// @Param Break body database.CreateBreakParams true "Create break"
// @Success 200 {object} database.Break
// @Router /vehicles/{vehicle_id}/breaks [post]
func (server *Server) createBreak(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&userInput)

	// Add the vehicle_id path variable
	vars := mux.Vars(r)
	userInput["vehicle_id"] = vars["vehicle_id"]

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateBreakParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	v_break := database.CreateBreakParams{}
	json.Unmarshal(userInputString, &v_break)

	// Validate the struct
	if err := server.validate.Struct(v_break); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_break, err := server.DBCreateBreak(ctx, v_break)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_break)
}
