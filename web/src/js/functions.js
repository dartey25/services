import $ from "jquery";
import Notify, { NotifyType } from "@mdoffice/md-component/Notify";

/**
 * Initialize any select2 by given jquery selector
 *
 * @param {string} selector - jquery selector
 */
export function initSelect2(selector) {
  try {
    //@ts-ignore
    $(selector).select2();
  } catch (e) {
    console.error(`failed to init select2 on ${selector}: ${e}`);
  }
}

/**
 *
 * @param {string} msg - message to display
 */
export function NotifySuccess(msg) {
  Notify.show(msg, NotifyType.success);
}

/**
 *
 * @param {string} msg - message to display
 */
export function NotifyInfo(msg) {
  Notify.show(msg, NotifyType.info);
}

/**
 *
 * @param {string} msg - message to display
 */
export function NotifyError(msg) {
  Notify.show(msg, NotifyType.error);
}

/**
 *
 * @param {string} text - text to copy
 */
export function copyToClipboard(text) {
  Notify.show("success", NotifyType.success);
}

/**
 * Scroll to top of content box
 *
 */
export function scrollToTop() {
  document.getElementsByClassName("content")[0]?.scrollIntoView();
}
