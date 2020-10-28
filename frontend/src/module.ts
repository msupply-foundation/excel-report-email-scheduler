import { ComponentClass } from 'react';
import { AppPlugin, AppRootProps } from '@grafana/data';

import { ConfigPageFactory } from './containers';
import { AppConfiguration, RootPage } from './components';

/**
 * Grafana requires an export "plugin" which is a react component which will have
 * props (AppRootProps) injected in.
 *
 * Each config page added will be a tab in the "Info" page (As well as the readme page)
 * when navigating to the plugin from either the plugin list page or from a user-defined link.
 */
export const plugin = new AppPlugin()

  // Root page navigated to through {GrafanaURL}/a/pluginId
  .setRootPage((RootPage as unknown) as ComponentClass<AppRootProps>)

  // Navigated through {GrafanaURL}/plugins/{pluginId} or {GrafanaURL}/plugins/{pluginId}?page=page1
  .addConfigPage({
    title: '\tApp Configuration',
    icon: 'fa fa-cogs',
    body: ConfigPageFactory(AppConfiguration),
    id: 'page1',
  });
