/*GRP-GNU-AGPL******************************************************************

File: vroom_examples_test.go

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

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func postRequest(mux *mux.Router, url string, data map[string]interface{}) (response map[string]interface{}, statusCode int, err error) {
	response = map[string]interface{}{}
	m, b := data, new(bytes.Buffer)
	if err = json.NewEncoder(b).Encode(m); err != nil {
		return
	}
	request, err := http.NewRequest("POST", url, b)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)

	resp := recorder.Result()
	statusCode = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		err = fmt.Errorf("Unexpected response type")
		return
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return
	}
	if responseData, ok := response["data"].(map[string]interface{}); ok {
		response = responseData
	} else {
		err = fmt.Errorf("Could not get the data field")
	}

	return
}

func createProject(mux *mux.Router, data map[string]interface{}) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/projects")
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the project")
	}
	return
}

func createJob(mux *mux.Router, data map[string]interface{}, projectID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/projects/%s/jobs", projectID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the job")
	}
	return
}

func createJobTimeWindow(mux *mux.Router, data map[string]interface{}, jobID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/jobs/%s/time_windows", jobID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the job time window")
	}
	return
}

func createShipment(mux *mux.Router, data map[string]interface{}, projectID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/projects/%s/shipments", projectID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the shipment")
	}
	return
}

func createShipmentTimeWindow(mux *mux.Router, data map[string]interface{}, shipmentID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/shipments/%s/time_windows", shipmentID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the shipment time window")
	}
	return
}

func createVehicle(mux *mux.Router, data map[string]interface{}, projectID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/projects/%s/vehicles", projectID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the vehicle")
	}
	return
}

func createBreak(mux *mux.Router, data map[string]interface{}, vehicleID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/vehicles/%s/breaks", vehicleID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the break")
	}
	return
}

func createBreakTimeWindow(mux *mux.Router, data map[string]interface{}, breakID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/breaks/%s/time_windows", breakID)
	response, statusCode, err := postRequest(mux, url, data)
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the break time window")
	}
	return
}

func createSchedule(mux *mux.Router, projectID string) (response map[string]interface{}, err error) {
	url := fmt.Sprintf("/projects/%s/schedule", projectID)
	response, statusCode, err := postRequest(mux, url, map[string]interface{}{})
	if statusCode != 201 {
		err = fmt.Errorf("Could not create the schedule")
	}
	return
}

func TestVroomExample1(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "")
	defer conn.Close(context.Background())
	mux := server.Router

	t.Run("VROOM Example 1", func(t *testing.T) {
		project, err := createProject(mux, map[string]interface{}{
			"name": "VROOM Example 1",
		})
		require.NoError(t, err)
		projectID, ok := project["id"].(string)
		require.Equal(t, true, ok)

		// Vehicles
		_, err = createVehicle(mux, map[string]interface{}{
			"start_location": map[string]interface{}{
				"longitude": 2.35044,
				"latitude":  48.71764,
			},
			"end_location": map[string]interface{}{
				"longitude": 2.35044,
				"latitude":  48.71764,
			},
			"capacity":     []interface{}{4},
			"skills":       []interface{}{1, 14},
			"tw_open":      "2020-09-18 08:00:00",
			"tw_close":     "2020-09-18 12:00:00",
			"speed_factor": 12,
		}, projectID)
		require.NoError(t, err)

		vehicle2, err := createVehicle(mux, map[string]interface{}{
			"start_location": map[string]interface{}{
				"longitude": 2.35044,
				"latitude":  48.71764,
			},
			"end_location": map[string]interface{}{
				"longitude": 2.35044,
				"latitude":  48.71764,
			},
			"capacity":     []interface{}{4},
			"skills":       []interface{}{2, 14},
			"tw_open":      "2020-09-18 08:00:00",
			"tw_close":     "2020-09-18 12:00:00",
			"speed_factor": 12,
		}, projectID)
		require.NoError(t, err)
		vehicle2ID, ok := vehicle2["id"].(string)
		require.Equal(t, true, ok)

		break2, err := createBreak(mux, map[string]interface{}{
			"service": "00:05:00",
		}, vehicle2ID)
		require.NoError(t, err)
		break2ID, ok := break2["id"].(string)
		require.Equal(t, true, ok)

		_, err = createBreakTimeWindow(mux, map[string]interface{}{
			"tw_open":  "2020-09-18 10:00:00",
			"tw_close": "2020-09-18 10:30:00",
		}, break2ID)
		require.NoError(t, err)

		// Jobs
		job1, err := createJob(mux, map[string]interface{}{
			"service":  "00:05:00",
			"delivery": []interface{}{1},
			"location": map[string]interface{}{
				"longitude": 1.98935,
				"latitude":  48.701,
			},
			"skills": []interface{}{1},
		}, projectID)
		require.NoError(t, err)
		job1ID, ok := job1["id"].(string)
		require.Equal(t, true, ok)

		_, err = createJobTimeWindow(mux, map[string]interface{}{
			"tw_open":  "2020-09-18 09:00:00",
			"tw_close": "2020-09-18 10:00:00",
		}, job1ID)
		require.NoError(t, err)

		_, err = createJob(mux, map[string]interface{}{
			"service": "00:05:00",
			"pickup":  []interface{}{1},
			"location": map[string]interface{}{
				"longitude": 2.03655,
				"latitude":  48.61128,
			},
			"skills": []interface{}{1},
		}, projectID)
		require.NoError(t, err)

		_, err = createJob(mux, map[string]interface{}{
			"service":  "00:05:00",
			"delivery": []interface{}{1},
			"location": map[string]interface{}{
				"longitude": 2.28325,
				"latitude":  48.5958,
			},
			"skills": []interface{}{14},
		}, projectID)
		require.NoError(t, err)

		_, err = createJob(mux, map[string]interface{}{
			"service":  "00:05:00",
			"delivery": []interface{}{1},
			"location": map[string]interface{}{
				"longitude": 2.89357,
				"latitude":  48.90736,
			},
			"skills": []interface{}{14},
		}, projectID)
		require.NoError(t, err)

		_, err = createShipment(mux, map[string]interface{}{
			"amount":    []interface{}{1},
			"p_service": "00:05:00",
			"p_location": map[string]interface{}{
				"longitude": 2.41808,
				"latitude":  49.22619,
			},
			"d_service": "00:05:00",
			"d_location": map[string]interface{}{
				"longitude": 2.39719,
				"latitude":  49.07611,
			},
			"skills": []interface{}{2},
		}, projectID)
		require.NoError(t, err)

		schedule, err := createSchedule(mux, projectID)

		expectedResponse := map[string]interface{}{
			"metadata": map[string]interface{}{
				"summary": []interface{}{
					map[string]interface{}{
						"service_time": "00:20:00",
						"setup_time":   "00:00:00",
						"travel_time":  "02:48:24",
						"waiting_time": "00:00:00",
						"vehicle_data": map[string]interface{}{},
						"vehicle_id":   "3377044399950218541",
					}, map[string]interface{}{
						"service_time": "00:15:00",
						"setup_time":   "00:00:00",
						"travel_time":  "01:17:38",
						"waiting_time": "00:00:00",
						"vehicle_data": map[string]interface{}{},
						"vehicle_id":   "3621811169251080866",
					},
				},
				"total_service": "00:35:00",
				"total_setup":   "00:00:00",
				"total_travel":  "04:06:02",
				"total_waiting": "00:00:00",
				"unassigned":    []interface{}{},
			},
			"project_id": "5858246155675874240",
			"schedule": []interface{}{
				map[string]interface{}{
					"vehicle_id":   "3621811169251080866",
					"vehicle_data": map[string]interface{}{},
					"route": []interface{}{
						map[string]interface{}{
							"type": "start",
							"location": map[string]interface{}{
								"latitude":  48.7176,
								"longitude": 2.3504,
							},
							"task_id":   "-1",
							"task_data": map[string]interface{}{},
							"arrival":   "2020-09-18 08:00:00",
							"departure": "2020-09-18 08:00:00",
							"load": []interface{}{
								1,
							},
							"service_time": "00:00:00",
							"setup_time":   "00:00:00",
							"travel_time":  "00:00:00",
							"waiting_time": "00:00:00",
							"created_at":   "2022-01-11 13:06:49",
							"updated_at":   "2022-01-11 13:06:49",
						}, map[string]interface{}{
							"arrival": "2020-09-18 08:50:09", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 08:55:09", "load": []interface{}{
								0,
							}, "location": map[string]interface{}{
								"latitude": 48.9074, "longitude": 2.8936,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "4319466355054430705", "travel_time": "00:50:09", "type": "job", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 09:50:18", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 09:55:18", "load": []interface{}{
								1,
							}, "location": map[string]interface{}{
								"latitude": 49.2262, "longitude": 2.4181,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "6606974347783782672", "travel_time": "00:55:09", "type": "pickup", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 10:13:56", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 10:18:56", "load": []interface{}{
								0,
							}, "location": map[string]interface{}{
								"latitude": 49.0761, "longitude": 2.3972,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "6606974347783782672", "travel_time": "00:18:38", "type": "delivery", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 10:18:56", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 10:23:56", "load": []interface{}{
								0,
							}, "location": map[string]interface{}{
								"latitude": 48.5958, "longitude": 2.2832,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "7123356154219576727", "travel_time": "00:00:00", "type": "break", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 11:08:24", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 11:08:24", "load": []interface{}{
								0,
							}, "location": map[string]interface{}{
								"latitude": 48.7176, "longitude": 2.3504,
							}, "service_time": "00:00:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "-1", "travel_time": "00:44:28", "type": "end", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 08:01:57", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 08:01:57", "load": []interface{}{
								2,
							}, "location": map[string]interface{}{
								"latitude": 48.7176, "longitude": 2.3504,
							}, "service_time": "00:00:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "-1", "travel_time": "00:00:00", "type": "start", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 08:17:58", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 08:22:58", "load": []interface{}{
								1,
							}, "location": map[string]interface{}{
								"latitude": 48.5958, "longitude": 2.2832,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "7315653272110479537", "travel_time": "00:16:01", "type": "job", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 08:43:16", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 08:48:16", "load": []interface{}{
								2,
							}, "location": map[string]interface{}{
								"latitude": 48.6113, "longitude": 2.0366,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "7032679603033425113", "travel_time": "00:20:18", "type": "job", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 09:00:00", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 09:05:00", "load": []interface{}{
								1,
							}, "location": map[string]interface{}{
								"latitude": 48.701, "longitude": 1.9894,
							}, "service_time": "00:05:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "4914873261911329386", "travel_time": "00:11:44", "type": "job", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						}, map[string]interface{}{
							"arrival": "2020-09-18 09:34:35", "created_at": "2022-01-11 13:06:49", "departure": "2020-09-18 09:34:35", "load": []interface{}{
								1,
							}, "location": map[string]interface{}{
								"latitude": 48.7176, "longitude": 2.3504,
							}, "service_time": "00:00:00", "setup_time": "00:00:00", "task_data": map[string]interface{}{}, "task_id": "-1", "travel_time": "00:29:35", "type": "end", "updated_at": "2022-01-11 13:06:49", "waiting_time": "00:00:00",
						},
					},
				},
			},
		}
		assert.Equal(t, expectedResponse, schedule)
	})
}
