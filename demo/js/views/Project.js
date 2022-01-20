import AbstractView from "./AbstractView.js";
import ProjectAPI from "../api/ProjectAPI.js";
import JobView from "./JobView.js";
import MapView from "./MapView.js";

export default class extends AbstractView {
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
        this.mapView.createMap().then(() => {
          // call JobAPI to get jobs and pass them to JobView
          this.projectAPI.getJobs(params.id).then((jobs) => {
            var jobView = new JobView({
              jobs: jobs,
              projectID: params.id,
              mapView: this.mapView,
            });
            jobView.render();
          });
        });
      })

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
