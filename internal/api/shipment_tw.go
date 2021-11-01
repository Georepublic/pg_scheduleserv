/*GRP-GNU-AGPL******************************************************************

File: shipment_tw.go

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

// CreateShipmentTimeWindow godoc
// @Summary Create a new shipment time window
// @Description Create a new shipment time window with the input payload
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Param ShipmentTimeWindow body database.CreateShipmentTimeWindowParams true "Create shipment time window"
// @Success 200 {object} database.ShipmentTimeWindow
// @Router /shipments/{shipment_id}/time_windows [post]
func (server *Server) createShipmentTimeWindow(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	logrus.Debug(r.Body)
	json.NewDecoder(r.Body).Decode(&userInput)

	// Add the shipment_id path variable
	vars := mux.Vars(r)
	userInput["id"] = vars["shipment_id"]

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateShipmentTimeWindowParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	shipment := database.CreateShipmentTimeWindowParams{}
	json.Unmarshal(userInputString, &shipment)

	logrus.Debugf("%v", userInput)
	logrus.Debugf("%+v", shipment)

	// Validate the struct
	if err := server.validate.Struct(shipment); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_shipment, err := server.DBCreateShipmentTimeWindow(ctx, shipment)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_shipment)
}

// ListShipmentTimeWindows godoc
// @Summary List shipment time windows for a shipment
// @Description Get a list of shipment time windows for a shipment with shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Success 200 {object} database.ShipmentTimeWindow
// @Router /shipments/{shipment_id}/time_windows [get]
func (server *Server) listShipmentTimeWindows(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipment_id, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_vehicle, err := server.DBListShipmentTimeWindows(ctx, shipment_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_vehicle)
}

// DeleteShipmentTimeWindows godoc
// @Summary Delete shipment time windows
// @Description Delete shipment time windows for a shipment with shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Param tw_open path string true "Shipment opening Time Window"
// @Param tw_close path string true "Shipment closing Time Window"
// @Success 200 {object} database.ShipmentTimeWindow
// @Router /shipments/{shipment_id}/time_windows/{tw_open}/{tw_close} [delete]
func (server *Server) deleteShipmentTimeWindow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipment_id, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	userInput := map[string]interface{}{
		"id":       shipment_id,
		"tw_open":  vars["tw_open"],
		"tw_close": vars["tw_close"],
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateShipmentTimeWindowParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	shipment_tw := database.CreateShipmentTimeWindowParams{}
	json.Unmarshal(userInputString, &shipment_tw)

	// Validate the struct
	if err := server.validate.Struct(shipment_tw); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	err = server.DBDeleteShipmentTimeWindow(ctx, shipment_tw)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, nil)
}
