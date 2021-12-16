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

type Schedule struct {
	ID          int64          `json:"-" example:"1234567812345678"`
	Type        string         `json:"type" example:"job"`
	ProjectID   int64          `json:"project_id,string" example:"1234567812345678"`
	VehicleID   int64          `json:"vehicle_id,string" example:"1234567812345678"`
	JobID       int64          `json:"job_id,string,omitempty" example:"1234567812345678"`
	ShipmentID  int64          `json:"shipment_id,string,omitempty" example:"1234567812345678"`
	BreakID     int64          `json:"break_id,string,omitempty" example:"1234567812345678"`
	Location    LocationParams `json:"location"`
	Arrival     string         `json:"arrival" example:"2021-12-01 13:00:00"`
	Departure   string         `json:"departure" example:"2021-12-01 13:00:00"`
	TravelTime  int64          `json:"travel_time" example:"1000"`
	ServiceTime int64          `json:"service_time" example:"120"`
	WaitingTime int64          `json:"waiting_time" example:"0"`
	StartLoad   []int64        `json:"start_load" example:"0,0"`
	EndLoad     []int64        `json:"end_load" example:"50,25"`
	CreatedAt   string         `json:"created_at" example:"2021-12-01 13:00:00"`
	UpdatedAt   string         `json:"updated_at" example:"2021-12-01 13:00:00"`
}

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

func secondsToTime(totalSecs int64) string {
	hours := totalSecs / 3600
	minutes := (totalSecs % 3600) / 60
	seconds := totalSecs % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getSummary(schedule Schedule) string {
	return fmt.Sprintf("%s - Vehicle %d", strings.Title(schedule.Type), schedule.VehicleID)
}

func getLocation(schedule Schedule) string {
	return fmt.Sprintf("(%.4f, %.4f)", *schedule.Location.Latitude, *schedule.Location.Longitude)
}

func getDescription(schedule Schedule) string {
	desc := fmt.Sprintf("Project ID: %d\n", schedule.ProjectID)
	desc += fmt.Sprintf("Vehicle ID: %d\n", schedule.VehicleID)
	switch schedule.Type {
	case "job":
		desc += fmt.Sprintf("Job ID: %d\n", schedule.JobID)
	case "pickup":
		desc += fmt.Sprintf("Shipment ID: %d\n", schedule.ShipmentID)
	case "delivery":
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

func (r *Formatter) GetScheduleICal(schedule []Schedule) ([]ICal, string) {
	var calendar []ICal

	var projectID int64
	if len(schedule) > 0 {
		projectID = schedule[0].ProjectID
	} else {
		projectID = 0
	}
	filename := fmt.Sprintf("schedule-%d.ics", projectID)
	for i := 0; i < len(schedule); i++ {
		entry := ICal{
			ID:          fmt.Sprintf("%d", schedule[i].ID),
			CreatedTime: parseTime(schedule[i].CreatedAt),
			DtStampTime: time.Now(),
			ModifiedAt:  parseTime(schedule[i].UpdatedAt),
			StartAt:     parseTime(schedule[i].Arrival),
			EndAt:       parseTime(schedule[i].Departure),
			Summary:     getSummary(schedule[i]),
			Location:    getLocation(schedule[i]),
			Description: getDescription(schedule[i]),
		}
		calendar = append(calendar, entry)
	}
	return calendar, filename
}
