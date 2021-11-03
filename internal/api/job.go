/*GRP-GNU-AGPL******************************************************************

File: job.go

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

// CreateJob godoc
// @Summary Create a new job
// @Description Create a new job with the input payload
// @Tags Job
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Param Job body database.UpdateJobParams true "Update job"
// @Success 200 {object} database.Job
// @Router /projects/{project_id}/jobs [post]
func (server *Server) CreateJob(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&userInput)
	}

	// Add the project_id path variable
	vars := mux.Vars(r)
	userInput["project_id"] = vars["project_id"]

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateJobParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	job := database.CreateJobParams{}
	json.Unmarshal(userInputString, &job)

	// Validate the struct
	if err := server.validate.Struct(job); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_job, err := server.DBCreateJob(ctx, job)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_job)
}

// ListJobs godoc
// @Summary List jobs for a project
// @Description Get a list of jobs for a project with project_id
// @Tags Job
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} database.Job
// @Router /projects/{project_id}/jobs [get]
func (server *Server) ListJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_job, err := server.DBListJobs(ctx, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_job)
}

// GetJob godoc
// @Summary Fetch a job
// @Description Fetch a job with its job_id
// @Tags Job
// @Accept application/json
// @Produce application/json
// @Param job_id path int true "Job ID"
// @Success 200 {object} database.Job
// @Router /jobs/{job_id} [get]
func (server *Server) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job_id, err := strconv.ParseInt(vars["job_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_job, err := server.DBGetJob(ctx, job_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_job)
}

// GetJob godoc
// @Summary Update a job
// @Description Update a job (partial update) with its job_id
// @Tags Job
// @Accept application/json
// @Produce application/json
// @Param job_id path int true "Job ID"
// @Success 200 {object} database.Job
// @Router /jobs/{job_id} [patch]
func (server *Server) UpdateJob(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&userInput)
	}

	vars := mux.Vars(r)
	job_id, err := strconv.ParseInt(vars["job_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateJobParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	job := database.UpdateJobParams{}
	json.Unmarshal(userInputString, &job)

	// Validate the struct
	if err := server.validate.Struct(job); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_job, err := server.DBUpdateJob(ctx, job, job_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_job)
}

// GetJob godoc
// @Summary Delete a job
// @Description Delete a job with its job_id
// @Tags Job
// @Accept application/json
// @Produce application/json
// @Param job_id path int true "Job ID"
// @Success 200 {object} database.Job
// @Router /jobs/{job_id} [delete]
func (server *Server) DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	job_id, err := strconv.ParseInt(vars["job_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	err = server.DBDeleteJob(ctx, job_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusNoContent, nil)
}
