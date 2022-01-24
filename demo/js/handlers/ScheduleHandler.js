import ScheduleAPI from "../api/ScheduleAPI.js";

export default class ScheduleHandler {
  constructor(data, projectID, { onScheduleCreate }) {
    this.data = data;
    this.projectID = projectID;
    this.scheduleAPI = new ScheduleAPI();

    this.main = document.querySelector("#main");
    this.schedule = document.querySelector("#schedule");

    this.handleScheduleCreate(onScheduleCreate);
  }

  handleScheduleCreate(onScheduleCreate) {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="schedule-create"]`);
      if (el) {
        console.log(this.data);
        this.scheduleAPI.createSchedule(this.projectID).then((data) => {
          onScheduleCreate(data);
        });
      }
    });
  }
}
