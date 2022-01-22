import Home from "./views/Home.js";
import ProjectList from "./views/ProjectList.js";
import Project from "./views/Project.js";

const pathToRegex = (path) =>
  new RegExp(
    "^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "([^/]+)") + "$"
  );

const getParams = (match) => {
  const values = match.result.slice(1);
  const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(
    (result) => result[1]
  );

  return Object.fromEntries(
    keys.map((key, i) => {
      return [key, values[i]];
    })
  );
};

const navigateTo = (url) => {
  history.pushState(null, null, url);
  router();
};

const router = async () => {
  const routes = [
    { path: "/", view: Home },
    { path: "/projects", view: ProjectList },
    { path: "/projects/:id", view: Project },
  ];

  const potentialMatches = routes.map((route) => {
    return {
      route: route,
      result: location.pathname.match(pathToRegex(route.path)),
    };
  });

  let match = potentialMatches.find(
    (potentialMatch) => potentialMatch.result !== null
  );

  if (!match) {
    match = {
      route: routes[0],
      result: [location.pathname],
    };
  }

  const view = new match.route.view(getParams(match));

  document.querySelectorAll(".nav-link").forEach((el) => {
    if (el.getAttribute("href") === match.route.path) {
      el.classList.add("active");
    } else {
      el.classList.remove("active");
    }
  });
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
  document.addEventListener("click", (e) => {
    const el = e.target.closest("[data-link]");
    if (el) {
      e.preventDefault();
      navigateTo(el.href);
    }
  });
  router();
});
