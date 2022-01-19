import BaseAPI from "./BaseAPI.js";

export default class VehicleAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  listVehicles(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/vehicles`);
  }

  createVehicle(projectID, data) {
    return this.baseAPI.post(`/projects/${projectID}/vehicles`, data);
  }

  getVehicle(vehicleID) {
    return this.baseAPI.get(`/vehicles/${vehicleID}`);
  }

  updateVehicle(vehicleID, data) {
    return this.baseAPI.patch(`/vehicles/${vehicleID}`, data);
  }

  deleteVehicle(vehicleID) {
    return this.baseAPI.delete(`/vehicles/${vehicleID}`);
  }
}
