import BaseAPI from "./BaseAPI.js";

export default class ScheduleAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  createSchedule(projectID, type) {
    var queryParam = "";
    if (type == "fresh") {
      queryParam = "?fresh=true";
    } else if (type == "normal") {
      queryParam = "?fresh=false";
    }
    return this.baseAPI.post(`/projects/${projectID}/schedule${queryParam}`);
  }

  getSchedule(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/schedule`);
  }

  getScheduleIcal(projectID) {
    return this.baseAPI.getIcal(`/projects/${projectID}/schedule`);
  }

  getJobScheduleIcal(jobID) {
    return this.baseAPI.getIcal(`/jobs/${jobID}/schedule`);
  }

  getShipmentScheduleIcal(shipmentID) {
    return this.baseAPI.getIcal(`/shipments/${shipmentID}/schedule`);
  }

  getVehicleScheduleIcal(vehicleID) {
    return this.baseAPI.getIcal(`/vehicles/${vehicleID}/schedule`);
  }

  deleteSchedule(projectID) {
    return this.baseAPI.delete(`/projects/${projectID}/schedule`);
  }
}
