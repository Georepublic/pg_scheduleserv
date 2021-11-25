/*GRP-GNU-AGPL******************************************************************

File: json.go

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

package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func (r *Formatter) FormatJSON(w http.ResponseWriter, respCode int, data interface{}) {
	if data == pgx.ErrNoRows {
		respCode = http.StatusNotFound
		data = nil
	}

	// Set the content-type and response code in the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respCode)

	if data == nil {
		if respCode >= 200 && respCode < 300 {
			fmt.Fprint(w, jsonSuccessResp)
			return
		}
		fmt.Fprintf(w, jsonErrResp, http.StatusText(respCode))
		return
	}

	// Convert validation errors to multi errors
	if typ, ok := data.(validator.ValidationErrors); ok {
		data = &MultiError{Errors: getErrorMsg(typ)}
	}

	// Handle multi errors
	if typ, ok := data.(*multierror.Error); ok {

		errs := typ.WrappedErrors()
		msgs := make([]string, 0, len(errs))
		for _, err := range errs {
			msgs = append(msgs, err.Error())
		}
		data = &MultiError{Errors: msgs}
	}

	// Handle single error
	if typ, ok := data.(error); ok {
		data = &MultiError{Errors: []string{typ.Error()}}
	}

	b := r.pool.Get().(*bytes.Buffer)
	b.Reset()
	defer r.pool.Put(b)

	if err := json.NewEncoder(b).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonErrResp, http.StatusText(http.StatusInternalServerError))
		return
	}

	_, err := b.WriteTo(w)
	if err != nil {
		logrus.Error(err)
	}
}

const jsonSuccessResp = `{"success": true}`

const jsonErrResp = `{"error": "%s"}`

type MultiError struct {
	Errors []string `json:"errors,omitempty" example:"Error message1,Error message2"`
}

type NotFound struct {
	Error string `json:"error" example:"Not Found"`
}

type Success struct {
	Success string `json:"success" example:"true"`
}
