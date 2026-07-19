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
    body: "Everything lives in this sidebar — dynamic modules, forms, tables, and data tools.",
    target: '[data-tour="sidebar-nav"]',
    route: "/dashboard",
    placement: "right",
  },
  {
    key: "forms",
    title: "Dynamic forms",
    body: "Create records with forms generated from module metadata and a backend validation schema.",
    target: '[data-tour="nav-forms"]',
    placement: "right",
  },
  {
    key: "tables",
    title: "Dynamic tables",
    body: "Metadata-driven tables with sorting, filtering, and saved views that persist per module.",
    target: '[data-tour="nav-tables"]',
    placement: "right",
  },
  {
    key: "imports",
    title: "Import engine",
    body: "Bring in CSV or Excel files. Columns are auto-mapped and rows are validated and processed in the background.",
    target: '[data-tour="nav-imports"]',
    placement: "right",
  },
  {
    key: "exports",
    title: "Export engine",
    body: "Export any module to CSV or Excel — instantly or as a background job — and reuse saved export templates.",
    target: '[data-tour="nav-exports"]',
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
