


# pg_scheduleserv API
This is an API for scheduling VRP tasks. Source code can be found on https://github.com/Georepublic/pg_scheduleserv
  

## Informations

### Version

0.1.0

### License

[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html)

### Contact

Team Georepublic info@georepublic.de 

### Terms Of Service

https://swagger.io/terms/

## Content negotiation

### URI Schemes
  * http
  * https

### Consumes
  * application/json

### Produces
  * application/json
  * text/calendar

## All endpoints

###  break

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /breaks/{break_id} | [delete breaks break ID](#delete-breaks-break-id) | Delete a break |
| DELETE | /breaks/{break_id}/time_windows | [delete breaks break ID time windows](#delete-breaks-break-id-time-windows) | Delete break time windows |
| GET | /breaks/{break_id} | [get breaks break ID](#get-breaks-break-id) | Fetch a break |
| GET | /breaks/{break_id}/time_windows | [get breaks break ID time windows](#get-breaks-break-id-time-windows) | List break time windows for a break |
| GET | /vehicles/{vehicle_id}/breaks | [get vehicles vehicle ID breaks](#get-vehicles-vehicle-id-breaks) | List breaks |
| PATCH | /breaks/{break_id} | [patch breaks break ID](#patch-breaks-break-id) | Update a break |
| POST | /breaks/{break_id}/time_windows | [post breaks break ID time windows](#post-breaks-break-id-time-windows) | Create a new break time window |
| POST | /vehicles/{vehicle_id}/breaks | [post vehicles vehicle ID breaks](#post-vehicles-vehicle-id-breaks) | Create a new break |
  


###  job

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /jobs/{job_id} | [delete jobs job ID](#delete-jobs-job-id) | Delete a job |
| DELETE | /jobs/{job_id}/time_windows | [delete jobs job ID time windows](#delete-jobs-job-id-time-windows) | Delete job time windows |
| GET | /jobs/{job_id} | [get jobs job ID](#get-jobs-job-id) | Fetch a job |
| GET | /jobs/{job_id}/schedule | [get jobs job ID schedule](#get-jobs-job-id-schedule) | Get the schedule for a job |
| GET | /jobs/{job_id}/time_windows | [get jobs job ID time windows](#get-jobs-job-id-time-windows) | List job time windows for a job |
| GET | /projects/{project_id}/jobs | [get projects project ID jobs](#get-projects-project-id-jobs) | List jobs for a project |
| PATCH | /jobs/{job_id} | [patch jobs job ID](#patch-jobs-job-id) | Update a job |
| POST | /jobs/{job_id}/time_windows | [post jobs job ID time windows](#post-jobs-job-id-time-windows) | Create a new job time window |
| POST | /projects/{project_id}/jobs | [post projects project ID jobs](#post-projects-project-id-jobs) | Create a new job |
  


###  project

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /projects/{project_id} | [delete projects project ID](#delete-projects-project-id) | Delete a project |
| GET | /projects | [get projects](#get-projects) | List projects |
| GET | /projects/{project_id} | [get projects project ID](#get-projects-project-id) | Fetch a project |
| PATCH | /projects/{project_id} | [patch projects project ID](#patch-projects-project-id) | Update a project |
| POST | /projects | [post projects](#post-projects) | Create a new project |
  


###  schedule

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /projects/{project_id}/schedule | [delete projects project ID schedule](#delete-projects-project-id-schedule) | Delete the schedule |
| GET | /projects/{project_id}/schedule | [get projects project ID schedule](#get-projects-project-id-schedule) | Get the schedule |
| POST | /projects/{project_id}/schedule | [post projects project ID schedule](#post-projects-project-id-schedule) | Schedule the tasks |
  


###  shipment

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /shipments/{shipment_id} | [delete shipments shipment ID](#delete-shipments-shipment-id) | Delete a shipment |
| DELETE | /shipments/{shipment_id}/time_windows | [delete shipments shipment ID time windows](#delete-shipments-shipment-id-time-windows) | Delete shipment time windows |
| GET | /projects/{project_id}/shipments | [get projects project ID shipments](#get-projects-project-id-shipments) | List shipments for a project |
| GET | /shipments/{shipment_id} | [get shipments shipment ID](#get-shipments-shipment-id) | Fetch a shipment |
| GET | /shipments/{shipment_id}/schedule | [get shipments shipment ID schedule](#get-shipments-shipment-id-schedule) | Get the schedule for a shipment |
| GET | /shipments/{shipment_id}/time_windows | [get shipments shipment ID time windows](#get-shipments-shipment-id-time-windows) | List shipment time windows for a shipment |
| PATCH | /shipments/{shipment_id} | [patch shipments shipment ID](#patch-shipments-shipment-id) | Update a shipment |
| POST | /projects/{project_id}/shipments | [post projects project ID shipments](#post-projects-project-id-shipments) | Create a new shipment |
| POST | /shipments/{shipment_id}/time_windows | [post shipments shipment ID time windows](#post-shipments-shipment-id-time-windows) | Create a new shipment time window |
  


###  vehicle

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /vehicles/{vehicle_id} | [delete vehicles vehicle ID](#delete-vehicles-vehicle-id) | Delete a vehicle |
| GET | /projects/{project_id}/vehicles | [get projects project ID vehicles](#get-projects-project-id-vehicles) | List vehicles for a project |
| GET | /vehicles/{vehicle_id} | [get vehicles vehicle ID](#get-vehicles-vehicle-id) | Fetch a vehicle |
| GET | /vehicles/{vehicle_id}/schedule | [get vehicles vehicle ID schedule](#get-vehicles-vehicle-id-schedule) | Get the schedule for a vehicle |
| PATCH | /vehicles/{vehicle_id} | [patch vehicles vehicle ID](#patch-vehicles-vehicle-id) | Update a vehicle |
| POST | /projects/{project_id}/vehicles | [post projects project ID vehicles](#post-projects-project-id-vehicles) | Create a new vehicle |
  


## Paths

### <span id="delete-breaks-break-id"></span> Delete a break (*DeleteBreaksBreakID*)

```
DELETE /breaks/{break_id}
```

Delete a break with its break_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| break_id | `path` | integer | `int64` |  | ✓ |  | Break ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-breaks-break-id-200) | OK | OK |  | [schema](#delete-breaks-break-id-200-schema) |
| [400](#delete-breaks-break-id-400) | Bad Request | Bad Request |  | [schema](#delete-breaks-break-id-400-schema) |
| [404](#delete-breaks-break-id-404) | Not Found | Not Found |  | [schema](#delete-breaks-break-id-404-schema) |

#### Responses


##### <span id="delete-breaks-break-id-200"></span> 200 - OK
Status: OK

###### <span id="delete-breaks-break-id-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-breaks-break-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-breaks-break-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-breaks-break-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-breaks-break-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-breaks-break-id-time-windows"></span> Delete break time windows (*DeleteBreaksBreakIDTimeWindows*)

```
DELETE /breaks/{break_id}/time_windows
```

Delete all break time windows for a break with break_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| break_id | `path` | integer | `int64` |  | ✓ |  | Break ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-breaks-break-id-time-windows-200) | OK | OK |  | [schema](#delete-breaks-break-id-time-windows-200-schema) |
| [400](#delete-breaks-break-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#delete-breaks-break-id-time-windows-400-schema) |
| [404](#delete-breaks-break-id-time-windows-404) | Not Found | Not Found |  | [schema](#delete-breaks-break-id-time-windows-404-schema) |

#### Responses


##### <span id="delete-breaks-break-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="delete-breaks-break-id-time-windows-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-breaks-break-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-breaks-break-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-breaks-break-id-time-windows-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-breaks-break-id-time-windows-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-jobs-job-id"></span> Delete a job (*DeleteJobsJobID*)

```
DELETE /jobs/{job_id}
```

Delete a job with its job_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-jobs-job-id-200) | OK | OK |  | [schema](#delete-jobs-job-id-200-schema) |
| [400](#delete-jobs-job-id-400) | Bad Request | Bad Request |  | [schema](#delete-jobs-job-id-400-schema) |
| [404](#delete-jobs-job-id-404) | Not Found | Not Found |  | [schema](#delete-jobs-job-id-404-schema) |

#### Responses


##### <span id="delete-jobs-job-id-200"></span> 200 - OK
Status: OK

###### <span id="delete-jobs-job-id-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-jobs-job-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-jobs-job-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-jobs-job-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-jobs-job-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-jobs-job-id-time-windows"></span> Delete job time windows (*DeleteJobsJobIDTimeWindows*)

```
DELETE /jobs/{job_id}/time_windows
```

Delete all job time windows for a job with job_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-jobs-job-id-time-windows-200) | OK | OK |  | [schema](#delete-jobs-job-id-time-windows-200-schema) |
| [400](#delete-jobs-job-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#delete-jobs-job-id-time-windows-400-schema) |
| [404](#delete-jobs-job-id-time-windows-404) | Not Found | Not Found |  | [schema](#delete-jobs-job-id-time-windows-404-schema) |

#### Responses


##### <span id="delete-jobs-job-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="delete-jobs-job-id-time-windows-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-jobs-job-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-jobs-job-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-jobs-job-id-time-windows-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-jobs-job-id-time-windows-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-projects-project-id"></span> Delete a project (*DeleteProjectsProjectID*)

```
DELETE /projects/{project_id}
```

Delete a project with its project_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-projects-project-id-200) | OK | OK |  | [schema](#delete-projects-project-id-200-schema) |
| [400](#delete-projects-project-id-400) | Bad Request | Bad Request |  | [schema](#delete-projects-project-id-400-schema) |
| [404](#delete-projects-project-id-404) | Not Found | Not Found |  | [schema](#delete-projects-project-id-404-schema) |

#### Responses


##### <span id="delete-projects-project-id-200"></span> 200 - OK
Status: OK

###### <span id="delete-projects-project-id-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-projects-project-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-projects-project-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-projects-project-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-projects-project-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-projects-project-id-schedule"></span> Delete the schedule (*DeleteProjectsProjectIDSchedule*)

```
DELETE /projects/{project_id}/schedule
```

Delete the schedule for a project

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-projects-project-id-schedule-200) | OK | OK |  | [schema](#delete-projects-project-id-schedule-200-schema) |
| [400](#delete-projects-project-id-schedule-400) | Bad Request | Bad Request |  | [schema](#delete-projects-project-id-schedule-400-schema) |
| [404](#delete-projects-project-id-schedule-404) | Not Found | Not Found |  | [schema](#delete-projects-project-id-schedule-404-schema) |

#### Responses


##### <span id="delete-projects-project-id-schedule-200"></span> 200 - OK
Status: OK

###### <span id="delete-projects-project-id-schedule-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-projects-project-id-schedule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-projects-project-id-schedule-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-projects-project-id-schedule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-projects-project-id-schedule-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-shipments-shipment-id"></span> Delete a shipment (*DeleteShipmentsShipmentID*)

```
DELETE /shipments/{shipment_id}
```

Delete a shipment with its shipment_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-shipments-shipment-id-200) | OK | OK |  | [schema](#delete-shipments-shipment-id-200-schema) |
| [400](#delete-shipments-shipment-id-400) | Bad Request | Bad Request |  | [schema](#delete-shipments-shipment-id-400-schema) |
| [404](#delete-shipments-shipment-id-404) | Not Found | Not Found |  | [schema](#delete-shipments-shipment-id-404-schema) |

#### Responses


##### <span id="delete-shipments-shipment-id-200"></span> 200 - OK
Status: OK

###### <span id="delete-shipments-shipment-id-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-shipments-shipment-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-shipments-shipment-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-shipments-shipment-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-shipments-shipment-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-shipments-shipment-id-time-windows"></span> Delete shipment time windows (*DeleteShipmentsShipmentIDTimeWindows*)

```
DELETE /shipments/{shipment_id}/time_windows
```

Delete all shipment time windows for a shipment with shipment_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-shipments-shipment-id-time-windows-200) | OK | OK |  | [schema](#delete-shipments-shipment-id-time-windows-200-schema) |
| [400](#delete-shipments-shipment-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#delete-shipments-shipment-id-time-windows-400-schema) |
| [404](#delete-shipments-shipment-id-time-windows-404) | Not Found | Not Found |  | [schema](#delete-shipments-shipment-id-time-windows-404-schema) |

#### Responses


##### <span id="delete-shipments-shipment-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="delete-shipments-shipment-id-time-windows-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-shipments-shipment-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-shipments-shipment-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-shipments-shipment-id-time-windows-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-shipments-shipment-id-time-windows-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="delete-vehicles-vehicle-id"></span> Delete a vehicle (*DeleteVehiclesVehicleID*)

```
DELETE /vehicles/{vehicle_id}
```

Delete a vehicle with its vehicle_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| vehicle_id | `path` | integer | `int64` |  | ✓ |  | Vehicle ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-vehicles-vehicle-id-200) | OK | OK |  | [schema](#delete-vehicles-vehicle-id-200-schema) |
| [400](#delete-vehicles-vehicle-id-400) | Bad Request | Bad Request |  | [schema](#delete-vehicles-vehicle-id-400-schema) |
| [404](#delete-vehicles-vehicle-id-404) | Not Found | Not Found |  | [schema](#delete-vehicles-vehicle-id-404-schema) |

#### Responses


##### <span id="delete-vehicles-vehicle-id-200"></span> 200 - OK
Status: OK

###### <span id="delete-vehicles-vehicle-id-200-schema"></span> Schema
   
  

[UtilSuccess](#util-success)

##### <span id="delete-vehicles-vehicle-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="delete-vehicles-vehicle-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="delete-vehicles-vehicle-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="delete-vehicles-vehicle-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

### <span id="get-breaks-break-id"></span> Fetch a break (*GetBreaksBreakID*)

```
GET /breaks/{break_id}
```

Fetch a break with its break_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| break_id | `path` | integer | `int64` |  | ✓ |  | Break ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-breaks-break-id-200) | OK | OK |  | [schema](#get-breaks-break-id-200-schema) |
| [400](#get-breaks-break-id-400) | Bad Request | Bad Request |  | [schema](#get-breaks-break-id-400-schema) |
| [404](#get-breaks-break-id-404) | Not Found | Not Found |  | [schema](#get-breaks-break-id-404-schema) |

#### Responses


##### <span id="get-breaks-break-id-200"></span> 200 - OK
Status: OK

###### <span id="get-breaks-break-id-200-schema"></span> Schema
   
  

[GetBreaksBreakIDOKBody](#get-breaks-break-id-o-k-body)

##### <span id="get-breaks-break-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-breaks-break-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-breaks-break-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-breaks-break-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-breaks-break-id-o-k-body"></span> GetBreaksBreakIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getBreaksBreakIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseBreak](#database-break)| `models.DatabaseBreak` |  | |  |  |



### <span id="get-breaks-break-id-time-windows"></span> List break time windows for a break (*GetBreaksBreakIDTimeWindows*)

```
GET /breaks/{break_id}/time_windows
```

Get a list of break time windows for a break with break_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| break_id | `path` | integer | `int64` |  | ✓ |  | Break ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-breaks-break-id-time-windows-200) | OK | OK |  | [schema](#get-breaks-break-id-time-windows-200-schema) |
| [400](#get-breaks-break-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#get-breaks-break-id-time-windows-400-schema) |

#### Responses


##### <span id="get-breaks-break-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="get-breaks-break-id-time-windows-200-schema"></span> Schema
   
  

[GetBreaksBreakIDTimeWindowsOKBody](#get-breaks-break-id-time-windows-o-k-body)

##### <span id="get-breaks-break-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-breaks-break-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-breaks-break-id-time-windows-o-k-body"></span> GetBreaksBreakIDTimeWindowsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getBreaksBreakIdTimeWindowsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseBreakTimeWindow](#database-break-time-window)| `[]*models.DatabaseBreakTimeWindow` |  | |  |  |



### <span id="get-jobs-job-id"></span> Fetch a job (*GetJobsJobID*)

```
GET /jobs/{job_id}
```

Fetch a job with its job_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-jobs-job-id-200) | OK | OK |  | [schema](#get-jobs-job-id-200-schema) |
| [400](#get-jobs-job-id-400) | Bad Request | Bad Request |  | [schema](#get-jobs-job-id-400-schema) |
| [404](#get-jobs-job-id-404) | Not Found | Not Found |  | [schema](#get-jobs-job-id-404-schema) |

#### Responses


##### <span id="get-jobs-job-id-200"></span> 200 - OK
Status: OK

###### <span id="get-jobs-job-id-200-schema"></span> Schema
   
  

[GetJobsJobIDOKBody](#get-jobs-job-id-o-k-body)

##### <span id="get-jobs-job-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-jobs-job-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-jobs-job-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-jobs-job-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-jobs-job-id-o-k-body"></span> GetJobsJobIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getJobsJobIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseJob](#database-job)| `models.DatabaseJob` |  | |  |  |



### <span id="get-jobs-job-id-schedule"></span> Get the schedule for a job (*GetJobsJobIDSchedule*)

```
GET /jobs/{job_id}/schedule
```

Get the schedule for a job using job_id

#### Consumes
  * application/json

#### Produces
  * application/json
  * text/calendar

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-jobs-job-id-schedule-200) | OK | OK |  | [schema](#get-jobs-job-id-schedule-200-schema) |
| [400](#get-jobs-job-id-schedule-400) | Bad Request | Bad Request |  | [schema](#get-jobs-job-id-schedule-400-schema) |
| [404](#get-jobs-job-id-schedule-404) | Not Found | Not Found |  | [schema](#get-jobs-job-id-schedule-404-schema) |

#### Responses


##### <span id="get-jobs-job-id-schedule-200"></span> 200 - OK
Status: OK

###### <span id="get-jobs-job-id-schedule-200-schema"></span> Schema
   
  

[GetJobsJobIDScheduleOKBody](#get-jobs-job-id-schedule-o-k-body)

##### <span id="get-jobs-job-id-schedule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-jobs-job-id-schedule-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-jobs-job-id-schedule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-jobs-job-id-schedule-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-jobs-job-id-schedule-o-k-body"></span> GetJobsJobIDScheduleOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getJobsJobIdScheduleOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][UtilScheduleDataTask](#util-schedule-data-task)| `[]*models.UtilScheduleDataTask` |  | |  |  |



### <span id="get-jobs-job-id-time-windows"></span> List job time windows for a job (*GetJobsJobIDTimeWindows*)

```
GET /jobs/{job_id}/time_windows
```

Get a list of job time windows for a job with job_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-jobs-job-id-time-windows-200) | OK | OK |  | [schema](#get-jobs-job-id-time-windows-200-schema) |
| [400](#get-jobs-job-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#get-jobs-job-id-time-windows-400-schema) |

#### Responses


##### <span id="get-jobs-job-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="get-jobs-job-id-time-windows-200-schema"></span> Schema
   
  

[GetJobsJobIDTimeWindowsOKBody](#get-jobs-job-id-time-windows-o-k-body)

##### <span id="get-jobs-job-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-jobs-job-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-jobs-job-id-time-windows-o-k-body"></span> GetJobsJobIDTimeWindowsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getJobsJobIdTimeWindowsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseJobTimeWindow](#database-job-time-window)| `[]*models.DatabaseJobTimeWindow` |  | |  |  |



### <span id="get-projects"></span> List projects (*GetProjects*)

```
GET /projects
```

Get a list of projects

#### Consumes
  * application/json

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-projects-200) | OK | OK |  | [schema](#get-projects-200-schema) |
| [400](#get-projects-400) | Bad Request | Bad Request |  | [schema](#get-projects-400-schema) |

#### Responses


##### <span id="get-projects-200"></span> 200 - OK
Status: OK

###### <span id="get-projects-200-schema"></span> Schema
   
  

[GetProjectsOKBody](#get-projects-o-k-body)

##### <span id="get-projects-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-projects-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-projects-o-k-body"></span> GetProjectsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getProjectsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseProject](#database-project)| `[]*models.DatabaseProject` |  | |  |  |



### <span id="get-projects-project-id"></span> Fetch a project (*GetProjectsProjectID*)

```
GET /projects/{project_id}
```

Fetch a project with its project_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-projects-project-id-200) | OK | OK |  | [schema](#get-projects-project-id-200-schema) |
| [400](#get-projects-project-id-400) | Bad Request | Bad Request |  | [schema](#get-projects-project-id-400-schema) |
| [404](#get-projects-project-id-404) | Not Found | Not Found |  | [schema](#get-projects-project-id-404-schema) |

#### Responses


##### <span id="get-projects-project-id-200"></span> 200 - OK
Status: OK

###### <span id="get-projects-project-id-200-schema"></span> Schema
   
  

[GetProjectsProjectIDOKBody](#get-projects-project-id-o-k-body)

##### <span id="get-projects-project-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-projects-project-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-projects-project-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-projects-project-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-projects-project-id-o-k-body"></span> GetProjectsProjectIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getProjectsProjectIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseProject](#database-project)| `models.DatabaseProject` |  | |  |  |



### <span id="get-projects-project-id-jobs"></span> List jobs for a project (*GetProjectsProjectIDJobs*)

```
GET /projects/{project_id}/jobs
```

Get a list of jobs for a project with project_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-projects-project-id-jobs-200) | OK | OK |  | [schema](#get-projects-project-id-jobs-200-schema) |
| [400](#get-projects-project-id-jobs-400) | Bad Request | Bad Request |  | [schema](#get-projects-project-id-jobs-400-schema) |

#### Responses


##### <span id="get-projects-project-id-jobs-200"></span> 200 - OK
Status: OK

###### <span id="get-projects-project-id-jobs-200-schema"></span> Schema
   
  

[GetProjectsProjectIDJobsOKBody](#get-projects-project-id-jobs-o-k-body)

##### <span id="get-projects-project-id-jobs-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-projects-project-id-jobs-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-projects-project-id-jobs-o-k-body"></span> GetProjectsProjectIDJobsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getProjectsProjectIdJobsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseJob](#database-job)| `[]*models.DatabaseJob` |  | |  |  |



### <span id="get-projects-project-id-schedule"></span> Get the schedule (*GetProjectsProjectIDSchedule*)

```
GET /projects/{project_id}/schedule
```

Get the schedule for a project.

**For JSON content type**: When overview = true, only the metadata is returned. Default value is false, which also returns the summary object.

#### Consumes
  * application/json

#### Produces
  * application/json
  * text/calendar

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |
| overview | `query` | boolean | `bool` |  |  |  | Overview |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-projects-project-id-schedule-200) | OK | OK |  | [schema](#get-projects-project-id-schedule-200-schema) |
| [400](#get-projects-project-id-schedule-400) | Bad Request | Bad Request |  | [schema](#get-projects-project-id-schedule-400-schema) |
| [404](#get-projects-project-id-schedule-404) | Not Found | Not Found |  | [schema](#get-projects-project-id-schedule-404-schema) |

#### Responses


##### <span id="get-projects-project-id-schedule-200"></span> 200 - OK
Status: OK

###### <span id="get-projects-project-id-schedule-200-schema"></span> Schema
   
  

[GetProjectsProjectIDScheduleOKBody](#get-projects-project-id-schedule-o-k-body)

##### <span id="get-projects-project-id-schedule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-projects-project-id-schedule-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-projects-project-id-schedule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-projects-project-id-schedule-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-projects-project-id-schedule-o-k-body"></span> GetProjectsProjectIDScheduleOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getProjectsProjectIdScheduleOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [UtilScheduleData](#util-schedule-data)| `models.UtilScheduleData` |  | |  |  |



### <span id="get-projects-project-id-shipments"></span> List shipments for a project (*GetProjectsProjectIDShipments*)

```
GET /projects/{project_id}/shipments
```

Get a list of shipments for a project with project_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-projects-project-id-shipments-200) | OK | OK |  | [schema](#get-projects-project-id-shipments-200-schema) |
| [400](#get-projects-project-id-shipments-400) | Bad Request | Bad Request |  | [schema](#get-projects-project-id-shipments-400-schema) |

#### Responses


##### <span id="get-projects-project-id-shipments-200"></span> 200 - OK
Status: OK

###### <span id="get-projects-project-id-shipments-200-schema"></span> Schema
   
  

[GetProjectsProjectIDShipmentsOKBody](#get-projects-project-id-shipments-o-k-body)

##### <span id="get-projects-project-id-shipments-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-projects-project-id-shipments-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-projects-project-id-shipments-o-k-body"></span> GetProjectsProjectIDShipmentsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getProjectsProjectIdShipmentsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseShipment](#database-shipment)| `[]*models.DatabaseShipment` |  | |  |  |



### <span id="get-projects-project-id-vehicles"></span> List vehicles for a project (*GetProjectsProjectIDVehicles*)

```
GET /projects/{project_id}/vehicles
```

Get a list of vehicles for a project with project_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-projects-project-id-vehicles-200) | OK | OK |  | [schema](#get-projects-project-id-vehicles-200-schema) |
| [400](#get-projects-project-id-vehicles-400) | Bad Request | Bad Request |  | [schema](#get-projects-project-id-vehicles-400-schema) |

#### Responses


##### <span id="get-projects-project-id-vehicles-200"></span> 200 - OK
Status: OK

###### <span id="get-projects-project-id-vehicles-200-schema"></span> Schema
   
  

[GetProjectsProjectIDVehiclesOKBody](#get-projects-project-id-vehicles-o-k-body)

##### <span id="get-projects-project-id-vehicles-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-projects-project-id-vehicles-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-projects-project-id-vehicles-o-k-body"></span> GetProjectsProjectIDVehiclesOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getProjectsProjectIdVehiclesOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseVehicle](#database-vehicle)| `[]*models.DatabaseVehicle` |  | |  |  |



### <span id="get-shipments-shipment-id"></span> Fetch a shipment (*GetShipmentsShipmentID*)

```
GET /shipments/{shipment_id}
```

Fetch a shipment with its shipment_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-shipments-shipment-id-200) | OK | OK |  | [schema](#get-shipments-shipment-id-200-schema) |
| [400](#get-shipments-shipment-id-400) | Bad Request | Bad Request |  | [schema](#get-shipments-shipment-id-400-schema) |
| [404](#get-shipments-shipment-id-404) | Not Found | Not Found |  | [schema](#get-shipments-shipment-id-404-schema) |

#### Responses


##### <span id="get-shipments-shipment-id-200"></span> 200 - OK
Status: OK

###### <span id="get-shipments-shipment-id-200-schema"></span> Schema
   
  

[GetShipmentsShipmentIDOKBody](#get-shipments-shipment-id-o-k-body)

##### <span id="get-shipments-shipment-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-shipments-shipment-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-shipments-shipment-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-shipments-shipment-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-shipments-shipment-id-o-k-body"></span> GetShipmentsShipmentIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getShipmentsShipmentIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseShipment](#database-shipment)| `models.DatabaseShipment` |  | |  |  |



### <span id="get-shipments-shipment-id-schedule"></span> Get the schedule for a shipment (*GetShipmentsShipmentIDSchedule*)

```
GET /shipments/{shipment_id}/schedule
```

Get the schedule for a shipment using shipment_id

#### Consumes
  * application/json

#### Produces
  * application/json
  * text/calendar

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-shipments-shipment-id-schedule-200) | OK | OK |  | [schema](#get-shipments-shipment-id-schedule-200-schema) |
| [400](#get-shipments-shipment-id-schedule-400) | Bad Request | Bad Request |  | [schema](#get-shipments-shipment-id-schedule-400-schema) |
| [404](#get-shipments-shipment-id-schedule-404) | Not Found | Not Found |  | [schema](#get-shipments-shipment-id-schedule-404-schema) |

#### Responses


##### <span id="get-shipments-shipment-id-schedule-200"></span> 200 - OK
Status: OK

###### <span id="get-shipments-shipment-id-schedule-200-schema"></span> Schema
   
  

[GetShipmentsShipmentIDScheduleOKBody](#get-shipments-shipment-id-schedule-o-k-body)

##### <span id="get-shipments-shipment-id-schedule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-shipments-shipment-id-schedule-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-shipments-shipment-id-schedule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-shipments-shipment-id-schedule-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-shipments-shipment-id-schedule-o-k-body"></span> GetShipmentsShipmentIDScheduleOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getShipmentsShipmentIdScheduleOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][UtilScheduleDataTask](#util-schedule-data-task)| `[]*models.UtilScheduleDataTask` |  | |  |  |



### <span id="get-shipments-shipment-id-time-windows"></span> List shipment time windows for a shipment (*GetShipmentsShipmentIDTimeWindows*)

```
GET /shipments/{shipment_id}/time_windows
```

Get a list of shipment time windows for a shipment with shipment_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-shipments-shipment-id-time-windows-200) | OK | OK |  | [schema](#get-shipments-shipment-id-time-windows-200-schema) |
| [400](#get-shipments-shipment-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#get-shipments-shipment-id-time-windows-400-schema) |

#### Responses


##### <span id="get-shipments-shipment-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="get-shipments-shipment-id-time-windows-200-schema"></span> Schema
   
  

[GetShipmentsShipmentIDTimeWindowsOKBody](#get-shipments-shipment-id-time-windows-o-k-body)

##### <span id="get-shipments-shipment-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-shipments-shipment-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-shipments-shipment-id-time-windows-o-k-body"></span> GetShipmentsShipmentIDTimeWindowsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getShipmentsShipmentIdTimeWindowsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseShipmentTimeWindow](#database-shipment-time-window)| `[]*models.DatabaseShipmentTimeWindow` |  | |  |  |



### <span id="get-vehicles-vehicle-id"></span> Fetch a vehicle (*GetVehiclesVehicleID*)

```
GET /vehicles/{vehicle_id}
```

Fetch a vehicle with its vehicle_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| vehicle_id | `path` | integer | `int64` |  | ✓ |  | Vehicle ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-vehicles-vehicle-id-200) | OK | OK |  | [schema](#get-vehicles-vehicle-id-200-schema) |
| [400](#get-vehicles-vehicle-id-400) | Bad Request | Bad Request |  | [schema](#get-vehicles-vehicle-id-400-schema) |
| [404](#get-vehicles-vehicle-id-404) | Not Found | Not Found |  | [schema](#get-vehicles-vehicle-id-404-schema) |

#### Responses


##### <span id="get-vehicles-vehicle-id-200"></span> 200 - OK
Status: OK

###### <span id="get-vehicles-vehicle-id-200-schema"></span> Schema
   
  

[GetVehiclesVehicleIDOKBody](#get-vehicles-vehicle-id-o-k-body)

##### <span id="get-vehicles-vehicle-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-vehicles-vehicle-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-vehicles-vehicle-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-vehicles-vehicle-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-vehicles-vehicle-id-o-k-body"></span> GetVehiclesVehicleIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getVehiclesVehicleIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseVehicle](#database-vehicle)| `models.DatabaseVehicle` |  | |  |  |



### <span id="get-vehicles-vehicle-id-breaks"></span> List breaks (*GetVehiclesVehicleIDBreaks*)

```
GET /vehicles/{vehicle_id}/breaks
```

Get a list of breaks

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| vehicle_id | `path` | integer | `int64` |  | ✓ |  | Vehicle ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-vehicles-vehicle-id-breaks-200) | OK | OK |  | [schema](#get-vehicles-vehicle-id-breaks-200-schema) |
| [400](#get-vehicles-vehicle-id-breaks-400) | Bad Request | Bad Request |  | [schema](#get-vehicles-vehicle-id-breaks-400-schema) |

#### Responses


##### <span id="get-vehicles-vehicle-id-breaks-200"></span> 200 - OK
Status: OK

###### <span id="get-vehicles-vehicle-id-breaks-200-schema"></span> Schema
   
  

[GetVehiclesVehicleIDBreaksOKBody](#get-vehicles-vehicle-id-breaks-o-k-body)

##### <span id="get-vehicles-vehicle-id-breaks-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-vehicles-vehicle-id-breaks-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="get-vehicles-vehicle-id-breaks-o-k-body"></span> GetVehiclesVehicleIDBreaksOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getVehiclesVehicleIdBreaksOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][DatabaseBreak](#database-break)| `[]*models.DatabaseBreak` |  | |  |  |



### <span id="get-vehicles-vehicle-id-schedule"></span> Get the schedule for a vehicle (*GetVehiclesVehicleIDSchedule*)

```
GET /vehicles/{vehicle_id}/schedule
```

Get the schedule for a vehicle using vehicle_id

**For JSON content type**: When overview = true, only the metadata is returned. Default value is false, which also returns the summary object.

#### Consumes
  * application/json

#### Produces
  * application/json
  * text/calendar

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| vehicle_id | `path` | integer | `int64` |  | ✓ |  | Vehicle ID |
| overview | `query` | boolean | `bool` |  |  |  | Overview |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-vehicles-vehicle-id-schedule-200) | OK | OK |  | [schema](#get-vehicles-vehicle-id-schedule-200-schema) |
| [400](#get-vehicles-vehicle-id-schedule-400) | Bad Request | Bad Request |  | [schema](#get-vehicles-vehicle-id-schedule-400-schema) |
| [404](#get-vehicles-vehicle-id-schedule-404) | Not Found | Not Found |  | [schema](#get-vehicles-vehicle-id-schedule-404-schema) |

#### Responses


##### <span id="get-vehicles-vehicle-id-schedule-200"></span> 200 - OK
Status: OK

###### <span id="get-vehicles-vehicle-id-schedule-200-schema"></span> Schema
   
  

[GetVehiclesVehicleIDScheduleOKBody](#get-vehicles-vehicle-id-schedule-o-k-body)

##### <span id="get-vehicles-vehicle-id-schedule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="get-vehicles-vehicle-id-schedule-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="get-vehicles-vehicle-id-schedule-404"></span> 404 - Not Found
Status: Not Found

###### <span id="get-vehicles-vehicle-id-schedule-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="get-vehicles-vehicle-id-schedule-o-k-body"></span> GetVehiclesVehicleIDScheduleOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*getVehiclesVehicleIdScheduleOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][UtilScheduleDB](#util-schedule-d-b)| `[]*models.UtilScheduleDB` |  | |  |  |



### <span id="patch-breaks-break-id"></span> Update a break (*PatchBreaksBreakID*)

```
PATCH /breaks/{break_id}
```

Update a break with its break_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| break_id | `path` | integer | `int64` |  | ✓ |  | Break ID |
| Break | `body` | [DatabaseCreateBreakParams](#database-create-break-params) | `models.DatabaseCreateBreakParams` | | ✓ | | Update break |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#patch-breaks-break-id-200) | OK | OK |  | [schema](#patch-breaks-break-id-200-schema) |
| [400](#patch-breaks-break-id-400) | Bad Request | Bad Request |  | [schema](#patch-breaks-break-id-400-schema) |
| [404](#patch-breaks-break-id-404) | Not Found | Not Found |  | [schema](#patch-breaks-break-id-404-schema) |

#### Responses


##### <span id="patch-breaks-break-id-200"></span> 200 - OK
Status: OK

###### <span id="patch-breaks-break-id-200-schema"></span> Schema
   
  

[PatchBreaksBreakIDOKBody](#patch-breaks-break-id-o-k-body)

##### <span id="patch-breaks-break-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="patch-breaks-break-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="patch-breaks-break-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="patch-breaks-break-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="patch-breaks-break-id-o-k-body"></span> PatchBreaksBreakIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*patchBreaksBreakIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseBreak](#database-break)| `models.DatabaseBreak` |  | |  |  |



### <span id="patch-jobs-job-id"></span> Update a job (*PatchJobsJobID*)

```
PATCH /jobs/{job_id}
```

Update a job (partial update) with its job_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#patch-jobs-job-id-200) | OK | OK |  | [schema](#patch-jobs-job-id-200-schema) |
| [400](#patch-jobs-job-id-400) | Bad Request | Bad Request |  | [schema](#patch-jobs-job-id-400-schema) |
| [404](#patch-jobs-job-id-404) | Not Found | Not Found |  | [schema](#patch-jobs-job-id-404-schema) |

#### Responses


##### <span id="patch-jobs-job-id-200"></span> 200 - OK
Status: OK

###### <span id="patch-jobs-job-id-200-schema"></span> Schema
   
  

[PatchJobsJobIDOKBody](#patch-jobs-job-id-o-k-body)

##### <span id="patch-jobs-job-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="patch-jobs-job-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="patch-jobs-job-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="patch-jobs-job-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="patch-jobs-job-id-o-k-body"></span> PatchJobsJobIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*patchJobsJobIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseJob](#database-job)| `models.DatabaseJob` |  | |  |  |



### <span id="patch-projects-project-id"></span> Update a project (*PatchProjectsProjectID*)

```
PATCH /projects/{project_id}
```

Update a project with its project_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |
| Project | `body` | [DatabaseCreateProjectParams](#database-create-project-params) | `models.DatabaseCreateProjectParams` | | ✓ | | Update project |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#patch-projects-project-id-200) | OK | OK |  | [schema](#patch-projects-project-id-200-schema) |
| [400](#patch-projects-project-id-400) | Bad Request | Bad Request |  | [schema](#patch-projects-project-id-400-schema) |
| [404](#patch-projects-project-id-404) | Not Found | Not Found |  | [schema](#patch-projects-project-id-404-schema) |

#### Responses


##### <span id="patch-projects-project-id-200"></span> 200 - OK
Status: OK

###### <span id="patch-projects-project-id-200-schema"></span> Schema
   
  

[PatchProjectsProjectIDOKBody](#patch-projects-project-id-o-k-body)

##### <span id="patch-projects-project-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="patch-projects-project-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="patch-projects-project-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="patch-projects-project-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="patch-projects-project-id-o-k-body"></span> PatchProjectsProjectIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*patchProjectsProjectIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseProject](#database-project)| `models.DatabaseProject` |  | |  |  |



### <span id="patch-shipments-shipment-id"></span> Update a shipment (*PatchShipmentsShipmentID*)

```
PATCH /shipments/{shipment_id}
```

Update a shipment with its shipment_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |
| Shipment | `body` | [DatabaseUpdateShipmentParams](#database-update-shipment-params) | `models.DatabaseUpdateShipmentParams` | | ✓ | | Update shipment |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#patch-shipments-shipment-id-200) | OK | OK |  | [schema](#patch-shipments-shipment-id-200-schema) |
| [400](#patch-shipments-shipment-id-400) | Bad Request | Bad Request |  | [schema](#patch-shipments-shipment-id-400-schema) |
| [404](#patch-shipments-shipment-id-404) | Not Found | Not Found |  | [schema](#patch-shipments-shipment-id-404-schema) |

#### Responses


##### <span id="patch-shipments-shipment-id-200"></span> 200 - OK
Status: OK

###### <span id="patch-shipments-shipment-id-200-schema"></span> Schema
   
  

[PatchShipmentsShipmentIDOKBody](#patch-shipments-shipment-id-o-k-body)

##### <span id="patch-shipments-shipment-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="patch-shipments-shipment-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="patch-shipments-shipment-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="patch-shipments-shipment-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="patch-shipments-shipment-id-o-k-body"></span> PatchShipmentsShipmentIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*patchShipmentsShipmentIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseShipment](#database-shipment)| `models.DatabaseShipment` |  | |  |  |



### <span id="patch-vehicles-vehicle-id"></span> Update a vehicle (*PatchVehiclesVehicleID*)

```
PATCH /vehicles/{vehicle_id}
```

Update a vehicle with its vehicle_id

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| vehicle_id | `path` | integer | `int64` |  | ✓ |  | Vehicle ID |
| Vehicle | `body` | [DatabaseUpdateVehicleParams](#database-update-vehicle-params) | `models.DatabaseUpdateVehicleParams` | | ✓ | | Update vehicle |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#patch-vehicles-vehicle-id-200) | OK | OK |  | [schema](#patch-vehicles-vehicle-id-200-schema) |
| [400](#patch-vehicles-vehicle-id-400) | Bad Request | Bad Request |  | [schema](#patch-vehicles-vehicle-id-400-schema) |
| [404](#patch-vehicles-vehicle-id-404) | Not Found | Not Found |  | [schema](#patch-vehicles-vehicle-id-404-schema) |

#### Responses


##### <span id="patch-vehicles-vehicle-id-200"></span> 200 - OK
Status: OK

###### <span id="patch-vehicles-vehicle-id-200-schema"></span> Schema
   
  

[PatchVehiclesVehicleIDOKBody](#patch-vehicles-vehicle-id-o-k-body)

##### <span id="patch-vehicles-vehicle-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="patch-vehicles-vehicle-id-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

##### <span id="patch-vehicles-vehicle-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="patch-vehicles-vehicle-id-404-schema"></span> Schema
   
  

[UtilNotFound](#util-not-found)

###### Inlined models

**<span id="patch-vehicles-vehicle-id-o-k-body"></span> PatchVehiclesVehicleIDOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*patchVehiclesVehicleIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseVehicle](#database-vehicle)| `models.DatabaseVehicle` |  | |  |  |



### <span id="post-breaks-break-id-time-windows"></span> Create a new break time window (*PostBreaksBreakIDTimeWindows*)

```
POST /breaks/{break_id}/time_windows
```

Create a new break time window with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| break_id | `path` | integer | `int64` |  | ✓ |  | Break ID |
| BreakTimeWindow | `body` | [DatabaseCreateBreakTimeWindowParams](#database-create-break-time-window-params) | `models.DatabaseCreateBreakTimeWindowParams` | | ✓ | | Create break time window |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-breaks-break-id-time-windows-200) | OK | OK |  | [schema](#post-breaks-break-id-time-windows-200-schema) |
| [400](#post-breaks-break-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#post-breaks-break-id-time-windows-400-schema) |

#### Responses


##### <span id="post-breaks-break-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="post-breaks-break-id-time-windows-200-schema"></span> Schema
   
  

[PostBreaksBreakIDTimeWindowsOKBody](#post-breaks-break-id-time-windows-o-k-body)

##### <span id="post-breaks-break-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-breaks-break-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-breaks-break-id-time-windows-o-k-body"></span> PostBreaksBreakIDTimeWindowsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postBreaksBreakIdTimeWindowsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseBreakTimeWindow](#database-break-time-window)| `models.DatabaseBreakTimeWindow` |  | |  |  |



### <span id="post-jobs-job-id-time-windows"></span> Create a new job time window (*PostJobsJobIDTimeWindows*)

```
POST /jobs/{job_id}/time_windows
```

Create a new job time window with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| job_id | `path` | integer | `int64` |  | ✓ |  | Job ID |
| JobTimeWindow | `body` | [DatabaseCreateJobTimeWindowParams](#database-create-job-time-window-params) | `models.DatabaseCreateJobTimeWindowParams` | | ✓ | | Create job time window |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-jobs-job-id-time-windows-200) | OK | OK |  | [schema](#post-jobs-job-id-time-windows-200-schema) |
| [400](#post-jobs-job-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#post-jobs-job-id-time-windows-400-schema) |

#### Responses


##### <span id="post-jobs-job-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="post-jobs-job-id-time-windows-200-schema"></span> Schema
   
  

[PostJobsJobIDTimeWindowsOKBody](#post-jobs-job-id-time-windows-o-k-body)

##### <span id="post-jobs-job-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-jobs-job-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-jobs-job-id-time-windows-o-k-body"></span> PostJobsJobIDTimeWindowsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postJobsJobIdTimeWindowsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseJobTimeWindow](#database-job-time-window)| `models.DatabaseJobTimeWindow` |  | |  |  |



### <span id="post-projects"></span> Create a new project (*PostProjects*)

```
POST /projects
```

Create a new project with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| Project | `body` | [DatabaseCreateProjectParams](#database-create-project-params) | `models.DatabaseCreateProjectParams` | | ✓ | | Create project |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-projects-200) | OK | OK |  | [schema](#post-projects-200-schema) |
| [400](#post-projects-400) | Bad Request | Bad Request |  | [schema](#post-projects-400-schema) |

#### Responses


##### <span id="post-projects-200"></span> 200 - OK
Status: OK

###### <span id="post-projects-200-schema"></span> Schema
   
  

[PostProjectsOKBody](#post-projects-o-k-body)

##### <span id="post-projects-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-projects-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-projects-o-k-body"></span> PostProjectsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postProjectsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseProject](#database-project)| `models.DatabaseProject` |  | |  |  |



### <span id="post-projects-project-id-jobs"></span> Create a new job (*PostProjectsProjectIDJobs*)

```
POST /projects/{project_id}/jobs
```

Create a new job with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |
| Job | `body` | [DatabaseCreateJobParams](#database-create-job-params) | `models.DatabaseCreateJobParams` | | ✓ | | Job object |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-projects-project-id-jobs-200) | OK | OK |  | [schema](#post-projects-project-id-jobs-200-schema) |
| [400](#post-projects-project-id-jobs-400) | Bad Request | Bad Request |  | [schema](#post-projects-project-id-jobs-400-schema) |

#### Responses


##### <span id="post-projects-project-id-jobs-200"></span> 200 - OK
Status: OK

###### <span id="post-projects-project-id-jobs-200-schema"></span> Schema
   
  

[PostProjectsProjectIDJobsOKBody](#post-projects-project-id-jobs-o-k-body)

##### <span id="post-projects-project-id-jobs-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-projects-project-id-jobs-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-projects-project-id-jobs-o-k-body"></span> PostProjectsProjectIDJobsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postProjectsProjectIdJobsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseJob](#database-job)| `models.DatabaseJob` |  | |  |  |



### <span id="post-projects-project-id-schedule"></span> Schedule the tasks (*PostProjectsProjectIDSchedule*)

```
POST /projects/{project_id}/schedule
```

Schedule the tasks present in a project, deleting any previous schedule and return the new schedule.

**For JSON content type**: When overview = true, only the metadata is returned. Default value is false, which also returns the summary object.

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |
| overview | `query` | boolean | `bool` |  |  |  | Overview |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#post-projects-project-id-schedule-201) | Created | Created |  | [schema](#post-projects-project-id-schedule-201-schema) |
| [400](#post-projects-project-id-schedule-400) | Bad Request | Bad Request |  | [schema](#post-projects-project-id-schedule-400-schema) |

#### Responses


##### <span id="post-projects-project-id-schedule-201"></span> 201 - Created
Status: Created

###### <span id="post-projects-project-id-schedule-201-schema"></span> Schema
   
  

[PostProjectsProjectIDScheduleCreatedBody](#post-projects-project-id-schedule-created-body)

##### <span id="post-projects-project-id-schedule-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-projects-project-id-schedule-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-projects-project-id-schedule-created-body"></span> PostProjectsProjectIDScheduleCreatedBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postProjectsProjectIdScheduleCreatedBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [UtilScheduleData](#util-schedule-data)| `models.UtilScheduleData` |  | |  |  |



### <span id="post-projects-project-id-shipments"></span> Create a new shipment (*PostProjectsProjectIDShipments*)

```
POST /projects/{project_id}/shipments
```

Create a new shipment with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |
| Shipment | `body` | [DatabaseCreateShipmentParams](#database-create-shipment-params) | `models.DatabaseCreateShipmentParams` | | ✓ | | Create shipment |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-projects-project-id-shipments-200) | OK | OK |  | [schema](#post-projects-project-id-shipments-200-schema) |
| [400](#post-projects-project-id-shipments-400) | Bad Request | Bad Request |  | [schema](#post-projects-project-id-shipments-400-schema) |

#### Responses


##### <span id="post-projects-project-id-shipments-200"></span> 200 - OK
Status: OK

###### <span id="post-projects-project-id-shipments-200-schema"></span> Schema
   
  

[PostProjectsProjectIDShipmentsOKBody](#post-projects-project-id-shipments-o-k-body)

##### <span id="post-projects-project-id-shipments-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-projects-project-id-shipments-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-projects-project-id-shipments-o-k-body"></span> PostProjectsProjectIDShipmentsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postProjectsProjectIdShipmentsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseShipment](#database-shipment)| `models.DatabaseShipment` |  | |  |  |



### <span id="post-projects-project-id-vehicles"></span> Create a new vehicle (*PostProjectsProjectIDVehicles*)

```
POST /projects/{project_id}/vehicles
```

Create a new vehicle with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| project_id | `path` | integer | `int64` |  | ✓ |  | Project ID |
| Vehicle | `body` | [DatabaseCreateVehicleParams](#database-create-vehicle-params) | `models.DatabaseCreateVehicleParams` | | ✓ | | Create vehicle |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-projects-project-id-vehicles-200) | OK | OK |  | [schema](#post-projects-project-id-vehicles-200-schema) |
| [400](#post-projects-project-id-vehicles-400) | Bad Request | Bad Request |  | [schema](#post-projects-project-id-vehicles-400-schema) |

#### Responses


##### <span id="post-projects-project-id-vehicles-200"></span> 200 - OK
Status: OK

###### <span id="post-projects-project-id-vehicles-200-schema"></span> Schema
   
  

[PostProjectsProjectIDVehiclesOKBody](#post-projects-project-id-vehicles-o-k-body)

##### <span id="post-projects-project-id-vehicles-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-projects-project-id-vehicles-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-projects-project-id-vehicles-o-k-body"></span> PostProjectsProjectIDVehiclesOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postProjectsProjectIdVehiclesOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseVehicle](#database-vehicle)| `models.DatabaseVehicle` |  | |  |  |



### <span id="post-shipments-shipment-id-time-windows"></span> Create a new shipment time window (*PostShipmentsShipmentIDTimeWindows*)

```
POST /shipments/{shipment_id}/time_windows
```

Create a new shipment time window with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| shipment_id | `path` | integer | `int64` |  | ✓ |  | Shipment ID |
| ShipmentTimeWindow | `body` | [DatabaseCreateShipmentTimeWindowParams](#database-create-shipment-time-window-params) | `models.DatabaseCreateShipmentTimeWindowParams` | | ✓ | | Create shipment time window |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-shipments-shipment-id-time-windows-200) | OK | OK |  | [schema](#post-shipments-shipment-id-time-windows-200-schema) |
| [400](#post-shipments-shipment-id-time-windows-400) | Bad Request | Bad Request |  | [schema](#post-shipments-shipment-id-time-windows-400-schema) |

#### Responses


##### <span id="post-shipments-shipment-id-time-windows-200"></span> 200 - OK
Status: OK

###### <span id="post-shipments-shipment-id-time-windows-200-schema"></span> Schema
   
  

[PostShipmentsShipmentIDTimeWindowsOKBody](#post-shipments-shipment-id-time-windows-o-k-body)

##### <span id="post-shipments-shipment-id-time-windows-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-shipments-shipment-id-time-windows-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-shipments-shipment-id-time-windows-o-k-body"></span> PostShipmentsShipmentIDTimeWindowsOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postShipmentsShipmentIdTimeWindowsOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseShipmentTimeWindow](#database-shipment-time-window)| `models.DatabaseShipmentTimeWindow` |  | |  |  |



### <span id="post-vehicles-vehicle-id-breaks"></span> Create a new break (*PostVehiclesVehicleIDBreaks*)

```
POST /vehicles/{vehicle_id}/breaks
```

Create a new break with the input payload

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| vehicle_id | `path` | integer | `int64` |  | ✓ |  | Vehicle ID |
| Break | `body` | [DatabaseCreateBreakParams](#database-create-break-params) | `models.DatabaseCreateBreakParams` | | ✓ | | Create break |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-vehicles-vehicle-id-breaks-200) | OK | OK |  | [schema](#post-vehicles-vehicle-id-breaks-200-schema) |
| [400](#post-vehicles-vehicle-id-breaks-400) | Bad Request | Bad Request |  | [schema](#post-vehicles-vehicle-id-breaks-400-schema) |

#### Responses


##### <span id="post-vehicles-vehicle-id-breaks-200"></span> 200 - OK
Status: OK

###### <span id="post-vehicles-vehicle-id-breaks-200-schema"></span> Schema
   
  

[PostVehiclesVehicleIDBreaksOKBody](#post-vehicles-vehicle-id-breaks-o-k-body)

##### <span id="post-vehicles-vehicle-id-breaks-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-vehicles-vehicle-id-breaks-400-schema"></span> Schema
   
  

[UtilErrorResponse](#util-error-response)

###### Inlined models

**<span id="post-vehicles-vehicle-id-breaks-o-k-body"></span> PostVehiclesVehicleIDBreaksOKBody**


  


* composed type [UtilSuccessResponse](#util-success-response)
* inlined member (*postVehiclesVehicleIdBreaksOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [DatabaseBreak](#database-break)| `models.DatabaseBreak` |  | |  |  |



## Models

### <span id="database-break"></span> database.Break


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| id | string| `string` |  | |  | `1234567812345678` |
| service | string| `string` |  | |  | `00:02:00` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| vehicle_id | string| `string` |  | |  | `1234567812345678` |



### <span id="database-break-time-window"></span> database.BreakTimeWindow


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| id | string| `string` |  | |  | `1234567812345678` |
| tw_close | string| `string` |  | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` |  | |  | `2021-12-31 23:00:00` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="database-create-break-params"></span> database.CreateBreakParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| service | string| `string` |  | |  | `00:02:00` |



### <span id="database-create-break-time-window-params"></span> database.CreateBreakTimeWindowParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| tw_close | string| `string` | ✓ | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` | ✓ | |  | `2021-12-31 23:00:00` |



### <span id="database-create-job-params"></span> database.CreateJobParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| delivery | []integer| `[]int64` |  | |  | `[10,20]` |
| location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` | ✓ | |  |  |
| pickup | []integer| `[]int64` |  | |  | `[5,15]` |
| priority | integer| `int64` |  | |  | `10` |
| service | string| `string` |  | |  | `00:02:00` |
| setup | string| `string` |  | |  | `00:00:00` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |



### <span id="database-create-job-time-window-params"></span> database.CreateJobTimeWindowParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| tw_close | string| `string` | ✓ | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` | ✓ | |  | `2021-12-31 23:00:00` |



### <span id="database-create-project-params"></span> database.CreateProjectParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| name | string| `string` | ✓ | |  | `Sample Project` |



### <span id="database-create-shipment-params"></span> database.CreateShipmentParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| amount | []integer| `[]int64` |  | |  | `[5,15]` |
| d_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` | ✓ | |  |  |
| d_service | string| `string` |  | |  | `00:02:00` |
| d_setup | string| `string` |  | |  | `00:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| p_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` | ✓ | |  |  |
| p_service | string| `string` |  | |  | `00:02:00` |
| p_setup | string| `string` |  | |  | `00:00:00` |
| priority | integer| `int64` |  | |  | `10` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |



### <span id="database-create-shipment-time-window-params"></span> database.CreateShipmentTimeWindowParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| kind | string| `string` | ✓ | |  | `p` |
| tw_close | string| `string` | ✓ | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` | ✓ | |  | `2021-12-31 23:00:00` |



### <span id="database-create-vehicle-params"></span> database.CreateVehicleParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| capacity | []integer| `[]int64` |  | |  | `[50,25]` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| end_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` | ✓ | |  |  |
| max_tasks | integer| `int64` |  | |  | `20` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |
| speed_factor | number| `float64` |  | |  | `1` |
| start_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` | ✓ | |  |  |
| tw_close | string| `string` |  | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` |  | |  | `2021-12-31 23:00:00` |



### <span id="database-job"></span> database.Job


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| delivery | []integer| `[]int64` |  | |  | `[10,20]` |
| id | string| `string` |  | |  | `1234567812345678` |
| location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| pickup | []integer| `[]int64` |  | |  | `[5,15]` |
| priority | integer| `int64` |  | |  | `10` |
| project_id | string| `string` |  | |  | `1234567812345678` |
| service | string| `string` |  | |  | `00:02:00` |
| setup | string| `string` |  | |  | `00:00:00` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="database-job-time-window"></span> database.JobTimeWindow


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| id | string| `string` |  | |  | `1234567812345678` |
| tw_close | string| `string` |  | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` |  | |  | `2021-12-31 23:00:00` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="database-project"></span> database.Project


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| id | string| `string` |  | |  | `1234567812345678` |
| name | string| `string` |  | |  | `Sample Project` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="database-shipment"></span> database.Shipment


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| amount | []integer| `[]int64` |  | |  | `[5,15]` |
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| d_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| d_service | string| `string` |  | |  | `00:02:00` |
| d_setup | string| `string` |  | |  | `00:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| id | string| `string` |  | |  | `1234567812345678` |
| p_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| p_service | string| `string` |  | |  | `00:02:00` |
| p_setup | string| `string` |  | |  | `00:00:00` |
| priority | integer| `int64` |  | |  | `10` |
| project_id | string| `string` |  | |  | `1234567812345678` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="database-shipment-time-window"></span> database.ShipmentTimeWindow


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| id | string| `string` |  | |  | `1234567812345678` |
| kind | string| `string` |  | |  | `p` |
| tw_close | string| `string` |  | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` |  | |  | `2021-12-31 23:00:00` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="database-update-shipment-params"></span> database.UpdateShipmentParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| amount | []integer| `[]int64` |  | |  | `[5,15]` |
| d_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| d_service | string| `string` |  | |  | `00:02:00` |
| d_setup | string| `string` |  | |  | `00:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| p_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| p_service | string| `string` |  | |  | `00:02:00` |
| p_setup | string| `string` |  | |  | `00:00:00` |
| priority | integer| `int64` |  | |  | `10` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |



### <span id="database-update-vehicle-params"></span> database.UpdateVehicleParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| capacity | []integer| `[]int64` |  | |  | `[50,25]` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| end_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| max_tasks | integer| `int64` |  | |  | `20` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |
| speed_factor | number| `float64` |  | |  | `1` |
| start_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| tw_close | string| `string` |  | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` |  | |  | `2021-12-31 23:00:00` |



### <span id="database-vehicle"></span> database.Vehicle


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| capacity | []integer| `[]int64` |  | |  | `[50,25]` |
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| end_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| id | string| `string` |  | |  | `1234567812345678` |
| max_tasks | integer| `int64` |  | |  | `20` |
| project_id | string| `string` |  | |  | `1234567812345678` |
| skills | []integer| `[]int64` |  | |  | `[1,5]` |
| speed_factor | number| `float64` |  | |  | `1` |
| start_location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| tw_close | string| `string` |  | |  | `2021-12-31 23:59:00` |
| tw_open | string| `string` |  | |  | `2021-12-31 23:00:00` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |



### <span id="util-error-response"></span> util.ErrorResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | string| `string` |  | |  | `400` |
| errors | []string| `[]string` |  | |  | `["Error message1","Error message2"]` |
| message | string| `string` |  | |  | `Bad Request` |



### <span id="util-location-params"></span> util.LocationParams


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| latitude | number| `float64` | ✓ | |  | `2.0365` |
| longitude | number| `float64` | ✓ | |  | `48.6113` |



### <span id="util-metadata-response"></span> util.MetadataResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| summary | [][UtilScheduleSummary](#util-schedule-summary)| `[]*UtilScheduleSummary` |  | |  |  |
| total_service | string| `string` |  | |  | `00:10:00` |
| total_setup | string| `string` |  | |  | `00:05:00` |
| total_travel | string| `string` |  | |  | `01:00:00` |
| total_waiting | string| `string` |  | |  | `00:30:00` |
| unassigned | [][UtilScheduleUnassigned](#util-schedule-unassigned)| `[]*UtilScheduleUnassigned` |  | |  |  |



### <span id="util-not-found"></span> util.NotFound


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | string| `string` |  | |  | `404` |
| error | string| `string` |  | |  | `Not Found` |



### <span id="util-schedule-d-b"></span> util.ScheduleDB


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| arrival | string| `string` |  | |  | `2021-12-01 13:00:00` |
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| departure | string| `string` |  | |  | `2021-12-01 13:00:00` |
| load | []integer| `[]int64` |  | |  | `[0,0]` |
| location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| project_id | string| `string` |  | |  | `1234567812345678` |
| service_time | string| `string` |  | |  | `00:02:00` |
| setup_time | string| `string` |  | |  | `00:00:00` |
| task_data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| task_id | string| `string` |  | |  | `1234567812345678` |
| travel_time | string| `string` |  | |  | `00:16:40` |
| type | string| `string` |  | |  | `job` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| vehicle_data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| vehicle_id | string| `string` |  | |  | `1234567812345678` |
| waiting_time | string| `string` |  | |  | `00:00:00` |



### <span id="util-schedule-data"></span> util.ScheduleData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| metadata | [UtilMetadataResponse](#util-metadata-response)| `UtilMetadataResponse` |  | |  |  |
| project_id | string| `string` |  | |  | `1234567812345678` |
| schedule | [][UtilScheduleResponse](#util-schedule-response)| `[]*UtilScheduleResponse` |  | |  |  |



### <span id="util-schedule-data-task"></span> util.ScheduleDataTask


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| project_id | string| `string` |  | |  | `1234567812345678` |
| schedule | [][UtilScheduleResponse](#util-schedule-response)| `[]*UtilScheduleResponse` |  | |  |  |



### <span id="util-schedule-response"></span> util.ScheduleResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| route | [][UtilScheduleRoute](#util-schedule-route)| `[]*UtilScheduleRoute` |  | |  |  |
| vehicle_data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| vehicle_id | string| `string` |  | |  | `0` |



### <span id="util-schedule-route"></span> util.ScheduleRoute


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| arrival | string| `string` |  | |  | `2021-12-01 13:00:00` |
| created_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| departure | string| `string` |  | |  | `2021-12-01 13:00:00` |
| load | []integer| `[]int64` |  | |  | `[0,0]` |
| location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| service_time | string| `string` |  | |  | `00:02:00` |
| setup_time | string| `string` |  | |  | `00:00:00` |
| task_data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| task_id | string| `string` |  | |  | `1234567812345678` |
| travel_time | string| `string` |  | |  | `00:16:40` |
| type | string| `string` |  | |  | `job` |
| updated_at | string| `string` |  | |  | `2021-12-01 13:00:00` |
| waiting_time | string| `string` |  | |  | `00:00:00` |



### <span id="util-schedule-summary"></span> util.ScheduleSummary


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| service_time | string| `string` |  | |  | `00:02:00` |
| setup_time | string| `string` |  | |  | `00:00:00` |
| travel_time | string| `string` |  | |  | `00:16:40` |
| vehicle_data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| vehicle_id | string| `string` |  | |  | `1234567812345678` |
| waiting_time | string| `string` |  | |  | `00:00:00` |



### <span id="util-schedule-unassigned"></span> util.ScheduleUnassigned


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| location | [UtilLocationParams](#util-location-params)| `UtilLocationParams` |  | |  |  |
| task_data | map of string| `map[string]string` |  | |  | `{"key1":"value1","key2":"value2"}` |
| task_id | string| `string` |  | |  | `1234567812345678` |
| type | string| `string` |  | |  | `job` |



### <span id="util-success"></span> util.Success


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | string| `string` |  | |  | `200` |
| message | string| `string` |  | |  | `OK` |



### <span id="util-success-response"></span> util.SuccessResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | string| `string` |  | |  | `200` |
| data | [interface{}](#interface)| `interface{}` |  | |  |  |
| message | string| `string` |  | |  | `OK` |


