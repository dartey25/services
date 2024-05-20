// md core
import "@mdoffice/md-ui/css/bootstrap.css";
import "@mdoffice/md-ui/css/bootstrap_limitless.css";
import "@mdoffice/md-ui/css/cci.color.css";
import "@mdoffice/md-ui/css/icomoon.css";
import "@mdoffice/md-ui/css/layout.css";
import "@mdoffice/md-ui/css/components.css";
import "@mdoffice/md-ui/css/colors.css";
import "@mdoffice/md-ui/css/limitless.datepicker.css";

// libs
import $ from "jquery";
import "bootstrap";
import htmx from "htmx.org";

// custom scripts
import "../css/style.css";
import { fabShow, fabHide, fabInit } from "./fab";
import { scrollToTop, initSelect2 } from "./util";
import { NotifyError, NotifyInfo, NotifySuccess } from "./notify";

// Different event handlers and initializations for the page
$(function () {
  htmx.on("htmx:sendError", () => {
    NotifyError("Сервер не відповідає");
  });

  htmx.on("htmx:responseError", (ev) => {
    //@ts-ignore
    if (ev.detail.xhr.status >= 500) {
      NotifyError("Помилка на сервері");
    }
  });
  fabInit();
});

// Global window functions (make sure to add their types in .d.ts file to make linter happy)
window.$;
window.htmx = htmx;
window.initSelect2 = initSelect2;
window.NotifySuccess = NotifySuccess;
window.NotifyInfo = NotifyInfo;
window.NotifyError = NotifyError;
window.scrollToTop = scrollToTop;
window.fabShow = fabShow;
window.fabHide = fabHide;
