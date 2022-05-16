import { getBackendSrv } from '@grafana/runtime';
import { ScheduleType } from 'types';

const getSchedules = () => getBackendSrv().get('api/plugins/msupplyfoundation-datasource/resources/schedule');

const createSchedule = (schedule: ScheduleType) => {
  return getBackendSrv().post(`/api/plugins/msupplyfoundation-datasource/resources/schedule`, schedule);
};

const deleteSchedule = async (scheduleID: string) => {
  return getBackendSrv().delete(`./api/plugins/msupplyfoundation-datasource/resources/schedule/${scheduleID}`);
};

const getScheduleByID = (scheduleID: string) => {
  return getBackendSrv().get(`/api/plugins/msupplyfoundation-datasource/resources/schedule/${scheduleID}`);
};

export { createSchedule, getSchedules, deleteSchedule, getScheduleByID };
