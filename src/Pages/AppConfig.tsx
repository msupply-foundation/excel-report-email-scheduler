import { AppPluginMeta, PluginConfigPageProps } from '@grafana/data';
import React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { AppConfigForm } from '../components';
import { AppConfigProps } from '../types';

const queryClient = new QueryClient();

interface Props extends PluginConfigPageProps<AppPluginMeta<AppConfigProps>> {}

const AppConfig = ({ plugin, ...rest }: Props) => {
  return (
    <QueryClientProvider client={queryClient}>
      <AppConfigForm plugin={plugin} {...rest} />
    </QueryClientProvider>
  );
};

export { AppConfig };
