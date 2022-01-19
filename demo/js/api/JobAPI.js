import BaseAPI from "./BaseAPI.js";

export default class JobAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  processData(data) {
    // convert location to {latitude, longitude} by removing () if present and splitting by comma and trim
    const location = data["location"].replace(/[()]/g, "").split(",").map(function(item) {
      return item.trim();
    });
    data["location"] = {
      // convert string to float64
      latitude: parseFloat(location[0]),
      longitude: parseFloat(location[1])
    }
    data["delivery"] = data["delivery"].split(",").map(function(item) {
      return parseInt(item);
    });
    data["pickup"] = data["pickup"].split(",").map(function (item) {
      return parseInt(item);
    });
    data["skills"] = data["skills"].split(",").map(function (item) {
      return parseInt(item);
    });

    data["priority"] = parseInt(data["priority"]);

    return data;
  }

  saveJob(data) {
    data = this.processData(data);

    if (data["id"]) {
      return this.baseAPI.patch(`/jobs/${data["id"]}`, data);
    } else {
      return this.baseAPI.post(`/projects/${data["project_id"]}/jobs`, data);
    }
  }

  getJob(jobID) {
    return this.baseAPI.get(`/jobs/${jobID}`);
  }

  deleteJob(jobID) {
    return this.baseAPI.delete(`/jobs/${jobID}`);
  }
}
