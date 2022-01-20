import JobAPI from "../api/JobAPI.js";

export default class {
  constructor(jobs, emptyJob, { onJobView, onJobCreateClick, onJobEditClick, onJobDelete, onJobSave, onJobClose }) {
    // get the jobs from the params
    this.jobs = jobs;
    this.emptyJob = emptyJob;
    this.jobAPI = new JobAPI();

    this.handleJobView(onJobView);
    this.handleJobCreateClick(onJobCreateClick);
    this.handleJobEditClick(onJobEditClick);
    this.handleJobDelete(onJobDelete);
    this.handleJobSave(onJobSave);
    this.handleJobClose(onJobClose);
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
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-view"]`);
      if (el) {
        let jobID = el.dataset.id;
        onJobView(this.getJob(jobID));
      }
    });
  }

  handleJobCreateClick(onJobCreateClick) {
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-create"]`);
      if (el) {
        onJobCreateClick();
      }
    });
  }

  handleJobEditClick(onJobEditClick) {
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-edit"]`);
      if (el) {
        let jobID = el.dataset.id;
        onJobEditClick(this.getJob(jobID));
      }
    });
  }

  handleJobDelete(onJobDelete) {
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-delete"]`);
      if (el) {
        let jobID = el.dataset.id;
        // call the job api to delete the job
        this.jobAPI.deleteJob(jobID)
          .then(() => {
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
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-save"]`);
      if (el) {
        const form = el.closest("form");
        const formData = new FormData(form);
        const job = {};
        for (const [key, value] of formData.entries()) {
          job[key] = value;
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
        }).catch(err => {
          console.log(err);
        })
      }
    });
  }

  handleJobClose(onJobClose) {
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="job-close"]`);
      if (el) {
        onJobClose();
      }
    });
  }
}
