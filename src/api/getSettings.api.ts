import { getBackendSrv } from '@grafana/runtime';

const getSettings = async (id: any) => {
  return getBackendSrv().get(`/api/plugins/${id}/settings`);
};

export { getSettings };
