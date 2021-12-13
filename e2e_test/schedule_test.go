/*GRP-GNU-AGPL******************************************************************

File: schedule_test.go

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
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSchedule(t *testing.T) {
	test_db := NewTestDatabase(t)
	server, conn := setup(test_db, "testdata.sql")
	defer conn.Close(context.Background())
	mux := server.Router

	testCases := []struct {
		name       string
		statusCode int
		projectID  int
		resBody    string
	}{
		{
			name:       "Invalid ID",
			statusCode: 200,
			projectID:  123,
			resBody:    "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//arran4//Golang ICS Library\r\nMETHOD:REQUEST\r\nEND:VCALENDAR\r\n",
		},
		{
			name:       "Valid ID",
			statusCode: 200,
			projectID:  3909655254191459782,
			resBody: "BEGIN:VCALENDAR\r\nVERSION:2.0\r\nPRODID:-//arran4//Golang ICS Library\r\nMETHOD:REQUEST\r\n" +
				"BEGIN:VEVENT\r\nUID:4341723776417023483\r\nCREATED:20211208T200416Z\r\nLAST-MODIFIED:20211208T200416Z\r\nDTSTART:20200101T101000Z\r\nDTEND:20200101T101000Z\r\nSUMMARY:start 7300272137290532980\r\nLOCATION:(-32.2340\\, -23.2342)\r\n" +
				"DESCRIPTION:Project ID: 3909655254191459782\\nVehicle ID:\r\n  7300272137290532980\\nTravel Time: 00:00:00\\nService Time:\r\n  00:00:00\\nWaiting Time: 00:00:00\\nLoad: [0 0] - [0 0]\\n\r\nORGANIZER;CN=This Machine:sender@domain\r\nEND:VEVENT\r\n" +
				"BEGIN:VEVENT\r\nUID:6390629987209858272\r\nCREATED:20211208T200416Z\r\nLAST-MODIFIED:20211208T200416Z\r\nDTSTART:20200101T101000Z\r\nDTEND:20200107T100531Z\r\nSUMMARY:pickup 7300272137290532980\r\nLOCATION:(-32.2340\\, -23.2342)\r\n" +
				"DESCRIPTION:Project ID: 3909655254191459782\\nVehicle ID:\r\n  7300272137290532980\\nTravel Time: 00:00:00\\nService Time:\r\n  00:00:01\\nWaiting Time: 00:00:00\\nLoad: [0 0] - [3 5]\\n\r\nORGANIZER;CN=This Machine:sender@domain\r\nEND:VEVENT\r\n" +
				"BEGIN:VEVENT\r\nUID:5021753332863055108\r\nCREATED:20211208T200416Z\r\nLAST-MODIFIED:20211208T200416Z\r\nDTSTART:20200107T100531Z\r\nDTEND:20200107T100534Z\r\nSUMMARY:delivery 7300272137290532980\r\nLOCATION:(23.3458\\, 2.3242)\r\n" +
				"DESCRIPTION:Project ID: 3909655254191459782\\nVehicle ID:\r\n  7300272137290532980\\nTravel Time: 143:55:30\\nService Time:\r\n  00:00:03\\nWaiting Time: 00:00:00\\nLoad: [3 5] - [0 0]\\n\r\nORGANIZER;CN=This Machine:sender@domain\r\nEND:VEVENT\r\n" +
				"BEGIN:VEVENT\r\nUID:682344376747508512\r\nCREATED:20211208T200416Z\r\nLAST-MODIFIED:20211208T200416Z\r\nDTSTART:20200107T100534Z\r\nDTEND:20200107T101058Z\r\nSUMMARY:break 7300272137290532980\r\nLOCATION:(23.3458\\, 2.3242)\r\n" +
				"DESCRIPTION:Project ID: 3909655254191459782\\nVehicle ID:\r\n  7300272137290532980\\nBreak ID: 2349284092384902582\\nTravel Time:\r\n  143:55:30\\nService Time: 00:05:24\\nWaiting Time: 00:00:00\\nLoad: [0 0] -\r\n  [0 0]\\n\r\nORGANIZER;CN=This Machine:sender@domain\r\nEND:VEVENT\r\n" +
				"BEGIN:VEVENT\r\nUID:3799072960370619615\r\nCREATED:20211208T200416Z\r\nLAST-MODIFIED:20211208T200416Z\r\nDTSTART:20200107T101058Z\r\nDTEND:20200107T101058Z\r\nSUMMARY:end 7300272137290532980\r\nLOCATION:(23.3458\\, 2.3242)\r\n" +
				"DESCRIPTION:Project ID: 3909655254191459782\\nVehicle ID:\r\n  7300272137290532980\\nTravel Time: 143:55:30\\nService Time:\r\n  00:00:00\\nWaiting Time: 00:00:00\\nLoad: [0 0] - [0 0]\\n\r\nORGANIZER;CN=This Machine:sender@domain\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/projects/%d/schedule", tc.projectID)
			request, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)

			resp := recorder.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			// Removing the current Date Time Stamp from the ical file
			bodyStr := string(body)
			m1 := regexp.MustCompile("DTSTAMP.*?\n")
			bodyStr = m1.ReplaceAllString(bodyStr, "")

			assert.Equal(t, tc.statusCode, resp.StatusCode)
			assert.Equal(t, "text/calendar", resp.Header.Get("Content-Type"))
			assert.Equal(t, tc.resBody, bodyStr)
		})
	}
}
