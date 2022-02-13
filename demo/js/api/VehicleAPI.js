import Parser from "../utils/Parser.js";
import Toast from "../utils/Toast.js";
import BaseAPI from "./BaseAPI.js";

export default class VehicleAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  parseVehicle(data) {
    return {
      id: data.id,
      start_location: Parser.parseLocation(data.start_location),
      end_location: Parser.parseLocation(data.end_location),
      capacity: Parser.parseAmount(data.capacity),
      skills: Parser.parseAmount(data.skills),
      tw_open: Parser.parseDateTime(data.tw_open),
      tw_close: Parser.parseDateTime(data.tw_close),
      speed_factor: Parser.parseSpeedFactor(data.speed_factor),
      max_tasks: Parser.parseMaxTasks(data.max_tasks),
      project_id: data.project_id,
      data: Parser.parseJSON(data.data),
    };
  }

  saveVehicle(data) {
    try {
      data = this.parseVehicle(data);
    } catch (e) {
      return Promise.reject(e);
    }

    if (data["id"]) {
      return this.baseAPI.patch(`/vehicles/${data["id"]}`, data);
    } else {
      return this.baseAPI.post(
        `/projects/${data["project_id"]}/vehicles`,
        data
      );
    }
  }

  listVehicles(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/vehicles`);
  }

  getVehicle(vehicleID) {
    return this.baseAPI.get(`/vehicles/${vehicleID}`);
  }

  deleteVehicle(vehicleID) {
    return this.baseAPI.delete(`/vehicles/${vehicleID}`);
  }
}
