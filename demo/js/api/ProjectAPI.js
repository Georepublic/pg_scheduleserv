import BaseAPI from "./BaseAPI.js";

export default class ProjectAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  listProjects() {
    return this.baseAPI.get(`/projects`);
  }

  createProject(data) {
    return this.baseAPI.post(`/projects`, data);
  }

  getProject(projectID) {
    return this.baseAPI.get(`/projects/${projectID}`);
  }

  editProject(projectID, data) {
    return this.baseAPI.patch(`/projects/${projectID}`, data);
  }

  deleteProject(projectID) {
    return this.baseAPI.delete(`/projects/${projectID}`);
  }
}
