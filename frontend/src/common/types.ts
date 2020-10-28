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
  reportRecipientID: string;
  reportGroupID: string;
};

export type Schedule = {
  id?: string;
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
  storeID: string;
  lookback: number;
};
