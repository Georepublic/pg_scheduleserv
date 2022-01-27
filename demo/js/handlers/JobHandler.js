import JobAPI from "../api/JobAPI.js";

export default class JobHandler {
  constructor(
    jobs,
    emptyJob,
    {
      onJobView,
      onJobCreateClick,
      onJobEditClick,
      onJobDelete,
      onJobSave,
      onJobClose,
    }
  ) {
    // get the jobs from the params
    this.jobs = jobs;
    this.emptyJob = emptyJob;
    this.jobAPI = new JobAPI();

    this.main = document.querySelector("#main");
    this.schedule = document.querySelector("#schedule");

    this.handleJobView(onJobView);
    this.handleJobCreateClick(onJobCreateClick);
    this.handleJobEditClick(onJobEditClick);
    this.handleJobDelete(onJobDelete);
    this.handleJobSave(onJobSave);
    this.handleJobClose(onJobClose);
    this.handleJobScheduleClick(onJobView);
    this.handleJobTwForm();
  }

  // get job from id
  getJob(jobID) {
    if (!jobID) {
      return this.emptyJob;
    }
    return this.jobs.find((job) => {
      return job.id === jobID;
    });
  }

  handleJobView(onJobView) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-view"]`);
      if (el) {
        let jobID = el.dataset.id;
        onJobView(this.getJob(jobID));
      }
    });
  }

  handleJobScheduleClick(onJobView) {
    this.schedule.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-view"]`);
      if (el) {
        let jobID = el.dataset.id;
        onJobView(this.getJob(jobID));
      }
    });
  }

  handleJobCreateClick(onJobCreateClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-create"]`);
      if (el) {
        onJobCreateClick();
      }
    });
  }

  handleJobEditClick(onJobEditClick) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-edit"]`);
      if (el) {
        let jobID = el.dataset.id;
        onJobEditClick(this.getJob(jobID));
      }
    });
  }

  handleJobDelete(onJobDelete) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-delete"]`);
      if (el) {
        let jobID = el.dataset.id;
        // call the job api to delete the job
        this.jobAPI.deleteJob(jobID).then(() => {
          // remove the job from the list
          this.jobs = this.jobs.filter((job) => {
            return job.id !== jobID;
          });
          // call the onJobDelete callback
          onJobDelete(this.jobs);
        });
      }
    });
  }

  handleJobSave(onJobSave) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-save"]`);
      if (el) {
        const form = el.closest("form");
        const formData = new FormData(form);
        const job = {};
        for (const [key, value] of formData.entries()) {
          // if key contains [] then it is an array
          if (key.includes("[]")) {
            if (!job[key]) {
              job[key] = [];
            }
            job[key].push(value);
          } else {
            job[key] = value;
          }
        }
        const id = job["id"];

        this.jobAPI.saveJob(job).then((job) => {
          // edit the job in the list, or append a new job to the list depending on the id
          if (id) {
            // update the job
            this.jobs = this.jobs.map((oldJob) => {
              if (oldJob.id === job.id) {
                return job;
              }
              return oldJob;
            });
          } else {
            // append the new job to the list
            this.jobs.push(job);
          }
          onJobSave(job, this.jobs);
        });
      }
    });
  }

  handleJobClose(onJobClose) {
    this.main.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-close"]`);
      if (el) {
        onJobClose();
      }
    });
  }

  handleJobTwForm() {
    this.main.addEventListener("click", (event) => {
      let el = event.target.closest(`[data-action="job-tw-form-create"]`);
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
      el = event.target.closest(`[data-action="job-tw-form-delete"]`);
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
