/*GRP-GNU-AGPL******************************************************************

File: job_tw_test.go

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

package e2etest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateJobTimeWindow(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobID      int
		body       map[string]interface{}
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Empty Body",
			statusCode: 400,
			jobID:      3324729385723589729,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'tw_open' of type 'string' is required",
					"Field 'tw_close' of type 'string' is required",
				},
			},
		},
		{
			name:       "Only tw_open",
			statusCode: 400,
			jobID:      3324729385723589729,
			body: map[string]interface{}{
				"tw_open": "2021-10-26T21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'tw_close' of type 'string' is required",
				},
			},
		},
		{
			name:       "Only tw_close",
			statusCode: 400,
			jobID:      3324729385723589729,
			body: map[string]interface{}{
				"tw_close": "2021-10-26T21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'tw_open' of type 'string' is required",
				},
			},
		},
		{
			name:       "Opening time greater than closing time",
			statusCode: 400,
			jobID:      3324729385723589729,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26T21:24:39",
				"tw_close": "2021-10-26T21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'tw_open' must be less than or equal to field 'tw_close'",
				},
			},
		},
		{
			name:       "Invalid JobID",
			statusCode: 400,
			jobID:      100,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26T21:24:38",
				"tw_close": "2021-10-26T21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Job with the given 'job_id' does not exist",
				},
			},
		},
		{
			name:       "All fields",
			statusCode: 201,
			jobID:      3324729385723589729,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26T21:20:20",
				"tw_close": "2021-10-26T21:24:38",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"tw_open":  "2021-10-26T21:20:20",
					"tw_close": "2021-10-26T21:24:38",
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Primary key violation",
			statusCode: 400,
			jobID:      3324729385723589729,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26T21:20:20",
				"tw_close": "2021-10-26T21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Jobs time window with given values already exist",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.todo == true {
				t.Skip("TODO")
			}
			m, b := tc.body, new(bytes.Buffer)
			if err := json.NewEncoder(b).Encode(m); err != nil {
				t.Error(err)
			}
			url := fmt.Sprintf("/jobs/%d/time_windows", tc.jobID)
			request, err := http.NewRequest("POST", url, b)
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)

			resp := recorder.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			m = map[string]interface{}{}
			if err = json.Unmarshal(body, &m); err != nil {
				t.Error(err)
			}
			if mData, ok := m["data"].(map[string]interface{}); ok {
				delete(mData, "id")
				delete(mData, "created_at")
				delete(mData, "updated_at")
				m["data"] = mData
			}
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestListJobTimeWindows(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobID      int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			jobID:      100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			jobID:      6362411701075685873,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":         "6362411701075685873",
						"tw_open":    "2020-10-10T00:00:00",
						"tw_close":   "2020-10-10T00:00:10",
						"created_at": "2021-10-26T21:25:41",
						"updated_at": "2021-10-26T21:25:41",
					},
					map[string]interface{}{
						"id":         "6362411701075685873",
						"tw_open":    "2020-10-11T00:00:00",
						"tw_close":   "2020-10-12T00:00:00",
						"created_at": "2021-10-26T21:25:51",
						"updated_at": "2021-10-26T21:25:51",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID but no time window",
			statusCode: 200,
			jobID:      3324729385723589729,
			resBody: map[string]interface{}{
				"data":    []interface{}{},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d/time_windows", tc.jobID)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)

			resp := recorder.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			m := map[string]interface{}{}
			if err = json.Unmarshal(body, &m); err != nil {
				t.Error(err)
			}
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestDeleteJobTimeWindow(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobID      int
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			jobID:      100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			jobID:      6362411701075685873,
			resBody: map[string]interface{}{
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Correct ID but no time window",
			statusCode: 404,
			jobID:      3324729385723589729,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.todo == true {
				t.Skip("TODO")
			}
			url := fmt.Sprintf("/jobs/%d/time_windows", tc.jobID)
			request, err := http.NewRequest("DELETE", url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)

			resp := recorder.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			m := map[string]interface{}{}
			if err = json.Unmarshal(body, &m); err != nil {
				t.Error(err)
			}
			assert.Equal(t, tc.resBody, m)
		})
	}
}
