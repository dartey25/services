import $ from "jquery";
import Notify, { NotifyType } from "@mdoffice/md-component/Notify";

/**
 *
 * @param {string} msg - message to display
 * @returns {void}
 */
export function NotifySuccess(msg) {
  Notify.show(msg, NotifyType.success);
}

/**
 *
 * @param {string} msg - message to display
 * @returns {void}
 */
export function NotifyInfo(msg) {
  Notify.show(msg, NotifyType.info);
}

/**
 *
 * @param {string} msg - message to display
 * @returns {void}
 */
export function NotifyError(msg) {
  Notify.show(msg, NotifyType.error);
}

/**
 *
 * @param {string} text - text to copy
 * @returns {void}
 */
export function copyToClipboard(text) {
  Notify.show("success", NotifyType.success);
}
