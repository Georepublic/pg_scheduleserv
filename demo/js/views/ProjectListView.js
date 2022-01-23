import ProjectAPI from "../api/ProjectAPI.js";
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

    this.refresh();
  }

  refresh() {
    this.projectAPI
      .listProjects()
      .then((projects) => {
        this.setHtml(this.getCompleteHtml(projects));
      })
      .then(() => {
        this.handleEditButton();
        this.handleDeleteButton();
        this.handleCreateButton();
      });
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
      <button class="btn btn-outline-primary project__create">
        <i class="fas fa-plus-circle"></i>
      </button>
    </div>
    `;
  }

  getProjectHtml(projects) {
    var html = "";
    projects.forEach((project) => {
      html += `
        <div id="project-${project.id}" class="list-group-item flex-column align-items-start">
          <div class="d-flex w-100 justify-content-between">
            <h5 class="mb-1">${project.name}</h5>
            <small class="text-muted">Created: ${project.created_at}</small>
          </div>
          <small>
            Exploration Level: ..., Timeout: ...
            <p style="float: right">
            <a style="margin-right:5px;" type="button" class="btn btn-outline-info" href="/projects/${project.id}" data-link>
              <i class="fas fa-folder-open"></i>
            </a>
            <button style="margin-right:5px;" type="button" class="btn btn-outline-warning project__edit">
              <i class="fa-solid fa-pen-to-square"></i>
            </button>
            <button type="button" class="btn btn-outline-danger project__delete">
              <i class="fa-solid fa-trash"></i>
            </button>
            </p>
          </small>
        </div>
      `;
    });
    return html;
  }

  // When the create button is clicked, append a div to existing div with the form
  handleCreateButton() {
    const btnCreateProject = document.querySelector(".project__create");
    btnCreateProject.addEventListener("click", (e) => {
      const div = e.target.closest("div");
      div.innerHTML = this.getProjectForm(
        // empty project
        {
          id: "",
          name: "",
          exploration_level: "",
          timeout: "",
        }
      );
      this.handleSaveButton(div);
      this.handleResetButton(div);
    });
  }

  // Edit the div when edit button is clicked and transform it into a form
  handleEditButton() {
    const btnEditProject = document.querySelectorAll(".project__edit");
    btnEditProject.forEach((btn) => {
      btn.addEventListener("click", (e) => {
        const div = e.target.closest("div");
        const id = div.id.split("-")[1];
        this.projectAPI
          .getProject(id)
          .then((project) => {
            div.innerHTML = this.getProjectForm(project);
          })
          .then(() => {
            this.handleSaveButton(div);
            this.handleResetButton(div);
          });
      });
    });
  }

  getProjectForm(project) {
    return `
      <form>
        <div class="form-group">
          <label for="name">Name</label>
          <input type="hidden" id="id" name="id" value="${project.id}">
          <input type="text" class="form-control" id="name" name="name" value="${project.name}">
        </div>
        <div class="form-group">
          <label for="explorationLevel">Exploration Level</label>
          <input type="text" class="form-control" id="explorationLevel" name="explorationLevel" value="${project.exploration_level}">
        </div>
        <div class="form-group">
          <label for="timeout">Timeout</label>
          <input type="text" class="form-control" id="timeout" name="timeout" value="${project.timeout}">
        </div>
        <button type="button" class="btn btn-outline-success project__save">Save</button>
        <button type="button" class="btn btn-outline-success project__reset">Reset</button>
      </form>
    `;
  }

  handleSaveButton(div) {
    const btnSaveProject = div.querySelector(".project__save");
    btnSaveProject.addEventListener("click", (e) => {
      const form = e.target.closest("form");
      const formData = new FormData(form);
      const project = {};
      for (const [key, value] of formData.entries()) {
        project[key] = value;
      }
      const id = project["id"];

      this.projectAPI
        .saveProject(project)
        .then((project) => {
          if (id) {
            // update the project
            div.outerHTML = this.getProjectHtml([project]);
          } else {
            // append the new project to the list
            const projects = document.querySelector(".list-group");
            projects.innerHTML += this.getProjectHtml([project]);

            // replace the div with create button and handler
            div.outerHTML = this.getProjectCreateButton();
            this.handleCreateButton();
          }
        })
        .then(() => {
          this.handleEditButton();
          this.handleDeleteButton();
        });
    });
  }

  handleResetButton(div) {
    const btnResetProject = div.querySelector(".project__reset");
    btnResetProject.addEventListener("click", (e) => {
      const form = e.target.closest("form");
      const formData = new FormData(form);
      const id = formData.get("id");
      if (id) {
        // reset the project
        this.projectAPI
          .getProject(id)
          .then((project) => {
            div.outerHTML = this.getProjectHtml([project]);
          })
          .then(() => {
            this.handleEditButton();
            this.handleDeleteButton();
          });
      } else {
        // remove the new project create form
        div.outerHTML = this.getProjectCreateButton();
        this.handleCreateButton();
      }
    });
  }

  handleDeleteButton() {
    const btnDeleteProject = document.querySelectorAll(".project__delete");
    btnDeleteProject.forEach((btn) => {
      btn.addEventListener("click", (e) => {
        const div = e.target.closest("div");
        const id = div.id.split("-")[1];
        this.projectAPI
          .deleteProject(id)
          .then(() => {
            div.remove();
          })
          .catch((err) => {
            console.log(err);
          });
      });
    });
  }
}
