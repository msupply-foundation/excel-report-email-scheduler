import pluginJson from './plugin.json';
import { NavItem } from './types';
import intl from 'react-intl-universal';
import { SelectableValue } from '@grafana/data';

export const PLUGIN_ID = `${pluginJson.id}`;
export const PLUGIN_BASE_URL = `/a/${PLUGIN_ID}`;

export enum ROUTES {
  REPORT_GROUP = 'report-groups',
  SCHEDULES = 'schedules',
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
  [ROUTES.SCHEDULES]: {
    id: ROUTES.SCHEDULES,
    text: 'Schedules',
    sub: "Schedules to send emails with selected panels's data to the selected user groups",
    icon: 'schedule',
    url: `${PLUGIN_BASE_URL}/schedules`,
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

export const getLookbacks = (): Array<SelectableValue<String>> => [
  { label: intl.get('1day'), value: 'now-1d' },
  { label: intl.get('2days'), value: 'now-2d' },
  { label: intl.get('3days'), value: 'now-3d' },
  { label: intl.get('1week'), value: 'now-1w' },
  { label: intl.get('2weeks'), value: 'now-2w' },
  { label: intl.get('4weeks'), value: 'now-4w' },
  { label: intl.get('3months'), value: 'now-3M' },
  { label: intl.get('6months'), value: 'now-6M' },
  { label: intl.get('1year'), value: 'now-1y' },
];
