export default class Toast {
  static success(message) {
    iziToast.success({
      title: "Success",
      message: message,
    });
  }

  static error(message) {
    iziToast.error({
      title: "Error",
      message: message,
    });
  }

  static info(message) {
    iziToast.info({
      title: "Info",
      message: message,
    });
  }
}
