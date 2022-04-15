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

export { createReportGroup, getReportGroups };
