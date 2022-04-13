import { getBackendSrv } from '@grafana/runtime';
import { ReportGroupType } from 'types';

const createReportGroup = (reportGroup: ReportGroupType) => {
  return getBackendSrv().post(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-datasource/resources/report-group`,
    reportGroup
  );
};

export { createReportGroup };
