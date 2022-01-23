import AbstractView from "./AbstractView.js";
import ProjectAPI from "../api/ProjectAPI.js";
import JobView from "./JobView.js";
import ShipmentView from "./ShipmentView.js";
import MapView from "./MapView.js";
import VehicleView from "./VehicleView.js";

export default class ProjectView extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Projects");

    this.setHtml(this.getLoadingHtml());
    this.setSubHeading("");

    this.projectAPI = new ProjectAPI();

    this.projectAPI
      .getProject(params.id)
      .then((project) => {
        this.setHeading(project.name);
        this.setHtml(this.getHtml(project));
      })
      .then(() => {
        this.mapView = new MapView();
        return this.mapView.createMap();
      })
      .then(() => {
        // call JobAPI to get jobs and pass them to JobView
        return this.projectAPI.getJobs(params.id).then((jobs) => {
          var jobView = new JobView({
            jobs: jobs,
            projectID: params.id,
            mapView: this.mapView,
          });
          jobView.render();
        });
      })
      .then(() => {
        // call ShipmentAPI to get shipments and pass them to ShipmentView
        return this.projectAPI.getShipments(params.id).then((shipments) => {
          var shipmentView = new ShipmentView({
            shipments: shipments,
            projectID: params.id,
            mapView: this.mapView,
          });
          shipmentView.render();
        });
      })
      .then(() => {
        // call VehicleAPI to get vehicles and pass them to VehicleView
        return this.projectAPI.getVehicles(params.id).then((vehicles) => {
          var vehicleView = new VehicleView({
            vehicles: vehicles,
            projectID: params.id,
            mapView: this.mapView,
          });
          vehicleView.render();
        });
      });

    this.setHeading("Projects");
  }

  getLoadingHtml() {
    return `
      <div class="list-group">
        <div class="list-group-item flex-column align-items-start">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-1">Loading...</h5>
          </div>
        </div>
      </div>
    `;
  }

  getHtml(project) {
    var html = `
    <div class="list-group">
      <div id="project-${project.id}" class="list-group-item flex-column align-items-start">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">${project.name}</h5>
          <small class="text-muted">Created: ${project.created_at}</small>
        </div>
        <small>Exploration Level: ..., Timeout: ...</small>
      </div>
    </div>
    <div id="map"></div>
    `;
    return html;
  }
}
