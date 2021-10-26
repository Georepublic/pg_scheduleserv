/*GRP-GNU-AGPL******************************************************************

File: vehicle.go

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

// CreateVehicles godoc
// @Summary Create a new vehicle
// @Description Create a new vehicle with the input payload
// @Tags Vehicle
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Param Vehicle body database.CreateVehicleParams true "Create vehicle"
// @Success 200 {object} database.Vehicle
// @Router /projects/{project_id}/vehicles [post]
func (server *Server) createVehicle(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&userInput)

	// Add the project_id path variable
	vars := mux.Vars(r)
	userInput["project_id"] = vars["project_id"]

	logrus.Debugf("%v", userInput)

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateVehicleParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	vehicle := database.CreateVehicleParams{}
	json.Unmarshal(userInputString, &vehicle)

	logrus.Debugf("%+v", vehicle)

	// Validate the struct
	if err := server.validate.Struct(vehicle); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_vehicle, err := server.DBCreateVehicle(ctx, vehicle)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_vehicle)
}
