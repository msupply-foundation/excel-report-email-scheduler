import * as React from 'react';
import { AppRootProps } from '@grafana/data';
import { QueryClient, QueryClientProvider } from 'react-query';

import { PluginPropsContext } from '../context';
import { AppRoutes } from './AppRoutes';

import { AppSettings } from 'types';

const queryClient = new QueryClient();
export class App extends React.PureComponent<AppRootProps<AppSettings>> {
  render() {
    return (
      <PluginPropsContext.Provider value={this.props}>
        <QueryClientProvider client={queryClient}>
          <AppRoutes />
        </QueryClientProvider>
      </PluginPropsContext.Provider>
    );
  }
}
