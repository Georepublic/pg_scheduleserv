export default class MapHandler {
  constructor({ onLocationTextChange }) {
    this.appLeft = document.querySelector("#app-left");
    this.app = document.querySelector("#app");
    this.appRight = document.querySelector("#app-right");

    this.handleLocationTextChange(onLocationTextChange);
  }

  handleLocationTextChange(onLocationTextChange) {
    this.appRight.addEventListener("input", (event) => {
      const el = event.target.closest(`[data-action="location-change"]`);
      if (el) {
        let location = el.value;

        const locationArray = location
          .replace(/[()]/g, "")
          .split(",")
          .map(function (item) {
            return item.trim();
          });

        if (locationArray.length === 2) {
          // check if both the elements of array are float
          if (
            !isNaN(parseFloat(locationArray[0])) &&
            !isNaN(parseFloat(locationArray[1]))
          ) {
            onLocationTextChange(locationArray[0], locationArray[1]);
          }
        }
      }
    });
  }
}
