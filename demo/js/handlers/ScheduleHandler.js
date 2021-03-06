import ScheduleAPI from "../api/ScheduleAPI.js";
import Toast from "../utils/Toast.js";

export default class ScheduleHandler {
  constructor(
    data,
    projectID,
    { onScheduleCreate, onScheduleDelete, onPlayRoute, onStopPlayRoute }
  ) {
    this.data = data;
    this.projectID = projectID;
    this.scheduleAPI = new ScheduleAPI();

    this.main = document.querySelector("#main");
    this.schedule = document.querySelector("#schedule");

    this.handleScheduleCreate(onScheduleCreate);
    this.handleScheduleDownload();
    this.handleVehicleScheduleDownload();
    this.handleScheduleDelete(onScheduleDelete);
    this.handleScheduleDropdownClick();
    this.handlePlayRoute(onPlayRoute);
    this.handleStopPlayRoute(onStopPlayRoute);
  }

  downloadIcal(data, filename) {
    const url = window.URL.createObjectURL(new Blob([data]));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", filename);
    document.body.appendChild(link);
    link.click();
    Toast.success("Schedule downloaded");
  }

  handleScheduleCreate(onScheduleCreate) {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="schedule-create"]`);
      if (el) {
        this.scheduleAPI
          .createSchedule(this.projectID, el.dataset.type)
          .then((data) => {
            onScheduleCreate(data);
          });
      }
    });
  }

  handleScheduleDownload() {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="schedule-download"]`);
      if (el) {
        this.scheduleAPI.getScheduleIcal(this.projectID).then((data) => {
          let filename = `schedule-${this.projectID}.ical`;
          this.downloadIcal(data, filename);
        });
      }
    });
  }

  handleVehicleScheduleDownload() {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(
        `[data-action="vehicle-schedule-download"]`
      );
      if (el) {
        let vehicleID = el.dataset.id;
        this.scheduleAPI.getVehicleScheduleIcal(vehicleID).then((data) => {
          let filename = `schedule-vehicle-${vehicleID}.ical`;
          this.downloadIcal(data, filename);
        });
      }
    });
  }

  handleScheduleDelete(onScheduleDelete) {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="schedule-delete"]`);
      if (el) {
        this.scheduleAPI.deleteSchedule(this.projectID).then(() => {
          onScheduleDelete();
        });
      }
    });
  }

  handleScheduleDropdownClick() {
    this.schedule.addEventListener("click", (event) => {
      const el = document.querySelector(`[data-action="schedule-create"]`);
      if (event.target.closest(`[data-action="schedule-normal"]`)) {
        el.dataset.type = "normal";
        el.innerText = "Create Schedule (Normal)";
        [el, el.nextElementSibling].forEach((el) => {
          el.classList.remove("btn-primary");
          el.classList.add("btn-success");
        });
      } else if (event.target.closest(`[data-action="schedule-fresh"]`)) {
        el.dataset.type = "fresh";
        el.innerText = "Create Schedule (Fresh)";
        [el, el.nextElementSibling].forEach((el) => {
          el.classList.remove("btn-success");
          el.classList.add("btn-primary");
        });
      }
    });
  }

  handlePlayRoute(onPlayRoute) {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="play-route"]`);
      if (el) {
        let vehicleID = el.dataset.id;
        onPlayRoute(vehicleID);
      }
    });
  }

  handleStopPlayRoute(onStopPlayRoute) {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="stop-play-route"]`);
      if (el) {
        let vehicleID = el.dataset.id;
        onStopPlayRoute(vehicleID);
      }
    });
  }
}
