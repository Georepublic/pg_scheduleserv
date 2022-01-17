export default class LocationAPI {
  constructor() {}

  // Get the coordinates of a user using ip address and axios
  getLocation() {
    return axios
      .get(`https://ipapi.co/json/`)
      .then((response) => {
        const data = response.data
        return {
          latitude: data.latitude,
          longitude: data.longitude
        }
      })
      .catch((error) => {
        throw error;
      });
  }
}
