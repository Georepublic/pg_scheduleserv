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
	"strconv"

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
// @Success 200 {object} util.SuccessResponse{data=database.Break}
// @Failure 400 {object} util.ErrorResponse
// @Router /vehicles/{vehicle_id}/breaks [post]
func (server *Server) CreateBreak(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

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
	if err = json.Unmarshal(userInputString, &v_break); err != nil {
		logrus.Error(err)
	}

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

// ListBreaks godoc
// @Summary List breaks
// @Description Get a list of breaks
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param vehicle_id path int true "Vehicle ID"
// @Success 200 {object} util.SuccessResponse{data=[]database.Break}
// @Failure 400 {object} util.ErrorResponse
// @Router /vehicles/{vehicle_id}/breaks [get]
func (server *Server) ListBreaks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicle_id, err := strconv.ParseInt(vars["vehicle_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_break, err := server.DBListBreaks(ctx, vehicle_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_break)
}

// GetBreak godoc
// @Summary Fetch a break
// @Description Fetch a break with its break_id
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Success 200 {object} util.SuccessResponse{data=database.Break}
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.NotFound
// @Router /breaks/{break_id} [get]
func (server *Server) GetBreak(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	break_id, err := strconv.ParseInt(vars["break_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_break, err := server.DBGetBreak(ctx, break_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_break)
}

// UpdateBreak godoc
// @Summary Update a break
// @Description Update a break with its break_id
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Param Break body database.CreateBreakParams true "Update break"
// @Success 200 {object} util.SuccessResponse{data=database.Break}
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.NotFound
// @Router /breaks/{break_id} [patch]
func (server *Server) UpdateBreak(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	vars := mux.Vars(r)
	break_id, err := strconv.ParseInt(vars["break_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.UpdateBreakParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	v_break := database.UpdateBreakParams{}
	if err = json.Unmarshal(userInputString, &v_break); err != nil {
		logrus.Error(err)
	}

	// Validate the struct
	if err := server.validate.Struct(v_break); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_break, err := server.DBUpdateBreak(ctx, v_break, break_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_break)
}

// DeleteBreak godoc
// @Summary Delete a break
// @Description Delete a break with its break_id
// @Tags Break
// @Accept application/json
// @Produce application/json
// @Param break_id path int true "Break ID"
// @Success 200 {object} util.Success
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.NotFound
// @Router /breaks/{break_id} [delete]
func (server *Server) DeleteBreak(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	break_id, err := strconv.ParseInt(vars["break_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	_, err = server.DBDeleteBreak(ctx, break_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, nil)
}
