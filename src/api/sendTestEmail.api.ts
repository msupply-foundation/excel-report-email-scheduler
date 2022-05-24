import { getBackendSrv } from '@grafana/runtime';

const sendTestEmail = (scheduleID: string) =>
  getBackendSrv().get(`api/plugins/msupplyfoundation-datasource/resources/test-email?schedule-id=${scheduleID}`);

export { sendTestEmail };
