export default class Random {
  static getRandomColor(id) {
    id = id.slice(-15);
    const hex = parseInt(id, 10).toString(16).slice(-6);
    const hexPadded = "000000".substring(0, 6 - hex.length) + hex;
    return `#${hexPadded}`;
  }
}
