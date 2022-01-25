export default class OSRMAPI {
  constructor() {
    this.serverUrl = "https://router.project-osrm.org";
  }

  // get the route between several coordinates by calling the OSRM API using axios
  getRoute(coordinates) {
    let coordinatesString = coordinates.map((coordinate) => {
      return `${coordinate.longitude},${coordinate.latitude}`;
    });
    coordinatesString = coordinatesString.join(";");

    // call the OSRM API
    return axios
      .get(`${this.serverUrl}/route/v1/driving/${coordinatesString}`, {
        params: {
          geometries: "geojson",
          overview: "full",
        },
      })
      .then((response) => {
        const data = response.data;
        return {
          geometry: data.routes[0].geometry,
        };
      })
      .catch((error) => {
        throw error;
      });
  }
}
