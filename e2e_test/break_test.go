/*GRP-GNU-AGPL******************************************************************

File: break_test.go

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

func TestCreateBreak(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		vehicleID  int
		body       map[string]interface{}
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Empty Body",
			statusCode: 201,
			vehicleID:  2550908592071787332,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"service":    "00:00:00",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{},
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Only service",
			statusCode: 201,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"service": "00:01:40",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"service":    "00:01:40",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{},
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Negative service",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"service": "-00:01:40",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'service' must be non-negative with the format 'HH:MM:SS'"},
			},
		},
		{
			name:       "Only data",
			statusCode: 201,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"service":    "00:00:00",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{"key": "value"},
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Invalid Vehicle ID",
			statusCode: 400,
			vehicleID:  123,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Vehicle with the given 'vehicle_id' does not exist",
				},
			},
		},
		{
			name:       "All fields",
			statusCode: 201,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"service": "00:03:35",
				"data":    map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"service":    "00:03:35",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{"key": "value"},
				},
				"code":    "201",
				"message": "Created",
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
			url := fmt.Sprintf("/vehicles/%d/breaks", tc.vehicleID)
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

func TestListBreaks(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		vehicleID  int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			vehicleID:  100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":         "4668767710686035977",
						"service":    "00:00:01",
						"vehicle_id": "2550908592071787332",
						"data":       map[string]interface{}{"key": "value"},
						"created_at": "2021-10-26 21:24:38",
						"updated_at": "2021-10-26 21:24:38",
					},
					map[string]interface{}{
						"id":         "3990300682121424906",
						"service":    "00:05:24",
						"vehicle_id": "2550908592071787332",
						"data":       map[string]interface{}{"s": float64(1)},
						"created_at": "2021-10-26 21:24:52",
						"updated_at": "2021-10-26 21:24:52",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID but no break",
			statusCode: 200,
			vehicleID:  150202809001685363,
			resBody: map[string]interface{}{
				"data":    []interface{}{},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/vehicles/%d/breaks", tc.vehicleID)
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

func TestGetBreak(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		breakID    int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			breakID:    100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			breakID:    4668767710686035977,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":         "4668767710686035977",
					"service":    "00:00:01",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{"key": "value"},
					"created_at": "2021-10-26 21:24:38",
					"updated_at": "2021-10-26 21:24:38",
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/breaks/%d", tc.breakID)
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

func TestUpdateBreak(t *testing.T) {
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
			statusCode: 200,
			breakID:    4668767710686035977,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":         "4668767710686035977",
					"service":    "00:00:01",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{"key": "value"},
					"created_at": "2021-10-26 21:24:38",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			breakID:    100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Only service",
			statusCode: 200,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"service": "00:01:40",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":         "4668767710686035977",
					"service":    "00:01:40",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{"key": "value"},
					"created_at": "2021-10-26 21:24:38",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Negative service",
			statusCode: 400,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"service": "-00:01:40",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'service' must be non-negative with the format 'HH:MM:SS'"},
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"data": map[string]interface{}{},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":         "4668767710686035977",
					"service":    "00:01:40",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{},
					"created_at": "2021-10-26 21:24:38",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "All fields",
			statusCode: 200,
			breakID:    4668767710686035977,
			body: map[string]interface{}{
				"service":    "00:01:41",
				"vehicle_id": "2550908592071787332",
				"data":       map[string]interface{}{"s": 1},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":         "4668767710686035977",
					"service":    "00:01:41",
					"vehicle_id": "2550908592071787332",
					"data":       map[string]interface{}{"s": float64(1)},
					"created_at": "2021-10-26 21:24:38",
				},
				"code":    "200",
				"message": "OK",
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
			url := fmt.Sprintf("/breaks/%d", tc.breakID)
			request, err := http.NewRequest("PATCH", url, b)
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
				delete(mData, "updated_at")
				m["data"] = mData
			}
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestDeleteBreak(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		breakID    int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			breakID:    100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			breakID:    4668767710686035977,
			resBody: map[string]interface{}{
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/breaks/%d", tc.breakID)
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
