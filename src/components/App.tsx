import React, { useState, useEffect } from 'react';

import { QueryClient, QueryClientProvider } from 'react-query';
import intl from 'react-intl-universal';
import { css } from '@emotion/css';

import { LoadingPlaceholder, useStyles2 } from '@grafana/ui';
import { AppRootProps, GrafanaTheme2 } from '@grafana/data';

import { PluginPropsContext } from '../context';
import { AppRoutes } from './AppRoutes';

import { AppSettings } from 'types';

import { locales } from '../locales';

const queryClient = new QueryClient();

const App = (props: AppRootProps<AppSettings>) => {
  const style = useStyles2(getStyles);
  const [initDone, setInitDone] = useState(false);

  useEffect(() => {
    loadLocales();
  }, []);

  const loadLocales = () => {
    // init method will load CLDR locale data according to currentLocale
    // react-intl-universal is singleton, so you should init it only once in your app
    intl
      .init({
        currentLocale: 'en', // TODO: determine locale here
        locales,
      })
      .then(() => {
        // After loading CLDR locale data, start to render
        setInitDone(true);
      });
  };

  return initDone ? (
    <PluginPropsContext.Provider value={props}>
      <QueryClientProvider client={queryClient}>
        <AppRoutes />
      </QueryClientProvider>
    </PluginPropsContext.Provider>
  ) : (
    <div className={style.loadingWrapper}>
      <LoadingPlaceholder text="Loading..." />
    </div>
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  loadingWrapper: css`
    display: flex;
    height: 50vh;
    align-items: center;
    justify-content: center;
  `,
});

export { App };
