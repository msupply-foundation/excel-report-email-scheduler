import { getBackendSrv } from '@grafana/runtime';
import { ScheduleType } from 'types';

const createSchedule = (schedule: ScheduleType) => {
  return getBackendSrv().post(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/schedule`,
    schedule
  );
};

export { createSchedule };
