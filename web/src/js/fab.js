import { getElementById } from "./util";

const FabID = "fab";
const ScrollThreshold = 300;

/**
 * Init FAB Element to toggle visibility on scroll
 *
 */
export function fabInit() {
  const content = document.querySelector(".content-inner");
  if (!content) {
    console.warn("Fab initialization failed. Reason: no .content");
    return;
  }

  content.addEventListener("scroll", () => {
    if (content.scrollTop > ScrollThreshold) {
      fabShow();
    } else {
      fabHide();
    }
  });
}

/**
 * Show FAB Element
 *
 */
export function fabShow() {
  getElementById(FabID)?.classList.remove("d-none");
}

/**
 * Hide FAB Element
 *
 */
export function fabHide() {
  getElementById(FabID)?.classList.add("d-none");
}
