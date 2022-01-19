export default class {
  constructor(params) {
    this.params = params;
  }

  setTitle(title) {
    document.title = title;
  }

  setHtml(html) {
    document.querySelector("#app").innerHTML = html;
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
