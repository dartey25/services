import $ from "jquery";
/**
 * Retrieves an element by its ID.
 *
 * @param {string} id - The ID of the element to retrieve.
 * @returns {HTMLElement | null} The element with the specified ID, or null if no element is found.
 */
export function getElementById(id) {
  return document.getElementById(id);
}

/**
 * Scroll to top of content box
 *
 * @returns {void}
 */
export function scrollToTop() {
  document.getElementsByClassName("content")[0]?.scrollIntoView();
}

/**
 * Initialize any select2 by given jquery selector
 *
 * @param {string} selector - jquery selector
 * @returns {void}
 */
export function initSelect2(selector) {
  try {
    $(selector).select2();
  } catch (e) {
    console.error(`failed to init select2 on ${selector}: ${e}`);
  }
}
