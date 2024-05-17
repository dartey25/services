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
