import LocationAPI from "../api/LocationAPI.js";
import MapHandler from "../handlers/MapHandler.js";

export default class MapView {
  constructor() {
    this.locationAPI = new LocationAPI();
    this.handler = new MapHandler(this.handlers());
  }

  createMap() {
    var latitude = 35.7127;
    var longitude = 139.762;

    return this.locationAPI.getLocation().then((location) => {
      latitude = location.latitude;
      longitude = location.longitude;
    }).catch((error) => {
      console.log(error);
    }).then(() => {
      // create new leafletjs map
      this.map = L.map("map").setView([latitude, longitude], 14);

      // add tile layer
      L.tileLayer("https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png", {
        attribution:
          '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(this.map);
    });
  }

  activateMap() {
    this.map.getContainer().classList.add("active");

    if (this.newMarker && this.newMarker.dragging) {
      this.newMarker.dragging.enable();
    }

    this.map.on("click", (e) => {
      this.removeMapPointer();

      var latitude = e.latlng.lat.toFixed(4);
      var longitude = e.latlng.lng.toFixed(4);
      this.addMarkerOnClick(latitude, longitude);

      if (this.newMarker && this.newMarker.dragging) {
        this.newMarker.dragging.enable();
      }
    });
  }

  deactivateMap() {
    this.map.off("mouseover");
    this.map.getContainer().classList.remove("active");
    this.map.off("click");

    if (this.newMarker && this.newMarker.dragging) {
      this.newMarker.dragging.disable();
    }
  }

  // click on map to add the marker, show the popup and fill the textbox with latitude and longitude
  addMarkerOnClick(latitude, longitude) {
    this.addMapPointer(latitude, longitude);
    this.setCenter(latitude, longitude);

    // show popup
    this.newMarker.bindPopup("Latitude: " + latitude + "<br>Longitude: " + longitude).openPopup();

    // select input name=location and fill with latitude and longitude
    document.querySelector("input[name=location]").value = latitude + ", " + longitude;

    // add event listener to marker
    this.newMarker.on("dragend", (e) => {
      // get latitude and longitude by rounding off to 4 decimal places
      var latitude = e.target._latlng.lat.toFixed(4);
      var longitude = e.target._latlng.lng.toFixed(4);

      // update popup
      this.newMarker.setPopupContent("Latitude: " + latitude + "<br>Longitude: " + longitude).openPopup();

      this.setCenter(latitude, longitude);

      // update textbox
      document.querySelector("input[name=location]").value = latitude + ", " + longitude;
    });
  }

  addMapPointer(latitude, longitude) {
    this.newMarker = L.marker([latitude, longitude]).addTo(this.map);
  }

  removeMapPointer() {
    if (this.newMarker) {
      this.map.removeLayer(this.newMarker);
    }
  }

  addMarker(latitude, longitude, icon, color) {
    // icons: truck, industry, warehouse, house
     var myIcon = L.divIcon({
       html: `<i class="fa-solid fa-${icon} fa-2x" style="color: ${color};"></i>`,
       iconSize: [20, 20],
       className: "myDivIcon",
     });

    // create new marker from latitude and longitude
    var marker = L.marker([latitude, longitude], { icon: myIcon }).addTo(
      this.map
    );
    return marker;
  }

  removeMarkers(markers) {
    for (var key in markers) {
      this.map.removeLayer(markers[key]);
    }
  }

  getCenter() {
    // get latitude and longitude of map centre
    var latitude = this.map.getCenter().lat.toFixed(4);
    var longitude = this.map.getCenter().lng.toFixed(4);

    // return as an array
    return [latitude, longitude];
  }

  setCenter(latitude, longitude) {
    // set map centre to latitude and longitude
    this.map.flyTo([latitude, longitude], 15, {
      animate: true,
    });
  }

  fitMarkers(markers) {
    if (Object.keys(markers).length > 0) {
      var bounds = new L.LatLngBounds();
      for (var key in markers) {
        var marker = markers[key];
        bounds.extend(marker.getLatLng());
      }
      this.map.flyToBounds(bounds, {
        animate: true,
      });
    }
  }

  removeAllMarkers() {
    // remove all markers from map
    this.map.eachLayer((layer) => {
      if (layer instanceof L.Marker) {
        this.map.removeLayer(layer);
      }
    });
  }

  handlers() {
    return {
      onLocationTextChange: (latitude, longitude) => {
        this.removeMapPointer();
        this.addMarkerOnClick(latitude, longitude);
        this.activateMap();
      }
    };
  }
}
