import pluginJson from './plugin.json';
import { NavItem } from './types';
import intl from 'react-intl-universal';
import { SelectableValue } from '@grafana/data';

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

export const getIntervals = () => [
  { label: intl.get('daily'), value: 0 },
  { label: intl.get('weekly'), value: 1 },
  { label: intl.get('fortnightly'), value: 2 },
  { label: intl.get('monthly'), value: 3 },
  { label: intl.get('quarterly'), value: 4 },
  { label: intl.get('yearly'), value: 5 },
];

export const getLookbacks = (): Array<SelectableValue<Number>> => [
  { label: intl.get('1day'), value: 24 * 60 * 1000 * 60 },
  { label: intl.get('2days'), value: 2 * 24 * 60 * 1000 * 60 },
  { label: intl.get('3days'), value: 3 * 24 * 60 * 1000 * 60 },
  { label: intl.get('1week'), value: 7 * 24 * 60 * 1000 * 60 },
  { label: intl.get('2weeks'), value: 14 * 24 * 60 * 1000 * 60 },
  { label: intl.get('4weeks'), value: 28 * 24 * 60 * 1000 * 60 },
  { label: intl.get('3months'), value: 28 * 3 * 24 * 60 * 1000 * 60 },
  { label: intl.get('6months'), value: 28 * 6 * 24 * 60 * 1000 * 60 },
  { label: intl.get('1year'), value: 356 * 24 * 60 * 1000 * 60 },
];
