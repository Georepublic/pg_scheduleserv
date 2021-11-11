/*GRP-GNU-AGPL******************************************************************

File: shipment_test.go

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

func TestCreateShipment(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
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
				"errors": []interface{}{
					"Field 'p_location' of type 'util.LocationParams' is required",
					"Field 'd_location' of type 'util.LocationParams' is required",
				},
			},
		},
		{
			name:       "Only Location - Wrong parameters 1",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": "Sample Location",
				"d_location": "Sample Location",
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'p_location' must be of 'util.LocationParams' type.",
					"Field 'd_location' must be of 'util.LocationParams' type.",
				},
			},
		},
		{
			name:       "Only Location - Wrong parameters 2",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  "12.34567",
					"longitude": "56.78",
				},
				"d_location": map[string]interface{}{
					"latitude":  "12.34567",
					"longitude": "56.78",
				},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'latitude' of type 'float64' is required",
					"Field 'longitude' of type 'float64' is required",
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
				"p_location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
				"d_location": map[string]interface{}{
					"latitude":  112.34567,
					"longitude": 256.78,
				},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{
					"Field 'latitude' must be between -90 and 90",
					"Field 'longitude' must be between -180 and 180",
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
				"p_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"d_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
			},
			resBody: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.3457,
					"longitude": 56.78,
				},
				"p_service": float64(0),
				"d_location": map[string]interface{}{
					"latitude":  -12.3457,
					"longitude": -56.78,
				},
				"d_service":  float64(0),
				"amount":     []interface{}{},
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
				"errors": []interface{}{
					"Field 'p_location' of type 'util.LocationParams' is required",
					"Field 'd_location' of type 'util.LocationParams' is required",
				},
			},
		},
		{
			name:       "Priority Min Range incorrect",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"d_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"priority": -1,
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'priority' must be between 1 and 100"},
			},
			todo: true,
		},
		{
			name:       "Priority Max Range incorrect",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"d_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"priority": 101,
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'priority' must be between 1 and 100"},
			},
			todo: true,
		},
		{
			name:       "Negative skills",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"d_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"skills": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'skills' must have non-negative values"},
			},
			todo: true,
		},
		{
			name:       "Negative amount",
			statusCode: 400,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"d_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"amount": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'amount' must have non-negative values"},
			},
			todo: true,
		},
		{
			name:       "All fields",
			statusCode: 201,
			projectID:  3909655254191459782,
			body: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.34567,
					"longitude": 56.78,
				},
				"p_service": 115,
				"d_location": map[string]interface{}{
					"latitude":  -12.34567,
					"longitude": -56.78,
				},
				"d_service": 215,
				"amount":    []interface{}{15, 16},
				"skills":    []interface{}{5, 50, 100},
				"priority":  10,
				"data":      map[string]interface{}{"key": "value"},
			},
			resBody: map[string]interface{}{
				"p_location": map[string]interface{}{
					"latitude":  12.3457,
					"longitude": 56.78,
				},
				"p_service": float64(115),
				"d_location": map[string]interface{}{
					"latitude":  -12.3457,
					"longitude": -56.78,
				},
				"d_service":  float64(215),
				"amount":     []interface{}{float64(15), float64(16)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
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
			url := fmt.Sprintf("/projects/%d/shipments", tc.projectID)
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

func TestListShipments(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectID  int
		resBody    []map[string]interface{}
	}{
		// TODO: Check this
		{
			name:       "Invalid ID",
			statusCode: 200,
			projectID:  100,
			resBody:    []map[string]interface{}{},
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			projectID:  2593982828701335033,
			resBody: []map[string]interface{}{
				{
					"id": "7794682317520784480",
					"p_location": map[string]interface{}{
						"latitude":  32.234,
						"longitude": -23.2342,
					},
					"d_location": map[string]interface{}{
						"latitude":  23.3458,
						"longitude": 2.3242,
					},
					"service":    float64(145),
					"delivery":   []interface{}{float64(10), float64(20)},
					"pickup":     []interface{}{float64(20), float64(30)},
					"skills":     []interface{}{float64(5), float64(50), float64(100)},
					"priority":   float64(11),
					"project_id": "2593982828701335033",
					"data":       map[string]interface{}{"key": "value"},
					"created_at": "2021-10-24 20:31:25",
					"updated_at": "2021-10-24 20:31:25",
				},
				{
					"id": "3329730179111013588",
					"p_location": map[string]interface{}{
						"latitude":  -32.234,
						"longitude": -23.2342,
					},
					"d_location": map[string]interface{}{
						"latitude":  23.3458,
						"longitude": 2.3242,
					},
					"service":    float64(61),
					"delivery":   []interface{}{float64(5), float64(6)},
					"pickup":     []interface{}{float64(7), float64(8)},
					"skills":     []interface{}{},
					"priority":   float64(0),
					"project_id": "2593982828701335033",
					"data":       map[string]interface{}{"data": []interface{}{"value1", float64(2)}},
					"created_at": "2021-10-24 21:12:24",
					"updated_at": "2021-10-24 21:12:24",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d/shipments", tc.projectID)
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
			err = json.Unmarshal(body, &m)
			assert.Equal(t, tc.resBody, m)
		})
	}
}

func TestGetShipment(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		shipmentID int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			shipmentID: 100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			shipmentID: 7794682317520784480,
			resBody: map[string]interface{}{
				"id": "7794682317520784480",
				"location": map[string]interface{}{
					"latitude":  32.234,
					"longitude": -23.2342,
				},
				"service":    float64(145),
				"delivery":   []interface{}{float64(10), float64(20)},
				"pickup":     []interface{}{float64(20), float64(30)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
				"updated_at": "2021-10-24 20:31:25",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/shipments/%d", tc.shipmentID)
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

func TestUpdateShipment(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		shipmentID int
		body       map[string]interface{}
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Empty Body",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  32.234,
					"longitude": -23.2342,
				},
				"service":    float64(145),
				"delivery":   []interface{}{float64(10), float64(20)},
				"pickup":     []interface{}{float64(20), float64(30)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Invalid ID",
			statusCode: 404,
			shipmentID: 100,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Only Location - Wrong parameters 1",
			statusCode: 400,
			shipmentID: 6362411701075685873,
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
			shipmentID: 6362411701075685873,
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
			shipmentID: 6362411701075685873,
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
			name:       "Priority Min Range incorrect",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"priority": -1,
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'priority' must be between 1 and 100"},
			},
			todo: true,
		},
		{
			name:       "Priority Max Range incorrect",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"priority": 101,
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'priority' must be between 1 and 100"},
			},
			todo: true,
		},
		{
			name:       "Negative skills",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"skills": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'skills' must have non-negative values"},
			},
			todo: true,
		},
		{
			name:       "Negative pickup",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"pickup": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'pickup' must have non-negative values"},
			},
			todo: true,
		},
		{
			name:       "Negative delivery",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"pickup": []interface{}{-1, -2},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'delivery' must have non-negative values"},
			},
			todo: true,
		},
		{
			name:       "Inconsistent pickup and delivery length",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"pickup":   []interface{}{1},
				"delivery": []interface{}{2, 3},
			},
			resBody: map[string]interface{}{
				"errors": []interface{}{"Field 'pickup' and 'delivery' must have same length"},
			},
			todo: true,
		},
		{
			name:       "Only location",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(145),
				"delivery":   []interface{}{float64(10), float64(20)},
				"pickup":     []interface{}{float64(20), float64(30)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Only service",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"service": 1005,
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(10), float64(20)},
				"pickup":     []interface{}{float64(20), float64(30)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Only delivery",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"delivery": []interface{}{20, 30},
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(20), float64(30)},
				"pickup":     []interface{}{float64(20), float64(30)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Only pickup",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"pickup": []interface{}{10, 20},
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(20), float64(30)},
				"pickup":     []interface{}{float64(10), float64(20)},
				"skills":     []interface{}{float64(5), float64(50), float64(100)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Only skills",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"skills": []interface{}{5},
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(20), float64(30)},
				"pickup":     []interface{}{float64(10), float64(20)},
				"skills":     []interface{}{float64(5)},
				"priority":   float64(11),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Only priority",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"priority": 100,
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(20), float64(30)},
				"pickup":     []interface{}{float64(10), float64(20)},
				"skills":     []interface{}{float64(5)},
				"priority":   float64(100),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{"key": "value"},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Only data",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"data": map[string]interface{}{},
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(20), float64(30)},
				"pickup":     []interface{}{float64(10), float64(20)},
				"skills":     []interface{}{float64(5)},
				"priority":   float64(100),
				"project_id": "2593982828701335033",
				"data":       map[string]interface{}{},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "Invalid projectID type",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"project_id": 100,
			},
			resBody: map[string]interface{}{"errors": []interface{}{"Field 'project_id' must be of 'string' type."}},
		},
		{
			name:       "Invalid projectID",
			statusCode: 400,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"project_id": "100",
			},
			resBody: map[string]interface{}{"errors": []interface{}{"Project with the given 'project_id' does not exist."}},
			todo:    true,
		},
		{
			name:       "Valid projectID",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"project_id": "8943284028902589305",
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  23.4567,
					"longitude": -78.90,
				},
				"service":    float64(1005),
				"delivery":   []interface{}{float64(20), float64(30)},
				"pickup":     []interface{}{float64(10), float64(20)},
				"skills":     []interface{}{float64(5)},
				"priority":   float64(100),
				"project_id": "8943284028902589305",
				"data":       map[string]interface{}{},
				"created_at": "2021-10-24 20:31:25",
			},
		},
		{
			name:       "All fields",
			statusCode: 200,
			shipmentID: 6362411701075685873,
			body: map[string]interface{}{
				"location": map[string]interface{}{
					"latitude":  -23.4567,
					"longitude": 78.90,
				},
				"service":    105,
				"delivery":   []interface{}{20},
				"pickup":     []interface{}{4},
				"skills":     []interface{}{},
				"priority":   float64(0),
				"project_id": "3909655254191459782",
				"data":       map[string]interface{}{"key": 123.23},
			},
			resBody: map[string]interface{}{
				"id": "6362411701075685873",
				"location": map[string]interface{}{
					"latitude":  -23.4567,
					"longitude": 78.90,
				},
				"service":    float64(105),
				"delivery":   []interface{}{float64(20)},
				"pickup":     []interface{}{float64(4)},
				"skills":     []interface{}{},
				"priority":   float64(0),
				"project_id": "3909655254191459782",
				"data":       map[string]interface{}{"key": 123.23},
				"created_at": "2021-10-24 20:31:25",
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
			url := fmt.Sprintf("/shipments/%d", tc.shipmentID)
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

func TestDeleteShipment(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		shipmentID int
		resBody    map[string]interface{}
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			shipmentID: 100,
			resBody: map[string]interface{}{
				"error": "Not Found",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			shipmentID: 7794682317520784480,
			resBody: map[string]interface{}{
				"success": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/shipments/%d", tc.shipmentID)
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
