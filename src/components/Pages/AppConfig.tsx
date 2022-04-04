import { AppPluginMeta, PluginConfigPageProps, SelectableValue } from '@grafana/data';
import React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { AppConfigForm } from 'components/AppConfigForm';
const queryClient = new QueryClient();

export type JsonData = {
  grafanaUsername?: string;
  isGrafanaPasswordSet?: boolean;
  senderEmailAddress?: string;
  senderEmailPassword?: string;
  isSenderEmailPasswordSet?: boolean;
  senderEmailHost?: string;
  senderEmailPort?: number;
  datasourceID?: number;
  selectedDatasource?: SelectableValue | null;
};

interface Props extends PluginConfigPageProps<AppPluginMeta<JsonData>> {}

const AppConfig = ({ plugin, ...rest }: Props) => {
  return (
    <QueryClientProvider client={queryClient}>
      <AppConfigForm plugin={plugin} {...rest} />
    </QueryClientProvider>
  );
};

export { AppConfig };
