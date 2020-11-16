import { useQuery } from 'react-query';
import { getStores } from './../api';
import { Store } from 'common/types';
import { useContext } from 'react';
import { AppDataContext } from 'containers';

export const useStores = () => {
  const { datasourceID } = useContext(AppDataContext);
  const { data: stores } = useQuery<Store[]>(['stores'], () => getStores(datasourceID));

  return stores ?? [];
};
