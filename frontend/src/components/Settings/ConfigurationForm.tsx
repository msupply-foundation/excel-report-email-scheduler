import React, { useState, useEffect } from 'react';
import intl from 'react-intl-universal';
import { Field, FieldSet, Form, Input } from '@grafana/ui';
import { FC } from 'react';

import { useQuery } from 'react-query';
import { getDatasources } from 'api';
import { SelectableValue } from '@grafana/data';
import { FieldInput } from 'components/FieldInput';
import { FieldSelect } from 'components/FieldSelect';
import { FormValues } from 'common/types';

type OnSubmit<FormValues> = (data: FormValues) => void;

type FormProps = {
  formValues: FormValues;
  onSubmit: OnSubmit<FormValues>;
};

export const ConfigurationForm: FC<FormProps> = ({ formValues, onSubmit }) => {
  const { data: datasources } = useQuery('datasources', getDatasources);
  const [selectedDatasource, selectDatasource] = useState<SelectableValue | null>(null);

  const wrappedSubmit = (params: any) => onSubmit({ ...params, datasourceID: selectedDatasource?.value?.id });

  useEffect(() => {
    const found = datasources?.find((ds: any) => ds.id === formValues.datasourceID);
    if (found) {
      selectDatasource({ label: found.name, value: found });
    }
  }, [datasources, formValues.datasourceID]);

  return (
    <Form<FormValues> defaultValues={formValues} onSubmit={wrappedSubmit}>
      {({ register, errors, getValues, formState }) => {
        const {
          grafanaUsername = '',
          email = '',
          grafanaPassword = '',
          grafanaURL = '',
          emailPassword = '',
          emailHost,
          emailPort,
        } = getValues();

        return (
          <>
            <FieldSet label={intl.get('grafana_details')}>
              <FieldInput
                tooltip={intl.get('grafana_username_tooltip')}
                label={intl.get('grafana_username')}
                defaultValue={grafanaUsername}
                placeholder={intl.get('grafana_username')}
                inputName="grafanaUsername"
                invalid={!!errors.grafanaUsername}
                errorMessage={errors.grafanaUsername?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />
              <FieldInput
                type="password"
                tooltip={intl.get('grafana_password_tooltip')}
                label={intl.get('grafana_password')}
                defaultValue={grafanaPassword}
                placeholder={intl.get('grafana_password')}
                inputName="grafanaPassword"
                invalid={!!errors.grafanaPassword}
                errorMessage={errors.grafanaPassword?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />
              <FieldInput
                tooltip={intl.get('grafana_url_tooltip')}
                label={intl.get('grafana_url')}
                defaultValue={grafanaURL}
                placeholder={intl.get('grafana_url')}
                inputName="grafanaURL"
                invalid={!!errors.grafanaURL}
                errorMessage={errors.grafanaURL?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />
            </FieldSet>

            <FieldSet label={intl.get('email_details')}>
              <FieldInput
                tooltip={intl.get('email_tooltip')}
                label={intl.get('email_address')}
                defaultValue={email}
                placeholder={intl.get('email_address')}
                inputName="email"
                invalid={!!errors.email}
                errorMessage={errors.email?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />
              <FieldInput
                type="password"
                tooltip={intl.get('email_password_tooltip')}
                label={intl.get('email_password')}
                defaultValue={emailPassword}
                placeholder={intl.get('email_password')}
                inputName="emailPassword"
                invalid={!!errors.emailPassword}
                errorMessage={errors.emailPassword?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />
              <FieldInput
                tooltip={intl.get('email_host_tooltip')}
                label={intl.get('email_host')}
                defaultValue={emailHost}
                placeholder={intl.get('email_host')}
                inputName="emailHost"
                invalid={!!errors.emailHost}
                errorMessage={errors.emailHost?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />

              <FieldInput
                tooltip={intl.get('email_port_tooltip')}
                label={intl.get('email_port')}
                defaultValue={emailPort}
                placeholder={intl.get('email_port')}
                inputName="emailPort"
                invalid={!!errors.emailPort}
                errorMessage={errors.emailPort?.message ?? ''}
                register={() => register({ required: intl.get('required') })}
              />
            </FieldSet>

            <FieldSet label={intl.get('datasource_details')}>
              <FieldSelect
                label={intl.get('datasource')}
                tooltip={intl.get('datasource_tooltip')}
                value={selectedDatasource}
                options={datasources?.map((datasource: any) => ({ label: datasource.name, value: datasource })) ?? []}
                onChange={(selectedDatasource: SelectableValue) => {
                  selectDatasource(selectedDatasource);
                }}
              />
            </FieldSet>

            <Field label={intl.get('save_details')} description={intl.get('save_details_description')}>
              <Input
                value={intl.get('submit')}
                type="submit"
                placeholder={intl.get('submit')}
                name="submit"
                css=""
                disabled={formState.dirty}
              />
            </Field>
          </>
        );
      }}
    </Form>
  );
};
