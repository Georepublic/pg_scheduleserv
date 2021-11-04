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
	ID          int64
	Type        int
	ProjectID   int64
	VehicleID   int64
	JobID       int64
	ShipmentID  int64
	BreakID     int64
	Arrival     string
	Departure   string
	TravelTime  int64
	ServiceTime int64
	WaitingTime int64
	StartLoad   []int64
	EndLoad     []int64
	CreatedAt   string
	UpdatedAt   string
}

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

func getTaskType(schedule Schedule) string {
	switch schedule.Type {
	case 1:
		return "Start"
	case 2:
		return "Job"
	case 3:
		return "Pickup"
	case 4:
		return "Delivery"
	case 5:
		return "Break"
	case 6:
		return "End"
	}
	logrus.Error("Invalid task type")
	return "None"
}

func getSummary(schedule Schedule) string {
	return fmt.Sprintf("%s %d", getTaskType(schedule), schedule.VehicleID)
}

func getLocation(schedule Schedule) string {
	return "TODO"
}

func getDescription(schedule Schedule) string {
	desc := fmt.Sprintf("Project ID: %d\n", schedule.ProjectID)
	desc += fmt.Sprintf("Vehicle ID: %d\n", schedule.VehicleID)
	switch getTaskType(schedule) {
	case "Job":
		desc += fmt.Sprintf("Job ID: %d\n", schedule.JobID)
	case "Shipment":
		desc += fmt.Sprintf("Shipment ID: %d\n", schedule.ShipmentID)
	case "Break":
		desc += fmt.Sprintf("Break ID: %d\n", schedule.BreakID)
	}
	desc += fmt.Sprintf("Travel Time: %s\n", secondsToTime(schedule.TravelTime))
	desc += fmt.Sprintf("Service Time: %s\n", secondsToTime(schedule.ServiceTime))
	desc += fmt.Sprintf("Waiting Time: %s\n", secondsToTime(schedule.WaitingTime))
	desc += fmt.Sprintf("Load: %d - %d\n", schedule.StartLoad, *&schedule.EndLoad)
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
		event.SetURL("https://github.com/Georepublic") // TODO
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
