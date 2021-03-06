/*GRP-GNU-AGPL******************************************************************

File: vehicle_test.go

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
	"regexp"
	"testing"
	"time"

	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateVehicle(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectID  int
		body       map[string]interface{}
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Empty Body",
			statusCode: 400,
			projectID:  3909655254191459782,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'start_location' of type 'util.LocationParams' is required",
					"Field 'end_location' of type 'util.LocationParams' is required",
				},
			},
		},
		{
			name:       "Only Location - Wrong parameters 1",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": "Sample Location",
				"end_location":   "Sample Location",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'start_location' must be of 'util.LocationParams' type.",
					"Field 'end_location' must be of 'util.LocationParams' type.",
				},
			},
		},
		{
			name:       "Only Location - Wrong parameters 2",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  "12.34567",
					"longitude": "56.78",
				},
				"end_location": map[string]interface{}{
					"latitude":  "12.34567",
					"longitude": "56.78",
				},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'latitude' and 'longitude' of type 'float64' is required",
				},
			},
		},
		{
			name:       "Only Location - Wrong range of parameters - 1",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'latitude' must be less than or equal to 90",
					"Field 'longitude' must be less than or equal to 180",
				},
			},
		},
		{
			name:       "Only Location - Wrong range of parameters - 2",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  -112.34567,
					"longitude": -256.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -112.34567,
					"longitude": -256.78,
				},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'latitude' must be greater than or equal to -90",
					"Field 'longitude' must be greater than or equal to -180",
				},
			},
		},
		{
			name:       "Only Location",
			statusCode: 201,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"start_location": map[string]interface{}{
						"latitude":  12.3457,
						"longitude": 56.78,
					},
					"end_location": map[string]interface{}{
						"latitude":  -12.3457,
						"longitude": -56.78,
					},
					"capacity":     []interface{}{},
					"skills":       []interface{}{},
					"tw_open":      "1970-01-01T00:00:00",
					"tw_close":     "2038-01-19T03:14:07",
					"speed_factor": float64(1),
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{},
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Only data",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'start_location' of type 'util.LocationParams' is required",
					"Field 'end_location' of type 'util.LocationParams' is required",
				},
			},
		},
		{
			name:       "Negative skills",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"skills": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'skills[0]' must be non-negative",
					"Field 'skills[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Negative capacity",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"capacity": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'capacity[0]' must be non-negative",
					"Field 'capacity[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Zero speed factor",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"speed_factor": 0,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'speed_factor' must be greater than 0",
				},
			},
		},
		{
			name:       "Negative speed factor",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"speed_factor": -1.1,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'speed_factor' must be greater than 0",
				},
			},
		},
		{
			name:       "All fields",
			statusCode: 201,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"capacity":     []interface{}{15, 16},
				"skills":       []interface{}{5, 50, 100},
				"tw_open":      "2021-01-01T01:01:01",
				"tw_close":     "2021-01-09T03:14:07",
				"speed_factor": 10.45,
				"max_tasks":    float64(25),
				"data":         map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"start_location": map[string]interface{}{
						"latitude":  12.3457,
						"longitude": 56.78,
					},
					"end_location": map[string]interface{}{
						"latitude":  -12.3457,
						"longitude": -56.78,
					},
					"capacity":     []interface{}{float64(15), float64(16)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"tw_open":      "2021-01-01T01:01:01",
					"tw_close":     "2021-01-09T03:14:07",
					"speed_factor": 10.45,
					"max_tasks":    float64(25),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
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
			url := fmt.Sprintf("/projects/%d/vehicles", tc.projectID)
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

func TestListVehicles(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectID  int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			projectID:  100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			projectID:  3909655254191459782,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id": "2550908592071787332",
						"start_location": map[string]interface{}{
							"latitude":  32.234,
							"longitude": -23.2342,
						},
						"end_location": map[string]interface{}{
							"latitude":  23.3458,
							"longitude": 2.3242,
						},
						"capacity":     []interface{}{float64(10), float64(30)},
						"skills":       []interface{}{float64(10)},
						"tw_open":      "2020-01-01T00:00:00",
						"tw_close":     "2020-01-10T07:14:07",
						"speed_factor": 10.5,
						"max_tasks":    float64(2147483647),
						"project_id":   "3909655254191459782",
						"data":         map[string]interface{}{"key": "value"},
						"created_at":   "2021-10-26T10:46:41",
						"updated_at":   "2021-10-26T10:46:41",
					},
					map[string]interface{}{
						"id": "7300272137290532980",
						"start_location": map[string]interface{}{
							"latitude":  -32.234,
							"longitude": -23.2342,
						},
						"end_location": map[string]interface{}{
							"latitude":  23.3458,
							"longitude": 2.3242,
						},
						"capacity":     []interface{}{float64(30), float64(50)},
						"skills":       []interface{}{float64(1)},
						"tw_open":      "2020-01-01T10:10:00",
						"tw_close":     "2020-01-11T03:14:07",
						"speed_factor": 34.25,
						"max_tasks":    float64(2147483647),
						"project_id":   "3909655254191459782",
						"data":         map[string]interface{}{"s": float64(1)},
						"created_at":   "2021-10-26T10:47:54",
						"updated_at":   "2021-10-26T10:47:54",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d/vehicles", tc.projectID)
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

func TestGetVehicle(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
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
			name:       "Correct ID",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  32.234,
						"longitude": -23.2342,
					},
					"end_location": map[string]interface{}{
						"latitude":  23.3458,
						"longitude": 2.3242,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(10)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 10.5,
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
					"updated_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/vehicles/%d", tc.vehicleID)
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

func TestUpdateVehicle(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
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
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  32.234,
						"longitude": -23.2342,
					},
					"end_location": map[string]interface{}{
						"latitude":  23.3458,
						"longitude": 2.3242,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(10)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 10.5,
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			vehicleID:  100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Only Location - Wrong parameters 1",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"start_location": "Sample Location",
				"end_location":   "Sample Location",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'start_location' must be of 'util.LocationParams' type.",
					"Field 'end_location' must be of 'util.LocationParams' type.",
				},
			},
		},
		{
			name:       "Only Location - Wrong parameters 2",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  "12.34567",
					"longitude": "56.78",
				},
				"end_location": map[string]interface{}{
					"latitude":  "12.34567",
					"longitude": "56.78",
				},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'latitude' and 'longitude' of type 'float64' is required",
				},
			},
		},
		{
			name:       "Only Location - Wrong range of parameters - 1",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'latitude' must be less than or equal to 90",
					"Field 'longitude' must be less than or equal to 180",
				},
			},
		},
		{
			name:       "Only Location - Wrong range of parameters - 2",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  -112.34567,
					"longitude": -256.78,
				},
				"end_location": map[string]interface{}{
					"latitude":  -112.34567,
					"longitude": -256.78,
				},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'latitude' must be greater than or equal to -90",
					"Field 'longitude' must be greater than or equal to -180",
				},
			},
		},
		{
			name:       "Negative skills",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"skills": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'skills[0]' must be non-negative",
					"Field 'skills[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Negative capacity",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"capacity": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'capacity[0]' must be non-negative",
					"Field 'capacity[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Zero speed factor",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"speed_factor": 0,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'speed_factor' must be greater than 0",
				},
			},
		},
		{
			name:       "Negative speed factor",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"speed_factor": -1.1,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'speed_factor' must be greater than 0",
				},
			},
		},
		{
			name:       "Only location",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"end_location": map[string]interface{}{
					"latitude":  -23.4567,
					"longitude": 78.90,
				},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(10)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 10.5,
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only capacity",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"capacity": []interface{}{10, 30},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(10)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 10.5,
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only skills",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"skills": []interface{}{5},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(5)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 10.5,
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only speed_factor",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"speed_factor": 1.234,
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(5)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 1.234,
					"max_tasks":    float64(2147483647),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only max_tasks",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"max_tasks": float64(15),
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(5)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 1.234,
					"max_tasks":    float64(15),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"data": map[string]interface{}{},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(5)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 1.234,
					"max_tasks":    float64(15),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Invalid projectID type",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"project_id": 100,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'project_id' must be of 'string' type."}},
		},
		{
			name:       "Invalid projectID",
			statusCode: 400,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"project_id": "100",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Project with the given 'project_id' does not exist"}},
		},
		{
			name:       "Valid projectID",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"project_id": "8943284028902589305",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"capacity":     []interface{}{float64(10), float64(30)},
					"skills":       []interface{}{float64(5)},
					"tw_open":      "2020-01-01T00:00:00",
					"tw_close":     "2020-01-10T07:14:07",
					"speed_factor": 1.234,
					"max_tasks":    float64(15),
					"project_id":   "8943284028902589305",
					"data":         map[string]interface{}{},
					"created_at":   "2021-10-26T10:46:41",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "All fields",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			body: map[string]interface{}{
				"start_location": map[string]interface{}{
					"latitude":  3.4567,
					"longitude": -8.90,
				},
				"end_location": map[string]interface{}{
					"latitude":  -3.4567,
					"longitude": 8.90,
				},
				"capacity":     []interface{}{float64(21)},
				"skills":       []interface{}{float64(5), float64(6)},
				"tw_open":      "2021-11-01T00:00:00",
				"tw_close":     "2021-11-10T03:14:07",
				"speed_factor": 11.234,
				"max_tasks":    float64(35),
				"project_id":   "3909655254191459782",
				"data":         map[string]interface{}{"s": 1},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "2550908592071787332",
					"start_location": map[string]interface{}{
						"latitude":  3.4567,
						"longitude": -8.90,
					},
					"end_location": map[string]interface{}{
						"latitude":  -3.4567,
						"longitude": 8.90,
					},
					"capacity":     []interface{}{float64(21)},
					"skills":       []interface{}{float64(5), float64(6)},
					"tw_open":      "2021-11-01T00:00:00",
					"tw_close":     "2021-11-10T03:14:07",
					"speed_factor": 11.234,
					"max_tasks":    float64(35),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"s": float64(1)},
					"created_at":   "2021-10-26T10:46:41",
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
			url := fmt.Sprintf("/vehicles/%d", tc.vehicleID)
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

func TestDeleteVehicle(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
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
			name:       "Correct ID",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			resBody: map[string]interface{}{
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/vehicles/%d", tc.vehicleID)
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

func TestGetVehicleScheduleJson(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
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
			vehicleID:  123,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Valid ID, no schedule",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"schedule": []interface{}{},
					"metadata": map[string]interface{}{
						"summary":       []interface{}{},
						"unassigned":    []interface{}{},
						"total_setup":   "00:00:00",
						"total_service": "00:00:00",
						"total_travel":  "00:00:00",
						"total_waiting": "00:00:00",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			vehicleID:  7300272137290532980,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"schedule": []interface{}{
						map[string]interface{}{
							"vehicle_id": "7300272137290532980",
							"vehicle_data": map[string]interface{}{
								"s": float64(1),
							},
							"route": []interface{}{
								map[string]interface{}{
									"type":    "start",
									"task_id": "-1",
									"location": map[string]interface{}{
										"latitude":  -32.234,
										"longitude": -23.2342,
									},
									"arrival":      "2020-01-01T10:10:00",
									"departure":    "2020-01-01T10:10:00",
									"travel_time":  "00:00:00",
									"setup_time":   "00:00:00",
									"service_time": "00:00:00",
									"waiting_time": "00:00:00",
									"load": []interface{}{
										float64(0),
										float64(0),
									},
									"task_data":  map[string]interface{}{},
									"created_at": "2021-12-08T20:04:16",
									"updated_at": "2021-12-08T20:04:16",
								},
								map[string]interface{}{
									"type":    "pickup",
									"task_id": "3341766951177830852",
									"location": map[string]interface{}{
										"latitude":  -32.234,
										"longitude": -23.2342,
									},
									"arrival":      "2020-01-01T10:10:00",
									"departure":    "2020-01-01T10:10:01",
									"travel_time":  "00:00:00",
									"setup_time":   "00:00:00",
									"service_time": "00:00:01",
									"waiting_time": "00:00:00",
									"load": []interface{}{
										float64(3),
										float64(5),
									},
									"task_data":  map[string]interface{}{},
									"created_at": "2021-12-08T20:04:16",
									"updated_at": "2021-12-08T20:04:16",
								},
								map[string]interface{}{
									"type":    "delivery",
									"task_id": "3341766951177830852",
									"location": map[string]interface{}{
										"latitude":  23.3458,
										"longitude": 2.3242,
									},
									"arrival":      "2020-01-03T20:52:34",
									"departure":    "2020-01-03T20:52:37",
									"travel_time":  "58:42:33",
									"setup_time":   "00:00:00",
									"service_time": "00:00:03",
									"waiting_time": "00:00:00",
									"load": []interface{}{
										float64(0),
										float64(0),
									},
									"task_data":  map[string]interface{}{},
									"created_at": "2021-12-08T20:04:16",
									"updated_at": "2021-12-08T20:04:16",
								},
								map[string]interface{}{
									"type":    "break",
									"task_id": "2349284092384902582",
									"location": map[string]interface{}{
										"latitude":  23.3458,
										"longitude": 2.3242,
									},
									"arrival":      "2020-01-03T20:52:37",
									"departure":    "2020-01-03T20:58:01",
									"travel_time":  "00:00:00",
									"setup_time":   "00:00:00",
									"service_time": "00:05:24",
									"waiting_time": "00:00:00",
									"load": []interface{}{
										float64(0),
										float64(0),
									},
									"task_data":  map[string]interface{}{},
									"created_at": "2021-12-08T20:04:16",
									"updated_at": "2021-12-08T20:04:16",
								},
								map[string]interface{}{
									"type":    "end",
									"task_id": "-1",
									"location": map[string]interface{}{
										"latitude":  23.3458,
										"longitude": 2.3242,
									},
									"arrival":      "2020-01-03T20:58:01",
									"departure":    "2020-01-03T20:58:01",
									"travel_time":  "00:00:00",
									"setup_time":   "00:00:00",
									"service_time": "00:00:00",
									"waiting_time": "00:00:00",
									"load": []interface{}{
										float64(0),
										float64(0),
									},
									"task_data":  map[string]interface{}{},
									"created_at": "2021-12-08T20:04:16",
									"updated_at": "2021-12-08T20:04:16",
								},
							},
						},
					},
					"metadata": map[string]interface{}{
						"summary": []interface{}{
							map[string]interface{}{
								"vehicle_id":   "7300272137290532980",
								"setup_time":   "00:00:00",
								"service_time": "00:05:28",
								"travel_time":  "58:42:33",
								"waiting_time": "00:00:00",
								"vehicle_data": map[string]interface{}{
									"s": float64(1),
								},
							},
						},
						"unassigned":    []interface{}{},
						"total_setup":   "00:00:00",
						"total_service": "00:05:28",
						"total_travel":  "58:42:33",
						"total_waiting": "00:00:00",
					},
					"project_id": "3909655254191459782",
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/vehicles/%d/schedule", tc.vehicleID)
			request, err := http.NewRequest("GET", url, nil)
			// Set the Accept headers to return json
			request.Header.Set("Accept", "application/json")
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

func TestGetVehicleScheduleICal(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		vehicleID  int
		resBody    []util.ICal
		filename   string
	}{
		{
			name:       "Valid ID, no schedule",
			statusCode: 200,
			vehicleID:  2550908592071787332,
			resBody:    []util.ICal{},
			filename:   "schedule-0.ics",
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			vehicleID:  7300272137290532980,
			resBody: []util.ICal{
				{
					ID:          "start7300272137290532980@scheduleserv",
					CreatedTime: time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					ModifiedAt:  time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					StartAt:     time.Date(2020, time.Month(1), 1, 10, 10, 0, 0, time.UTC),
					EndAt:       time.Date(2020, time.Month(1), 1, 10, 10, 0, 0, time.UTC),
					Summary:     "Start - Vehicle 7300272137290532980",
					Location:    "(-32.2340, -23.2342)",
					Description: "Project ID: 3909655254191459782\nVehicle ID: 7300272137290532980\nTask ID: -1\nTravel Time: 00:00:00\nService Time: 00:00:00\nWaiting Time: 00:00:00\nLoad: [0 0]\n",
				},
				{
					ID:          "pickup3341766951177830852@scheduleserv",
					CreatedTime: time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					ModifiedAt:  time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					StartAt:     time.Date(2020, time.Month(1), 1, 10, 10, 0, 0, time.UTC),
					EndAt:       time.Date(2020, time.Month(1), 1, 10, 10, 1, 0, time.UTC),
					Summary:     "Pickup - Vehicle 7300272137290532980",
					Location:    "(-32.2340, -23.2342)",
					Description: "Project ID: 3909655254191459782\nVehicle ID: 7300272137290532980\nTask ID: 3341766951177830852\nTravel Time: 00:00:00\nService Time: 00:00:01\nWaiting Time: 00:00:00\nLoad: [3 5]\n",
				},
				{
					ID:          "delivery3341766951177830852@scheduleserv",
					CreatedTime: time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					ModifiedAt:  time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					StartAt:     time.Date(2020, time.Month(1), 3, 20, 52, 34, 0, time.UTC),
					EndAt:       time.Date(2020, time.Month(1), 3, 20, 52, 37, 0, time.UTC),
					Summary:     "Delivery - Vehicle 7300272137290532980",
					Location:    "(23.3458, 2.3242)",
					Description: "Project ID: 3909655254191459782\nVehicle ID: 7300272137290532980\nTask ID: 3341766951177830852\nTravel Time: 58:42:33\nService Time: 00:00:03\nWaiting Time: 00:00:00\nLoad: [0 0]\n",
				},
				{
					ID:          "break2349284092384902582@scheduleserv",
					CreatedTime: time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					ModifiedAt:  time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					StartAt:     time.Date(2020, time.Month(1), 3, 20, 52, 37, 0, time.UTC),
					EndAt:       time.Date(2020, time.Month(1), 3, 20, 58, 1, 0, time.UTC),
					Summary:     "Break - Vehicle 7300272137290532980",
					Location:    "(23.3458, 2.3242)",
					Description: "Project ID: 3909655254191459782\nVehicle ID: 7300272137290532980\nTask ID: 2349284092384902582\nTravel Time: 00:00:00\nService Time: 00:05:24\nWaiting Time: 00:00:00\nLoad: [0 0]\n",
				},
				{
					ID:          "end7300272137290532980@scheduleserv",
					CreatedTime: time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					ModifiedAt:  time.Date(2021, time.Month(12), 8, 20, 4, 16, 0, time.UTC),
					StartAt:     time.Date(2020, time.Month(1), 3, 20, 58, 1, 0, time.UTC),
					EndAt:       time.Date(2020, time.Month(1), 3, 20, 58, 1, 0, time.UTC),
					Summary:     "End - Vehicle 7300272137290532980",
					Location:    "(23.3458, 2.3242)",
					Description: "Project ID: 3909655254191459782\nVehicle ID: 7300272137290532980\nTask ID: -1\nTravel Time: 00:00:00\nService Time: 00:00:00\nWaiting Time: 00:00:00\nLoad: [0 0]\n",
				},
			},
			filename: "schedule-3909655254191459782.ics",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/vehicles/%d/schedule", tc.vehicleID)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)

			resp := recorder.Result()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			// Removing the current Date Time Stamp from the ical file
			bodyStr := string(body)
			regex := regexp.MustCompile("DTSTAMP.*?\n")
			bodyStr = regex.ReplaceAllString(bodyStr, "")

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, "text/calendar", resp.Header.Get("Content-Type"))
			assert.Equal(t, fmt.Sprintf("attachment; filename=%s", tc.filename), resp.Header.Get("Content-Disposition"))

			expectedIcal := util.SerializeICal(tc.resBody)
			expectedIcal = regex.ReplaceAllString(expectedIcal, "")
			assert.Equal(t, expectedIcal, bodyStr)
		})
	}
}
