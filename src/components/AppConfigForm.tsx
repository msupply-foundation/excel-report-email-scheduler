import React, { useState, ChangeEvent, useEffect } from 'react';
import { Field, Input, FieldSet, Button, useStyles2, Select } from '@grafana/ui';
import { PluginConfigPageProps, AppPluginMeta, GrafanaTheme2, PluginMeta, SelectableValue } from '@grafana/data';
import { getBackendSrv, locationService } from '@grafana/runtime';
import intl from 'react-intl-universal';
import { useQuery } from 'react-query';
import { Loading, SecretInput } from './common';
import { css } from '@emotion/css';
import { locales } from '../locales';
import { AppConfigProps, AppConfigStateType } from 'types';
import { getDatasources } from 'api/getDatasource.api';

interface Props extends PluginConfigPageProps<AppPluginMeta<AppConfigProps>> {}

const AppConfigForm = ({ plugin }: Props) => {
  const style = useStyles2(getStyles);

  const [loading, setLoading] = useState(true);
  const { data: datasources, isLoading: isDatasourceListLoading } = useQuery('datasources', getDatasources);

  const { enabled, pinned, jsonData } = plugin.meta;

  const [state, setState] = useState<AppConfigStateType>({
    grafanaUsername: jsonData?.grafanaUsername || '',
    grafanaPassword: '',
    grafanaURL: jsonData?.grafanaURL || '',
    isGrafanaPasswordSet: Boolean(jsonData?.isGrafanaPasswordSet),
    senderEmailAddress: jsonData?.senderEmailAddress || '',
    senderEmailPassword: '',
    isSenderEmailPasswordSet: Boolean(jsonData?.isSenderEmailPasswordSet),
    senderEmailHost: jsonData?.senderEmailHost || '',
    senderEmailPort: jsonData?.senderEmailPort || 0,
    datasourceID: jsonData?.datasourceID || 0,
  });

  useEffect(() => {
    intl
      .init({
        currentLocale: 'en-US',
        locales,
      })
      .then(() => {
        // After loading locale data, start to render
        setLoading(false);
      });
  }, []);

  const onResetGrafanaPassword = () =>
    setState({
      ...state,
      grafanaPassword: '',
      isGrafanaPasswordSet: false,
    });

  const onResetSenderEmailPassword = () =>
    setState({
      ...state,
      senderEmailPassword: '',
      isSenderEmailPasswordSet: false,
    });

  const onChangeGrafanaPassword = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      grafanaPassword: event.target.value.trim(),
    });
  };

  const onChangeSenderEmailPassword = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      senderEmailPassword: event.target.value.trim(),
    });
  };

  const onChangeGrafanaUsername = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      grafanaUsername: event.target.value.trim(),
    });
  };

  const onChangeGrafanaURL = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      grafanaURL: event.target.value.trim(),
    });
  };

  const onEmailAddressChange = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      senderEmailAddress: event.target.value.trim(),
    });
  };

  const onSenderEmailHost = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      senderEmailHost: event.target.value.trim(),
    });
  };

  const onSenderEmailPort = (event: ChangeEvent<HTMLInputElement>) => {
    setState({
      ...state,
      senderEmailPort: Number(event.target.value.trim()),
    });
  };

  if (loading || isDatasourceListLoading) {
    return <Loading text="Loading..." />;
  }

  return (
    !loading && (
      <div>
        {/* Grafana Username */}
        <FieldSet label={intl.get('grafana_details')}>
          <Field label={intl.get('grafana_url')} description="Full URL of your Grafana installation">
            <Input
              width={60}
              id="api-grafana-url"
              data-testid="api-grafana-url"
              label={intl.get('grafana_url')}
              value={state?.grafanaURL}
              placeholder={intl.get('grafana_url')}
              onChange={onChangeGrafanaURL}
            />
          </Field>

          <Field label={intl.get('grafana_username')}>
            <Input
              width={60}
              id="api-grafana-username"
              data-testid="api-grafana-username"
              label={intl.get('grafana_username')}
              value={state?.grafanaUsername}
              placeholder={intl.get('grafana_username')}
              onChange={onChangeGrafanaUsername}
            />
          </Field>

          <Field label={intl.get('grafana_password')} description={intl.get('grafana_password_tooltip')}>
            <SecretInput
              width={60}
              id="api-grafana-password"
              data-testid="api-grafana-password"
              label={intl.get('grafana_password')}
              value={state?.grafanaPassword}
              isConfigured={state.isGrafanaPasswordSet}
              placeholder={intl.get('grafana_password')}
              onChange={onChangeGrafanaPassword}
              onReset={onResetGrafanaPassword}
            />
          </Field>
        </FieldSet>

        <FieldSet label={intl.get('email_details')}>
          <Field label={intl.get('email_address')}>
            <Input
              width={60}
              id="api-email-address"
              data-testid="api-email-address"
              label={intl.get('email_address')}
              value={state?.senderEmailAddress}
              placeholder={intl.get('email_address')}
              onChange={onEmailAddressChange}
            />
          </Field>

          <Field label={intl.get('email_password')} description={intl.get('email_password_tooltip')}>
            <SecretInput
              width={60}
              id="api-email-password"
              data-testid="api-email-password"
              label={intl.get('email_password')}
              value={state?.senderEmailPassword}
              isConfigured={state.isSenderEmailPasswordSet}
              placeholder={intl.get('email_password')}
              onChange={onChangeSenderEmailPassword}
              onReset={onResetSenderEmailPassword}
            />
          </Field>

          <Field label={intl.get('email_host')} description={intl.get('email_host_tooltip')}>
            <Input
              width={60}
              id="api-email-host"
              data-testid="api-email-host"
              label={intl.get('email_host')}
              value={state?.senderEmailHost}
              placeholder={intl.get('email_host')}
              onChange={onSenderEmailHost}
            />
          </Field>

          <Field label={intl.get('email_port')} description={intl.get('email_port_tooltip')}>
            <Input
              width={60}
              id="api-email-port"
              data-testid="api-email-port"
              label={intl.get('email_port')}
              value={state?.senderEmailPort}
              placeholder={intl.get('email_port')}
              onChange={onSenderEmailPort}
            />
          </Field>
        </FieldSet>

        <FieldSet label={intl.get('datasource_details')}>
          <Field label={intl.get('datasource')}>
            <Select
              width={60}
              menuShouldPortal
              value={state?.datasourceID}
              options={datasources?.map((datasource: any) => ({ label: datasource.name, value: datasource.id })) ?? []}
              onChange={(selectedDatasource: SelectableValue) => {
                setState({
                  ...state,
                  datasourceID: Number(selectedDatasource.value),
                });
              }}
            ></Select>
          </Field>
        </FieldSet>

        <div className={style.marginTop}>
          <Button
            type="submit"
            onClick={() =>
              updatePluginAndReload(plugin.meta.id, {
                enabled,
                pinned,
                jsonData: {
                  grafanaUsername: state.grafanaUsername,
                  grafanaURL: state.grafanaURL,
                  isGrafanaPasswordSet: true,
                  senderEmailAddress: state.senderEmailAddress,
                  isSenderEmailPasswordSet: true,
                  senderEmailHost: state.senderEmailHost,
                  senderEmailPort: state.senderEmailPort,
                  datasourceID: state.datasourceID,
                },
                secureJsonData:
                  state.isGrafanaPasswordSet && state.isSenderEmailPasswordSet
                    ? undefined
                    : Object.fromEntries(
                        Object.entries({
                          grafanaPassword: state.grafanaPassword,
                          senderEmailPassword: state.senderEmailPassword,
                        }).filter(([_, v]) => v !== '')
                      ),
              })
            }
            disabled={Boolean(
              !state.grafanaUsername ||
                (!state.isGrafanaPasswordSet && !state.grafanaPassword) ||
                (!state.isSenderEmailPasswordSet && !state.senderEmailPassword) ||
                !state.senderEmailAddress ||
                !state.senderEmailHost ||
                !state.senderEmailPort ||
                !state.grafanaURL ||
                !state.datasourceID
            )}
          >
            Save settings
          </Button>
        </div>
      </div>
    )
  );
};

const getStyles = (theme: GrafanaTheme2) => ({
  colorWeak: css`
    color: ${theme.colors.text.secondary};
  `,
  marginTop: css`
    margin-top: ${theme.spacing(3)};
  `,
  marginTopXl: css`
    margin-top: ${theme.spacing(6)};
  `,
});

const updatePluginAndReload = async (pluginId: string, data: Partial<PluginMeta<AppConfigProps>>) => {
  try {
    await updatePlugin(pluginId, data);

    // Reloading the page as the changes made here wouldn't be propagated to the actual plugin otherwise.
    // This is not ideal, however unfortunately currently there is no supported way for updating the plugin state.
    locationService.reload();
  } catch (e) {
    console.error('Error while updating the plugin', e);
  }
};

export const updatePlugin = async (pluginId: string, data: Partial<PluginMeta>) => {
  await getBackendSrv().post(`/api/plugins/${pluginId}/resources/settings`, {
    ...data.jsonData,
    ...data.secureJsonData,
  });
  const response = await getBackendSrv().post(`/api/plugins/${pluginId}/settings`, data);
  return response?.data;
};

export { AppConfigForm };
