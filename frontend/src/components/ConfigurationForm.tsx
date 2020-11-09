import React, { useState, useEffect } from 'react';
import intl from 'react-intl-universal';
import { Field, Form, Input, Select } from '@grafana/ui';
import { FC } from 'react';

import { FormValues } from '../types';
import { useQuery } from 'react-query';
import { getDatasources } from 'api';
import { SelectableValue } from '@grafana/data';

type OnSubmit<FormValues> = (data: FormValues) => void;

type FormProps = {
  formValues: FormValues;
  onSubmit: OnSubmit<FormValues>;
};

export const ConfigurationForm: FC<FormProps> = ({ formValues, onSubmit }) => {
  const { data: datasources } = useQuery('datasources', getDatasources);
  const [selectedDatasource, selectDatasource] = useState<SelectableValue | null>(null);

  const wrappedSubmit = (params: any) => {
    onSubmit({ ...params, datasourceID: selectedDatasource?.value?.id });
  };

  useEffect(() => {
    const found = datasources?.find((ds: any) => ds.id === formValues.datasourceID);

    if (found) {
      selectDatasource({ label: found.name, value: found });
    }
  }, [datasources]);

  return (
    <Form<FormValues> defaultValues={formValues} onSubmit={wrappedSubmit}>
      {({ register, errors, getValues, formState }) => {
        const { grafanaUsername = '', email = '', grafanaPassword = '', emailPassword = '' } = getValues();
        return (
          <>
            <Field
              label={intl.get('grafanaUsername')}
              description={intl.get('grafanaDescription')}
              invalid={!!errors.grafanaUsername}
              error={errors.grafanaUsername?.message}
            >
              <Input
                defaultValue={grafanaUsername}
                placeholder={intl.get('grafanaUsername')}
                name="grafanaUsername"
                ref={register({ required: 'required!' })}
                css=""
                loading
              />
            </Field>
            <Field
              label={intl.get('grafanaPassword')}
              invalid={!!errors.grafanaPassword}
              error={errors.grafanaPassword?.message}
            >
              <Input
                type="password"
                defaultValue={grafanaPassword}
                placeholder={intl.get('grafanaPassword')}
                name="grafanaPassword"
                ref={register({ required: intl.get('required') })}
                css=""
                loading
              />
            </Field>
            <Field
              label={intl.get('emailUsername')}
              description={intl.get('emailDescription')}
              invalid={!!errors.email}
              error={errors.email?.message}
            >
              <Input
                defaultValue={email}
                placeholder={intl.get('email')}
                name="email"
                ref={register({ required: intl.get('required') })}
                css=""
                loading
              />
            </Field>
            <Field
              label={intl.get('emailPassword')}
              invalid={!!errors.emailPassword}
              error={errors.emailPassword?.message}
            >
              <Input
                type="password"
                defaultValue={emailPassword}
                placeholder={intl.get('emailPassword')}
                name="emailPassword"
                ref={register({ required: intl.get('required') })}
                css=""
                loading
              />
            </Field>
            <Field label={intl.get('saveDetails')} description={intl.get('saveDetailsDescription')}>
              <Input
                value={intl.get('submit')}
                type="submit"
                placeholder="Submit"
                name="submit"
                css=""
                disabled={formState.dirty}
              ></Input>
            </Field>
            <Field label={intl.get('datasource')} description={intl.get('datasource_details')}>
              <Select
                value={selectedDatasource}
                options={datasources?.map((datasource: any) => ({ label: datasource.name, value: datasource })) ?? []}
                onChange={(selectedDatasource: SelectableValue) => {
                  selectDatasource(selectedDatasource);
                }}
              />
            </Field>
          </>
        );
      }}
    </Form>
  );
};
