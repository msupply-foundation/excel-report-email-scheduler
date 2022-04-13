import React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import intl from 'react-intl-universal';
import { css } from '@emotion/css';

import { LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { AppRootProps, GrafanaTheme2, getLocale } from '@grafana/data';

import { PluginPropsContext } from '../context';
import { AppRoutes } from './AppRoutes';

import { AppSettings } from 'types';

import { locales } from '../locales';

const queryClient = new QueryClient();

type AppRootState = { initDone: boolean };
class App extends React.PureComponent<AppRootProps<AppSettings>, AppRootState> {
  state: AppRootState = {
    // optional second annotation for better type inference
    initDone: false,
  };

  style = getStyles;

  componentDidMount() {
    this.loadLocales();
  }

  loadLocales() {
    // init method will load CLDR locale data according to currentLocale
    // react-intl-universal is singleton, so you should init it only once in your app
    intl
      .init({
        currentLocale: getLocale(), // TODO: determine locale here
        locales,
      })
      .then(() => {
        // After loading CLDR locale data, start to render
        this.setState({
          initDone: true,
        });
      });
  }

  render() {
    return this.state.initDone ? (
      <PluginPropsContext.Provider value={this.props}>
        <QueryClientProvider client={queryClient}>
          <AppRoutes />
        </QueryClientProvider>
      </PluginPropsContext.Provider>
    ) : (
      <div className={this.style.loadingWrapper}>
        <LoadingPlaceholder text="Loading..." />
      </div>
    );
  }
}

const getStyles = (theme: GrafanaTheme2) => ({
  loadingWrapper: css`
    display: flex;
    height: 50vh;
    align-items: center;
    justify-content: center;
  `,
});

export { App };
