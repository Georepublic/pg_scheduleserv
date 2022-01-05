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
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/sirupsen/logrus"
)

type ScheduleDB struct {
	Type        string         `json:"type" example:"job"`
	ProjectID   int64          `json:"project_id,string" example:"1234567812345678"`
	VehicleID   int64          `json:"vehicle_id,string" example:"1234567812345678"`
	TaskID      int64          `json:"task_id,string" example:"1234567812345678"`
	Location    LocationParams `json:"location"`
	Arrival     string         `json:"arrival" example:"2021-12-01 13:00:00"`
	Departure   string         `json:"departure" example:"2021-12-01 13:00:00"`
	TravelTime  string         `json:"travel_time" example:"00:16:40"`
	SetupTime   string         `json:"setup_time" example:"00:00:00"`
	ServiceTime string         `json:"service_time" example:"00:02:00"`
	WaitingTime string         `json:"waiting_time" example:"00:00:00"`
	Load        []int64        `json:"load" example:"0,0"`
	VehicleData interface{}    `json:"vehicle_data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	TaskData    interface{}    `json:"task_data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt   string         `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt   string         `json:"updated_at" example:"2021-12-01 13:00:00"`
}

/*
-------------------------
Schedule Response
-------------------------
*/

type ScheduleRoute struct {
	Type        string         `json:"type" example:"job"`
	TaskID      int64          `json:"task_id,string" example:"1234567812345678"`
	Location    LocationParams `json:"location"`
	Arrival     string         `json:"arrival" example:"2021-12-01 13:00:00"`
	Departure   string         `json:"departure" example:"2021-12-01 13:00:00"`
	TravelTime  string         `json:"travel_time" example:"00:16:40"`
	SetupTime   string         `json:"setup_time" example:"00:00:00"`
	ServiceTime string         `json:"service_time" example:"00:02:00"`
	WaitingTime string         `json:"waiting_time" example:"00:00:00"`
	Load        []int64        `json:"load" example:"0,0"`
	TaskData    interface{}    `json:"task_data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	CreatedAt   string         `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt   string         `json:"updated_at" example:"2021-12-01 13:00:00"`
}

type ScheduleResponse struct {
	VehicleID   int64           `json:"vehicle_id,string"`
	VehicleData interface{}     `json:"vehicle_data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
	Route       []ScheduleRoute `json:"route"`
}

/*
-------------------------
Metadata Response
-------------------------
*/

type ScheduleSummary struct {
	VehicleID   int64       `json:"vehicle_id,string" example:"1234567812345678"`
	TravelTime  string      `json:"travel_time" example:"00:16:40"`
	SetupTime   string      `json:"setup_time" example:"00:00:00"`
	ServiceTime string      `json:"service_time" example:"00:02:00"`
	WaitingTime string      `json:"waiting_time" example:"00:00:00"`
	VehicleData interface{} `json:"vehicle_data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

type ScheduleUnassigned struct {
	Type     string         `json:"type" example:"job"`
	TaskID   int64          `json:"task_id,string" example:"1234567812345678"`
	Location LocationParams `json:"location"`
	TaskData interface{}    `json:"task_data" swaggertype:"object,string" example:"key1:value1,key2:value2"`
}

type MetadataResponse struct {
	Summary      []ScheduleSummary    `json:"summary"`
	Unassigned   []ScheduleUnassigned `json:"unassigned"`
	TotalTravel  string               `json:"total_travel" example:"01:00:00"`
	TotalSetup   string               `json:"total_setup" example:"00:05:00"`
	TotalService string               `json:"total_service" example:"00:10:00"`
	TotalWaiting string               `json:"total_waiting" example:"00:30:00"`
}

/*
-------------------------
Schedule Data
-------------------------
*/

type ScheduleData struct {
	Schedule  []ScheduleResponse `json:"schedule"`
	Metadata  MetadataResponse   `json:"metadata"`
	ProjectID int64              `json:"project_id,string,omitempty" example:"1234567812345678"`
}

type ScheduleDataOverview struct {
	Metadata  MetadataResponse `json:"metadata"`
	ProjectID int64            `json:"project_id,string,omitempty" example:"1234567812345678"`
}

type ScheduleDataTask struct {
	Schedule  []ScheduleResponse `json:"schedule"`
	ProjectID int64              `json:"project_id,string,omitempty" example:"1234567812345678"`
}

/*
-------------------------
ICal Format struct
-------------------------
*/

type ICal struct {
	ID          string
	CreatedTime time.Time
	DtStampTime time.Time
	ModifiedAt  time.Time
	StartAt     time.Time
	EndAt       time.Time
	Summary     string
	Location    string
	Description string
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

func getSummary(route ScheduleRoute, vehicleID int64) string {
	return fmt.Sprintf("%s - Vehicle %d", strings.Title(route.Type), vehicleID)
}

func getLocation(route ScheduleRoute) string {
	return fmt.Sprintf("(%.4f, %.4f)", *route.Location.Latitude, *route.Location.Longitude)
}

func getDescription(route ScheduleRoute, projectID int64, vehicleID int64) string {
	desc := fmt.Sprintf("Project ID: %d\n", projectID)
	desc += fmt.Sprintf("Vehicle ID: %d\n", vehicleID)
	desc += fmt.Sprintf("Task ID: %d\n", route.TaskID)
	desc += fmt.Sprintf("Travel Time: %s\n", route.TravelTime)
	desc += fmt.Sprintf("Service Time: %s\n", route.ServiceTime)
	desc += fmt.Sprintf("Waiting Time: %s\n", route.WaitingTime)
	desc += fmt.Sprintf("Load: %d\n", route.Load)
	return desc
}

func SerializeICal(calendar []ICal) string {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for i := 0; i < len(calendar); i++ {
		event := cal.AddEvent(calendar[i].ID)
		event.SetCreatedTime(calendar[i].CreatedTime)
		event.SetDtStampTime(calendar[i].DtStampTime)
		event.SetModifiedAt(calendar[i].ModifiedAt)
		event.SetStartAt(calendar[i].StartAt)
		event.SetEndAt(calendar[i].EndAt)
		event.SetSummary(calendar[i].Summary)
		event.SetLocation(calendar[i].Location)
		event.SetDescription(calendar[i].Description)
	}
	return cal.Serialize()
}

func (r *Formatter) FormatICAL(w http.ResponseWriter, respCode int, calendar []ICal, filename string) {
	// Set the content-type, content-disposition, and response code in the header
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.WriteHeader(respCode)

	b := r.pool.Get().(*bytes.Buffer)
	b.Reset()
	defer r.pool.Put(b)

	b = bytes.NewBufferString(SerializeICal(calendar))

	_, err := b.WriteTo(w)
	if err != nil {
		logrus.Error(err)
	}
}

func (r *Formatter) GetScheduleICal(scheduleData ScheduleData) ([]ICal, string) {
	var calendar []ICal

	projectID := scheduleData.ProjectID
	filename := fmt.Sprintf("schedule-%d.ics", projectID)
	schedule := scheduleData.Schedule
	for i := 0; i < len(schedule); i++ {
		route := schedule[i].Route
		vehicleID := schedule[i].VehicleID
		for j := 0; j < len(route); j++ {
			entry := ICal{
				ID:          fmt.Sprintf("%s%d@scheduleserv", route[j].Type, route[j].TaskID),
				CreatedTime: parseTime(route[j].CreatedAt),
				DtStampTime: time.Now(),
				ModifiedAt:  parseTime(route[j].UpdatedAt),
				StartAt:     parseTime(route[j].Arrival),
				EndAt:       parseTime(route[j].Departure),
				Summary:     getSummary(route[j], vehicleID),
				Location:    getLocation(route[j]),
				Description: getDescription(route[j], projectID, vehicleID),
			}
			calendar = append(calendar, entry)
		}
	}
	return calendar, filename
}
