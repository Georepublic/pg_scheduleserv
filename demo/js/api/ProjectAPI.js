import Parser from "../utils/Parser.js";
import BaseAPI from "./BaseAPI.js";

export default class ProjectAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  parseProject(data) {
    return {
      id: data.id,
      name: data.name,
      duration_calc: data.duration_calc,
      max_shift: Parser.parseDuration(data.max_shift),
      exploration_level: parseInt(data.exploration_level),
      timeout: Parser.parseDuration(data.timeout),
    };
  }

  listProjects() {
    return this.baseAPI.get(`/projects`);
  }

  saveProject(data) {
    try {
      data = this.parseProject(data);
    } catch (e) {
      return Promise.reject(e);
    }

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
