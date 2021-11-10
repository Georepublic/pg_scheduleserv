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

func TestCreateJob(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
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
				"errors": []interface{}{"Field 'location' of type 'util.LocationParams' is required"},
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
				"errors": []interface{}{"Field 'location' must be of 'util.LocationParams' type."},
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
				"errors": []interface{}{
					"Field 'latitude' of type 'float64' is required",
					"Field 'longitude' of type 'float64' is required",
				},
			},
			todo: true,
		},
		{
			name:       "Only Location - Wrong range of parameters",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'latitude' must be between -90 and 90",
					"Field 'longitude' must be between -180 and 180",
				},
			},
			todo: true,
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
				"location": map[string]interface{}{
					"latitude":  12.3457,
					"longitude": 56.78,
				},
				"service":    float64(0),
				"delivery":   []interface{}{},
				"pickup":     []interface{}{},
				"skills":     []interface{}{},
				"priority":   float64(0),
				"project_id": "3909655254191459782",
				"data":       map[string]interface{}{},
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
				"errors": []interface{}{"Field 'location' of type 'util.LocationParams' is required"},
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
				"errors": []interface{}{"Field 'priority' must be between 1 and 100"},
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
				"errors": []interface{}{"Field 'priority' must be between 1 and 100"},
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
				"service":  15,
				"delivery": []interface{}{float64(10), float64(20)},
				"pickup":   []interface{}{float64(15), float64(16)},
				"skills":   []interface{}{float64(5), float64(50)},
				"priority": 10,
				"data":     map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  12.3457,
					"longitude": 56.78,
				},
				"service":    float64(15),
				"delivery":   []interface{}{float64(10), float64(20)},
				"pickup":     []interface{}{float64(15), float64(16)},
				"skills":     []interface{}{float64(5), float64(50)},
				"priority":   float64(10),
				"project_id": "3909655254191459782",
				"data":       map[string]interface{}{"key": "value"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.todo == true {
				t.Skip("TODO")
			}
			m, b := tc.body, new(bytes.Buffer)
			json.NewEncoder(b).Encode(m)
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
			err = json.Unmarshal(body, &m)
			delete(m, "id")
			delete(m, "created_at")
			delete(m, "updated_at")
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestGetJob(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobId      int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			jobId:      100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			jobId:      6362411701075685873,
			resBody: map[string]interface{}{
				"id":         "6362411701075685873",
				"name":       "Sample Job",
				"data":       "random",
				"created_at": "2021-10-22 23:29:31",
				"updated_at": "2021-10-22 23:29:31",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d", tc.jobId)
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
			err = json.Unmarshal(body, &m)
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestListJobs(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		resBody    []map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 200,
			resBody: []map[string]interface{}{
				{
					"id":         "3909655254191459782",
					"name":       "Sample Job",
					"data":       "random",
					"created_at": "2021-10-22 23:29:31",
					"updated_at": "2021-10-22 23:29:31",
				},
				{
					"id":         "2593982828701335033",
					"name":       "",
					"data":       map[string]interface{}{"s": float64(1)},
					"created_at": "2021-10-24 19:52:52",
					"updated_at": "2021-10-24 19:52:52",
				},
				{
					"id":         "8943284028902589305",
					"name":       "",
					"data":       map[string]interface{}{"s": float64(1)},
					"created_at": "2021-10-24 19:52:52",
					"updated_at": "2021-10-24 19:52:52",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/jobs", nil)
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
			err = json.Unmarshal(body, &m)
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestUpdateJob(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobId      int
		body       map[string]interface{}
		resBody    map[string]interface{}
	}{
		{
			name:       "Empty Body",
			statusCode: 200,
			jobId:      3909655254191459782,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Sample Job",
				"data":       "random",
				"created_at": "2021-10-22 23:29:31",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			jobId:      100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Only Name",
			statusCode: 200,
			jobId:      3909655254191459782,
			body: map[string]interface{}{
				"name": "Another Sample Job",
			},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Another Sample Job",
				"data":       "random",
				"created_at": "2021-10-22 23:29:31",
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			jobId:      3909655254191459782,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Another Sample Job",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-22 23:29:31",
			},
		},
		{
			name:       "Integer name",
			statusCode: 400,
			jobId:      3909655254191459782,
			body: map[string]interface{}{
				"name": 123,
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'name' must be of string type."},
			},
		},
		{
			name:       "Integer data",
			statusCode: 200,
			jobId:      3909655254191459782,
			body: map[string]interface{}{
				"data": 123,
			},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Another Sample Job",
				"data":       float64(123),
				"created_at": "2021-10-22 23:29:31",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, b := tc.body, new(bytes.Buffer)
			json.NewEncoder(b).Encode(m)
			url := fmt.Sprintf("/jobs/%d", tc.jobId)
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
			err = json.Unmarshal(body, &m)
			delete(m, "updated_at")
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestDeleteJob(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		jobId      int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			jobId:      100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 204,
			jobId:      3909655254191459782,
			resBody: map[string]interface{}{
				"success": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/jobs/%d", tc.jobId)
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
			err = json.Unmarshal(body, &m)
			assert.Equal(t, tc.resBody, m)
		})
	}
}
