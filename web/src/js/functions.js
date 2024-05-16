import $ from "jquery";

export function initSelect2(selector) {
  try {
    $(selector).select2();
  } catch (e) {
    console.error(`failed to init select2 on ${selector}: ${e}`);
  }
}
