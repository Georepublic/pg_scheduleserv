import AbstractView from "./AbstractView.js";
import ProjectAPI from "../api/ProjectAPI.js";
import LocationAPI from "../api/LocationAPI.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Projects");

    this.setHtml(this.getLoadingHtml());

    this.projectAPI = new ProjectAPI();
    this.locationAPI = new LocationAPI();

    this.projectAPI
      .getProject(params.id)
      .then((project) => {
        this.setHtml(this.getHtml(project));
      })
      .then(() => {
        this.createMap();
      });
  }

  getLoadingHtml() {
    return `
      <div class="heading"><h2>Project</h2></div>
      <div class="list-group">
        <div class="list-group-item flex-column align-items-start">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-1">Loading...</h5>
          </div>
        </div>
      </div>
    `;
  }

  createMap() {
    var latitude = 35.7127;
    var longitude = 139.762;

    this.locationAPI.getLocation().then((location) => {
      latitude = location.latitude;
      longitude = location.longitude;
    }).catch((error) => {
      console.log(error);
    }).then(() => {
      // create new open layer map
      var map = new ol.Map({
        target: "map",
        layers: [
          new ol.layer.Tile({
            source: new ol.source.OSM(),
          }),
        ],
        view: new ol.View({
          center: ol.proj.fromLonLat([longitude, latitude]),
          zoom: 12,
        }),
      });
    })
  }

  getHtml(project) {
    var html = `
    <div class="heading"><h2>${project.name}</h2></div>
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
