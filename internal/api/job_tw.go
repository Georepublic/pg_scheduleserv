/*GRP-GNU-AGPL******************************************************************

File: job_tw.go

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

// CreateJobTimeWindows godoc
// @Summary Create a new job time window
// @Description Create a new job time window with the input payload
// @Tags Job
// @Accept application/json
// @Produce application/json
// @Param job_id path int true "Job ID"
// @Param JobTimeWindow body database.CreateJobTimeWindowParams true "Create job time window"
// @Success 200 {object} database.JobTimeWindow
// @Router /jobs/{job_id}/time_windows [post]
func (server *Server) createJobTimeWindow(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	logrus.Debug(r.Body)
	json.NewDecoder(r.Body).Decode(&userInput)

	// Add the job_id path variable
	vars := mux.Vars(r)
	userInput["id"] = vars["job_id"]

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateJobTimeWindowParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	job := database.CreateJobTimeWindowParams{}
	json.Unmarshal(userInputString, &job)

	logrus.Debugf("%v", userInput)
	logrus.Debugf("%+v", job)

	// Validate the struct
	if err := server.validate.Struct(job); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_job, err := server.DBCreateJobTimeWindow(ctx, job)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_job)
}
