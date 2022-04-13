import { getBackendSrv } from '@grafana/runtime';
import { ReportGroupType } from 'types';

const createReportGroup = async (reportGroup: ReportGroupType) => {
  return getBackendSrv().post(`./api/plugins/msupplyfoundation-datasource/resources/report-group`, reportGroup);
};

export { createReportGroup };
