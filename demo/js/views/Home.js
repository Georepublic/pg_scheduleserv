import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Home");
    this.setHtml(this.getHtml());
  }

  getHtml() {
    return `
        <center><h2>Welcome to pg_scheduleserv demo application</h2></center>
        `;
  }
}
