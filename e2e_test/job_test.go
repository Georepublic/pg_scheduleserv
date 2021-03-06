/*GRP-GNU-AGPL******************************************************************

File: job_test.go

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

func TestCreateJob(t *testing.T) {
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
				"errors":  []interface{}{"Field 'location' of type 'util.LocationParams' is required"},
			},
		},
		{
			name:       "Only Location - Wrong parameters 1",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": "Sample Location",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'location' must be of 'util.LocationParams' type."},
			},
		},
		{
			name:       "Only Location - Wrong parameters 2",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
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
				"location": map[string]interface{}{
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
				"location": map[string]interface{}{
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
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"location": map[string]interface{}{
						"latitude":  12.3457,
						"longitude": 56.78,
					},
					"setup":        "00:00:00",
					"service":      "00:00:00",
					"delivery":     []interface{}{},
					"pickup":       []interface{}{},
					"skills":       []interface{}{},
					"priority":     float64(0),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{},
					"time_windows": []interface{}{},
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
				"errors":  []interface{}{"Field 'location' of type 'util.LocationParams' is required"},
			},
		},
		{
			name:       "Priority Min Range incorrect",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"priority": -1,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'priority' must be non-negative"},
			},
		},
		{
			name:       "Priority Max Range incorrect",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"priority": 101,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'priority' must be less than or equal to 100"},
			},
		},
		{
			name:       "Negative skills",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
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
			name:       "Negative pickup",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"pickup": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'pickup[0]' must be non-negative",
					"Field 'pickup[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Negative delivery",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"delivery": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'delivery[0]' must be non-negative",
					"Field 'delivery[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Inconsistent pickup and delivery length",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"pickup":   []interface{}{1},
				"delivery": []interface{}{2, 3},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'pickup' and 'delivery' must have same length"},
			},
		},
		{
			name:       "All fields",
			statusCode: 201,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.3457,
					"longitude": 56.78,
				},
				"setup":    "00:00:10",
				"service":  "00:03:35",
				"delivery": []interface{}{10, 20},
				"pickup":   []interface{}{15, 16},
				"skills":   []interface{}{5, 50, 100},
				"priority": 10,
				"data":     map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"location": map[string]interface{}{
						"latitude":  12.3457,
						"longitude": 56.78,
					},
					"setup":        "00:00:10",
					"service":      "00:03:35",
					"delivery":     []interface{}{float64(10), float64(20)},
					"pickup":       []interface{}{float64(15), float64(16)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(10),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": "value"},
					"time_windows": []interface{}{},
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
			url := fmt.Sprintf("/projects/%d/jobs", tc.projectID)
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

func TestListJobs(t *testing.T) {
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
			projectID:  2593982828701335033,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id": "6362411701075685873",
						"location": map[string]interface{}{
							"latitude":  32.234,
							"longitude": -23.2342,
						},
						"setup":      "00:00:00",
						"service":    "00:02:25",
						"delivery":   []interface{}{float64(10), float64(20)},
						"pickup":     []interface{}{float64(20), float64(30)},
						"skills":     []interface{}{float64(5), float64(50), float64(100)},
						"priority":   float64(11),
						"project_id": "2593982828701335033",
						"data":       map[string]interface{}{"key": "value"},
						"created_at": "2021-10-24T20:31:25",
						"updated_at": "2021-10-24T20:31:25",
						"time_windows": []interface{}{
							[]interface{}{
								"2020-10-10T00:00:00",
								"2020-10-10T00:00:10",
							},
							[]interface{}{
								"2020-10-11T00:00:00",
								"2020-10-12T00:00:00",
							},
						},
					},
					map[string]interface{}{
						"id": "2229737119501208952",
						"location": map[string]interface{}{
							"latitude":  -81.23,
							"longitude": float64(12),
						},
						"setup":      "00:00:00",
						"service":    "00:01:01",
						"delivery":   []interface{}{float64(5), float64(6)},
						"pickup":     []interface{}{float64(7), float64(8)},
						"skills":     []interface{}{},
						"priority":   float64(0),
						"project_id": "2593982828701335033",
						"data":       map[string]interface{}{"data": []interface{}{"value1", float64(2)}},
						"created_at": "2021-10-24T21:12:24",
						"updated_at": "2021-10-24T21:12:24",
						"time_windows": []interface{}{
							[]interface{}{
								"2020-10-10T00:10:00",
								"2020-10-10T00:10:10",
							},
						},
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d/jobs", tc.projectID)
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

func TestGetJob(t *testing.T) {
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
			name:       "Correct ID",
			statusCode: 200,
			jobID:      6362411701075685873,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  32.234,
						"longitude": -23.2342,
					},
					"setup":      "00:00:00",
					"service":    "00:02:25",
					"delivery":   []interface{}{float64(10), float64(20)},
					"pickup":     []interface{}{float64(20), float64(30)},
					"skills":     []interface{}{float64(5), float64(50), float64(100)},
					"priority":   float64(11),
					"project_id": "2593982828701335033",
					"data":       map[string]interface{}{"key": "value"},
					"created_at": "2021-10-24T20:31:25",
					"updated_at": "2021-10-24T20:31:25",
					"time_windows": []interface{}{
						[]interface{}{
							"2020-10-10T00:00:00",
							"2020-10-10T00:00:10",
						},
						[]interface{}{
							"2020-10-11T00:00:00",
							"2020-10-12T00:00:00",
						},
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d", tc.jobID)
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

func TestUpdateJob(t *testing.T) {
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
			statusCode: 200,
			jobID:      6362411701075685873,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  32.234,
						"longitude": -23.2342,
					},
					"setup":        "00:00:00",
					"service":      "00:02:25",
					"delivery":     []interface{}{float64(10), float64(20)},
					"pickup":       []interface{}{float64(20), float64(30)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			jobID:      100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Only Location - Wrong parameters 1",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"location": "Sample Location",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'location' must be of 'util.LocationParams' type."},
			},
		},
		{
			name:       "Only Location - Wrong parameters 2",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
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
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
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
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
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
			name:       "Priority Min Range incorrect",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"priority": -1,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'priority' must be non-negative"},
			},
		},
		{
			name:       "Priority Max Range incorrect",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"priority": 101,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'priority' must be less than or equal to 100"},
			},
		},
		{
			name:       "Negative skills",
			statusCode: 400,
			jobID:      6362411701075685873,
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
			name:       "Negative pickup",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"pickup": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'pickup[0]' must be non-negative",
					"Field 'pickup[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Negative delivery",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"delivery": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'delivery[0]' must be non-negative",
					"Field 'delivery[1]' must be non-negative",
				},
			},
		},
		{
			name:       "Inconsistent pickup and delivery length",
			statusCode: 400,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"pickup":   []interface{}{1},
				"delivery": []interface{}{2, 3},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'pickup' and 'delivery' must have same length"},
			},
		},
		{
			name:       "Only location",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:00:00",
					"service":      "00:02:25",
					"delivery":     []interface{}{float64(10), float64(20)},
					"pickup":       []interface{}{float64(20), float64(30)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only setup",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"setup": "00:01:40",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:02:25",
					"delivery":     []interface{}{float64(10), float64(20)},
					"pickup":       []interface{}{float64(20), float64(30)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only service",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"service": "00:16:45",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(10), float64(20)},
					"pickup":       []interface{}{float64(20), float64(30)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only delivery",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"delivery": []interface{}{20, 30},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(20), float64(30)},
					"pickup":       []interface{}{float64(20), float64(30)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only pickup",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"pickup": []interface{}{10, 20},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(20), float64(30)},
					"pickup":       []interface{}{float64(10), float64(20)},
					"skills":       []interface{}{float64(5), float64(50), float64(100)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only skills",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"skills": []interface{}{5},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(20), float64(30)},
					"pickup":       []interface{}{float64(10), float64(20)},
					"skills":       []interface{}{float64(5)},
					"priority":     float64(11),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only priority",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"priority": 100,
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(20), float64(30)},
					"pickup":       []interface{}{float64(10), float64(20)},
					"skills":       []interface{}{float64(5)},
					"priority":     float64(100),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{"key": "value"},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"data": map[string]interface{}{},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(20), float64(30)},
					"pickup":       []interface{}{float64(10), float64(20)},
					"skills":       []interface{}{float64(5)},
					"priority":     float64(100),
					"project_id":   "2593982828701335033",
					"data":         map[string]interface{}{},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Invalid projectID type",
			statusCode: 400,
			jobID:      6362411701075685873,
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
			jobID:      6362411701075685873,
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
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"project_id": "8943284028902589305",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  23.4567,
						"longitude": -78.90,
					},
					"setup":        "00:01:40",
					"service":      "00:16:45",
					"delivery":     []interface{}{float64(20), float64(30)},
					"pickup":       []interface{}{float64(10), float64(20)},
					"skills":       []interface{}{float64(5)},
					"priority":     float64(100),
					"project_id":   "8943284028902589305",
					"data":         map[string]interface{}{},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "All fields",
			statusCode: 200,
			jobID:      6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  -23.4567,
					"longitude": 78.90,
				},
				"setup":      "00:00:10",
				"service":    "00:01:45",
				"delivery":   []interface{}{20},
				"pickup":     []interface{}{4},
				"skills":     []interface{}{},
				"priority":   float64(0),
				"project_id": "3909655254191459782",
				"data":       map[string]interface{}{"key": 123.23},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "6362411701075685873",
					"location": map[string]interface{}{
						"latitude":  -23.4567,
						"longitude": 78.90,
					},
					"setup":        "00:00:10",
					"service":      "00:01:45",
					"delivery":     []interface{}{float64(20)},
					"pickup":       []interface{}{float64(4)},
					"skills":       []interface{}{},
					"priority":     float64(0),
					"project_id":   "3909655254191459782",
					"data":         map[string]interface{}{"key": 123.23},
					"created_at":   "2021-10-24T20:31:25",
					"time_windows": []interface{}{},
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
			url := fmt.Sprintf("/jobs/%d", tc.jobID)
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

func TestDeleteJob(t *testing.T) {
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
		// TODO
		{
			name:       "Invalid ID",
			statusCode: 200,
			jobID:      100,
			resBody: map[string]interface{}{
				"code":    "200",
				"message": "OK",
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d", tc.jobID)
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

func TestGetJobScheduleJson(t *testing.T) {
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
			jobID:      123,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Valid ID, no schedule",
			statusCode: 200,
			jobID:      6362411701075685873,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"schedule": []interface{}{},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			jobID:      3324729385723589730,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"schedule": []interface{}{
						map[string]interface{}{
							"vehicle_id": "7300272137290532981",
							"vehicle_data": map[string]interface{}{
								"s": float64(1),
							},
							"route": []interface{}{
								map[string]interface{}{
									"type":    "job",
									"task_id": "3324729385723589730",
									"location": map[string]interface{}{
										"latitude":  23.3458,
										"longitude": 2.3242,
									},
									"arrival":      "2020-01-03T16:42:27",
									"departure":    "2020-01-03T16:47:27",
									"travel_time":  "54:32:27",
									"setup_time":   "00:00:00",
									"service_time": "00:05:00",
									"waiting_time": "00:00:00",
									"load": []interface{}{
										0.0,
										0.0,
									},
									"task_data": map[string]interface{}{
										"key": "value",
									},
									"created_at": "2021-12-29T01:05:34",
									"updated_at": "2021-12-29T01:05:34",
								},
							},
						},
					},
					"project_id": "3909655254191459783",
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d/schedule", tc.jobID)
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

func TestGetJobScheduleICal(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobID      int
		resBody    []util.ICal
		filename   string
	}{
		{
			name:       "Valid ID, no schedule",
			statusCode: 200,
			jobID:      2229737119501208952,
			resBody:    []util.ICal{},
			filename:   "schedule-0.ics",
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			jobID:      3324729385723589730,
			resBody: []util.ICal{
				{
					ID:          "job3324729385723589730@scheduleserv",
					CreatedTime: time.Date(2021, time.Month(12), 29, 1, 5, 34, 0, time.UTC),
					ModifiedAt:  time.Date(2021, time.Month(12), 29, 1, 5, 34, 0, time.UTC),
					StartAt:     time.Date(2020, time.Month(1), 3, 16, 42, 27, 0, time.UTC),
					EndAt:       time.Date(2020, time.Month(1), 3, 16, 47, 27, 0, time.UTC),
					Summary:     "Job - Vehicle 7300272137290532981",
					Location:    "(23.3458, 2.3242)",
					Description: "Project ID: 3909655254191459783\nVehicle ID: 7300272137290532981\nTask ID: 3324729385723589730\nTravel Time: 54:32:27\nService Time: 00:05:00\nWaiting Time: 00:00:00\nLoad: [0 0]\n",
				},
			},
			filename: "schedule-3909655254191459783.ics",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d/schedule", tc.jobID)
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
