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
