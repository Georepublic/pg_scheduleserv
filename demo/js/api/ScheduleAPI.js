import BaseAPI from "./BaseAPI.js";

export default class ScheduleAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  createSchedule(projectID) {
    return this.baseAPI.post(`/projects/${projectID}/schedule`);
  }

  getSchedule(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/schedule`);
  }

  getScheduleIcal(projectID) {
    return this.baseAPI.getIcal(`/projects/${projectID}/schedule`);
  }

  getJobSchedule(jobID) {
    return this.baseAPI.get(`/jobs/${jobID}/schedule`);
  }

  getShipmentSchedule(shipmentID) {
    return this.baseAPI.get(`/shipments/${shipmentID}/schedule`);
  }

  getVehicleSchedule(vehicleID) {
    return this.baseAPI.get(`/vehicles/${vehicleID}/schedule`);
  }

  deleteSchedule(projectID) {
    return this.baseAPI.delete(`/projects/${projectID}/schedule`);
  }
}
