import ProjectAPI from "../api/ProjectAPI.js";
import ProjectHandler from "../handlers/ProjectHandler.js";
import AbstractView from "./AbstractView.js";

export default class ProjectListView extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Projects");
    this.projectAPI = new ProjectAPI();

    // initialize the view with loading icon before data is loaded
    this.setHtml(this.getLoadingHtml());
    this.setHeading("Projects");
    this.setSubHeading("");

    this.durationType = {
      euclidean: "Euclidean",
      valhalla: "Valhalla",
      osrm: "OSRM",
    };

    this.projectAPI
      .listProjects()
      .then((projects) => {
        this.projects = projects;
      })
      .then(() => {
        this.handler = new ProjectHandler(
          this.projects,
          this.getEmptyProject(),
          this.handlers()
        );
        this.render(this.projects);
      });
  }

  render(projects) {
    this.setHtml(this.getCompleteHtml(projects));
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

  // get the html heading and list of projects by calling the getProjectHtml function
  getCompleteHtml(projects) {
    return `
      <div class="list-group">${this.getProjectHtml(projects)}</div>
      ${this.getProjectCreateButton()}
    `;
  }

  getProjectCreateButton() {
    return `
    <div class="text-center">
      <button class="btn btn-outline-primary" data-action="project-create">
        <i class="fas fa-plus-circle"></i>
      </button>
    </div>
    `;
  }

  getProjectHtml(projects) {
    var html = "";
    projects.forEach((project) => {
      html += `
        <div class="list-group-item flex-column align-items-start" data-id="${
          project.id
        }">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-1">${project.name}</h5>
            <small class="text-muted">Created: ${project.created_at}</small>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Duration Calculation: ${
              this.durationType[project.duration_calc]
            }</p>
          </div>
          <small>
            Max Shift in schedule: ${project.max_shift}<br/>
            Exploration Level: ${project.exploration_level}, Timeout: ${
        project.timeout
      }
            <p style="float: right">
            <a style="margin-right:5px;" type="button" class="btn btn-outline-info" href="/projects/${
              project.id
            }" data-link>
              <i class="fas fa-folder-open"></i>
            </a>
            <button style="margin-right:5px;" type="button" class="btn btn-outline-warning" data-action="project-edit" data-id="${
              project.id
            }">
              <i class="fa-solid fa-pen-to-square"></i>
            </button>
            <button type="button" class="btn btn-outline-danger" data-action="project-delete" data-id="${
              project.id
            }">
              <i class="fa-solid fa-trash"></i>
            </button>
            </p>
          </small>
        </div>
      `;
    });

    // if html is empty, return empty project
    if (html === "") {
      html = `
        <div class="list-group-item flex-column align-items-start" data-attribute="empty">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-3">No projects found</h5>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Click below to create a new project</p>
          </div>
        </div>
      `;
    }

    return html;
  }

  getProjectForm(project) {
    let durationTypesHtml = ["euclidean", "valhalla", "osrm"].map((type) => {
      let selected = type == project.duration_calc ? "selected" : "";
      return `
        <option value="${type}" ${selected}>${this.durationType[type]}</option>
      `;
    });

    return `
      <form>
        <div class="form-group">
          <label for="name">Name</label>
          <input type="hidden" id="id" name="id" value="${project.id}">
          <input type="text" class="form-control" name="name" value="${
            project.name
          }">
        </div>
        <div class="form-group">
          <label for="duration_calc">Duration calculation</label>
          <select class="form-select" name="duration_calc">
            ${durationTypesHtml.join("")}
          </select>
        </div>
        <div class="form-group">
          <label for="max_shift">Max Shift</label>
          <input type="time" class="form-control" name="max_shift" value="${
            project.max_shift
          }" step="1">
        </div>
        <div class="form-group">
          <label for="exploration_level">Exploration Level</label>
          <input type="number" class="form-control" name="exploration_level" min="0" max="5" value="${
            project.exploration_level
          }">
        </div>
        <div class="form-group">
          <label for="timeout">Timeout</label>
          <input type="time" class="form-control" name="timeout" value="${
            project.timeout
          }" step="1">
        </div>
        <button type="button" class="btn btn-outline-success" data-action="project-save">Save</button>
        <button type="button" class="btn btn-outline-success" data-action="project-reset">Reset</button>
      </form>
    `;
  }

  getEmptyProject() {
    return {
      id: "",
      name: "Sample Project",
      duration_calc: "euclidean",
      max_shift: "00:30:00",
      exploration_level: "5",
      timeout: "00:10:00",
    };
  }

  handlers() {
    return {
      onProjectResetClick: (projects) => {
        this.render(projects);
      },
      onProjectCreateClick: (el) => {
        el.innerHTML = this.getProjectForm(
          // empty project
          this.getEmptyProject()
        );
      },
      onProjectEditClick: (project) => {
        const el = document.querySelector(
          `.list-group-item[data-id="${project.id}"]`
        );
        el.innerHTML = this.getProjectForm(project);
      },
      onProjectSave: (newProjects) => {
        this.render(newProjects);
      },
      onProjectDelete: (newProjects) => {
        this.render(newProjects);
      },
    };
  }
}
