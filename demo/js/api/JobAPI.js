import Parser from "../utils/Parser.js";
import BaseAPI from "./BaseAPI.js";

export default class JobAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  parseJob(data) {
    return {
      id: data.id,
      location: Parser.parseLocation(data.location),
      setup: Parser.parseDuration(data.setup),
      service: Parser.parseDuration(data.service),
      delivery: Parser.parseAmount(data.delivery),
      pickup: Parser.parseAmount(data.pickup),
      skills: Parser.parseAmount(data.skills),
      priority: Parser.parsePriority(data.priority),
      project_id: data.project_id,
      data: Parser.parseJSON(data.data),
    };
  }

  saveJob(data) {
    try {
      data = this.parseJob(data);
    } catch (e) {
      return Promise.reject(e);
    }

    if (data["id"]) {
      return this.baseAPI.patch(`/jobs/${data["id"]}`, data);
    } else {
      return this.baseAPI.post(`/projects/${data["project_id"]}/jobs`, data);
    }
  }

  listJobs(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/jobs`);
  }

  getJob(jobID) {
    return this.baseAPI.get(`/jobs/${jobID}`);
  }

  deleteJob(jobID) {
    return this.baseAPI.delete(`/jobs/${jobID}`);
  }
}
