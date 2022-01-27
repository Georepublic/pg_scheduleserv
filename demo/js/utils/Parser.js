import Toast from "./Toast.js";

// function to raise error
function raiseError(error) {
  Toast.error(error);
  throw new Error(error);
}

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
      raiseError(
        "Invalid duration: " + duration + " must be in hh:mm:ss format"
      );
    }
    // also check for hours, minutes, seconds
    if (parseInt(durationArray[1]) > 59 || parseInt(durationArray[2]) > 59) {
      raiseError(
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
      raiseError("Invalid date: " + date + " must be in YYYY-MM-DD format");
    }
    // also check for year, month, day
    if (
      parseInt(dateArray[0]) < 1970 ||
      parseInt(dateArray[1]) < 1 ||
      parseInt(dateArray[1]) > 12 ||
      parseInt(dateArray[2]) < 1 ||
      parseInt(dateArray[2]) > 31
    ) {
      raiseError("Invalid date: " + date + " must be in YYYY-MM-DD format");
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
      raiseError(
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
      raiseError(
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
      raiseError(
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
          raiseError(`Invalid amount: ${item} must be an integer`);
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
      raiseError("Invalid Data field: " + e.message);
      throw e;
    }
  }

  static parseMaxTasks(max_tasks) {
    if (
      isNaN(parseInt(max_tasks)) ||
      parseInt(max_tasks) < 0 ||
      parseInt(max_tasks) > 2147483647
    ) {
      raiseError(
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
      raiseError(
        "Invalid speed_factor: " +
          speed_factor +
          " must be a float between 0.01 and 200 inclusive"
      );
    }
    return Parser._parseFloat(speed_factor);
  }

  static parseTimeWindows(tw_open_array, tw_close_array) {
    const time_windows = [];
    if (tw_open_array && tw_close_array) {
      if (tw_open_array.length !== tw_close_array.length) {
        raiseError(
          "Invalid time_windows: open and close array must be of same length"
        );
      }
      for (let i = 0; i < tw_open_array.length; i++) {
        if (tw_open_array[i] && tw_close_array[i]) {
          let parsed_tw_open = Parser.parseDateTime(tw_open_array[i]);
          let parsed_tw_close = Parser.parseDateTime(tw_close_array[i]);
          if (parsed_tw_open > parsed_tw_close) {
            raiseError(
              "Invalid time_windows: open must be less than or equal to close"
            );
          }
          time_windows.push([parsed_tw_open, parsed_tw_close]);
        } else if (tw_open_array[i] || tw_close_array[i]) {
          raiseError(
            "Invalid time_windows: both open and close must be present"
          );
        }
      }
    }
    return time_windows;
  }

  // parse and validate float
  static _parseFloat(item) {
    if (isNaN(parseFloat(item))) {
      raiseError("Invalid value: " + item + " must be a float");
    }
    return parseFloat(item.trim());
  }

  // parse and validate integer
  static _parseInteger(item) {
    if (isNaN(parseInt(item))) {
      raiseError("Invalid value: " + item + " must be an integer");
    }
    return parseInt(item.trim());
  }
}
