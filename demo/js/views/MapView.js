import LocationAPI from "../api/LocationAPI.js";
import MapHandler from "../handlers/MapHandler.js";
import Random from "../utils/Random.js";

export default class MapView {
  constructor() {
    this.locationAPI = new LocationAPI();
    this.mapPointers = [];
    this.jobMarkers = [];
    this.shipmentMarkers = [];
    this.vehicleMarkers = [];
    this.handler = new MapHandler(this.handlers());
  }

  createMap() {
    var latitude = 35.7127;
    var longitude = 139.762;

    return this.locationAPI
      .getLocation()
      .then((location) => {
        latitude = location.latitude;
        longitude = location.longitude;
      })
      .catch((error) => {
        console.log(error);
      })
      .then(() => {
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
    this.enableDraggingPointers();
  }

  deactivateMap() {
    this.map.off("mouseover");
    this.map.getContainer().classList.remove("active");
    this.map.off("click");
    this.disableDraggingPointers();
  }

  // click on map to add the marker, show the popup and fill the textbox with latitude and longitude
  addMapPointer(latitude, longitude, inputName, popupPrefix) {
    const marker = L.marker([latitude, longitude]).addTo(this.map);
    this.setCenter(latitude, longitude);

    // show popup
    marker
      .bindPopup(
        `${popupPrefix}Latitude: ${latitude}<br>Longitude: ${longitude}`,
        {
          autoClose: false,
        }
      )
      .openPopup();

    const inputElement = document.querySelector(`input[name=${inputName}]`);
    if (inputElement) {
      inputElement.value = `${latitude}, ${longitude}`;
    }

    // add event listener to marker
    marker.on("dragend", (e) => {
      // get latitude and longitude by rounding off to 4 decimal places
      var latitude = e.target._latlng.lat.toFixed(4);
      var longitude = e.target._latlng.lng.toFixed(4);

      // update popup
      marker
        .bindPopup(
          `${popupPrefix}Latitude: ${latitude}<br>Longitude: ${longitude}`
        )
        .openPopup();

      const inputElement = document.querySelector(`input[name=${inputName}]`);
      if (inputElement) {
        inputElement.value = `${latitude}, ${longitude}`;
      }
    });

    this.mapPointers.push(marker);
  }

  addJobMapPointer(latitude, longitude) {
    this.removeMapPointers();
    this.addMapPointer(latitude, longitude, "location", "");
    this.fitBounds([[latitude, longitude]]);
  }

  addShipmentMapPointer(p_latitude, p_longitude, d_latitude, d_longitude) {
    this.removeMapPointers();
    this.addMapPointer(
      p_latitude,
      p_longitude,
      "p_location",
      "<b>Pickup Location</b><br>"
    );
    this.addMapPointer(
      d_latitude,
      d_longitude,
      "d_location",
      "<b>Delivery Location</b><br>"
    );
    this.fitBounds([
      [p_latitude, p_longitude],
      [d_latitude, d_longitude],
    ]);
  }

  addVehicleMapPointer(
    start_latitude,
    start_longitude,
    end_latitude,
    end_longitude
  ) {
    this.removeMapPointers();
    this.addMapPointer(
      start_latitude,
      start_longitude,
      "start_location",
      "<b>Start Location</b><br>"
    );
    this.addMapPointer(
      end_latitude,
      end_longitude,
      "end_location",
      "<b>End Location</b><br>"
    );
    this.fitBounds([
      [start_latitude, start_longitude],
      [end_latitude, end_longitude],
    ]);
  }

  removeMapPointers() {
    this.mapPointers.forEach((pointer) => {
      this.map.removeLayer(pointer);
    });
    this.mapPointers = [];
  }

  enableDraggingPointers() {
    // for each pointers in this.mapPointers list, enable the dragging
    this.mapPointers.forEach((pointer) => {
      if (pointer.dragging) {
        pointer.dragging.enable();
      }
    });
  }

  disableDraggingPointers() {
    // for each pointers in this.mapPointers list, disable the dragging
    this.mapPointers.forEach((pointer) => {
      if (pointer.dragging) {
        pointer.dragging.disable();
      }
    });
  }

  addMarker(latitude, longitude, icon, color, iconPrefix="") {
    // icons: truck, industry, warehouse, house

    iconPrefix = "";  // not adding icon prefix for now
    var myIcon = L.divIcon({
      html: `${iconPrefix}<i class="fa-solid fa-${icon} fa-2x" style="color: ${color};"></i>`,
      iconSize: [20, 20],
      className: "myDivIcon",
    });

    // create new marker from latitude and longitude
    var marker = L.marker([latitude, longitude], { icon: myIcon }).addTo(
      this.map
    );
    return marker;
  }

  addJobMarkers(jobs) {
    this.removeJobMarkers();
    jobs.forEach((job) => {
      const randomColor = Random.getRandomColor(job.id);
      this.jobMarkers.push(
        this.addMarker(
          job.location.latitude,
          job.location.longitude,
          "house",
          randomColor,
          "<b>Job</b> "
        )
      );
    });
  }

  addShipmentMarkers(shipments) {
    this.removeShipmentMarkers();
    shipments.forEach((shipment) => {
      const randomColor = Random.getRandomColor(shipment.id);
      this.shipmentMarkers.push(
        this.addMarker(
          shipment.p_location.latitude,
          shipment.p_location.longitude,
          "warehouse",
          randomColor,
          "<b>Pickup</b> "
        )
      );
      this.shipmentMarkers.push(
        this.addMarker(
          shipment.d_location.latitude,
          shipment.d_location.longitude,
          "warehouse",
          randomColor,
          "<b>Delivery</b> "
        )
      );
    });
  }

  addVehicleMarkers(vehicles) {
    this.removeVehicleMarkers();
    vehicles.forEach((vehicle) => {
      const randomColor = Random.getRandomColor(vehicle.id);
      this.vehicleMarkers.push(
        this.addMarker(
          vehicle.start_location.latitude,
          vehicle.end_location.longitude,
          "industry",
          randomColor,
          "<b>Start</b> "
        )
      );
      this.vehicleMarkers.push(
        this.addMarker(
          vehicle.start_location.latitude,
          vehicle.end_location.longitude,
          "industry",
          randomColor,
          "<b>End</b> "
        )
      );
    });
  }

  removeMarkers(markers) {
    for (var key in markers) {
      this.map.removeLayer(markers[key]);
    }
  }

  removeJobMarkers() {
    this.removeMarkers(this.jobMarkers);
    this.jobMarkers = [];
  }

  removeShipmentMarkers() {
    this.removeMarkers(this.shipmentMarkers);
    this.shipmentMarkers = [];
  }

  removeVehicleMarkers() {
    this.removeMarkers(this.vehicleMarkers);
    this.vehicleMarkers = [];
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
      padding: [30, 30],
    });
  }

  fitBounds(bounds) {
    if (bounds) {
      this.map.flyToBounds(bounds, {
        animate: true,
        maxZoom: 15,
        padding: [30, 30],
      });
    }
  }

  fitMarkers(markers) {
    // for each marker in markers list, add to bounds
    var bounds = L.latLngBounds([]);
    markers.forEach((marker) => {
      bounds.extend(marker.getLatLng());
    });
    this.fitBounds(bounds);
  }

  fitAllMarkers() {
    var markers = [];
    markers = markers.concat(this.jobMarkers);
    markers = markers.concat(this.shipmentMarkers);
    markers = markers.concat(this.vehicleMarkers);
    this.fitMarkers(markers);
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
      onJobLocationTextChange: (location) => {
        this.removeMapPointers();
        this.addJobMapPointer(location[0], location[1]);
        this.activateMap();
      },
      onShipmentLocationTextChange: (pickup, delivery) => {
        this.removeMapPointers();
        this.addShipmentMapPointer(pickup[0], pickup[1], delivery[0], delivery[1]);
        this.activateMap();
      },
      onVehicleLocationTextChange: (start, end) => {
        this.removeMapPointers();
        this.addVehicleMapPointer(start[0], start[1], end[0], end[1]);
        this.activateMap();
      },
    };
  }
}
