import { getBackendSrv } from '@grafana/runtime';
import { ReportGroupType } from 'types';

const getReportGroups = () =>
  getBackendSrv().get('api/plugins/msupplyfoundation-excelreportemailscheduler-datasource/resources/report-group');

const createReportGroup = (reportGroup: ReportGroupType) => {
  return getBackendSrv().post(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-datasource/resources/report-group`,
    reportGroup
  );
};

export { createReportGroup, getReportGroups };
