import pluginJson from './plugin.json';
import { NavItem } from './types';

export const PLUGIN_BASE_URL = `/a/${pluginJson.id}`;

export enum ROUTES {
  REPORT_GROUP = 'report-groups',
  SCHEDULERS = 'schedulers',
}

export const NAVIGATION_TITLE = 'Excel report e-mail scheduler';
export const NAVIGATION_SUBTITLE = `The plugin takes data from panels of mSupply dashboard to generate excel reports. The reports are then emailed to a custom user group created with mSupply users pulled from mSupply Dashboard's datasource. The timing of the scheduler can be set in the plugin.`;

// Add a navigation item for each route you would like to display in the navigation bar
export const NAVIGATION: Record<string, NavItem> = {
  [ROUTES.REPORT_GROUP]: {
    id: ROUTES.REPORT_GROUP,
    text: 'Report Groups',
    icon: 'users-alt',
    url: `${PLUGIN_BASE_URL}/report-groups`,
  },
  [ROUTES.SCHEDULERS]: {
    id: ROUTES.SCHEDULERS,
    text: 'Schedulers',
    icon: 'schedule',
    url: `${PLUGIN_BASE_URL}/schedulers`,
  },
};
