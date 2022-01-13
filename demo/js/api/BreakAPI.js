import BaseAPI from "./BaseAPI.js";

export default class BreakAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  listBreaks(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/breaks`);
  }

  createBreak(projectID, data) {
    return this.baseAPI.post(`/projects/${projectID}/breaks`, data);
  }

  getBreak(breakID) {
    return this.baseAPI.get(`/${breakID}`);
  }

  editBreak(breakID, data) {
    return this.baseAPI.patch(`/${breakID}`, data);
  }

  deleteBreak(breakID) {
    return this.baseAPI.delete(`/${breakID}`);
  }
}
