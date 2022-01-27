import JobAPI from "../api/JobAPI.js";
import JobHandler from "../handlers/JobHandler.js";
import Random from "../utils/Random.js";
import AbstractView from "./AbstractView.js";

export default class JobView extends AbstractView {
  constructor(params) {
    super(params, false);

    // get the jobs from the params
    this.jobs = params.jobs;
    this.projectID = params.projectID;
    this.mapView = params.mapView;
    this.jobAPI = new JobAPI();

    this.jobLeftDiv = document.createElement("div");
    document.querySelector("#app-left").appendChild(this.jobLeftDiv);

    this.handler = new JobHandler(
      params.jobs,
      this.getEmptyJob(),
      this.handlers()
    );
  }

  // render the jobs for this project
  render() {
    // get the html for the jobs
    let jobsHtml = this.getJobsHtml();

    // set the html for the jobs
    this.jobLeftDiv.innerHTML = jobsHtml;

    this.mapView.addJobMarkers(this.jobs);
    this.mapView.fitAllMarkers();
  }

  // get the html for the jobs
  getJobsHtml() {
    // get the html for each job
    let jobsHtml = this.jobs.map((job) => {
      return this.getJobHtml(job);
    });

    if (jobsHtml.length === 0) {
      jobsHtml = [
        `
        <div class="list-group-item flex-column align-items-start">
          <p class="mb-1">No jobs found...</p>
        </div>
      `,
      ];
    }

    return `
      <div class="list-group">
        <div class="card">
          <div class="card-header job-view-heading">
            <h5 class="mb-0">
              Jobs
              <button type="button" class="btn btn-success" data-action="job-create" style="float: right">Add</button>
            </h5>
          </div>
          <div class="card-body-custom">
            ${jobsHtml.join("")}
          </div>
        </div>
      </div>
    `;
  }

  // get the html for the job
  getJobHtml(job) {
    const color = Random.getRandomColor(job.id);
    let html = `
      <div style="background-color: ${color}">
      <div class="list-group-item flex-column align-items-start" data-action="job-view" data-id="${
        job.id
      }">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">${job.id}</h5>
        </div>

        <div class="d-flex w-100 justify-content-between">
          <p class="mb-1">${JSON.stringify(job.data)}</p>
        </div>
      </div>
      </div>
    `;

    // return the html for the job
    return html;
  }

  // get the time windows html for the job, time windows is a 2D array of start and end times
  getTimeWindowsHtml(time_windows) {
    let timeWindowsHtml = time_windows.map((timeWindow) => {
      return `
          <li>${timeWindow[0]} - ${timeWindow[1]}</li>
      `;
    });
    let timeWindows = timeWindowsHtml.join("");
    if (timeWindows == "") {
      timeWindows = "<li>No time window</li>";
    }
    return `
      <div>
        <ul>
        ${timeWindows}
        </ul>
      </div>
    `;
  }

  getCompleteJobHtml(job) {
    let timeWindowsHtml = this.getTimeWindowsHtml(job.time_windows);
    let html = `
      <div class="card">
        <div class="card-header job-view-heading">
          <h5 class="mb-0">
            Job
            <button type="button" class="btn btn-danger" data-action="job-close">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">ID: ${job.id}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Location (Lat, Lon): ${job.location.latitude}, ${
      job.location.longitude
    }</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Setup: ${job.setup}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Service: ${job.service}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Delivery: [${job.delivery}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Pickup: [${job.pickup}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Skills: [${job.skills}]</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Priority: ${job.priority}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Time Windows:</p>
          </div>
          ${timeWindowsHtml}
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Project ID: ${job.project_id}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Data: ${JSON.stringify(job.data)}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Created At: ${job.created_at}</p>
          </div>
          <div class="d-flex w-100 justify-content-between">
            <p class="mb-1">Updated At: ${job.updated_at}</p>
          </div>
          <div class="d-flex w-100 justify-content-center">
            <button class="btn btn-primary mx-2" data-action="job-edit" data-id="${
              job.id
            }">Edit</button>
            <button class="btn btn-danger mx-2" data-action="job-delete" data-id="${
              job.id
            }">Delete</button>
          </div>
        </div>
      </div>
    `;

    // return the html for the job
    return html;
  }

  getTimeWindowsFormHtml(time_windows) {
    let timeWindowsHtml = time_windows.map((timeWindow) => {
      return `
        <div class="input-group">
          <input type="datetime-local" class="form-control" name="tw_open[]" value="${timeWindow[0]}" step="1" style="font-size: 13px;">
          <input type="datetime-local" class="form-control" name="tw_close[]" value="${timeWindow[1]}" step="1" style="font-size: 13px;">
        </div>
      `;
    });
    return `
      <div class="form-group">
      <label>Time Window (Open and Close)</label>
      ${timeWindowsHtml.join("")}
      <div class="text-center">
        <button type="button" class="btn btn-outline-primary" data-action="job-tw-form-create">
          <i class="fas fa-plus-circle"></i>
        </button>
        <button type="button" class="btn btn-outline-danger" data-action="job-tw-form-delete">
          <i class="fa-solid fa-trash"></i>
        </button>
      </div>
    </div>
    `;
  }

  getJobFormHtml(job) {
    // get time windows html for the job
    let timeWindowsHtml = this.getTimeWindowsFormHtml(job.time_windows);

    let html = `
      <div class="card">
        <div class="card-header job-view-heading">
          <h5 class="mb-0">
            Job
            <button type="button" class="btn btn-danger" data-action="job-close" data-id="${
              job.id
            }">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <form>
            <input type="hidden" name="id" value="${job.id}">
            <input type="hidden" name="project_id" value="${job.project_id}">
            <div class="form-group">
              <label>Location (Lat, Lon)</label>
              <input type="text" class="form-control" name="location" value="${
                job.location.latitude
              }, ${job.location.longitude}" data-action="job-location-change">
            </div>
            <div class="form-group">
              <label>Setup</label>
              <input type="time" class="form-control" name="setup" value="${
                job.setup
              }" step="1">
            </div>
            <div class="form-group">
              <label>Service</label>
              <input type="time" class="form-control" name="service" value="${
                job.service
              }" step="1">
            </div>
            <div class="form-group">
              <label>Delivery</label>
              <input type="text" class="form-control" name="delivery" value="${
                job.delivery
              }">
            </div>
            <div class="form-group">
              <label>Pickup</label>
              <input type="text" class="form-control" name="pickup" value="${
                job.pickup
              }">
            </div>
            <div class="form-group">
              <label>Skills</label>
              <input type="text" class="form-control" name="skills" value="${
                job.skills
              }">
            </div>
            <div class="form-group">
              <label>Priority</label>
              <input type="number" class="form-control" name="priority" min="0" max="100" value="${
                job.priority
              }">
            </div>
            ${timeWindowsHtml}
            <div class="form-group">
              <label>Data</label>
              <input type="text" class="form-control" name="data" value='${JSON.stringify(
                job.data
              )}'>
            </div>
            <div class="d-flex w-100 justify-content-center">
              <button type="button" class="btn btn-primary mx-2" data-action="job-save" data-id="${
                job.id
              }">Save</button>
              <button type="button" class="btn btn-warning mx-2" data-action="job-edit" data-id="${
                job.id
              }">Reset</button>
              <button type="button" class="btn btn-danger mx-2" data-action="job-view" data-id="${
                job.id
              }">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    `;

    // return the html for the job
    return html;
  }

  selectJob(jobID) {
    this.deselectAll();
    let jobViewElement = document.querySelector(
      `[data-action="job-view"][data-id="${jobID}"]`
    );
    jobViewElement.classList.add("active");

    // move the element into view
    jobViewElement.scrollIntoView({ behavior: "smooth", block: "nearest" });
  }

  deselectAll() {
    // for all elements in query selector, remove their active class
    document.querySelectorAll(`.list-group-item.active`).forEach((element) => {
      element.classList.remove("active");
    });
  }

  getEmptyJob() {
    // get map center
    let coordinates = this.mapView.getCenter();
    return {
      id: "",
      location: {
        latitude: coordinates[0],
        longitude: coordinates[1],
      },
      setup: "00:00:00",
      service: "00:00:00",
      delivery: "",
      pickup: "",
      skills: "",
      priority: "0",
      project_id: this.projectID,
      data: {},
      created_at: "",
      updated_at: "",
      time_windows: [],
    };
  }

  handlers() {
    return {
      onJobView: (job) => {
        if (!job.id) {
          this.setHtmlRight("");
          return;
        }
        // get the complete html for the job
        let jobHtml = this.getCompleteJobHtml(job);

        // set the html for the job
        this.setHtmlRight(jobHtml);

        // select the job
        this.selectJob(job.id);

        this.mapView.addJobMapPointer(
          job.location.latitude,
          job.location.longitude
        );
        this.mapView.deactivateMap();
      },
      onJobCreateClick: () => {
        this.deselectAll();

        const job = this.getEmptyJob();

        // create the job form html with empty job
        let jobHtml = this.getJobFormHtml(job);

        // set the html for the job
        this.setHtmlRight(jobHtml);

        this.mapView.addJobMapPointer(
          job.location.latitude,
          job.location.longitude
        );
        this.mapView.activateMap();
        this.mapView.fitAllMarkers();
      },
      onJobEditClick: (job) => {
        // get the complete html for the job
        let jobHtml = this.getJobFormHtml(job);

        // set the html for the job
        this.setHtmlRight(jobHtml);

        this.mapView.addJobMapPointer(
          job.location.latitude,
          job.location.longitude
        );
        this.mapView.activateMap();
      },
      onJobSave: (job, newJobs) => {
        this.jobs = newJobs;
        this.render();
        this.handlers().onJobView(job);
      },
      onJobDelete: (newJobs) => {
        this.deselectAll();
        this.setHtmlRight("");
        this.jobs = newJobs;
        this.render();
        this.mapView.removeMapPointers();
        this.mapView.deactivateMap();
        this.mapView.fitAllMarkers();
      },
      onJobClose: () => {
        this.deselectAll();
        this.setHtmlRight("");
        this.mapView.removeMapPointers();
        this.mapView.deactivateMap();
        this.mapView.fitAllMarkers();
      },
    };
  }
}
