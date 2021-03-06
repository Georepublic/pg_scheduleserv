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
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func getFinalData(respCode int, data interface{}) (int, interface{}) {
	if data == pgx.ErrNoRows {
		respCode = http.StatusNotFound
		data = nil
		return respCode, data
	}

	// Custom error message in case of conversion errors
	if typ, ok := data.(*strconv.NumError); ok {
		data = fmt.Errorf(`Invalid value "%s", %s`, typ.Num, typ.Err)
	}

	// Convert validation errors to []string
	if typ, ok := data.(validator.ValidationErrors); ok {
		data = getErrorMsg(typ)
	}

	// Handle multi errors and convert them to []string
	if typ, ok := data.(*multierror.Error); ok {
		errs := typ.WrappedErrors()
		msgs := make([]string, 0, len(errs))
		for _, err := range errs {
			msgs = append(msgs, err.Error())
		}
		data = msgs
	}

	// Handle single error
	if typ, ok := data.(error); ok {
		data = []string{typ.Error()}
	}

	return respCode, data
}

func (r *Formatter) FormatJSON(w http.ResponseWriter, respCode int, data interface{}) {
	respCode, data = getFinalData(respCode, data)

	// Set the content-type and response code in the header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respCode)

	if data == nil {
		if respCode >= 200 && respCode < 300 {
			fmt.Fprintf(w, jsonSuccessResp, http.StatusText(respCode), respCode)
		} else {
			fmt.Fprintf(w, jsonErrResp, http.StatusText(respCode), respCode)
		}
		return
	}

	b := r.pool.Get().(*bytes.Buffer)
	b.Reset()
	defer r.pool.Put(b)

	if respCode >= 200 && respCode < 300 {
		data = SuccessResponse{
			Data:    data,
			Message: http.StatusText(respCode),
			Code:    fmt.Sprintf("%d", respCode),
		}
	} else {
		data = ErrorResponse{
			Errors:  data,
			Message: http.StatusText(respCode),
			Code:    fmt.Sprintf("%d", respCode),
		}
	}

	if err := json.NewEncoder(b).Encode(data); err != nil {
		respCode = http.StatusInternalServerError
		w.WriteHeader(respCode)
		fmt.Fprintf(w, jsonErrResp, http.StatusText(respCode), respCode)
		return
	}

	_, err := b.WriteTo(w)
	if err != nil {
		logrus.Error(err)
	}
}

const jsonSuccessResp = `{"message": "%s", "code": "%d"}`

const jsonErrResp = `{"error": "%s", "code": "%d"}`

type NotFound struct {
	Error string `json:"error" example:"Not Found"`
	Code  string `json:"code" example:"404"`
}

type Success struct {
	Message string `json:"message" example:"OK"`
	Code    string `json:"code" example:"200"`
}

type ErrorResponse struct {
	Errors  interface{} `json:"errors" swaggertype:"array,string" example:"Error message1,Error message2"`
	Message string      `json:"message" example:"Bad Request"`
	Code    string      `json:"code" example:"400"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message" example:"OK"`
	Code    string      `json:"code" example:"200"`
}
