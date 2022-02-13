import Parser from "../utils/Parser.js";
import BaseAPI from "./BaseAPI.js";

export default class BreakAPI {
  constructor() {
    this.baseAPI = new BaseAPI();
  }

  parseBreak(data) {
    return {
      id: data.id,
      service: Parser.parseDuration(data.service),
      data: Parser.parseJSON(data.data),
      time_windows: Parser.parseTimeWindows(
        data["tw_open[]"],
        data["tw_close[]"]
      ),
      vehicle_id: data.vehicle_id,
    };
  }

  saveBreak(data) {
    try {
      data = this.parseBreak(data);
    } catch (e) {
      return Promise.reject(e);
    }

    if (data["id"]) {
      return this.baseAPI.patch(`/breaks/${data["id"]}`, data);
    } else {
      return this.baseAPI.post(`/vehicles/${data["vehicle_id"]}/breaks`, data);
    }
  }

  listBreaks(vehicleID) {
    return this.baseAPI.get(`/vehicles/${vehicleID}/breaks`);
  }

  getBreak(breakID) {
    return this.baseAPI.get(`/breaks/${breakID}`);
  }

  deleteBreak(breakID) {
    return this.baseAPI.delete(`/breaks/${breakID}`);
  }
}
