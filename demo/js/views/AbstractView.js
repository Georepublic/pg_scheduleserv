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
}
