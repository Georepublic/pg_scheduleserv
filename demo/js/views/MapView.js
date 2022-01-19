import LocationAPI from "../api/LocationAPI.js";
import MapHandler from "../handlers/MapHandler.js";

export default class {
  constructor() {
    this.locationAPI = new LocationAPI();
    this.mapActivated = false;

    this.handler = new MapHandler(this.mapActivated, this.handlers());
  }

  createMap() {
    var latitude = 35.7127;
    var longitude = 139.762;

    this.locationAPI.getLocation().then((location) => {
      latitude = location.latitude;
      longitude = location.longitude;
    }).catch((error) => {
      console.log(error);
    }).then(() => {
      // create new leafletjs map
      this.map = L.map("map").setView([latitude, longitude], 13);

      // add tile layer
      L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
        attribution:
          '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(this.map);
    });
  }

  activateMap() {
    // convert mouse to pointer on map hover
    this.map.on("mouseover", () => {
      this.map.getContainer().style.cursor = "pointer";
    });

    this.map.on("click", (e) => {
      // remove previous clicked marker and add the new marker
      if (this.newMarker) {
        this.map.removeLayer(this.newMarker);
      }
      this.addMarker(e);
    });
  }

  deactivateMap() {
    this.map.off("mouseover");
    this.map.on("mouseover", () => {
      // remove the pointer
      this.map.getContainer().style.cursor = "";
    });
    this.map.off("click");
  }

  // click on map to add the marker, show the popup and fill the textbox with latitude and longitude
  addMarker(e) {
    // get latitude and longitude by rounding off to 4 decimal places
    var latitude = e.latlng.lat.toFixed(4);
    var longitude = e.latlng.lng.toFixed(4);

    // create new marker from latitude and longitude
    this.newMarker = L.marker([latitude, longitude]).addTo(this.map);

    // show popup
    this.newMarker.bindPopup("Latitude: " + latitude + "<br>Longitude: " + longitude).openPopup();

    // select input name=location and fill with latitude and longitude
    document.querySelector("input[name=location]").value = latitude + ", " + longitude;
  }

  handlers() {
    return {
      onToggleMapClick: () => {
        if (!this.mapActivated) {
          this.activateMap();
          this.mapActivated = true;
        } else {
          this.deactivateMap();
          this.mapActivated = false;
        }
      },
    };
  }
}
