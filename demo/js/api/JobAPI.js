import BaseAPI from "./BaseAPI.js";

export default class JobAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  listJobs(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/jobs`);
  }

  createJob(projectID, data) {
    return this.baseAPI.post(`/projects/${projectID}/jobs`, data);
  }

  getJob(jobID) {
    return this.baseAPI.get(`/${jobID}`);
  }

  editJob(jobID, data) {
    return this.baseAPI.patch(`/${jobID}`, data);
  }

  deleteJob(jobID) {
    return this.baseAPI.delete(`/${jobID}`);
  }
}
