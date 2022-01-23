import VehicleAPI from "../api/VehicleAPI.js";

export default class VehicleHandler {
  constructor(
    vehicles,
    emptyVehicle,
    {
      onVehicleView,
      onVehicleCreateClick,
      onVehicleEditClick,
      onVehicleDelete,
      onVehicleSave,
      onVehicleClose,
    }
  ) {
    // get the vehicles from the params
    this.vehicles = vehicles;
    this.emptyVehicle = emptyVehicle;
    this.vehicleAPI = new VehicleAPI();

    this.main = document.querySelector("#main");

    this.handleVehicleView(onVehicleView);
    this.handleVehicleCreateClick(onVehicleCreateClick);
    this.handleVehicleEditClick(onVehicleEditClick);
    this.handleVehicleDelete(onVehicleDelete);
    this.handleVehicleSave(onVehicleSave);
    this.handleVehicleClose(onVehicleClose);
  }

  // get vehicle from id
  getVehicle(vehicleID) {
    if (!vehicleID) {
      return this.emptyVehicle;
    }
    return this.vehicles.find((vehicle) => {
      return vehicle.id === vehicleID;
    });
  }

  handleVehicleView(onVehicleView) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="vehicle-view"]`);
      if (el) {
        let vehicleID = el.dataset.id;
        onVehicleView(this.getVehicle(vehicleID));
      }
    });
  }

  handleVehicleCreateClick(onVehicleCreateClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="vehicle-create"]`);
      if (el) {
        onVehicleCreateClick();
      }
    });
  }

  handleVehicleEditClick(onVehicleEditClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="vehicle-edit"]`);
      if (el) {
        let vehicleID = el.dataset.id;
        onVehicleEditClick(this.getVehicle(vehicleID));
      }
    });
  }

  handleVehicleDelete(onVehicleDelete) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="vehicle-delete"]`);
      if (el) {
        let vehicleID = el.dataset.id;
        // call the vehicle api to delete the vehicle
        this.vehicleAPI.deleteVehicle(vehicleID).then(() => {
          // remove the vehicle from the list
          this.vehicles = this.vehicles.filter((vehicle) => {
            return vehicle.id !== vehicleID;
          });
          // call the onVehicleDelete callback
          onVehicleDelete(this.vehicles);
        });
      }
    });
  }

  handleVehicleSave(onVehicleSave) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="vehicle-save"]`);
      if (el) {
        const form = el.closest("form");
        const formData = new FormData(form);
        const vehicle = {};
        for (const [key, value] of formData.entries()) {
          vehicle[key] = value;
        }
        const id = vehicle["id"];

        this.vehicleAPI.saveVehicle(vehicle).then((vehicle) => {
          // edit the vehicle in the list, or append a new vehicle to the list depending on the id
          if (id) {
            // update the vehicle
            this.vehicles = this.vehicles.map((oldVehicle) => {
              if (oldVehicle.id === vehicle.id) {
                return vehicle;
              }
              return oldVehicle;
            });
          } else {
            // append the new vehicle to the list
            this.vehicles.push(vehicle);
          }
          onVehicleSave(vehicle, this.vehicles);
        });
      }
    });
  }

  handleVehicleClose(onVehicleClose) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="vehicle-close"]`);
      if (el) {
        onVehicleClose();
      }
    });
  }
}
