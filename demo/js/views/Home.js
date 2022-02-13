import AbstractView from "./AbstractView.js";

export default class Home extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Home");
    this.setHeading("Welcome to pg_scheduleserv demo application");
    this.setSubHeading(`
    <a href="/projects" type="button" class="btn btn-outline-primary" data-link>View all Projects</a>
    `);
  }
}
