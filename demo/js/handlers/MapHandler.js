export default class {
  constructor(mapActivated, { onToggleMapClick }) {
    this.handleToggleMapClick(mapActivated, onToggleMapClick);
  }

  handleToggleMapClick(mapActivated, onToggleMapClick) {
    document.addEventListener("click", (event) => {
      const el = event.target.closest(`[data-action="toggle-map-click"]`);
      if (el) {
        mapActivated = !mapActivated;
        if (mapActivated) {
          el.innerText = "Done";
        } else {
          el.innerText = "Choose on Map";
        }
        onToggleMapClick();
      }
    });
  }
}
