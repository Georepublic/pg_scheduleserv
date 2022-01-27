/*GRP-GNU-AGPL******************************************************************

File: osrm_api.go

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

package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// make get request to an url with content-type, and return the response body as json
func Get(url string, contentType string, target interface{}) (int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	return res.StatusCode, json.NewDecoder(res.Body).Decode(target)
}

func GetMatrix(locationIds []int64) (startIds []int64, endIds []int64, durations []int64, err error) {
	// iterate through locationIds, convert all the ids to latitude and longitude, and append [longitude, latitude] in an array
	coordinates := make([][]float64, 0)
	for _, id := range locationIds {
		latitude, longitude := GetCoordinates(id)
		coordinates = append(coordinates, []float64{longitude, latitude})
	}

	// call the osrm api function to get the matrix
	matrix, err := GetMatrixFromOSRM(coordinates)

	if err != nil {
		return nil, nil, nil, err
	}

	// iterate through the 2D matrix and append the start and end ids and durations. start id is locationIds[i], end id is locationIds[j], duration is matrix[i][j]
	for i := 0; i < len(locationIds); i++ {
		for j := 0; j < len(locationIds); j++ {
			startIds = append(startIds, locationIds[i])
			endIds = append(endIds, locationIds[j])
			durations = append(durations, matrix[i][j])
		}
	}

	return startIds, endIds, durations, nil
}

func GetMatrixFromOSRM(coordinates [][]float64) ([][]int64, error) {
	// convert the coordinates to a string
	coordinatesString := make([]string, 0)
	for _, coordinate := range coordinates {
		coordinatesString = append(coordinatesString, fmt.Sprintf("%.4f,%.4f", coordinate[0], coordinate[1]))
	}

	baseUrl := "http://router.project-osrm.org"

	// call the osrm api function to get the matrix
	url := fmt.Sprintf("%s/table/v1/driving/%s", baseUrl, strings.Join(coordinatesString, ";"))

	// decode the response body as json, pass json in Get() function
	response := make(map[string]interface{})
	statusCode, err := Get(url, "application/json", &response)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: %s", response["message"])
	}

	// get the matrix from the response
	matrix := response["durations"].([]interface{})

	// convert the matrix to int64
	matrixInt64 := make([][]int64, 0)
	for _, row := range matrix {
		rowInt64 := make([]int64, 0)
		for _, value := range row.([]interface{}) {
			if value == nil {
				// append max 16 bytes integer value
				rowInt64 = append(rowInt64, int64(1<<16-1))
			} else {
				rowInt64 = append(rowInt64, int64(value.(float64)))
			}
		}
		matrixInt64 = append(matrixInt64, rowInt64)
	}

	return matrixInt64, nil
}
