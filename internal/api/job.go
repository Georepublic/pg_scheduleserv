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
	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/sirupsen/logrus"
)

type CreateJobRequest struct {
	Latitude  float64      `json:"latitude"   validate:"required"`
	Longitude float64      `json:"longitude"  validate:"required"`
	Service   int64        `json:"service"`
	Delivery  []int64      `json:"delivery"`
	Pickup    []int64      `json:"pickup"`
	Skills    []int32      `json:"skills"`
	Priority  int32        `json:"priority"`
	ProjectID int64        `json:"project_id"`
	Data      pgtype.JSONB `json:"data"`
}

func (server *Server) createJob(w http.ResponseWriter, r *http.Request) {
	job := CreateJobRequest{}
	json.NewDecoder(r.Body).Decode(&job)
	if err := server.validate.Struct(job); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		logrus.Error(err)
	}
	arg := database.CreateJobParams{
		Latitude:  job.Latitude,
		Longitude: job.Longitude,
		Service:   job.Service,
		Delivery:  job.Delivery,
		Pickup:    job.Pickup,
		Skills:    job.Skills,
		Priority:  job.Priority,
		ProjectID: project_id,
		Data:      job.Data,
	}

	ctx := r.Context()
	created_job, err := server.CreateJob(ctx, arg)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_job)
}
