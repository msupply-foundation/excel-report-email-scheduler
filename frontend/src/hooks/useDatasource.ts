import { getSettings } from './../api';
import { useQuery } from 'react-query';
import { AppDataContext } from 'containers';
import { useContext } from 'react';

export const useDatasourceID = (): number => {
  const { pluginID } = useContext(AppDataContext);
  const { data: settings } = useQuery('settings', () => getSettings(pluginID));

  const { jsonData } = settings ?? {};
  const { datasourceID } = jsonData ?? {};

  return datasourceID;
};
