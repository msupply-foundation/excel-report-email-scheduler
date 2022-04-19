import { getBackendSrv } from '@grafana/runtime';
import { ReportGroupType } from 'types';

const getReportGroups = () =>
  getBackendSrv().get('api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/report-group');

const createReportGroup = (reportGroup: ReportGroupType) => {
  return getBackendSrv().post(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/report-group`,
    reportGroup
  );
};

const getReportGroupByID = (reportGroupID: string) => {
  return getBackendSrv().get(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/report-group/${reportGroupID}`
  );
};

const getReportGroupMembersByGroupID = (reportGroupID: string) => {
  return getBackendSrv().get(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/report-group-membership?&group-id=${reportGroupID}`
  );
};

export { createReportGroup, getReportGroups, getReportGroupByID, getReportGroupMembersByGroupID };
