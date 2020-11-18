export type ReportGroup = {
  id: string;
  name: string;
  description: string;
};

export type User = {
  id: string;
  e_mail: string;
  name: string;
};

export type ReportGroupMember = {
  id: string;
  userID: string;
  reportGroupID: string;
};

export type Schedule = {
  id: string;
  name?: string;
  description?: string;
  interval?: number;
  nextReportTime?: number;
  lookback?: number;
  reportGroupID: string;
};

export type ReportContent = {
  id: string;
  panelID: number;
  lookback: number;
  dashboardID: string;
  variables: string;
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

export type Store = {
  name: string;
  id: string;
};

export interface CreateContentVars {
  scheduleID: string;
  panelID: number;
  dashboardID: string;
  variables: string;
}
export interface CreateGroupMemberVars {
  user: User;
  reportGroupID: string;
}

export interface CreateReportGroupVariables {}

export type AppData = {
  pluginID: number | string | undefined;
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
  refresh: number;
  datasource: string;
  options: VariableOption[];
  multi: boolean;
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
};

export type ContentVariables = {
  [key: string]: string[];
};

export type SelectableVariable = {
  name: string;
  value: string;
};

export type FormValues = {
  grafanaUsername: string;
  grafanaPassword: string;
  email: string;
  emailPassword: string;
  datasourceID: number;
  emailHost: string;
  emailPort: number;
  grafanaURL: string;
};
