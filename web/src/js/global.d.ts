interface Window {
  $: typeof import("jquery");
  htmx: typeof import("htmx.org");
  initSelect2: (selector: string) => void;
  NotifySuccess: (message: string) => void;
  NotifyInfo: (message: string) => void;
  NotifyError: (message: string) => void;
}
