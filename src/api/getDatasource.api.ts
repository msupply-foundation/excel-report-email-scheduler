import { getBackendSrv } from '@grafana/runtime';

export const getDatasources = async () => {
  return getBackendSrv().get(`./api/datasources`);
};

export const getDatasource = async (id: number) => {
  return getBackendSrv().get(`./api/datasources/${id}`);
};
