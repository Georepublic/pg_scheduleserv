import VehicleAPI from "../api/VehicleAPI.js";
import BreakAPI from "../api/BreakAPI.js";
import VehicleHandler from "../handlers/VehicleHandler.js";
import BreakHandler from "../handlers/BreakHandler.js";
import Random from "../utils/Random.js";
import AbstractView from "./AbstractView.js";

export default class VehicleView extends AbstractView {
  constructor(params) {
    super(params, false);

    // get the vehicles from the params
    this.vehicles = params.vehicles;
    this.projectID = params.projectID;
    this.mapView = params.mapView;
    this.vehicleAPI = new VehicleAPI();
    this.breakAPI = new BreakAPI();

    this.vehicleLeftDiv = document.createElement("div");
    document.querySelector("#app-left").appendChild(this.vehicleLeftDiv);

    this.handler = new VehicleHandler(
      params.vehicles,
      this.getEmptyVehicle(),
      this.handlers()
    );
    this.breakHandler = new BreakHandler(this.breakHandlers());
  }

  // render the vehicles for this project
  render() {
    // get the html for the vehicles
    let vehiclesHtml = this.getVehiclesHtml();

    // set the html for the vehicles
    this.vehicleLeftDiv.innerHTML = vehiclesHtml;

    this.mapView.addVehicleMarkers(this.vehicles);
    this.mapView.fitAllMarkers();
  }

  // get the html for the vehicles
  getVehiclesHtml() {
    // get the html for each vehicle
    let vehiclesHtml = this.vehicles.map((vehicle) => {
      return this.getVehicleHtml(vehicle);
    });

    if (vehiclesHtml.length === 0) {
      vehiclesHtml = [
        `
        <div class="list-group-item flex-column align-items-start">
          <p class="mb-1">No vehicles found...</p>
        </div>
      `,
      ];
    }

    return `
      <div class="list-group">
        <div class="card">
          <div class="card-header vehicle-view-heading">
            <h5 class="mb-0">
              Vehicles
              <button type="button" class="btn btn-success" data-action="vehicle-create" style="float: right">Add</button>
            </h5>
          </div>
          <div class="card-body-custom">
            ${vehiclesHtml.join("")}
          </div>
        </div>
      </div>
    `;
  }

  // get the html for the vehicle
  getVehicleHtml(vehicle) {
    const color = Random.getRandomColor(vehicle.id);
    let html = `
      <div style="background-color: ${color}">
      <div class="list-group-item flex-column align-items-start" data-action="vehicle-view" data-id="${
        vehicle.id
      }">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">${vehicle.id}</h5>
        </div>

        <div class="d-flex w-100 justify-content-between">
          <p class="mb-1">${JSON.stringify(vehicle.data)}</p>
        </div>
      </div>
      </div>
    `;

    // return the html for the vehicle
    return html;
  }

  getCompleteVehicleHtml(vehicle) {
    let html = `
      <div class="card">
        <div class="card-header vehicle-view-heading">
          <h5 class="mb-0">
            Vehicle
            <button type="button" class="btn btn-danger" data-action="vehicle-close">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">ID: ${vehicle.id}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Start Location (Lat, Lon): ${
              vehicle.start_location.latitude
            }, ${vehicle.start_location.longitude}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">End Location (Lat, Lon): ${
              vehicle.end_location.latitude
            }, ${vehicle.end_location.longitude}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Capacity: [${vehicle.capacity}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Skills: [${vehicle.skills}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Time Window open: ${vehicle.tw_open}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Time Window close: ${vehicle.tw_close}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Speed Factor: ${vehicle.speed_factor}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Max Tasks: ${vehicle.max_tasks}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Project ID: ${vehicle.project_id}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Data: ${JSON.stringify(vehicle.data)}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Created At: ${vehicle.created_at}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Updated At: ${vehicle.updated_at}</p>
          </div>
          <div class="d-flex w-100 justify-content-center">
            <button class="btn btn-primary mx-2" data-action="vehicle-edit" data-id="${
              vehicle.id
            }">Edit</button>
            <button class="btn btn-danger mx-2" data-action="vehicle-delete" data-id="${
              vehicle.id
            }">Delete</button>
          </div>
        </div>
      </div>
    `;

    // return the html for the vehicle
    return html;
  }

  getVehicleFormHtml(vehicle) {
    let html = `
      <div class="card">
        <div class="card-header vehicle-view-heading">
          <h5 class="mb-0">
            Vehicle
            <button type="button" class="btn btn-danger" data-action="vehicle-close" data-id="${
              vehicle.id
            }">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <form>
            <input type="hidden" name="id" value="${vehicle.id}">
            <input type="hidden" name="project_id" value="${
              vehicle.project_id
            }">
            <div class="form-group">
              <label>Start Location (Lat, Lon)</label>
              <input type="text" class="form-control" name="start_location" value="${
                vehicle.start_location.latitude
              }, ${
      vehicle.start_location.longitude
    }" data-action="start_vehicle-location-change">
            </div>
            <div class="form-group">
              <label>End Location (Lat, Lon)</label>
              <input type="text" class="form-control" name="end_location" value="${
                vehicle.end_location.latitude
              }, ${
      vehicle.end_location.longitude
    }" data-action="end_vehicle-location-change">
            </div>
            <div class="form-group">
              <label>Capacity</label>
              <input type="text" class="form-control" name="capacity" value="${
                vehicle.capacity
              }">
            </div>
            <div class="form-group">
              <label>Skills</label>
              <input type="text" class="form-control" name="skills" value="${
                vehicle.skills
              }">
            </div>

            <div class="form-group">
              <label>Time Window (Open and Close)</label>
              <div class="input-group">
                <input type="datetime-local" class="form-control" name="tw_open" value="${
                  vehicle.tw_open
                }" step="1" style="font-size: 13px;">
                <span class="input-group-addon"></span>
                <input type="datetime-local" class="form-control" name="tw_close" value="${
                  vehicle.tw_close
                }" step="1" style="font-size: 13px;">
              </div>
            </div>
            <div class="form-group">
              <label>Speed Factor</label>
              <input type="number" class="form-control" name="speed_factor" min="0.01" max="200" value="${
                vehicle.speed_factor
              }" step="0.01">
            </div>
            <div class="form-group">
              <label>Max Tasks</label>
              <input type="number" class="form-control" name="max_tasks" min="0" max="2147483647" value="${
                vehicle.max_tasks
              }">
            </div>
            <div class="form-group">
              <label>Data</label>
              <input type="text" class="form-control" name="data" value='${JSON.stringify(
                vehicle.data
              )}'>
            </div>
            <div class="d-flex w-100 justify-content-center">
              <button type="button" class="btn btn-primary mx-2" data-action="vehicle-save" data-id="${
                vehicle.id
              }">Save</button>
              <button type="button" class="btn btn-warning mx-2" data-action="vehicle-edit" data-id="${
                vehicle.id
              }">Reset</button>
              <button type="button" class="btn btn-danger mx-2" data-action="vehicle-view" data-id="${
                vehicle.id
              }">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    `;

    // return the html for the vehicle
    return html;
  }

  // ------------------
  // BREAKS START
  // -------------------

  // get the time windows html for the break, time windows is a 2D array of start and end times
  getTimeWindowsHtml(time_windows) {
    let timeWindowsHtml = time_windows.map((timeWindow) => {
      return `
          <li>${timeWindow[0]} - ${timeWindow[1]}</li>
      `;
    });
    let timeWindows = timeWindowsHtml.join("");
    if (timeWindows == "") {
      timeWindows = "<li>No time window</li>";
    }
    return `
      <div>
        <ul class="mb-0">
        ${timeWindows}
        </ul>
      </div>
    `;
  }

  getTimeWindowsFormHtml(time_windows) {
    let timeWindowsHtml = time_windows.map((timeWindow) => {
      return `
        <div class="input-group">
          <input type="datetime-local" class="form-control" name="tw_open[]" value="${timeWindow[0]}" step="1" style="font-size: 12px;">
          <input type="datetime-local" class="form-control" name="tw_close[]" value="${timeWindow[1]}" step="1" style="font-size: 12px;">
        </div>
      `;
    });
    return `
      <div class="form-group">
      <label>Time Window (Open and Close)</label>
      ${timeWindowsHtml.join("")}
      <div class="text-center">
        <button type="button" class="btn btn-outline-primary" data-action="break-tw-form-create">
          <i class="fas fa-plus-circle"></i>
        </button>
        <button type="button" class="btn btn-outline-danger" data-action="break-tw-form-delete">
          <i class="fa-solid fa-trash"></i>
        </button>
      </div>
    </div>
    `;
  }

  getBreakHtml(break_) {
    let timeWindowsHtml = this.getTimeWindowsHtml(break_.time_windows);
    let html = `
      <div class="list-group-item flex-column align-items-start" data-id="${
        break_.id
      }">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">${break_.id}</h5>
        </div>
        <div>
          Service Time: ${break_.service}<br/>
          Data: ${JSON.stringify(break_.data)}<br/>
          Time Windows: <br/>
          ${timeWindowsHtml}
          <div style="text-align: center;">
            <button style="margin-right:5px;" type="button" class="btn btn-outline-warning" data-action="break-edit" data-id="${
              break_.id
            }">
              <i class="fa-solid fa-pen-to-square"></i>
            </button>
            <button type="button" class="btn btn-outline-danger" data-action="break-delete" data-id="${
              break_.id
            }" data-vehicle-id="${break_.vehicle_id}">
              <i class="fa-solid fa-trash"></i>
            </button>
          </div>
        </div>
      </div>
    `;

    // return the html for the break
    return html;
  }

  getBreaksHtml(breaksList, vehicleID) {
    let breakHtmlList = breaksList.map((break_) => {
      return this.getBreakHtml(break_);
    });

    let breakHtml = breakHtmlList.join("");
    if (breakHtml == "") {
      breakHtml = `
        <div class="list-group-item flex-column align-items-start" data-attribute="empty" data-id="">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-1">No Breaks</h5>
          </div>
        </div>
      `;
    }

    // return the html for the breaks
    let html = `
      <div class="card">
        <div class="card-header break-view-heading">
          <h5 class="mb-0">
            Break
          </h5>
        </div>
        <div class="card-body card-body-custom" style="max-height: 35vh;">
          ${breakHtml}
          <div class="text-center">
            <button type="button" class="btn btn-outline-primary" data-action="break-create" data-id="${vehicleID}">
              <i class="fas fa-plus-circle"></i>
            </button>
          </div>
        </div>
      </div>
    `;

    // return the html for the break
    return html;
  }

  getBreakFormHtml(break_) {
    let timeWindowsHtml = this.getTimeWindowsFormHtml(break_.time_windows);
    let html = `
      <form>
        <input type="hidden" name="id" value="${break_.id}">
        <input type="hidden" name="vehicle_id" value="${break_.vehicle_id}">
        <div class="form-group">
          <label>Service Time</label>
          <input type="time" class="form-control" name="service" value="${
            break_.service
          }" step="1">
        </div>
        <div class="form-group">
          <label>Data</label>
          <input type="text" class="form-control" name="data" value='${JSON.stringify(
            break_.data
          )}'>
        </div>
        ${timeWindowsHtml}
        <div class="d-flex w-100 justify-content-center">
          <button type="button" class="btn btn-primary mx-2" data-action="break-save" data-id="${
            break_.id
          }">Save</button>
          <button type="button" class="btn btn-danger mx-2" data-action="vehicle-view" data-id="${
            break_.vehicle_id
          }">Cancel</button>
        </div>
      </form>
    `;

    // return the html for the break
    return html;
  }

  // --------------------
  // BREAK END
  // --------------------

  selectVehicle(vehicleID) {
    this.deselectAll();
    let vehicleViewElement = document.querySelector(
      `[data-action="vehicle-view"][data-id="${vehicleID}"]`
    );
    vehicleViewElement.classList.add("active");

    // move the element into view
    vehicleViewElement.scrollIntoView({
      behavior: "smooth",
      block: "nearest",
    });
  }

  deselectAll() {
    // for all elements in query selector, remove their active class
    document.querySelectorAll(`.list-group-item.active`).forEach((element) => {
      element.classList.remove("active");
    });
  }

  getEmptyBreak(vehicle_id) {
    return {
      id: "",
      service: "00:00:00",
      data: {},
      time_windows: [],
      vehicle_id: vehicle_id,
    };
  }

  getEmptyVehicle() {
    // get todays date as string in the format YYYY-MM-DD with prepended 0 if required

    // get map center
    let coordinates = this.mapView.getCenter();
    let p_longitude = parseFloat(coordinates[1]) - 0.001;
    let d_longitude = parseFloat(coordinates[1]) + 0.001;
    return {
      id: "",
      start_location: {
        latitude: coordinates[0],
        longitude: p_longitude.toFixed(4),
      },
      end_location: {
        latitude: coordinates[0],
        longitude: d_longitude.toFixed(4),
      },
      capacity: "",
      skills: "",
      tw_open: ``,
      tw_close: ``,
      speed_factor: "1.00",
      max_tasks: "2147483647",
      project_id: this.projectID,
      data: {},
      created_at: "",
      updated_at: "",
    };
  }

  breakHandlers() {
    return {
      onBreakView: (vehicleID) => {
        // click on vehicle view button
        let vehicleViewElement = document.querySelector(
          `[data-action="vehicle-view"][data-id="${vehicleID}"]`
        );
        vehicleViewElement.click();
      },
      onBreakCreateClick: (vehicleID, sibling) => {
        // create an empty break form
        let break_ = this.getEmptyBreak(vehicleID);
        let breakFormHtml = this.getBreakFormHtml(break_);

        // if sibling has data-attribute=empty, then replace it with the form, else append the form
        if (sibling.dataset.attribute == "empty") {
          sibling.removeAttribute("data-attribute");
          sibling.innerHTML = breakFormHtml;
        } else {
          sibling.insertAdjacentHTML("afterend", breakFormHtml);
        }
      },
      onBreakEditClick: (breakID) => {
        let breakElement = document.querySelector(
          `.list-group-item[data-id="${breakID}"]`
        );

        this.breakAPI.getBreak(breakID).then((break_) => {
          let breakFormHtml = this.getBreakFormHtml(break_);
          breakElement.innerHTML = breakFormHtml;
        });
      },
    };
  }

  handlers() {
    return {
      onVehicleView: (vehicle) => {
        if (!vehicle.id) {
          this.setHtmlRight("");
          return;
        }
        // get the complete html for the vehicle
        let vehicleHtml = this.getCompleteVehicleHtml(vehicle);

        // set the html for the vehicle
        this.setHtmlRight(vehicleHtml);

        // call the break get api to get the break for the vehicle
        this.breakAPI.listBreaks(vehicle.id).then((breakList) => {
          let breakHtml = this.getBreaksHtml(breakList, vehicle.id);
          // append the break html to the vehicle html
          this.setHtmlRight(vehicleHtml + breakHtml);
        });

        // select the vehicle
        this.selectVehicle(vehicle.id);

        this.mapView.addVehicleMapPointer(
          vehicle.start_location.latitude,
          vehicle.start_location.longitude,
          vehicle.end_location.latitude,
          vehicle.end_location.longitude
        );
        this.mapView.deactivateMap();
      },
      onVehicleCreateClick: () => {
        this.deselectAll();

        const vehicle = this.getEmptyVehicle();

        // create the vehicle form html with empty vehicle
        let vehicleHtml = this.getVehicleFormHtml(vehicle);

        // set the html for the vehicle
        this.setHtmlRight(vehicleHtml);

        this.mapView.addVehicleMapPointer(
          vehicle.start_location.latitude,
          vehicle.start_location.longitude,
          vehicle.end_location.latitude,
          vehicle.end_location.longitude
        );
        this.mapView.activateMap();
        this.mapView.fitAllMarkers();
      },
      onVehicleEditClick: (vehicle) => {
        // get the complete html for the vehicle
        let vehicleHtml = this.getVehicleFormHtml(vehicle);

        // set the html for the vehicle
        this.setHtmlRight(vehicleHtml);

        this.mapView.addVehicleMapPointer(
          vehicle.start_location.latitude,
          vehicle.start_location.longitude,
          vehicle.end_location.latitude,
          vehicle.end_location.longitude
        );
        this.mapView.activateMap();
      },
      onVehicleSave: (vehicle, newVehicles) => {
        this.vehicles = newVehicles;
        this.render();
        this.handlers().onVehicleView(vehicle);
      },
      onVehicleDelete: (newVehicles) => {
        this.deselectAll();
        this.setHtmlRight("");
        this.vehicles = newVehicles;
        this.render();
        this.mapView.removeMapPointers();
        this.mapView.deactivateMap();
        this.mapView.fitAllMarkers();
      },
      onVehicleClose: () => {
        this.deselectAll();
        this.setHtmlRight("");
        this.mapView.removeMapPointers();
        this.mapView.deactivateMap();
        this.mapView.fitAllMarkers();
      },
    };
  }
}
