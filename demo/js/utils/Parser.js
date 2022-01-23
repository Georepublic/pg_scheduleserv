import Toast from "./Toast.js";

export default class Parser {
  // parse and validate duration to be of hh:mm:ss format, with proper time
  static parseDuration(duration) {
    const durationArray = duration.split(":").map(function (item) {
      return item.trim();
    });
    if (
      durationArray.length !== 3 ||
      isNaN(parseInt(durationArray[0])) ||
      isNaN(parseInt(durationArray[1])) ||
      isNaN(parseInt(durationArray[2]))
    ) {
      Toast.error(
        "Invalid duration: " + duration + " must be in hh:mm:ss format"
      );
      throw new Error(
        "Invalid duration: " + duration + " must be in hh:mm:ss format"
      );
    }
    // also check for hours, minutes, seconds
    if (parseInt(durationArray[1]) > 59 || parseInt(durationArray[2]) > 59) {
      Toast.error(
        "Invalid duration: " + duration + " must be in hh:mm:ss format"
      );
      throw new Error(
        "Invalid duration: " + duration + " must be in hh:mm:ss format"
      );
    }
    return durationArray.join(":");
  }

  // parse and validate date to be of YYYY-MM-DD format, with proper date
  static parseDate(date) {
    const dateArray = date.split("-").map(function (item) {
      return item.trim();
    });
    if (
      dateArray.length !== 3 ||
      isNaN(parseInt(dateArray[0])) ||
      isNaN(parseInt(dateArray[1])) ||
      isNaN(parseInt(dateArray[2]))
    ) {
      Toast.error("Invalid date: " + date + " must be in YYYY-MM-DD format");
      throw new Error(
        "Invalid date: " + date + " must be in YYYY-MM-DD format"
      );
    }
    // also check for year, month, day
    if (
      parseInt(dateArray[0]) < 1970 ||
      parseInt(dateArray[1]) < 1 ||
      parseInt(dateArray[1]) > 12 ||
      parseInt(dateArray[2]) < 1 ||
      parseInt(dateArray[2]) > 31
    ) {
      Toast.error("Invalid date: " + date + " must be in YYYY-MM-DD format");
      throw new Error(
        "Invalid date: " + date + " must be in YYYY-MM-DD format"
      );
    }
    return dateArray.join("-");
  }

  // parse and validate datetime to be of YYYY-MM-DDThh:mm:ss format, with proper date and time
  static parseDateTime(datetime) {
    const datetimeArray = datetime.split("T").map(function (item) {
      return item.trim();
    });
    if (
      datetimeArray.length !== 2 ||
      isNaN(parseInt(datetimeArray[0])) ||
      isNaN(parseInt(datetimeArray[1]))
    ) {
      Toast.error(
        "Invalid datetime: " +
          datetime +
          " must be in YYYY-MM-DDThh:mm:ss format"
      );
      throw new Error(
        "Invalid datetime: " +
          datetime +
          " must be in YYYY-MM-DDThh:mm:ss format"
      );
    }
    Parser.parseDate(datetimeArray[0]);
    Parser.parseDuration(datetimeArray[1]);
    return datetimeArray.join("T");
  }

  // parse and validate location to be of {latitude, longitude} by removing () if present and splitting by comma and trim
  static parseLocation(location) {
    const locationArray = location
      .replace(/[()]/g, "")
      .split(",")
      .map(function (item) {
        return item.trim();
      });
    if (locationArray.length !== 2) {
      Toast.error(
        "Invalid location: " +
          location +
          " must be in (latitude, longitude) format"
      );
      throw new Error(
        "Invalid location: " +
          location +
          " must be in (latitude, longitude) format"
      );
    }

    // return location as {latitude, longitude}
    return {
      latitude: Parser._parseFloat(locationArray[0]),
      longitude: Parser._parseFloat(locationArray[1]),
    };
  }

  // parse and validate priority to be an integer between 0 and 100 inclusive
  static parsePriority(priority) {
    if (
      isNaN(parseInt(priority)) ||
      parseInt(priority) < 0 ||
      parseInt(priority) > 100
    ) {
      Toast.error(
        "Invalid priority: " +
          priority +
          " must be an integer between 0 and 100 inclusive"
      );
      throw new Error(
        "Invalid priority: " +
          priority +
          " must be an integer between 0 and 100 inclusive"
      );
    }
    return Parser._parseInteger(priority);
  }

  // parse delivery and pickup arrays
  static parseAmount(items) {
    // split by comma, trim and convert to integer array.
    const itemsArray = items
      .split(",")
      .map((item) => {
        item = item.trim();
        if (item === "") {
          return null;
        }
        if (isNaN(parseInt(item))) {
          throw new Error(`Invalid amount: ${item} must be an integer`);
        }
        return parseInt(item.trim());
      })
      .filter((item) => {
        return item !== null;
      });
    return itemsArray;
  }

  // parse json
  static parseJSON(json) {
    try {
      return JSON.parse(json);
    } catch (e) {
      Toast.error("Invalid Data field: " + e.message);
      throw e;
    }
  }

  static parseMaxTasks(max_tasks) {
    if (
      isNaN(parseInt(max_tasks)) ||
      parseInt(max_tasks) < 0 ||
      parseInt(max_tasks) > 2147483647
    ) {
      Toast.error(
        "Invalid max_tasks: " +
          max_tasks +
          " must be an integer between 0 and 2147483647 inclusive"
      );
      throw new Error(
        "Invalid max_tasks: " +
          max_tasks +
          " must be an integer between 0 and 2147483647 inclusive"
      );
    }
    return Parser._parseInteger(max_tasks);
  }

  static parseSpeedFactor(speed_factor) {
    if (
      isNaN(parseFloat(speed_factor)) ||
      parseFloat(speed_factor) < 0.01 ||
      parseFloat(speed_factor) > 200
    ) {
      Toast.error(
        "Invalid speed_factor: " +
          speed_factor +
          " must be a float between 0.01 and 200 inclusive"
      );
      throw new Error(
        "Invalid speed_factor: " +
          speed_factor +
          " must be a float between 0.01 and 200 inclusive"
      );
    }
    return Parser._parseFloat(speed_factor);
  }

  // parse and validate float
  static _parseFloat(item) {
    if (isNaN(parseFloat(item))) {
      Toast.error("Invalid value: " + item + " must be a float");
      throw new Error("Invalid value: " + item + " must be a float");
    }
    return parseFloat(item.trim());
  }

  // parse and validate integer
  static _parseInteger(item) {
    if (isNaN(parseInt(item))) {
      Toast.error("Invalid value: " + item + " must be an integer");
      throw new Error("Invalid value: " + item + " must be an integer");
    }
    return parseInt(item.trim());
  }
}
