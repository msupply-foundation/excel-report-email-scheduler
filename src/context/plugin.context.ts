import React, { useContext } from 'react';
import { AppRootProps } from '@grafana/data';

// This is used to be able to retrieve the root plugin props anywhere inside the app.
const PluginPropsContext = React.createContext<AppRootProps | null>(null);

const usePluginProps = () => {
  const pluginProps = useContext(PluginPropsContext);
  return pluginProps;
};

const usePluginMeta = () => {
  const pluginProps = usePluginProps();
  return pluginProps?.meta;
};

export { PluginPropsContext, usePluginProps, usePluginMeta };
