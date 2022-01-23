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
      return this.baseAPI.post(`/projects/${data["project_id"]}/shipments`, data);
    }
  }

  getShipment(shipmentID) {
    return this.baseAPI.get(`/shipments/${shipmentID}`);
  }

  deleteShipment(shipmentID) {
    return this.baseAPI.delete(`/shipments/${shipmentID}`);
  }
}
