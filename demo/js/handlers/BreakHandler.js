import BreakAPI from "../api/BreakAPI.js";

export default class BreakHandler {
  constructor({ onBreakView, onBreakCreateClick, onBreakEditClick }) {
    this.breakAPI = new BreakAPI();

    this.main = document.querySelector("#main");
    this.schedule = document.querySelector("#schedule");

    this.handleBreakDelete(onBreakView);
    this.handleBreakSave(onBreakView);
    this.handleBreakEditClick(onBreakEditClick);
    this.handleBreakCreateClick(onBreakCreateClick);
    this.handleBreakTwForm();
  }

  handleBreakDelete(onBreakView) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="break-delete"]`);
      if (el) {
        let breakID = el.dataset.id;
        let vehicleID = el.dataset.vehicleId;

        // call the break api to delete the break
        this.breakAPI.deleteBreak(breakID).then(() => {
          onBreakView(vehicleID);
        });
      }
    });
  }

  handleBreakCreateClick(onBreakCreateClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="break-create"]`);
      if (el) {
        const vehicleID = el.dataset.id;
        let sibling = el.parentElement.previousElementSibling;
        onBreakCreateClick(vehicleID, sibling);
        el.outerHTML = "";
      }
    });
  }

  handleBreakEditClick(onBreakEditClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="break-edit"]`);
      if (el) {
        let breakID = el.dataset.id;
        onBreakEditClick(breakID);
      }
    });
  }

  handleBreakSave(onBreakView) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="break-save"]`);
      if (el) {
        const form = el.closest("form");
        const formData = new FormData(form);
        const break_ = {};
        for (const [key, value] of formData.entries()) {
          // if key contains [] then it is an array
          if (key.includes("[]")) {
            if (!break_[key]) {
              break_[key] = [];
            }
            break_[key].push(value);
          } else {
            break_[key] = value;
          }
        }

        this.breakAPI.saveBreak(break_).then((break_) => {
          onBreakView(break_["vehicle_id"]);
        });
      }
    });
  }

  handleBreakTwForm() {
    this.main.addEventListener("click", (event) => {
      let el = event.target.closest(`[data-action="break-tw-form-create"]`);
      if (el) {
        // select the parent element of el, and just before it append the html
        const parent = el.parentElement;
        const html = `
        <div class="input-group">
          <input type="datetime-local" class="form-control" name="tw_open[]" step="1" style="font-size: 13px;">
          <input type="datetime-local" class="form-control" name="tw_close[]" step="1" style="font-size: 13px;">
        </div>`;

        parent.insertAdjacentHTML("beforebegin", html);
      }
      el = event.target.closest(`[data-action="break-tw-form-delete"]`);
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
