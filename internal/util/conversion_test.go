/*GRP-GNU-AGPL******************************************************************

File: conversion_test.go

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocationId(t *testing.T) {
	var cases = []struct {
		name      string
		latitude  float64
		longitude float64
		id        int64
	}{
		{"zero_value", 0, 0, 0},
		{"max_value", 90, 180, 90000001800000},
		{"min_value", -90, -180, 1090000011800000},
		{"decimal_one_digit1", 1.23456789, -8.90123456, 1234610089012},
		{"decimal_one_digit2", -8.90123456, 1.23456789, 1008901200012346},
		{"decimal_two_digits1", 12.3456789, -89.0123456, 12345710890123},
		{"decimal_two_digits2", -89.0123456, 12.3456789, 1089012300123457},
		{"decimal_three_digits1", -12.3456789, 170.123456, 1012345701701235},
		{"decimal_three_digits2", 12.3456789, -170.123456, 12345711701235},
	}

	assert := assert.New(t)

	for _, tc := range cases {
		output := GetLocationId(tc.latitude, tc.longitude)
		assert.Equal(tc.id, output, fmt.Sprintf("%s: %v, %v", tc.name, tc.latitude, tc.longitude))
	}
}

func TestGetCoordinates(t *testing.T) {
	var cases = []struct {
		name      string
		latitude  float64
		longitude float64
		id        int64
	}{
		{"zero_value", 0, 0, 0},
		{"max_value", 90, 180, 90000001800000},
		{"min_value", -90, -180, 1090000011800000},
		{"decimal_one_digit1", 1.2346, -8.9012, 1234610089012},
		{"decimal_one_digit2", -8.9012, 1.2346, 1008901200012346},
		{"decimal_two_digits1", 12.3457, -89.0123, 12345710890123},
		{"decimal_two_digits2", -89.0123, 12.3457, 1089012300123457},
		{"decimal_three_digits1", -12.3457, 170.1235, 1012345701701235},
		{"decimal_three_digits2", 12.3457, -170.1235, 12345711701235},
	}

	assert := assert.New(t)

	for _, tc := range cases {
		latitude, longitude := GetCoordinates(tc.id)
		assert.Equal(tc.latitude, latitude, fmt.Sprintf("%s: %v", tc.name, tc.id))
		assert.Equal(tc.longitude, longitude, fmt.Sprintf("%s: %v", tc.name, tc.id))
	}
}
