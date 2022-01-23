export default class AbstractView {
  constructor(params, init = true) {
    this.params = params;
    if (init) {
      this.initHtml();
    }
  }

  initHtml() {
    let html = `
    <div class="row" id="main">
      <div class="col-md-3" id="app-left"></div>
      <div class="col-md-6" id="app"></div>
      <div class="col-md-3" id="app-right"></div>
    </div>
    `;
    document.querySelector("#main").outerHTML = html;
  }

  setTitle(title) {
    document.title = title;
  }

  setHtml(html) {
    document.querySelector("#app").innerHTML = html;
  }

  appendHtmlLeft(html) {
    document.querySelector("#app-left").innerHTML += html;
  }

  setHtmlLeft(html) {
    document.querySelector("#app-left").innerHTML = html;
  }

  setHtmlRight(html) {
    document.querySelector("#app-right").innerHTML = html;
  }

  setHeading(heading) {
    let html = `<h2>${heading}</h2>`;
    document.querySelector("#heading").innerHTML = html;
  }

  setSubHeading(subHeading) {
    document.querySelector("#sub-heading").innerHTML = subHeading;
  }
}
