/*GRP-GNU-AGPL******************************************************************

File: conversion.go

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
	"math"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GetShipmentTimeWindows(kind []*string, timeWindows [][]*string) ([][]string, [][]string) {
	// If no time windows exist, then [[NULL, NULL]] is stored in timeWindows when LEFT JOIN is performed
	p_parsedTimeWindows := make([][]string, 0)
	d_parsedTimeWindows := make([][]string, 0)
	for i, tw := range timeWindows {
		if tw[0] != nil && tw[1] != nil && kind[i] != nil {
			// depending on the kind of shipment, the time window is either the pickup or the delivery time window
			if *kind[i] == "p" {
				p_parsedTimeWindows = append(p_parsedTimeWindows, []string{*tw[0], *tw[1]})
			} else {
				d_parsedTimeWindows = append(d_parsedTimeWindows, []string{*tw[0], *tw[1]})
			}
		}
	}
	return p_parsedTimeWindows, d_parsedTimeWindows
}

func GetTimeWindows(timeWindows [][]*string) [][]string {
	// If no time windows exist, then [[NULL, NULL]] is stored in timeWindows when LEFT JOIN is performed
	parsedTimeWindows := make([][]string, 0)
	for _, tw := range timeWindows {
		if tw[0] != nil && tw[1] != nil {
			parsedTimeWindows = append(parsedTimeWindows, []string{*tw[0], *tw[1]})
		}
	}
	return parsedTimeWindows
}

func GetLocationId(latitude float64, longitude float64) int64 {
	lat_prefix := '0'
	lon_prefix := '0'
	if latitude < 0 {
		lat_prefix = '1'
	}
	if longitude < 0 {
		lon_prefix = '1'
	}
	s := fmt.Sprintf(
		"%c%07d%c%07d",
		lat_prefix, int(math.Abs(latitude)*10000+0.5),
		lon_prefix, int(math.Abs(longitude)*10000+0.5),
	)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logrus.Errorf("Unable to parse %v to int64", s)
	}
	return i
}

func GetCoordinates(id int64) (latitude float64, longitude float64) {
	latitude = float64(id/100000000) / 10000.0
	if latitude >= 1000 {
		latitude = -(latitude - 1000)
	}
	longitude = float64(id-id/100000000*100000000) / 10000.0
	if longitude >= 1000 {
		longitude = -(longitude - 1000)
	}

	// Rounding to 4 decimal places
	return math.Round(latitude*10000) / 10000, math.Round(longitude*10000) / 10000
}
