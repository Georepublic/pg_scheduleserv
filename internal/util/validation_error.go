/*GRP-GNU-AGPL******************************************************************

File: validation_error.go

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

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func getErrorMsg(ve validator.ValidationErrors) []string {
	var msg []string
	for i := 0; i < len(ve); i++ {
		var err string
		logrus.Info(ve[i].Value())
		switch field := ve[i].Tag(); field {
		case "required":
			err = fmt.Sprintf("Field '%s' of type '%s' is required", ve[i].Field(), ve[i].Type().Elem())
		case "datetime":
			err = fmt.Sprintf("Field '%s' must be of '%s' format", ve[i].Field(), ve[i].Param())
		default:
			err = fmt.Sprintf("Validation of Field '%s' failed on '%s' tag", ve[i].Field(), field)
		}
		msg = append(msg, err)
	}
	return msg
}
