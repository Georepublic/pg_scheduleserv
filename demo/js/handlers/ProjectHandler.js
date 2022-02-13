import ProjectAPI from "../api/ProjectAPI.js";

export default class ProjectHandler {
  constructor(
    projects,
    emptyProject,
    {
      onProjectResetClick,
      onProjectCreateClick,
      onProjectEditClick,
      onProjectDelete,
      onProjectSave,
    }
  ) {
    // get the projects from the params
    this.projects = projects;
    this.emptyProject = emptyProject;
    this.projectAPI = new ProjectAPI();

    this.main = document.querySelector("#main");
    this.schedule = document.querySelector("#schedule");

    this.handleProjectResetClick(onProjectResetClick);
    this.handleProjectCreateClick(onProjectCreateClick);
    this.handleProjectEditClick(onProjectEditClick);
    this.handleProjectDelete(onProjectDelete);
    this.handleProjectSave(onProjectSave);
  }

  // get project from id
  getProject(projectID) {
    if (!projectID) {
      return this.emptyProject;
    }
    return this.projects.find((project) => {
      return project.id === projectID;
    });
  }

  handleProjectResetClick(onProjectResetClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="project-reset"]`);
      if (el) {
        onProjectResetClick(this.projects);
      }
    });
  }

  handleProjectCreateClick(onProjectCreateClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="project-create"]`);
      if (el) {
        onProjectCreateClick(el.closest("div"));
      }
    });
  }

  handleProjectEditClick(onProjectEditClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="project-edit"]`);
      if (el) {
        let projectID = el.dataset.id;
        onProjectEditClick(this.getProject(projectID));
      }
    });
  }

  handleProjectDelete(onProjectDelete) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="project-delete"]`);
      if (el) {
        let projectID = el.dataset.id;
        // call the project api to delete the project
        this.projectAPI.deleteProject(projectID).then(() => {
          // remove the project from the list
          this.projects = this.projects.filter((project) => {
            return project.id !== projectID;
          });
          // call the onProjectDelete callback
          onProjectDelete(this.projects);
        });
      }
    });
  }

  handleProjectSave(onProjectSave) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="project-save"]`);
      if (el) {
        const form = el.closest("form");
        const formData = new FormData(form);
        const project = {};
        for (const [key, value] of formData.entries()) {
          // if key contains [] then it is an array
          if (key.includes("[]")) {
            if (!project[key]) {
              project[key] = [];
            }
            project[key].push(value);
          } else {
            project[key] = value;
          }
        }
        const id = project["id"];

        this.projectAPI.saveProject(project).then((project) => {
          // edit the project in the list, or append a new project to the list depending on the id
          if (id) {
            // update the project
            this.projects = this.projects.map((oldProject) => {
              if (oldProject.id === project.id) {
                return project;
              }
              return oldProject;
            });
          } else {
            // append the new project to the list
            this.projects.push(project);
          }
          onProjectSave(this.projects);
        });
      }
    });
  }
}
