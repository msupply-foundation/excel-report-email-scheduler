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

export { AppConfigProps, AppConfigStateType, ReportGroupType, ReportGroupTypeWithMembersDetail };
