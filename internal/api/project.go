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

	"github.com/Georepublic/pg_scheduleserv/internal/database"
	"github.com/jackc/pgtype"
)

type CreateProjectRequest struct {
	Name string       `json:"name" validate:"required"`
	Data pgtype.JSONB `json:"data"`
}

func (server *Server) createProject(w http.ResponseWriter, r *http.Request) {
	project := CreateProjectRequest{}
	json.NewDecoder(r.Body).Decode(&project)
	if err := server.validate.Struct(project); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}
	arg := database.CreateProjectParams{
		Name: project.Name,
		Data: project.Data,
	}

	ctx := r.Context()
	created_project, err := server.CreateProject(ctx, arg)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_project)
}
