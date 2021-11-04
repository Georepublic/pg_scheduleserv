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

func TestCreateProject(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
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
				"errors": []interface{}{"Field 'name' of type 'string' is required"},
			},
		},
		{
			name:       "Only Name",
			statusCode: 201,
			body: map[string]interface{}{
				"name": "Sample Project",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{},
				"name": "Sample Project",
			},
		},
		{
			name:       "Only data",
			statusCode: 400,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'name' of type 'string' is required"},
			},
		},
		{
			name:       "Integer name",
			statusCode: 400,
			body: map[string]interface{}{
				"name": 123,
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'name' must be of string type."},
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
				"name": "123",
				"data": float64(123),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, b := tc.body, new(bytes.Buffer)
			json.NewEncoder(b).Encode(m)
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
			err = json.Unmarshal(body, &m)
			delete(m, "id")
			delete(m, "created_at")
			delete(m, "updated_at")
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestGetProject(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectId  int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			projectId:  100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			projectId:  3909655254191459782,
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Sample Project",
				"data":       "random",
				"created_at": "2021-10-22 23:29:31",
				"updated_at": "2021-10-22 23:29:31",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d", tc.projectId)
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

func TestListProjects(t *testing.T) {
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
					"name":       "Sample Project",
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
			m := []map[string]interface{}{}
			err = json.Unmarshal(body, &m)
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectId  int
		body       map[string]interface{}
		resBody    map[string]interface{}
	}{
		{
			name:       "Empty Body",
			statusCode: 200,
			projectId:  3909655254191459782,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Sample Project",
				"data":       "random",
				"created_at": "2021-10-22 23:29:31",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			projectId:  100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Only Name",
			statusCode: 200,
			projectId:  3909655254191459782,
			body: map[string]interface{}{
				"name": "Another Sample Project",
			},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Another Sample Project",
				"data":       "random",
				"created_at": "2021-10-22 23:29:31",
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			projectId:  3909655254191459782,
			body: map[string]interface{}{
				"data": map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Another Sample Project",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-22 23:29:31",
			},
		},
		{
			name:       "Integer name",
			statusCode: 400,
			projectId:  3909655254191459782,
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
			projectId:  3909655254191459782,
			body: map[string]interface{}{
				"data": 123,
			},
			resBody: map[string]interface{}{
				"id":         "3909655254191459782",
				"name":       "Another Sample Project",
				"data":       float64(123),
				"created_at": "2021-10-22 23:29:31",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, b := tc.body, new(bytes.Buffer)
			json.NewEncoder(b).Encode(m)
			url := fmt.Sprintf("/projects/%d", tc.projectId)
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

func TestDeleteProject(t *testing.T) {
	db := NewDatabase(t)
	server, conn := setup(db)
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectId  int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			projectId:  100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 204,
			projectId:  3909655254191459782,
			resBody: map[string]interface{}{
				"success": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d", tc.projectId)
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
