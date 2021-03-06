import Toast from "../utils/Toast.js";
import { config } from "../config.js";

export default class BaseAPI {
  constructor() {
    this.baseURL = config.serverUrl;
  }

  showError(error) {
    if (error.response) {
      if (error.response.data.error) {
        Toast.error(error.response.data.error);
      } else if (error.response.data.errors) {
        Toast.error(error.response.data.errors.join("<br>"));
      } else {
        Toast.error("Unknown error occured");
      }
    } else {
      Toast.error(error.message);
    }
  }

  get(url) {
    return axios
      .get(this.baseURL + url, {
        headers: {
          Accept: "application/json",
        },
      })
      .then((response) => {
        return response.data.data;
      })
      .catch((error) => {
        this.showError(error);
        throw error;
      });
  }

  getIcal(url) {
    Toast.info("Processing...");
    return axios
      .get(this.baseURL + url, {
        headers: {
          "Content-Type": "text/calendar",
        },
      })
      .then((response) => {
        return response.data;
      })
      .catch((error) => {
        this.showError(error);
        throw error;
      });
  }

  post(url, data) {
    Toast.info("Processing...");
    return axios
      .post(this.baseURL + url, data)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        this.showError(error);
        throw error;
      });
  }

  patch(url, data) {
    Toast.info("Processing...");
    return axios
      .patch(this.baseURL + url, data)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        this.showError(error);
        throw error;
      });
  }

  delete(url) {
    Toast.info("Processing...");
    return axios
      .delete(this.baseURL + url)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        this.showError(error);
        throw error;
      });
  }
}
