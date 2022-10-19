import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';

export type NavItem = {
  id: string;
  text: string;
  sub?: string;
  icon?: string;
  url?: string;
};

export type AppData = {
  pluginID: number | string | undefined;
};

export type User = {
  id: string;
  e_mail: string;
  name: string;
};

export interface MyQuery extends DataQuery {
  queryText?: string;
  constant: number;
  withStreaming: boolean;
}

export const defaultQuery: Partial<MyQuery> = {
  constant: 6.5,
  withStreaming: false,
};

export interface AppSettings {}

type ReportGroupType = {
  id: string;
  name: string;
  description?: string;
  members: string[];
};

type ReportGroupTypeWithMembersDetail = {
  id: string;
  name: string;
  description?: string;
  members: User[];
};

type ScheduleType = {
  id: string;
  name: string;
  description?: string;
  interval: number;
  time: string;
  day: number;
  reportGroupID: string;
  nextReportTime?: number;
  panels: PanelListSelectedType[];
  panelDetails: PanelDetails[];
  dateFormat: string;
  datePosition: string;
};

export type VariableOption = {
  text: string;
  value: string;
};

export type Variable = {
  label: string;
  type: string;
  name: string;
  definition: string;
  includeAll: boolean;
  refresh: number;
  datasource: {
    type: string;
    uid: string;
  };
  options: VariableOption[];
  multi: boolean;
};

export type RawPanelTarget = {
  rawSql: string;
};

export type RawPanel = {
  targets: RawPanelTarget[];
  title: string;
  description: string;
  id: number;
  dashboardID: string;
  type: string;
  error?: string;
};

export type Panel = {
  id: number;
  dashboardID: string;
  variables: Variable[];
  description: string;
  title: string;
  rawSql: string;
  type: string;
  error?: string;
};

export type Templating = {
  list: Variable[];
};

export type Dashboard = {
  uid: string;
  panels: RawPanel[];
  templating: Templating;
};

export type DashboardResponse = {
  dashboard: Dashboard;
};

export type DashboardMeta = {
  uid: string;
  folderTitle: string;
  type: string;
};

/**
 * These are options configured for each DataSource instance.
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}

type AppConfigProps = {
  grafanaUsername?: string;
  isGrafanaPasswordSet?: boolean;
  grafanaURL?: string;
  senderEmailAddress?: string;
  senderEmailPassword?: string;
  isSenderEmailPasswordSet?: boolean;
  senderEmailHost?: string;
  senderEmailPort?: number;
  datasourceID?: number;
};

type AppConfigStateType = Required<AppConfigProps> & {
  grafanaPassword: string;
  senderEmailPassword: string;
  selectedDatasource?: SelectableValue | null;
};

export type ContentVariables = {
  [key: string]: string[];
};

export type SelectableVariable = {
  name: string;
  value: string;
};

export type JSONValue = string | number | boolean | { [x: string]: JSONValue } | JSONValue[];

export type PanelDetails = {
  id: string;
  scheduleID: string;
  panelID: number;
  lookback: string;
  dashboardID: string;
  variables: string | null;
};

export type PanelListSelectedType = {
  panelID: number;
  dashboardID: string;
};

export { AppConfigProps, ScheduleType, AppConfigStateType, ReportGroupType, ReportGroupTypeWithMembersDetail };
