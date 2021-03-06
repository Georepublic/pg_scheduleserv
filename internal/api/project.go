/*GRP-GNU-AGPL******************************************************************

File: project.go

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

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project with the input payload
// @Description The "duration_calc" parameter must be either "euclidean", "valhalla" or "osrm"
// @Tags Project
// @Accept application/json
// @Produce application/json
// @Param Project body database.CreateProjectParams true "Create project"
// @Success 200 {object} util.SuccessResponse{data=database.Project}
// @Failure 400 {object} util.ErrorResponse
// @Router /projects [post]
func (server *Server) CreateProject(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateProjectParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	project := database.CreateProjectParams{}
	if err = json.Unmarshal(userInputString, &project); err != nil {
		logrus.Error(err)
	}

	// Validate the struct
	if err := server.validate.Struct(project); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_project, err := server.DBCreateProject(ctx, project)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_project)
}

// ListProjects godoc
// @Summary List projects
// @Description Get a list of projects
// @Tags Project
// @Accept application/json
// @Produce application/json
// @Success 200 {object} util.SuccessResponse{data=[]database.Project}
// @Failure 400 {object} util.ErrorResponse
// @Router /projects [get]
func (server *Server) ListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	created_project, err := server.DBListProjects(ctx)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_project)
}

// GetProject godoc
// @Summary Fetch a project
// @Description Fetch a project with its project_id
// @Tags Project
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} util.SuccessResponse{data=database.Project}
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.NotFound
// @Router /projects/{project_id} [get]
func (server *Server) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_project, err := server.DBGetProject(ctx, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_project)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update a project with its project_id
// @Description The "duration_calc" parameter must be either "euclidean", "valhalla" or "osrm"
// @Tags Project
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Param Project body database.CreateProjectParams true "Update project"
// @Success 200 {object} util.SuccessResponse{data=database.Project}
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.NotFound
// @Router /projects/{project_id} [patch]
func (server *Server) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.UpdateProjectParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	project := database.UpdateProjectParams{}
	if err = json.Unmarshal(userInputString, &project); err != nil {
		logrus.Error(err)
	}

	// Validate the struct
	if err := server.validate.Struct(project); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_project, err := server.DBUpdateProject(ctx, project, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_project)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project with its project_id
// @Tags Project
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} util.Success
// @Failure 400 {object} util.ErrorResponse
// @Failure 404 {object} util.NotFound
// @Router /projects/{project_id} [delete]
func (server *Server) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	_, err = server.DBDeleteProject(ctx, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, nil)
}
