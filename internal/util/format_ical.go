/*GRP-GNU-AGPL******************************************************************

File: format_ical.go

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
	"bytes"
	"fmt"
	"net/http"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/sirupsen/logrus"
)

type Schedule struct {
	ID          int64   `json:"-" example:"1234567812345678"`
	Type        string  `json:"type" example:"job"`
	ProjectID   int64   `json:"project_id,string" example:"1234567812345678"`
	VehicleID   int64   `json:"vehicle_id,string" example:"1234567812345678"`
	JobID       int64   `json:"job_id,string" example:"1234567812345678"`
	ShipmentID  int64   `json:"shipment_id,string" example:"1234567812345678"`
	BreakID     int64   `json:"break_id,string" example:"1234567812345678"`
	LocationID  int64   `json:"location_id,string" example:"1234567812345678"`
	Arrival     string  `json:"arrival" example:"2021-12-01 13:00:00"`
	Departure   string  `json:"departure" example:"2021-12-01 13:00:00"`
	TravelTime  int64   `json:"travel_time" example:"1000"`
	ServiceTime int64   `json:"service_time" example:"120"`
	WaitingTime int64   `json:"waiting_time" example:"0"`
	StartLoad   []int64 `json:"start_load" example:"0,0"`
	EndLoad     []int64 `json:"end_load" example:"50,25"`
	CreatedAt   string  `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt   string  `json:"updated_at" example:"2021-12-01 13:00:00"`
}

// Example for Schedule in ical format
const ScheduleIcal = `BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//arran4//Golang ICS Library\r\nMETHOD:REQUEST\r\n
BEGIN:VEVENT\r\nUID:1234567812345678\r\nCREATED:20211201T130000Z\r\nLAST-MODIFIED:20211201T130000Z\r\nDTSTART:20211201T130000Z\r\nDTEND:20211201T130000Z\r\nSUMMARY:Summary here\r\nLOCATION:(2.0365\\, 48.6113)\r\n
DESCRIPTION:Description here\\n\r\nORGANIZER;CN=This Machine:sender@domain\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n`

func parseTime(time_str string) time.Time {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, time_str)

	if err != nil {
		logrus.Error(err)
	}
	return t
}

func secondsToTime(totalSecs int64) string {
	hours := totalSecs / 3600
	minutes := (totalSecs % 3600) / 60
	seconds := totalSecs % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getSummary(schedule Schedule) string {
	return fmt.Sprintf("%s %d", schedule.Type, schedule.VehicleID)
}

func getLocation(schedule Schedule) string {
	latitude, longitude := GetCoordinates(schedule.LocationID)
	return fmt.Sprintf("(%.4f, %.4f)", latitude, longitude)
}

func getDescription(schedule Schedule) string {
	desc := fmt.Sprintf("Project ID: %d\n", schedule.ProjectID)
	desc += fmt.Sprintf("Vehicle ID: %d\n", schedule.VehicleID)
	switch schedule.Type {
	case "job":
		desc += fmt.Sprintf("Job ID: %d\n", schedule.JobID)
	case "shipment":
		desc += fmt.Sprintf("Shipment ID: %d\n", schedule.ShipmentID)
	case "break":
		desc += fmt.Sprintf("Break ID: %d\n", schedule.BreakID)
	}
	desc += fmt.Sprintf("Travel Time: %s\n", secondsToTime(schedule.TravelTime))
	desc += fmt.Sprintf("Service Time: %s\n", secondsToTime(schedule.ServiceTime))
	desc += fmt.Sprintf("Waiting Time: %s\n", secondsToTime(schedule.WaitingTime))
	desc += fmt.Sprintf("Load: %d - %d\n", schedule.StartLoad, schedule.EndLoad)
	return desc
}

func (r *Formatter) FormatICAL(w http.ResponseWriter, respCode int, schedule []Schedule) {
	// Set the content-type and response code in the header
	w.Header().Set("Content-Type", "text/calendar")
	w.WriteHeader(respCode)

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for i := 0; i < len(schedule); i++ {
		event := cal.AddEvent(fmt.Sprintf("%d", schedule[i].ID))
		event.SetCreatedTime(parseTime(schedule[i].CreatedAt))
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(parseTime(schedule[i].UpdatedAt))
		event.SetStartAt(parseTime(schedule[i].Arrival))
		event.SetEndAt(parseTime(schedule[i].Departure))
		event.SetSummary(getSummary(schedule[i]))
		event.SetLocation(getLocation(schedule[i]))
		event.SetDescription(getDescription(schedule[i]))
		event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
	}

	b := r.pool.Get().(*bytes.Buffer)
	b.Reset()
	defer r.pool.Put(b)

	b = bytes.NewBufferString(cal.Serialize())

	_, err := b.WriteTo(w)
	if err != nil {
		logrus.Error(err)
	}
}
