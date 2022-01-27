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
    this.handleShipmentTwForm();
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
          // if key contains [] then it is an array
          if (key.includes("[]")) {
            if (!shipment[key]) {
              shipment[key] = [];
            }
            shipment[key].push(value);
          } else {
            shipment[key] = value;
          }
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

  handleShipmentTwForm() {
    this.main.addEventListener("click", (event) => {
      let el = event.target.closest(`[data-action="shipment-tw-form-create"]`);
      if (el) {
        let value = el.dataset.value;

        // select the parent element of el, and just before it append the html
        const parent = el.parentElement;
        const html = `
        <div class="input-group">
          <input type="datetime-local" class="form-control" name="${value}_tw_open[]" step="1" style="font-size: 13px;">
          <span class="input-group-addon"></span>
          <input type="datetime-local" class="form-control" name="${value}_tw_close[]" step="1" style="font-size: 13px;">
        </div>`;

        parent.insertAdjacentHTML("beforebegin", html);
      }
      el = event.target.closest(`[data-action="shipment-tw-form-delete"]`);
      if (el) {
        const parent = el.parentElement;

        // select the adjacent element of parent and remove it
        const sibling = parent.previousElementSibling;

        // if sibling has the class input-group, remove it
        if (sibling.classList.contains("input-group")) {
          sibling.remove();
        }
      }
    });
  }
}
