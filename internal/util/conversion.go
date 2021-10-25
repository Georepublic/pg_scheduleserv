package util

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GetLocationIndex(latitude float64, longitude float64) int64 {
	lat_prefix := '0'
	lon_prefix := '0'
	if latitude < 0 {
		lat_prefix = 1
	}
	if longitude < 0 {
		lon_prefix = 1
	}
	s := fmt.Sprintf("%c%07d%c%07d", lat_prefix, int(latitude*10000+0.5), lon_prefix, int(longitude*10000+0.5))
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logrus.Error("Invalid values")
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
	return
}
