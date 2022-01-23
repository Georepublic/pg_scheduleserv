import ShipmentAPI from "../api/ShipmentAPI.js";

export default class ShipmentHandler {
  constructor(
    shipments,
    emptyShipment,
    {
      onShipmentView,
      onShipmentCreateClick,
      onShipmentEditClick,
      onShipmentDelete,
      onShipmentSave,
      onShipmentClose,
    }
  ) {
    // get the shipments from the params
    this.shipments = shipments;
    this.emptyShipment = emptyShipment;
    this.shipmentAPI = new ShipmentAPI();

    this.main = document.querySelector("#main");

    this.handleShipmentView(onShipmentView);
    this.handleShipmentCreateClick(onShipmentCreateClick);
    this.handleShipmentEditClick(onShipmentEditClick);
    this.handleShipmentDelete(onShipmentDelete);
    this.handleShipmentSave(onShipmentSave);
    this.handleShipmentClose(onShipmentClose);
  }

  // get shipment from id
  getShipment(shipmentID) {
    if (!shipmentID) {
      return this.emptyShipment;
    }
    return this.shipments.find((shipment) => {
      return shipment.id === shipmentID;
    });
  }

  handleShipmentView(onShipmentView) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="shipment-view"]`);
      if (el) {
        let shipmentID = el.dataset.id;
        onShipmentView(this.getShipment(shipmentID));
      }
    });
  }

  handleShipmentCreateClick(onShipmentCreateClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="shipment-create"]`);
      if (el) {
        onShipmentCreateClick();
      }
    });
  }

  handleShipmentEditClick(onShipmentEditClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="shipment-edit"]`);
      if (el) {
        let shipmentID = el.dataset.id;
        onShipmentEditClick(this.getShipment(shipmentID));
      }
    });
  }

  handleShipmentDelete(onShipmentDelete) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="shipment-delete"]`);
      if (el) {
        let shipmentID = el.dataset.id;
        // call the shipment api to delete the shipment
        this.shipmentAPI.deleteShipment(shipmentID).then(() => {
          // remove the shipment from the list
          this.shipments = this.shipments.filter((shipment) => {
            return shipment.id !== shipmentID;
          });
          // call the onShipmentDelete callback
          onShipmentDelete(this.shipments);
        });
      }
    });
  }

  handleShipmentSave(onShipmentSave) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="shipment-save"]`);
      if (el) {
        const form = el.closest("form");
        const formData = new FormData(form);
        const shipment = {};
        for (const [key, value] of formData.entries()) {
          shipment[key] = value;
        }
        const id = shipment["id"];

        this.shipmentAPI.saveShipment(shipment).then((shipment) => {
          // edit the shipment in the list, or append a new shipment to the list depending on the id
          if (id) {
            // update the shipment
            this.shipments = this.shipments.map((oldShipment) => {
              if (oldShipment.id === shipment.id) {
                return shipment;
              }
              return oldShipment;
            });
          } else {
            // append the new shipment to the list
            this.shipments.push(shipment);
          }
          onShipmentSave(shipment, this.shipments);
        });
      }
    });
  }

  handleShipmentClose(onShipmentClose) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="shipment-close"]`);
      if (el) {
        onShipmentClose();
      }
    });
  }
}
