export default class Random {
  static getRandomColor(id) {
    const colors = [
      "#FF0000",
      "#FF8000",
      "#FFFF00",
      "#80FF00",
      "#00FF00",
      "#00FF80",
      "#00FFFF",
      "#0080FF",
      "#0000FF",
      "#7F00FF",
      "#FF00FF",
      "#FF007F",
    ];
    return colors[id % colors.length];
  }
}
