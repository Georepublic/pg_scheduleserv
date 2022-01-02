/*GRP-GNU-AGPL******************************************************************

File: shipment_tw_test.go

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

func TestCreateShipmentTimeWindow(t *testing.T) {
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
			statusCode: 400,
			shipmentID: 3341766951177830852,
			body:       map[string]interface{}{},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'kind' of type 'string' is required",
					"Field 'tw_open' of type 'string' is required",
					"Field 'tw_close' of type 'string' is required",
				},
			},
		},
		{
			name:       "Only kind",
			statusCode: 400,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"kind": "p",
			},
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
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"tw_open": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'kind' of type 'string' is required",
					"Field 'tw_close' of type 'string' is required",
				},
			},
		},
		{
			name:       "Only tw_close",
			statusCode: 400,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'kind' of type 'string' is required",
					"Field 'tw_open' of type 'string' is required",
				},
			},
		},
		{
			name:       "Opening time greater than closing time",
			statusCode: 400,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"kind":     "p",
				"tw_open":  "2021-10-26 21:24:39",
				"tw_close": "2021-10-26 21:24:38",
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
			name:       "Invalid ShipmentID",
			statusCode: 400,
			shipmentID: 100,
			body: map[string]interface{}{
				"kind":     "p",
				"tw_open":  "2021-10-26 21:24:38",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Shipment with the given 'shipment_id' does not exist",
				},
			},
		},
		{
			name:       "Invalid kind",
			statusCode: 400,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"kind":     "invalid",
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Field 'kind' must be one out of p, d",
				},
			},
		},
		{
			name:       "All fields - Pickup",
			statusCode: 201,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"kind":     "p",
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"kind":     "p",
					"tw_open":  "2021-10-26 21:20:20",
					"tw_close": "2021-10-26 21:24:38",
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "All fields - Delivery",
			statusCode: 201,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"kind":     "d",
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"data": map[string]interface{}{
					"kind":     "d",
					"tw_open":  "2021-10-26 21:20:20",
					"tw_close": "2021-10-26 21:24:38",
				},
				"code":    "201",
				"message": "Created",
			},
		},
		{
			name:       "Primary key violation",
			statusCode: 400,
			shipmentID: 3341766951177830852,
			body: map[string]interface{}{
				"kind":     "d",
				"tw_open":  "2021-10-26 21:20:20",
				"tw_close": "2021-10-26 21:24:38",
			},
			resBody: map[string]interface{}{
				"code":    "400",
				"message": "Bad Request",
				"errors": []interface{}{
					"Shipments time window with given values already exist",
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
			url := fmt.Sprintf("/shipments/%d/time_windows", tc.shipmentID)
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

func TestListShipmentTimeWindows(t *testing.T) {
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
			statusCode: 200,
			shipmentID: 100,
			resBody: map[string]interface{}{
				"data":    []interface{}{},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID - 1",
			statusCode: 200,
			shipmentID: 7794682317520784480,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":         "7794682317520784480",
						"kind":       "p",
						"tw_open":    "2020-10-10 00:00:00",
						"tw_close":   "2020-10-10 00:00:00",
						"created_at": "2021-10-26 20:45:31",
						"updated_at": "2021-10-26 20:45:31",
					},
					map[string]interface{}{
						"id":         "7794682317520784480",
						"kind":       "d",
						"tw_open":    "2020-10-10 00:00:00",
						"tw_close":   "2020-10-11 00:00:00",
						"created_at": "2021-10-26 20:45:31",
						"updated_at": "2021-10-26 20:45:31",
					},
					map[string]interface{}{
						"id":         "7794682317520784480",
						"kind":       "p",
						"tw_open":    "2020-10-10 00:00:10",
						"tw_close":   "2020-10-12 00:00:00",
						"created_at": "2021-10-26 20:45:31",
						"updated_at": "2021-10-26 20:45:31",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID - 2",
			statusCode: 200,
			shipmentID: 3329730179111013588,
			resBody: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id":         "3329730179111013588",
						"kind":       "d",
						"tw_open":    "2020-10-10 00:00:00",
						"tw_close":   "2020-10-10 00:00:00",
						"created_at": "2021-10-26 20:45:31",
						"updated_at": "2021-10-26 20:45:31",
					},
				},
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Valid ID but no time window",
			statusCode: 200,
			shipmentID: 3341766951177830852,
			resBody: map[string]interface{}{
				"data":    []interface{}{},
				"code":    "200",
				"message": "OK",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/shipments/%d/time_windows", tc.shipmentID)
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

func TestDeleteShipmentTimeWindow(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		shipmentID int
		resBody    map[string]interface{}
		todo       bool
	}{
		{
			name:       "Invalid ID",
			statusCode: 404,
			shipmentID: 100,
			resBody: map[string]interface{}{
				"error": "Not Found",
				"code":  "404",
			},
		},
		{
			name:       "Correct ID",
			statusCode: 200,
			shipmentID: 7794682317520784480,
			resBody: map[string]interface{}{
				"code":    "200",
				"message": "OK",
			},
		},
		{
			name:       "Correct ID but no time window",
			statusCode: 404,
			shipmentID: 3341766951177830852,
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
			url := fmt.Sprintf("/shipments/%d/time_windows", tc.shipmentID)
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
