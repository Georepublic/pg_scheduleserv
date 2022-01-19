import BaseAPI from "./BaseAPI.js";

export default class ShipmentAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  listShipments(projectID) {
    return this.baseAPI.get(`/projects/${projectID}/shipments`);
  }

  createShipment(projectID, data) {
    return this.baseAPI.post(`/projects/${projectID}/shipments`, data);
  }

  getShipment(shipmentID) {
    return this.baseAPI.get(`/shipments/${shipmentID}`);
  }

  updateShipment(shipmentID, data) {
    return this.baseAPI.patch(`/shipments/${shipmentID}`, data);
  }

  deleteShipment(shipmentID) {
    return this.baseAPI.delete(`/shipments/${shipmentID}`);
  }
}
