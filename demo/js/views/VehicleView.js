import VehicleAPI from "../api/VehicleAPI.js";
import VehicleHandler from "../handlers/VehicleHandler.js";
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

    this.vehicleLeftDiv = document.createElement("div");
    document.querySelector("#app-left").appendChild(this.vehicleLeftDiv);

    this.handler = new VehicleHandler(
      params.vehicles,
      this.getEmptyVehicle(),
      this.handlers()
    );
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
              <label>Time Window open</label>
              <input type="datetime-local" class="form-control" name="tw_open" value="${
                vehicle.tw_open
              }" step="1">
            </div>
            <div class="form-group">
              <label>Time Window close</label>
              <input type="datetime-local" class="form-control" name="tw_close" value="${
                vehicle.tw_close
              }" step="1">
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

  getEmptyVehicle() {
    // get todays date as string in the format YYYY-MM-DD with prepended 0 if required
    let today = new Date();
    let dd = today.getDate();
    let mm = today.getMonth() + 1; //January is 0!
    let yyyy = today.getFullYear();

    if (dd < 10) {
      dd = "0" + dd;
    }
    if (mm < 10) {
      mm = "0" + mm;
    }
    const date = `${yyyy}-${mm}-${dd}`;

    // get map center
    let coordinates = this.mapView.getCenter();
    return {
      id: "",
      start_location: {
        latitude: coordinates[0],
        longitude: coordinates[1],
      },
      end_location: {
        latitude: coordinates[0],
        longitude: coordinates[1],
      },
      capacity: "",
      skills: "",
      tw_open: `${date}T07:00:00`,
      tw_close: `${date}T19:00:00`,
      speed_factor: "1.00",
      max_tasks: "2147483647",
      project_id: this.projectID,
      data: {},
      created_at: "",
      updated_at: "",
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
