import * as React from 'react';
import { AppRootProps } from '@grafana/data';

import { AppRouter } from '../components/Routes';

import { PluginPropsContext } from 'context';
import { BrowserRouter } from 'react-router-dom';

export class App extends React.PureComponent<AppRootProps> {
  render() {
    return (
      <PluginPropsContext.Provider value={this.props}>
        <BrowserRouter>
          <AppRouter />
        </BrowserRouter>
      </PluginPropsContext.Provider>
    );
  }
}
