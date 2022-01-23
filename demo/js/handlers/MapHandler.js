export default class MapHandler {
  constructor({
    onJobLocationTextChange,
    onShipmentLocationTextChange,
    onVehicleLocationTextChange,
  }) {
    this.main = document.querySelector("#main");

    this.handleJobLocationTextChange(onJobLocationTextChange);
    this.handleShipmentLocationTextChange(onShipmentLocationTextChange);
    this.handleVehicleLocationTextChange(onVehicleLocationTextChange);
  }

  getLocationArray(location) {
    const locationArray = location
      .replace(/[()]/g, "")
      .split(",")
      .map(function (item) {
        return item.trim();
      });

    return locationArray;
  }

  checkLocationArray(locationArray) {
    if (locationArray.length === 2) {
      // check if both the elements of array are float
      if (
        !isNaN(parseFloat(locationArray[0])) &&
        !isNaN(parseFloat(locationArray[1]))
      ) {
        return true;
      }
    }
    return false;
  }

  handleJobLocationTextChange(onJobLocationTextChange) {
    this.main.addEventListener("input", (event) => {
      var location = event.target.closest(
        `[data-action="job-location-change"]`
      );
      if (location) {
        const locationArray = this.getLocationArray(
          document.querySelector(`[data-action="job-location-change"]`).value
        );
        if (this.checkLocationArray(locationArray)) {
          onJobLocationTextChange(locationArray);
        }
      }
    });
  }

  handleShipmentLocationTextChange(onShipmentLocationTextChange) {
    this.main.addEventListener("input", (event) => {
      var pickup = event.target.closest(
        `[data-action="p_shipment-location-change"]`
      );
      var delivery = event.target.closest(
        `[data-action="d_shipment-location-change"]`
      );
      if (pickup || delivery) {
        const pickupLocationArray = this.getLocationArray(
          document.querySelector(`[data-action="p_shipment-location-change"]`)
            .value
        );
        const deliveryLocationArray = this.getLocationArray(
          document.querySelector(`[data-action="d_shipment-location-change"]`)
            .value
        );
        if (
          this.checkLocationArray(pickupLocationArray) &&
          this.checkLocationArray(deliveryLocationArray)
        ) {
          onShipmentLocationTextChange(
            pickupLocationArray,
            deliveryLocationArray
          );
        }
      }
    });
  }

  handleVehicleLocationTextChange(onVehicleLocationTextChange) {
    this.main.addEventListener("input", (event) => {
      var start = event.target.closest(
        `[data-action="start_vehicle-location-change"]`
      );
      var end = event.target.closest(
        `[data-action="end_vehicle-location-change"]`
      );
      if (start || end) {
        const startLocationArray = this.getLocationArray(
          document.querySelector(
            `[data-action="start_vehicle-location-change"]`
          ).value
        );
        const endLocationArray = this.getLocationArray(
          document.querySelector(`[data-action="end_vehicle-location-change"]`)
            .value
        );
        if (
          this.checkLocationArray(startLocationArray) &&
          this.checkLocationArray(endLocationArray)
        ) {
          onVehicleLocationTextChange(startLocationArray, endLocationArray);
        }
      }
    });
  }
}
