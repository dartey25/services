import "../css/style.css";
import {
  NotifyError,
  NotifyInfo,
  NotifySuccess,
  initSelect2,
} from "./functions";
import "@mdoffice/md-ui/css/bootstrap.css";
import "@mdoffice/md-ui/css/bootstrap_limitless.css";
import "@mdoffice/md-ui/css/cci.color.css";
import "@mdoffice/md-ui/css/icomoon.css";
import "@mdoffice/md-ui/css/layout.css";
import "@mdoffice/md-ui/css/components.css";
import "@mdoffice/md-ui/css/colors.css";
import "@mdoffice/md-ui/css/limitless.datepicker.css";

import $ from "jquery";
import "bootstrap";
import htmx from "htmx.org";

/**
 * Different event handlers for the page
 */
$(function () {
  htmx.on("htmx:sendError", () => {
    NotifyError("Сервер не відповідає");
  });
});

/**
 * Global window functions (make sure to add their types in .d.ts file to make linter happy)
 */
window.$;
window.htmx = htmx;
window.initSelect2 = initSelect2;
window.NotifySuccess = NotifySuccess;
window.NotifyInfo = NotifyInfo;
window.NotifyError = NotifyError;
