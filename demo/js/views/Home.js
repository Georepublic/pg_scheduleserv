import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Home");
    this.setHtml(this.getHtml());
  }

  getHtml() {
    return `
        <div class="heading">
          <h2>Welcome to pg_scheduleserv demo application</h2>
        </div>
        <div class="sub-heading">
          <a href="/projects" type="button" class="btn btn-outline-primary" data-link>View all Projects</a>
        </div>
        `;
  }
}
