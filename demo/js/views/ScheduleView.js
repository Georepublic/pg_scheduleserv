import ScheduleAPI from "../api/ScheduleAPI.js";
import ScheduleHandler from "../handlers/ScheduleHandler.js";
import Random from "../utils/Random.js";
import AbstractView from "./AbstractView.js";

export default class ScheduleView extends AbstractView {
  constructor(params) {
    super(params, false);

    this.data = params.data;
    this.projectID = params.projectID;
    this.mapView = params.mapView;
    this.scheduleAPI = new ScheduleAPI();

    this.scheduleDiv = document.createElement("div");
    document.querySelector("#schedule").appendChild(this.scheduleDiv);

    this.handler = new ScheduleHandler(
      params.data,
      params.projectID,
      this.handlers()
    );
  }

  // render the schedules for this project
  render() {
    // get the html for the schedules
    let schedulesHtml = this.getSchedulesHtml();

    // set the html for the schedules
    this.scheduleDiv.innerHTML = schedulesHtml;

    // this.mapView.addScheduleMarkers(this.schedules);
    // this.mapView.fitAllMarkers();

    var tooltipTriggerList = [].slice.call(
      document.querySelectorAll('[data-bs-toggle="tooltip"]')
    );
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
      return new bootstrap.Tooltip(tooltipTriggerEl);
    });
  }

  // get the html for the schedules
  getSchedulesHtml() {
    // get the timeline end points
    const minMaxHours = this.getMinMaxHours(this.data.schedule);
    const minHours = minMaxHours[0];
    const maxHours = minMaxHours[1];

    // get the vertical lines html
    let verticalLinesHtml = this.getVerticalLinesHtml(minHours, maxHours);

    let schedulesHtml = this.data.schedule.map((schedule) => {
      return this.getScheduleHtml(
        schedule,
        minHours,
        maxHours,
        verticalLinesHtml
      );
    });

    if (schedulesHtml.length === 0) {
      schedulesHtml = [
        `
        <div class="list-group-item flex-column align-items-start">
          <p class="mb-1">No schedules found...</p>
        </div>
      `,
      ];
    }

    return `
      <div class="list-group">
        <div class="card">
          <div class="card-header schedule-view-heading">
            <h5 class="mb-0">
              Schedules
              <button type="button" class="btn btn-success" data-action="schedule-create" style="float: right">Create Schedule</button>
            </h5>
          </div>
          <div class="card-body-custom">
            ${schedulesHtml.join("")}
          </div>
        </div>
      </div>
    `;
  }

  // get the html for the schedule
  getScheduleHtml(schedule, minHours, maxHours, verticalLinesHtml) {
    let color = Random.getRandomColor(schedule.vehicle_id);

    let width = (maxHours - minHours) * 10;
    let widthFactor = 1;
    if (width < 100) {
      widthFactor = 100 / width;
    }

    // get the tasks html
    let tasksHtml = this.getTasksHtml(
      schedule.route,
      minHours,
      maxHours,
      color
    );

    let html = `
      <div style="background-color: ${color};">
      <div class="list-group-item flex-column align-items-start" data-id="${
        schedule.vehicle_id
      }">
        <div class="container-fluid" style="margin: 0; padding: 0">
          <div class="row">
            <div class="col-2">
              <div class="d-flex w-100 justify-content-between">
                <p class="mb-1">${schedule.vehicle_id}</p>
              </div>
              <div class="d-flex w-100 justify-content-between">
                <p class="mb-1">${JSON.stringify(schedule.vehicle_data)}</p>
              </div>
            </div>
            <div class="col-10">
              <div class="timelines-container" style="padding-bottom: 20px;>
                <div class="timeline-item flex-container items-center even">
                    <div class="timeline flex-main">
                      <div class="line"></div>
                      <div class="value-line" style="background-color: ${color}; width: ${
      width * widthFactor
    }%; left: 0%;"></div>
                      ${tasksHtml}
                      <div class="label-vertical-lines">
                        ${verticalLinesHtml}
                      </div>
                    </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      </div>
    `;

    // return the html for the schedule
    return html;
  }

  // get the html for the tasks
  getTasksHtml(route, minHours, maxHours, defaultColor) {
    let width = (maxHours - minHours) * 10;
    let widthFactor = 1;
    if (width < 100) {
      widthFactor = 100 / width;
    }

    let tasksHtml = route.map((task) => {
      let color =
        task.task_id === "-1"
          ? defaultColor
          : Random.getRandomColor(task.task_id);

      let arrivalHours = new Date(task.arrival).getTime() / 1000 / 60 / 60;
      let departureHours = new Date(task.departure).getTime() / 1000 / 60 / 60;

      // take the difference between the start time and the min hours to get the left position
      let left = (arrivalHours - minHours) * 10;

      // duration is the difference between the start time and the end time
      let duration = (departureHours - arrivalHours) * 10;

      let widthString = `${duration * widthFactor}%;`;
      if (widthString === "0%;") {
        widthString = "10px;";
      }

      // get the html for the task
      let html = `
        <div class="task-item" style="left: ${
          left * widthFactor
        }%; width: ${widthString}">
          <a data-bs-toggle="tooltip" data-bs-placement="bottom" title="Tooltip Text">
            <div class="task-item-full" style="background-color: ${color};"></div>
          </a>
        </div>
      `;

      console.log(task);
      console.log(html);

      return html;
    });

    return tasksHtml.join("");
  }

  // get the html for the vertical lines
  // MAYBE: Change this to prevent overflow of timeline
  getVerticalLinesHtml(minHours, maxHours) {
    let html = "";

    let width = (maxHours - minHours) * 10;
    let widthFactor = 1;
    if (width < 100) {
      widthFactor = 100 / width;
    }

    // get the number of vertical lines
    const numberOfLines = maxHours - minHours;

    // get the html for the vertical lines
    for (let i = 0; i <= numberOfLines; i++) {
      html += `
        <div class="label-vertical-line" style="left: ${
          10 * i * widthFactor
        }%;"></div>
      `;
    }

    // return the html for the vertical lines
    return html;
  }

  getMinMaxHours(schedule) {
    if (schedule.length === 0) {
      return 0;
    }

    let minDate = new Date(9999, 1, 1, 0, 0, 0, 0);
    let maxDate = new Date(0, 1, 1, 0, 0, 0, 0);

    schedule.forEach((vehicle) => {
      if (new Date(vehicle.route[0].arrival) < minDate) {
        minDate = new Date(vehicle.route[0].arrival);
      }
      if (new Date(vehicle.route[vehicle.route.length - 1].arrival) > maxDate) {
        maxDate = new Date(vehicle.route[vehicle.route.length - 1].arrival);
      }
    });

    // get hours since epoch
    const minDateHours = Math.floor(minDate.getTime() / 1000 / 60 / 60);
    const maxDateHours = Math.ceil(maxDate.getTime() / 1000 / 60 / 60);

    return [minDateHours, maxDateHours];
  }

  getTimelineHtml() {}

  handlers() {
    return {
      onScheduleCreate: (data) => {
        this.data = data;
        this.render();
      },
    };
  }
}
