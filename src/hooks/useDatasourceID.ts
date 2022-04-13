import { getSettings } from '../api/getSettings.api';
import { useQuery } from 'react-query';
import { usePluginMeta } from 'context';

export const useDatasourceID = (): number => {
  const pluginMeta = usePluginMeta();

  const { data: settings } = useQuery('settings', () => getSettings(pluginMeta?.id), {
    refetchOnWindowFocus: false,
  });

  const { jsonData } = settings ?? {};
  const { datasourceID } = jsonData ?? {};

  return datasourceID;
};
