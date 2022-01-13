import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Projects");
        this.setHtml(this.getHtml());
    }

    getHtml() {
        return `
        You are viewing the project with ID ${this.params.id}
        `
    }
}
