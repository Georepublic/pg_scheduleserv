/*GRP-GNU-AGPL******************************************************************

File: project_test.go

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

func TestCreateProject(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		body       map[string]interface{}
		resBody    map[string]interface{}
	}{
		{
			name:       "Empty Body",
			statusCode: 400,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'name' of type 'string' is required"},
			},
		},
		{
			name:       "Only Name",
			statusCode: 201,
			body: map[string]interface{}{
				"name": "Sample Project",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"data":          map[string]interface{}{},
					"name":          "Sample Project",
					"distance_calc": "euclidean",
					"max_shift":     "00:30:00",
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Only data",
			statusCode: 400,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'name' of type 'string' is required"},
			},
		},
		{
			name:       "Integer name",
			statusCode: 400,
			body: map[string]interface{}{
				"name": 123,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'name' must be of 'string' type."},
			},
		},
		{
			name:       "Integer data",
			statusCode: 201,
			body: map[string]interface{}{
				"name": "123",
				"data": 123,
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"name":          "123",
					"data":          float64(123),
					"distance_calc": "euclidean",
					"max_shift":     "00:30:00",
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "All fields",
			statusCode: 201,
			body: map[string]interface{}{
				"name": "123",
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"name":          "123",
					"data":          map[string]interface{}{"key": "value"},
					"distance_calc": "euclidean",
					"max_shift":     "00:30:00",
				},
				"code":    "201",
				"message": "Created",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, b := tc.body, new(bytes.Buffer)
			if err := json.NewEncoder(b).Encode(m); err != nil {
				t.Error(err)
			}
			request, err := http.NewRequest("POST", "/projects", b)
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

func TestGetProject(t *testing.T) {
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
			name:       "Correct ID",
			statusCode: 200,
			projectID:  3909655254191459782,
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":            "3909655254191459782",
					"name":          "Sample Project",
					"data":          "random",
					"distance_calc": "osrm",
					"max_shift":     "00:30:00",
					"created_at":    "2021-10-22T23:29:31",
					"updated_at":    "2021-10-22T23:29:31",
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d", tc.projectID)
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

func TestListProjects(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close()
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		resBody    map[string]interface{}
	}{
		{
			name:       "All projects",
			statusCode: 200,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":            "3909655254191459782",
						"name":          "Sample Project",
						"data":          "random",
						"distance_calc": "osrm",
						"max_shift":     "00:30:00",
						"created_at":    "2021-10-22T23:29:31",
						"updated_at":    "2021-10-22T23:29:31",
					},
					map[string]interface{}{
						"id":            "3909655254191459783",
						"name":          "Sample Project2",
						"data":          "random",
						"distance_calc": "osrm",
						"max_shift":     "00:30:00",
						"created_at":    "2021-10-22T23:29:31",
						"updated_at":    "2021-10-22T23:29:31",
					},
					map[string]interface{}{
						"id":            "2593982828701335033",
						"name":          "",
						"data":          map[string]interface{}{"s": float64(1)},
						"distance_calc": "osrm",
						"max_shift":     "00:30:00",
						"created_at":    "2021-10-24T19:52:52",
						"updated_at":    "2021-10-24T19:52:52",
					},
					map[string]interface{}{
						"id":            "8943284028902589305",
						"name":          "",
						"data":          map[string]interface{}{"s": float64(1)},
						"distance_calc": "osrm",
						"max_shift":     "00:30:00",
						"created_at":    "2021-10-24T19:52:52",
						"updated_at":    "2021-10-24T19:52:52",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/projects", nil)
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

func TestUpdateProject(t *testing.T) {
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
	}{
		{
			name:       "Empty Body",
			statusCode: 200,
			projectID:  3909655254191459782,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":            "3909655254191459782",
					"name":          "Sample Project",
					"data":          "random",
					"distance_calc": "osrm",
					"max_shift":     "00:30:00",
					"created_at":    "2021-10-22T23:29:31",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			projectID:  100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Only Name",
			statusCode: 200,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"name": "Another Sample Project",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":            "3909655254191459782",
					"name":          "Another Sample Project",
					"data":          "random",
					"distance_calc": "osrm",
					"max_shift":     "00:30:00",
					"created_at":    "2021-10-22T23:29:31",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":            "3909655254191459782",
					"name":          "Another Sample Project",
					"data":          map[string]interface{}{"key": "value"},
					"distance_calc": "osrm",
					"max_shift":     "00:30:00",
					"created_at":    "2021-10-22T23:29:31",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Integer name",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"name": 123,
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors":  []interface{}{"Field 'name' must be of 'string' type."},
			},
		},
		{
			name:       "Integer data",
			statusCode: 200,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"data": 123,
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":            "3909655254191459782",
					"name":          "Another Sample Project",
					"data":          float64(123),
					"distance_calc": "osrm",
					"max_shift":     "00:30:00",
					"created_at":    "2021-10-22T23:29:31",
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "All fields",
			statusCode: 200,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"name": "Final Sample Project",
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"id":            "3909655254191459782",
					"name":          "Final Sample Project",
					"data":          map[string]interface{}{"key": "value"},
					"distance_calc": "osrm",
					"max_shift":     "00:30:00",
					"created_at":    "2021-10-22T23:29:31",
				},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, b := tc.body, new(bytes.Buffer)
			if err := json.NewEncoder(b).Encode(m); err != nil {
				t.Error(err)
			}
			url := fmt.Sprintf("/projects/%d", tc.projectID)
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

func TestDeleteProject(t *testing.T) {
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
			name:       "Correct ID",
			statusCode: 200,
			projectID:  3909655254191459782,
			resBody: map[string]interface{}{
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d", tc.projectID)
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
