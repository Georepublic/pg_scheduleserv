import BaseAPI from "./BaseAPI.js";

export default class ProjectAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  listProjects() {
    return this.baseAPI.get(`/projects`);
  }

  saveProject(data) {
    if (data["id"]) {
      return this.baseAPI.patch(`/projects/${data["id"]}`, data);
    } else {
      return this.baseAPI.post(`/projects`, data);
    }
  }

  getProject(projectID) {
    return this.baseAPI.get(`/projects/${projectID}`);
  }

  deleteProject(projectID) {
    return this.baseAPI.delete(`/projects/${projectID}`);
  }
}
