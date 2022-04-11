import pluginJson from './plugin.json';
import { NavItem } from './types';

export const PLUGIN_ID = `${pluginJson.id}`;
export const PLUGIN_BASE_URL = `/a/${PLUGIN_ID}`;

export enum ROUTES {
  REPORT_GROUP = 'report-groups',
  SCHEDULERS = 'schedulers',
}

export const NAVIGATION_TITLE = 'Excel report e-mail scheduler';
export const NAVIGATION_SUBTITLE = `Generate Excel reports from mSupply dashboard. Send the reports to custom created user-groups on pre-defined schedule.`;

// Add a navigation item for each route you would like to display in the navigation bar
export const NAVIGATION: Record<string, NavItem> = {
  [ROUTES.REPORT_GROUP]: {
    id: ROUTES.REPORT_GROUP,
    text: 'Report Groups',
    sub: 'Contain users to whom email reports would sent',
    icon: 'users-alt',
    url: `${PLUGIN_BASE_URL}/report-groups`,
  },
  [ROUTES.SCHEDULERS]: {
    id: ROUTES.SCHEDULERS,
    text: 'Schedulers',
    sub: "Schedules to send emails with selected panels's data to the selected user groups",
    icon: 'schedule',
    url: `${PLUGIN_BASE_URL}/schedulers`,
  },
};
