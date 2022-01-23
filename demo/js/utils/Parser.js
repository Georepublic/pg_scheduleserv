import Toast from "./Toast.js";

export default class Parser {
  // parse and validate duration to be of HH:MM:SS format, with proper time
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
