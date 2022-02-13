import OSRMAPI from "../api/OSRMAPI.js";
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
    this.osrmAPI = new OSRMAPI();

    this.scheduleDiv = document.createElement("div");
    document.querySelector("#schedule").appendChild(this.scheduleDiv);

    this.handler = new ScheduleHandler(
      params.data,
      params.projectID,
      this.handlers()
    );
  }

  // render the schedules for this project
  render(refresh) {
    // get the html for the schedules
    let schedulesHtml = this.getSchedulesHtml();

    // when refreshing, don't change the #schedule-view-heading div
    if (refresh) {
      let scheduleViewHeading = document.querySelector(
        "#schedule-view-heading"
      );
      // remove the right sibling of the schedule view heading
      scheduleViewHeading.nextElementSibling.remove();
      // add the new schedules html
      scheduleViewHeading.insertAdjacentHTML("afterend", schedulesHtml);
    } else {
      this.scheduleDiv.innerHTML = this.getScheduleBodyHtml();
      let scheduleViewHeading = document.querySelector(
        "#schedule-view-heading"
      );
      scheduleViewHeading.insertAdjacentHTML("afterend", schedulesHtml);
    }

    this.adjustTaskItemWidth();

    var tooltipTriggerList = [].slice.call(
      document.querySelectorAll('[data-bs-toggle="tooltip"]')
    );
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
      return new bootstrap.Tooltip(tooltipTriggerEl);
    });

    this.renderGeometry();
    this.renderNumberPointers();
    this.renderUnassignedPointers();
  }

  // if the width of any circle is less than 10px, then make it 10px
  adjustTaskItemWidth() {
    let taskItems = document.querySelectorAll(".task-item");
    taskItems.forEach((taskItem) => {
      if (taskItem.offsetWidth < 10) {
        taskItem.style.width = "10px";
      }
    });
  }

  // render the geometry on the map
  renderGeometry() {
    // get the schedules
    let schedules = this.data.schedule;

    // get the routes for each schedule (array) and add to map
    let routes = schedules.map((schedule) => {
      return this.getRoute(schedule);
    });

    // wait for all the route to be resolved
    Promise.all(routes).then((route) => {
      // iterate through geometry, add to map
      route.forEach((route, index) => {
        // create style
        let style = {
          color: Random.getRandomColor(schedules[index].vehicle_id),
          weight: 5,
          opacity: 1,
          lineJoin: "round",
          lineCap: "round",
        };
        this.mapView.addRouteLayer(route.geometry, style);
      });
      this.mapView.fitAllMarkers();
    });
  }

  // render the number pointers on the map
  renderNumberPointers() {
    // get the schedules
    let schedules = this.data.schedule;

    // iterate through schedules array, for each route add a number pointer to the map
    schedules.map((schedule) => {
      let number = 1;
      let vehicleID = schedule.vehicle_id;
      schedule.route.forEach((route) => {
        this.mapView.addNumberPointer(
          vehicleID,
          route.location.latitude,
          route.location.longitude,
          number
        );
        number++;
      });
      this.mapView.setStyle(vehicleID);
    });
  }

  // render the unassigned cross pointers on the map
  renderUnassignedPointers() {
    // get the unassigned
    let unassigned = this.data.metadata.unassigned;

    // iterate through schedules array, for each route add a number pointer to the map
    unassigned.map((task) => {
      this.mapView.addUnassignedPointer(
        task.location.latitude,
        task.location.longitude
      );
    });
  }

  // get the geometry for the schedule
  getRoute(schedule) {
    // get the coordinates
    let coordinates = [];
    schedule.route.forEach((route) => {
      coordinates.push(route.location);
    });

    // call the osrm api with the coordinates to get route geometry
    return this.osrmAPI.getRoute(coordinates).then((route) => {
      return route;
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
    let labelHtml = this.getLabelHtml(minHours, maxHours);

    let schedulesHtml = this.data.schedule.map((schedule) => {
      return this.getScheduleHtml(
        schedule,
        minHours,
        maxHours,
        verticalLinesHtml,
        labelHtml
      );
    });

    let timelineBodyHtml = `
        <div class="card-body-schedule">
          <div class="container-fluid">
            <div class="row">
              <div class="col-2"></div>
              <div class="col-10">
                <div class="timelines-container">
                  <div class="timeline-item flex-container items-center">
                    <div class="timeline flex-main">
                      <div>
                        ${labelHtml[0]}
                      </div>
                    </div>
                  </div>
                  <div class="timeline-item flex-container items-center">
                    <div class="timeline flex-main">
                      <div>
                        ${labelHtml[1]}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            ${schedulesHtml.join("")}
          </div>
        </div>
    `;

    if (schedulesHtml.length === 0) {
      timelineBodyHtml = [
        `
        <div class="card-body-schedule">
          <div class="container-fluid">
            <p class="my-4">No schedules found...</p>
          </div>
        </div>
      `,
      ];
    }
    return timelineBodyHtml;
  }

  getScheduleBodyHtml() {
    return `
    <div class="list-group">
        <div class="card">
          <div class="card-header" id="schedule-view-heading">
            <h5 class="mb-0">
              Schedules
              <button type="button" class="btn btn-danger mx-2" data-action="schedule-delete" style="float: right">Delete Schedule</button>
              <button type="button" class="btn btn-info mx-2" data-action="schedule-download" style="float: right;">Download Schedule</button>
              <div class="btn-group mx-2" style="float: right">
                <button type="button" class="btn btn-success" data-action="schedule-create" data-type="normal" style="float: right">Create Schedule (Normal)</button>
                <button type="button" class="btn btn-success dropdown-toggle dropdown-toggle-split" data-bs-toggle="dropdown" aria-expanded="false">
                  <span class="visually-hidden">Toggle Dropdown</span>
                </button>
                <ul class="dropdown-menu">
                  <li><a class="dropdown-item" data-action="schedule-normal">Create Schedule (Normal)</a></li>
                  <li><a class="dropdown-item" data-action="schedule-fresh">Create Schedule (Fresh)</a></li>
                </ul>
              </div>
            </h5>
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
          <div class="row my-3">
            <div class="col-2">
              <div class="d-flex w-100 justify-content-between">
                <p class="mb-1">${schedule.vehicle_id}</p>
              </div>
              <div class="d-flex w-100 justify-content-between">
                <p class="mb-1">${JSON.stringify(schedule.vehicle_data)}</p>
              </div>
              <div class="d-flex w-100 justify-content-between">
                <p class="mb-1">
                  <button class="btn btn-primary" data-action="vehicle-schedule-download" data-id="${
                    schedule.vehicle_id
                  }">Download Schedule</button>
                </p>
              </div>
              <!-- <div class="d-flex w-100 justify-content-between">
                <p class="mb-1">
                <button class="btn btn-primary" data-action="play-route" data-id="${
                  schedule.vehicle_id
                }">Play Route</button>
                </p>
              </div> -->
            </div>
            <div class="col-10">
              <div class="timelines-container" style="padding-bottom: 20px;">
                <div class="timeline-item flex-container items-center">
                  <div class="timeline flex-main">
                    <div class="line"></div>
                    <div class="value-line" style="background-color: ${color}; width: ${
      width * widthFactor
    }%; left: 0%;"></div>
                    ${tasksHtml}
                    <div>
                      ${verticalLinesHtml}
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

    let tasksHtml = route.map((task, index) => {
      let color =
        task.task_id === "-1"
          ? defaultColor
          : Random.getRandomColor(task.task_id);

      let arrivalHours =
        new Date(task.arrival + "Z").getTime() / 1000 / 60 / 60;
      let departureHours =
        new Date(task.departure + "Z").getTime() / 1000 / 60 / 60;

      // take the difference between the start time and the min hours to get the left position
      let left = (arrivalHours - minHours) * 10;

      // duration is the difference between the start time and the end time
      let duration = (departureHours - arrivalHours) * 10;

      // It is displayed centred from left position, so we need to add half the duration
      left = left + duration / 2;

      let widthString = `${duration * widthFactor}%;`;
      if (widthString === "0%;") {
        widthString = "10px";
      }

      let taskType = task.type.charAt(0).toUpperCase() + task.type.slice(1);
      let taskId = "";
      if (taskType != "Start" && taskType != "End") {
        taskId = `ID: ${task.task_id}<br/>`;
      }

      // set tooltip text for the task
      let tooltipText = `
        <b>${taskType}</b><br/>
        ${taskId}
        Arr: ${task.arrival}<br/>
        Dep: ${task.departure}<br/>
        Travel time: ${task.travel_time}<br/>
        Waiting time: ${task.waiting_time}<br/>
        Setup time: ${task.setup_time}<br/>
        Service time: ${task.service_time}<br/>
        Data: ${JSON.stringify(task.task_data)}<br/>
      `;

      let zIndex = 10 + route.length - index;

      // get the html for the task
      let html = `
        <div class="task-item" style="left: ${
          left * widthFactor
        }%; width: ${widthString}; z-index: ${zIndex}">
          <a data-bs-toggle="tooltip" data-bs-placement="bottom" data-bs-html="true" title='${tooltipText}'>
            <div class="task-item-full" style="background-color: ${color};"></div>
          </a>
        </div>
      `;

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

  getLabelHtml(minHours, maxHours) {
    let timeHtml = "";
    let dateHtml = "";

    let width = (maxHours - minHours) * 10;
    let widthFactor = 1;
    if (width < 100) {
      widthFactor = 100 / width;
    }

    // get the number of vertical lines
    const numberOfLines = maxHours - minHours;

    // get the html for the vertical lines
    for (let i = 0; i <= numberOfLines; i++) {
      // set time as minHours + i hours in epoch, and convert it to date
      let totalHours = minHours + i;
      let dateTime = new Date(totalHours * 3600 * 1000)
        .toISOString()
        .split("T");

      let date = dateTime[0];
      let time = dateTime[1];

      // convert time from hh:mm:ss.sss to hh:mm
      time = time.split(":");
      time = time[0] + ":" + time[1];

      let offset = 10 * i * widthFactor;

      dateHtml += `
        <div class="label-item" style="left: ${offset}%; white-space: nowrap;">${date}</div>
      `;

      timeHtml += `
        <div class="label-item" style="left: ${offset}%; white-space: nowrap;">${time}</div>
      `;
    }

    // return the html for the labels
    return [dateHtml, timeHtml];
  }

  getMinMaxHours(schedule) {
    if (schedule.length === 0) {
      return 0;
    }

    let minDate = new Date(9999, 1, 1, 0, 0, 0, 0);
    let maxDate = new Date(0, 1, 1, 0, 0, 0, 0);

    schedule.forEach((vehicle) => {
      if (new Date(vehicle.route[0].arrival + "Z") < minDate) {
        minDate = new Date(vehicle.route[0].arrival + "Z");
      }
      if (
        new Date(vehicle.route[vehicle.route.length - 1].arrival + "Z") >
        maxDate
      ) {
        maxDate = new Date(
          vehicle.route[vehicle.route.length - 1].arrival + "Z"
        );
      }
    });

    // get hours since epoch
    const minDateHours = Math.floor(minDate.getTime() / 1000 / 60 / 60);
    const maxDateHours = Math.ceil(maxDate.getTime() / 1000 / 60 / 60);

    return [minDateHours, maxDateHours];
  }

  getEmptySchedule() {
    return {
      metadata: {
        summary: [],
        total_service: "00:00:00",
        total_setup: "00:00:00",
        total_travel: "00:00:00",
        total_waiting: "00:00:00",
        unassigned: [],
      },
      schedule: [],
    };
  }

  handlers() {
    return {
      onScheduleCreate: (data) => {
        this.data = data;
        this.mapView.deleteAllNumberPointers();
        this.mapView.deleteAllUnassignedPointers();
        this.mapView.deleteAllRouteLayers();
        this.render(true);
      },
      onScheduleDelete: () => {
        this.data = this.getEmptySchedule();
        this.render(true);
        this.mapView.deleteAllNumberPointers();
        this.mapView.deleteAllUnassignedPointers();
        this.mapView.deleteAllRouteLayers();
      },
      onPlayRoute: (vehicleID) => {
        this.mapView.playRoute(vehicleID);
      },
      onStopPlayRoute: (vehicleID) => {
        this.mapView.stopPlayRoute(vehicleID);
      },
    };
  }
}
