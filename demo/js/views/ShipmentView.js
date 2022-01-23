import ShipmentAPI from "../api/ShipmentAPI.js";
import ShipmentHandler from "../handlers/ShipmentHandler.js";
import Random from "../utils/Random.js";
import AbstractView from "./AbstractView.js";

export default class ShipmentView extends AbstractView {
  constructor(params) {
    super(params, false);

    // get the shipments from the params
    this.shipments = params.shipments;
    this.projectID = params.projectID;
    this.mapView = params.mapView;
    this.shipmentAPI = new ShipmentAPI();

    this.shipmentLeftDiv = document.createElement("div");
    document.querySelector("#app-left").appendChild(this.shipmentLeftDiv);

    this.handler = new ShipmentHandler(
      params.shipments,
      this.getEmptyShipment(),
      this.handlers()
    );
  }

  // render the shipments for this project
  render() {
    // get the html for the shipments
    let shipmentsHtml = this.getShipmentsHtml();

    // set the html for the shipments
    this.shipmentLeftDiv.innerHTML = shipmentsHtml;

    this.mapView.addShipmentMarkers(this.shipments);
    this.mapView.fitAllMarkers();
  }

  // get the html for the shipments
  getShipmentsHtml() {
    // get the html for each shipment
    let shipmentsHtml = this.shipments.map((shipment) => {
      return this.getShipmentHtml(shipment);
    });

    if (shipmentsHtml.length === 0) {
      shipmentsHtml = [
        `
        <div class="list-group-item flex-column align-items-start">
          <p class="mb-1">No shipments found...</p>
        </div>
      `,
      ];
    }

    // return the html for the shipments, with card heading of shipments with max height of 30vh and scrolling
    return `
      <div class="list-group">
        <div class="card">
          <div class="card-header shipment-view-heading">
            <h5 class="mb-0">
              Shipments
              <button type="button" class="btn btn-success" data-action="shipment-create" style="float: right">Add</button>
            </h5>
          </div>
          <div style="max-height: 30vh; overflow-y: scroll;">
            ${shipmentsHtml.join("")}
          </div>
        </div>
      </div>
    `;
  }

  // get the html for the shipment
  getShipmentHtml(shipment) {
    const color = Random.getRandomColor(shipment.id);
    let html = `
      <div style="background-color: ${color}">
      <div class="list-group-item flex-column align-items-start" data-action="shipment-view" data-id="${
        shipment.id
      }">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">${shipment.id}</h5>
        </div>

        <div class="d-flex w-100 justify-content-between">
          <p class="mb-1">${JSON.stringify(shipment.data)}</p>
        </div>
      </div>
      </div>
    `;

    // return the html for the shipment
    return html;
  }

  getCompleteShipmentHtml(shipment) {
    let html = `
      <div class="card">
        <div class="card-header shipment-view-heading">
          <h5 class="mb-0">
            Shipment
            <button type="button" class="btn btn-danger" data-action="shipment-close">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">ID: ${shipment.id}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Pickup Location (Lat, Lon): ${
              shipment.p_location.latitude
            }, ${shipment.p_location.longitude}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Delivery Location (Lat, Lon): ${
              shipment.d_location.latitude
            }, ${shipment.d_location.longitude}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Pickup Setup: ${shipment.p_setup}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Delivery Setup: ${shipment.d_setup}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Pickup Service: ${shipment.p_service}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Delivery Service: ${shipment.d_service}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Amount: [${shipment.amount}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Skills: [${shipment.skills}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Priority: ${shipment.priority}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Project ID: ${shipment.project_id}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Data: ${JSON.stringify(shipment.data)}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Created At: ${shipment.created_at}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Updated At: ${shipment.updated_at}</p>
          </div>
          <div class="d-flex w-100 justify-content-center">
            <button class="btn btn-primary mx-2" data-action="shipment-edit" data-id="${
              shipment.id
            }">Edit</button>
            <button class="btn btn-danger mx-2" data-action="shipment-delete" data-id="${
              shipment.id
            }">Delete</button>
          </div>
        </div>
      </div>
    `;

    // return the html for the shipment
    return html;
  }

  getShipmentFormHtml(shipment) {
    let html = `
      <div class="card">
        <div class="card-header shipment-view-heading">
          <h5 class="mb-0">
            Shipment
            <button type="button" class="btn btn-danger" data-action="shipment-close" data-id="${
              shipment.id
            }">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <form>
            <input type="hidden" name="id" value="${shipment.id}">
            <input type="hidden" name="project_id" value="${
              shipment.project_id
            }">
            <div class="form-group">
              <label>Pickup Location (Lat, Lon)</label>
              <input type="text" class="form-control" name="p_location" value="${
                shipment.p_location.latitude
              }, ${
      shipment.p_location.longitude
    }" data-action="p_shipment-location-change">
            </div>
            <div class="form-group">
              <label>Delivery Location (Lat, Lon)</label>
              <input type="text" class="form-control" name="d_location" value="${
                shipment.d_location.latitude
              }, ${
      shipment.d_location.longitude
    }" data-action="d_shipment-location-change">
            </div>
            <div class="form-group">
              <label>Pickup Setup</label>
              <input type="time" class="form-control" name="p_setup" value="${
                shipment.p_setup
              }" step="1">
            </div>
            <div class="form-group">
              <label>Delivery Setup</label>
              <input type="time" class="form-control" name="d_setup" value="${
                shipment.d_setup
              }" step="1">
            </div>
            <div class="form-group">
              <label>Pickup Service</label>
              <input type="time" class="form-control" name="p_service" value="${
                shipment.p_service
              }" step="1">
            </div>
            <div class="form-group">
              <label>Delivery Service</label>
              <input type="time" class="form-control" name="d_service" value="${
                shipment.d_service
              }" step="1">
            </div>
            <div class="form-group">
              <label>Amount</label>
              <input type="text" class="form-control" name="amount" value="${
                shipment.amount
              }">
            </div>
            <div class="form-group">
              <label>Skills</label>
              <input type="text" class="form-control" name="skills" value="${
                shipment.skills
              }">
            </div>
            <div class="form-group">
              <label>Priority</label>
              <input type="number" class="form-control" name="priority" min="0" max="100" value="${
                shipment.priority
              }">
            </div>
            <div class="form-group">
              <label>Data</label>
              <input type="text" class="form-control" name="data" value='${JSON.stringify(
                shipment.data
              )}'>
            </div>
            <div class="d-flex w-100 justify-content-center">
              <button type="button" class="btn btn-primary mx-2" data-action="shipment-save" data-id="${
                shipment.id
              }">Save</button>
              <button type="button" class="btn btn-warning mx-2" data-action="shipment-edit" data-id="${
                shipment.id
              }">Reset</button>
              <button type="button" class="btn btn-danger mx-2" data-action="shipment-view" data-id="${
                shipment.id
              }">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    `;

    // return the html for the shipment
    return html;
  }

  selectShipment(shipmentID) {
    this.deselectAll();
    let shipmentViewElement = document.querySelector(
      `[data-action="shipment-view"][data-id="${shipmentID}"]`
    );
    shipmentViewElement.classList.add("active");

    // move the element into view
    shipmentViewElement.scrollIntoView({
      behavior: "smooth",
      block: "nearest",
    });
  }

  deselectAll() {
    // for all elements in query selector, remove their active class
    document.querySelectorAll(`.list-group-item.active`).forEach((element) => {
      element.classList.remove("active");
    });
  }

  getEmptyShipment() {
    // get map center
    let coordinates = this.mapView.getCenter();
    return {
      id: "",
      p_location: {
        latitude: coordinates[0],
        longitude: coordinates[1],
      },
      d_location: {
        latitude: coordinates[0],
        longitude: coordinates[1],
      },
      p_setup: "00:00:00",
      d_setup: "00:00:00",
      p_service: "00:00:00",
      d_service: "00:00:00",
      amount: "",
      skills: "",
      priority: "0",
      project_id: this.projectID,
      data: {},
      created_at: "",
      updated_at: "",
    };
  }

  handlers() {
    return {
      onShipmentView: (shipment) => {
        if (!shipment.id) {
          this.setHtmlRight("");
          return;
        }
        // get the complete html for the shipment
        let shipmentHtml = this.getCompleteShipmentHtml(shipment);

        // set the html for the shipment
        this.setHtmlRight(shipmentHtml);

        // select the shipment
        this.selectShipment(shipment.id);

        this.mapView.addShipmentMapPointer(
          shipment.p_location.latitude,
          shipment.p_location.longitude,
          shipment.d_location.latitude,
          shipment.d_location.longitude
        );
        this.mapView.deactivateMap();
      },
      onShipmentCreateClick: () => {
        this.deselectAll();

        const shipment = this.getEmptyShipment();

        // create the shipment form html with empty shipment
        let shipmentHtml = this.getShipmentFormHtml(shipment);

        // set the html for the shipment
        this.setHtmlRight(shipmentHtml);

        this.mapView.addShipmentMapPointer(
          shipment.p_location.latitude,
          shipment.p_location.longitude,
          shipment.d_location.latitude,
          shipment.d_location.longitude,
        );
        this.mapView.activateMap();
        this.mapView.fitAllMarkers();
      },
      onShipmentEditClick: (shipment) => {
        // get the complete html for the shipment
        let shipmentHtml = this.getShipmentFormHtml(shipment);

        // set the html for the shipment
        this.setHtmlRight(shipmentHtml);

        this.mapView.addShipmentMapPointer(
          shipment.p_location.latitude,
          shipment.p_location.longitude,
          shipment.d_location.latitude,
          shipment.d_location.longitude
        );
        this.mapView.activateMap();
      },
      onShipmentSave: (shipment, newShipments) => {
        this.shipments = newShipments;
        this.render();
        this.handlers().onShipmentView(shipment);
      },
      onShipmentDelete: (newShipments) => {
        this.deselectAll();
        this.setHtmlRight("");
        this.shipments = newShipments;
        this.render();
        this.mapView.removeMapPointers();
        this.mapView.deactivateMap();
        this.mapView.fitAllMarkers();
      },
      onShipmentClose: () => {
        this.deselectAll();
        this.setHtmlRight("");
        this.mapView.removeMapPointers();
        this.mapView.deactivateMap();
        this.mapView.fitAllMarkers();
      },
    };
  }
}
