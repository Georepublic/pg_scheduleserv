export default class {
  constructor({ onLocationTextChange }) {
    this.handleLocationTextChange(onLocationTextChange);
  }

  handleLocationTextChange(onLocationTextChange) {
    document.addEventListener("input", (event) => {
      console.log("Input event")
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
