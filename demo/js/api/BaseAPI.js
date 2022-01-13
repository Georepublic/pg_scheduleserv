import Toast from "../utils/Toast.js"

export default class BaseAPI {
  constructor() {
    this.baseURL = "http://localhost:9100";
  }

  get(url) {
    return axios
      .get(this.baseURL + url)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        if (error.response) {
          Toast.error(error.response.data.error);
        } else {
          Toast.error(error.message);
        }
      });
  }

  post(url, data) {
    return axios
      .post(this.baseURL + url, data)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        if (error.response) {
          Toast.error(error.response.data.error);
        } else {
          Toast.error(error.message);
        }
      });
  }

  post(url, data) {
    return axios
      .post(this.baseURL + url, data)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        if (error.response) {
          Toast.error(error.response.data.error);
        } else {
          Toast.error(error.message);
        }
      });
  }

  post(url, data) {
    return axios
      .patch(this.baseURL + url, data)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        if (error.response) {
          Toast.error(error.response.data.error);
        } else {
          Toast.error(error.message);
        }
      });
  }

  delete(url) {
    return axios
      .post(this.baseURL + url)
      .then((response) => {
        Toast.success(response.data.message);
        return response.data.data;
      })
      .catch((error) => {
        if (error.response) {
          Toast.error(error.response.data.error);
        } else {
          Toast.error(error.message);
        }
      });
  }
}
