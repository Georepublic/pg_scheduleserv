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
    this.selectedJob = null;
    this.jobAPI = new JobAPI();

    this.handler = new JobHandler(params.jobs, this.getEmptyJob(), this.handlers());
  }

  // render the jobs for this project
  render() {
    // get the html for the jobs
    let jobsHtml = this.getJobsHtml();

    // set the html for the jobs
    this.setHtmlLeft(jobsHtml);

    if (this.markers) {
      this.mapView.removeMarkers(this.markers);
    }

    this.markers = {};
    this.jobs.forEach((job) => {
      const randomColor = Random.getRandomColor(job.id);
      let marker = this.mapView.addMarker(job.location.latitude, job.location.longitude, "house", randomColor);

      marker.on('click', () => {
        this.handlers().onJobView(job);
      });

      this.markers[job.id] = marker;
      this.mapView.fitMarkers(this.markers);
    });
  }

  // get the html for the jobs
  getJobsHtml() {
    // get the html for each job
    let jobsHtml = this.jobs.map((job) => {
      return this.getJobHtml(job);
    });

    if (jobsHtml.length === 0) {
      jobsHtml = [`
        <div class="list-group-item flex-column align-items-start">
          <p class="mb-1">No jobs found...</p>
        </div>
      `];
    }

    // return the html for the jobs, with card heading of jobs with max height of 30vh and scrolling
    return `
      <div class="list-group">
        <div class="card">
          <div class="card-header job-view-heading">
            <h5 class="mb-0">
              Jobs
              <button type="button" class="btn btn-success" data-action="job-create" style="float: right">Add</button>
            </h5>
          </div>
          <div style="max-height: 30vh; overflow-y: scroll;">
            ${jobsHtml.join("")}
          </div>
        </div>
      </div>
    `;
  }

  // get the html for the job
  getJobHtml(job) {
    let html = `
      <div class="list-group-item flex-column align-items-start" data-action="job-view" data-id="${job.id}">
        <div class="d-flex w-100 justify-content-between">
          <h5 class="mb-1">${job.id}</h5>
        </div>

        <div class="d-flex w-100 justify-content-between">
          <p class="mb-1">${JSON.stringify(job.data)}</p>
        </div>
      </div>
    `;

    // return the html for the job
    return html;
  }

  getCompleteJobHtml(job) {
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
            <p class="mb-1">Location (Lat, Lon): ${job.location.latitude}, ${job.location.longitude}</p>
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
            <button class="btn btn-primary mx-2" data-action="job-edit" data-id="${job.id}">Edit</button>
            <button class="btn btn-danger mx-2" data-action="job-delete" data-id="${job.id}">Delete</button>
          </div>
        </div>
      </div>
    `;

    // return the html for the job
    return html;
  }

  getJobFormHtml(job) {
    let html = `
      <div class="card">
        <div class="card-header job-view-heading">
          <h5 class="mb-0">
            Edit Job
            <button type="button" class="btn btn-danger" data-action="job-close" data-id="${job.id}">
              <i class="fas fa-times"></i>
            </button>
          </h5>
        </div>
        <div class="card-body">
          <form>
            <input type="hidden" name="id" value="${job.id}">
            <input type="hidden" name="project_id" value="${job.project_id}">
            <div class="form-group">
              <label for="edit-job-location">Location (Lat, Lon)</label>
              <input type="text" class="form-control" name="location" value="${job.location.latitude}, ${job.location.longitude}" data-action="location-change">
            </div>
            <div class="form-group">
              <label for="edit-job-setup">Setup</label>
              <input type="time" class="form-control" name="setup" value="${job.setup}" step="1">
            </div>
            <div class="form-group">
              <label for="edit-job-service">Service</label>
              <input type="time" class="form-control" name="service" value="${job.service}" step="1">
            </div>
            <div class="form-group">
              <label for="edit-job-delivery">Delivery</label>
              <input type="text" class="form-control" name="delivery" value="${job.delivery}">
            </div>
            <div class="form-group">
              <label for="edit-job-pickup">Pickup</label>
              <input type="text" class="form-control" name="pickup" value="${job.pickup}">
            </div>
            <div class="form-group">
              <label for="edit-job-skills">Skills</label>
              <input type="text" class="form-control" name="skills" value="${job.skills}">
            </div>
            <div class="form-group">
              <label for="edit-job-priority">Priority</label>
              <input type="number" class="form-control" name="priority" min="0" max="100" value="${job.priority}">
            </div>
            <div class="form-group">
              <label for="edit-job-data">Data</label>
              <input type="text" class="form-control" name="data" value='${JSON.stringify(job.data)}'>
            </div>
            <div class="d-flex w-100 justify-content-center">
              <button type="button" class="btn btn-primary mx-2" data-action="job-save" data-id="${job.id}">Save</button>
              <button type="button" class="btn btn-warning mx-2" data-action="job-edit" data-id="${job.id}">Reset</button>
              <button type="button" class="btn btn-danger mx-2" data-action="job-view" data-id="${job.id}">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    `;

    // return the html for the job
    return html;
  }

  selectJob(jobID) {
    this.deselectJob();
    let jobViewElement = document.querySelector(`[data-action="job-view"][data-id="${jobID}"]`);
    jobViewElement.classList.add("active");
    this.selectedJob = jobID;

    // move the element into view
    jobViewElement.scrollIntoView({ behavior: "smooth", block: "nearest" });
  }

  deselectJob() {
    if (this.selectedJob) {
      let jobViewElement = document.querySelector(`[data-action="job-view"][data-id="${this.selectedJob}"]`);
      jobViewElement.classList.remove("active");
      this.selectedJob = null;
    }
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

        this.mapView.removeMapPointer();
        this.mapView.deactivateMap();
        this.mapView.setCenter(job.location.latitude, job.location.longitude);
        this.mapView.addMapPointer(job.location.latitude, job.location.longitude);
      },
      onJobCreateClick: () => {
        this.deselectJob();

        const job = this.getEmptyJob();

        // create the job edit html with empty job
        let jobHtml = this.getJobFormHtml(job);

        // set the html for the job
        this.setHtmlRight(jobHtml);

        this.mapView.removeMapPointer();
        this.mapView.fitMarkers(this.markers);
        this.mapView.addMarkerOnClick(job.location.latitude, job.location.longitude);
        this.mapView.activateMap();
      },
      onJobEditClick: (job) => {
        // get the complete html for the job
        let jobHtml = this.getJobFormHtml(job);

        // set the html for the job
        this.setHtmlRight(jobHtml);

        this.mapView.removeMapPointer();
        this.mapView.setCenter(job.location.latitude, job.location.longitude);
        this.mapView.addMarkerOnClick(job.location.latitude, job.location.longitude);
        this.mapView.activateMap();
      },
      onJobSave: (job, newJobs) => {
        this.jobs = newJobs;
        this.render();
        this.handlers().onJobView(job);
      },
      onJobDelete: (newJobs) => {
        this.deselectJob();
        this.setHtmlRight("");
        this.jobs = newJobs;
        this.render();
        this.mapView.removeMapPointer();
        this.mapView.deactivateMap();
      },
      onJobClose: () => {
        this.deselectJob();
        this.setHtmlRight("");
        this.mapView.removeMapPointer();
        this.mapView.deactivateMap();
        this.mapView.fitMarkers(this.markers);
      },
    };
  }
}
