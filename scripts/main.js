// import "./core";
import "./style.css";
// import $ from "jquery";

export function initSelect2(selector) {
  try {
    $(selector).select2();
  } catch (e) {
    console.error(`failed to init select2 on ${selector}: ${e}`);
  }
}

$(() => {
  initSelect2("#country-select");
  function hello() {
    alert("hello world");
  }
});

// script showSearchResults() {
//     htmx.removeClass(htmx.find("#results"), "d-none")
// }
//
// script hideResults() {
//     htmx.addClass(htmx.find("#results"), "d-none")
//     htmx.find("#results-body").innerHTML = ""
// }
