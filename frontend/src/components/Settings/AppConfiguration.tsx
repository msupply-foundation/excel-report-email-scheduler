import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { AppPluginMeta, PluginConfigPageProps } from '@grafana/data';
import { getBackendSrv } from '@grafana/runtime';
import { css } from 'emotion';
import { Button } from '@grafana/ui';
import { DisabledState } from './DisabledState';

import { ConfigurationForm } from './ConfigurationForm';
import { FormValues } from 'common/types';

interface Props extends PluginConfigPageProps<AppPluginMeta> {}

export const AppConfiguration: FC<Props> = (props: Props) => {
  const toggleAppState = (newState = true) => {
    getBackendSrv()
      .post(`/api/plugins/${props.plugin.meta.id}/settings`, {
        ...props.plugin.meta,
        enabled: newState,
        pinned: newState,
      })
      // Reload the current URL to update the app and show the sidebar
      // link and icon.
      .then(() => (window.location.href = window.location.href));
  };

  const onSubmit = (newJsonData: FormValues) => {
    const mapped = { ...newJsonData, emailPort: Number(newJsonData.emailPort) };
    getBackendSrv().post(`/api/plugins/msupply-datasource/resources/settings`, mapped);
    getBackendSrv().post(`/api/plugins/${props.plugin.meta.id}/settings`, {
      ...props.plugin.meta,
      jsonData: mapped,
    });
  };

  const defaultFormValues = {
    grafanaPassword: props.plugin.meta.jsonData?.grafanaPassword ?? '',
    grafanaUsername: props.plugin.meta.jsonData?.grafanaUsername ?? '',
    grafanaURL: props.plugin.meta.jsonData?.grafanaURL ?? '',
    email: props.plugin.meta.jsonData?.email ?? '',
    emailPassword: props.plugin.meta.jsonData?.emailPassword ?? '',
    datasourceID: props.plugin.meta.jsonData?.datasourceID ?? 1,
    datasourceName: props.plugin.meta.jsonData?.datasourceName ?? '',
    emailHost: props.plugin.meta.jsonData?.emailHost ?? 'smtp.gmail.com',
    emailPort: props.plugin.meta.jsonData?.emailPort ?? 587,
  };

  const isEnabled = props.plugin.meta.enabled;

  return isEnabled ? (
    <div
      className={css`
        margin: auto;
      `}
    >
      <ConfigurationForm formValues={defaultFormValues} onSubmit={onSubmit} />
      <Button variant="destructive" onClick={() => toggleAppState(false)}>
        {intl.get('disable')}
      </Button>
    </div>
  ) : (
    <DisabledState toggle={toggleAppState} />
  );
};
