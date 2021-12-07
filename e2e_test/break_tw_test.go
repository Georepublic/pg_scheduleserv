/*GRP-GNU-AGPL******************************************************************

File: break_tw_test.go

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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBreakTimeWindow(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		breakID    int
		body       map[string]interface{}
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Empty Body",
			statusCode: 400,
			breakID:    4668767710686035977,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'tw_open' of type 'string' is required",
					"Field 'tw_close' of type 'string' is required",
				},
			},
		},
		{
			name:       "Only tw_open",
			statusCode: 400,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"tw_open": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'tw_close' of type 'string' is required",
				},
			},
		},
		{
			name:       "Only tw_close",
			statusCode: 400,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'tw_open' of type 'string' is required",
				},
			},
		},
		{
			name:       "Opening time greater than closing time",
			statusCode: 400,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26 21:24:39",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'tw_open' must be less than or equal to field 'tw_close'",
				},
			},
		},
		{
			name:       "Invalid BreakID",
			statusCode: 400,
			breakID:    100,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26 21:24:38",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Break with the given 'break_id' does not exist",
				},
			},
		},
		{
			name:       "All fields",
			statusCode: 201,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
		},
		{
			name:       "Primary key violation",
			statusCode: 400,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Breaks time window with given values already exist",
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
			url := fmt.Sprintf("/breaks/%d/time_windows", tc.breakID)
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
			delete(m, "id")
			delete(m, "created_at")
			delete(m, "updated_at")
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestListBreakTimeWindows(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		breakID    int
		resBody    []map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 200,
			breakID:    100,
			resBody:    []map[string]interface{}{},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			breakID:    3990300682121424906,
			resBody: []map[string]interface{}{
				{
					"id":         "3990300682121424906",
					"tw_open":    "2020-10-10 00:00:00",
					"tw_close":   "2020-10-10 00:00:10",
					"created_at": "2021-10-26 21:25:41",
					"updated_at": "2021-10-26 21:25:41",
				},
				{
					"id":         "3990300682121424906",
					"tw_open":    "2020-10-11 00:00:00",
					"tw_close":   "2020-10-12 00:00:00",
					"created_at": "2021-10-26 21:25:51",
					"updated_at": "2021-10-26 21:25:51",
				},
			},
		},
		{
			name:       "Valid ID but no time window",
			statusCode: 200,
			breakID:    4668767710686035977,
			resBody:    []map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/breaks/%d/time_windows", tc.breakID)
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
			m := []map[string]interface{}{}
			if err = json.Unmarshal(body, &m); err != nil {
				t.Error(err)
			}
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestDeleteBreakTimeWindow(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		breakID    int
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			breakID:    100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			breakID:    3990300682121424906,
			resBody: map[string]interface{}{
				"success": true,
			},
		},
		{
			name:       "Correct ID but no time window",
			statusCode: 404,
			breakID:    4668767710686035977,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.todo == true {
				t.Skip("TODO")
			}
			url := fmt.Sprintf("/breaks/%d/time_windows", tc.breakID)
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
