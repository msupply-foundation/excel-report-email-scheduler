import { getBackendSrv } from '@grafana/runtime';
import { ScheduleType } from 'types';

const getSchedules = () =>
  getBackendSrv().get('api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/schedule');

const createSchedule = (schedule: ScheduleType) => {
  return getBackendSrv().post(
    `/api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/schedule`,
    schedule
  );
};

const deleteSchedule = async (scheduleID: string) => {
  return getBackendSrv().delete(
    `./api/plugins/msupplyfoundation-excelreportemailscheduler-app/resources/schedule/${scheduleID}`
  );
};

export { createSchedule, getSchedules, deleteSchedule };
