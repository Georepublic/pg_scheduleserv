import Parser from "../utils/Parser.js";
import Toast from "../utils/Toast.js";
import BaseAPI from "./BaseAPI.js";

export default class ShipmentAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  parseShipment(data) {
    return {
      id: data.id,
      p_location: Parser.parseLocation(data.p_location),
      d_location: Parser.parseLocation(data.d_location),
      p_setup: Parser.parseDuration(data.p_setup),
      d_setup: Parser.parseDuration(data.d_setup),
      p_service: Parser.parseDuration(data.p_service),
      d_service: Parser.parseDuration(data.d_service),
      amount: Parser.parseAmount(data.amount),
      skills: Parser.parseAmount(data.skills),
      priority: Parser.parsePriority(data.priority),
      project_id: data.project_id,
      data: Parser.parseJSON(data.data),
      p_time_windows: Parser.parseTimeWindows(
        data["p_tw_open[]"],
        data["p_tw_close[]"]
      ),
      d_time_windows: Parser.parseTimeWindows(
        data["d_tw_open[]"],
        data["d_tw_close[]"]
      ),
    };
  }

  saveShipment(data) {
    try {
      data = this.parseShipment(data);
    } catch (e) {
      return Promise.reject(e);
    }

    if (data["id"]) {
      return this.baseAPI.patch(`/shipments/${data["id"]}`, data);
    } else {
      return this.baseAPI.post(
        `/projects/${data["project_id"]}/shipments`,
        data
      );
    }
  }

  listShipments(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/shipments`);
  }

  getShipment(shipmentID) {
    return this.baseAPI.get(`/shipments/${shipmentID}`);
  }

  deleteShipment(shipmentID) {
    return this.baseAPI.delete(`/shipments/${shipmentID}`);
  }
}
