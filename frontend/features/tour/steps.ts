// Client fallback catalogue for the orientation tour. Prefer metadata workflow
// `crm_orientation_tour` via loadOrientationSteps(); progress still uses
// tour_progress (see features/tour/api.ts).

import type { GuidedPlacement } from "@/features/guided";

export type TourPlacement = GuidedPlacement;

export interface TourStep {
  // Stable identifier persisted in `completed_steps`. Never reorder-sensitive.
  key: string;
  title: string;
  body: string;
  // CSS selector of the element to spotlight. Omit for a centered step.
  target?: string;
  // Route to navigate to before showing the step.
  route?: string;
  // Preferred placement of the tooltip relative to the target.
  placement?: TourPlacement;
}

export const APP_TOUR_KEY = "app";

export const APP_TOUR_STEPS: TourStep[] = [
  {
    key: "welcome",
    title: "Welcome to CRM Lite",
    body: "Take a 60-second tour of the workspace. You can skip anytime and restart later from your profile menu.",
    route: "/dashboard",
    placement: "center",
  },
  {
    key: "navigation",
    title: "Your workspace",
    body: "Everything lives in this sidebar — Dashboard, your modules, notifications, and Settings.",
    target: '[data-tour="sidebar-nav"]',
    route: "/dashboard",
    placement: "right",
  },
  {
    key: "forms",
    title: "Form Designer",
    body: "Settings → Form Designer previews how create forms are built from field metadata. Real records are created from each module’s Add button.",
    target: '[data-tutorial-action="open-forms"]',
    route: "/settings/forms",
    placement: "right",
  },
  {
    key: "tables",
    title: "Your modules",
    body: "Each module in the workspace sidebar has its own page — view, edit, and delete records there.",
    target: '[data-tour="sidebar-nav"]',
    route: "/dashboard",
    placement: "right",
  },
  {
    key: "imports",
    title: "Import engine",
    body: "Bring in CSV or Excel files from Settings → Import. Columns are auto-mapped and rows are validated in the background.",
    target: '[data-tutorial-action="open-imports"]',
    route: "/settings/imports",
    placement: "right",
  },
  {
    key: "exports",
    title: "Export engine",
    body: "Export any module to CSV or Excel from Settings → Export — instantly or as a background job.",
    target: '[data-tutorial-action="open-exports"]',
    route: "/settings/exports",
    placement: "right",
  },
  {
    key: "search",
    title: "Global search",
    body: "Jump to any record fast. Search spans dynamic module data from anywhere in the app.",
    target: '[data-tour="global-search"]',
    route: "/dashboard",
    placement: "bottom",
  },
  {
    key: "notifications",
    title: "Notifications",
    body: "Delivery updates for WhatsApp and email automations show up here as they are sent.",
    target: '[data-tour="notification-bell"]',
    placement: "bottom",
  },
  {
    key: "restart",
    title: "Restart anytime",
    body: "Open your profile menu to take this tour again whenever you need a refresher.",
    target: '[data-tour="user-menu"]',
    placement: "bottom",
  },
  {
    key: "done",
    title: "You're all set",
    body: "That's the whirlwind tour. Explore freely — your progress is saved automatically.",
    placement: "center",
  },
];
