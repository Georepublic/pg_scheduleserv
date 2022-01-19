import Toast from "./Toast.js";

export default class Json {
  // function to parse json and return error
  static parseJson(json) {
    try {
      return JSON.parse(json);
    } catch (e) {
      Toast.error("Invalid Data field: " + e.message);
      throw e;
    }
  }
}
