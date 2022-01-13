import ProjectAPI from "../api/ProjectAPI.js";
import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Projects");

    const projectAPI = new ProjectAPI();
    projectAPI
      .listProjects()
      .then((projects) => {
        this.setHtml(this.getHtml(projects));
      });
  }

  getHtml(projects) {
    var html = "";
    projects.forEach((project) => {
      html += `
        <a href="/projects/${project.id}" class="list-group-item list-group-item-action flex-column align-items-start">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-1">${project.name}</h5>
            <small class="text-muted">Created: ${project.created_at}</small>
          </div>
          <small>Exploration Level: ..., Timeout: ...</small>
        </a>
      `;
    });
    return `
      <div class="heading"><h2>Projects</h2></div>
      <div class="list-group">${html}</div>
    `;
  }
}
