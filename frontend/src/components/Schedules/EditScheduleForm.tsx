import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Field, Input, Select } from '@grafana/ui';

import { ScheduleKey } from 'common/enums';
import { Schedule } from 'common/types';

const intervals = [
  { label: 'Daily', value: 60 * 1000 * 60 * 24 },
  { label: 'Weekly', value: 60 * 1000 * 60 * 24 * 7 },
  { label: 'Fortnightly', value: 60 * 1000 * 60 * 24 * 14 },
  { label: 'Monthly', value: 60 * 1000 * 60 * 24 * 30 },
  { label: 'Quarterly', value: 60 * 1000 * 60 * 24 * 30 * 6 },
  { label: 'Yearly', value: 60 * 1000 * 60 * 24 * 30 * 12 },
];

const container = css`
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  flex: 1;
  padding-right: 30px;
`;

const flexWrapping = css`
  flex-basis: calc(50% - 20px);
`;

type Props = {
  onUpdate: (key: ScheduleKey, value: string) => void;
  schedule: Schedule;
};

export const EditScheduleForm: FC<Props> = ({ onUpdate, schedule }) => (
  <div className={container}>
    <Field className={flexWrapping} label={intl.get('name')} description={intl.get('group_name')}>
      <Input
        onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.NAME, value)}
        name={intl.get('name')}
        defaultValue={schedule.name}
        placeholder={intl.get('name')}
        css=""
      />
    </Field>

    <Field className={flexWrapping} label={intl.get('description')} description={intl.get('group_description')}>
      <Input
        onChange={({ currentTarget: { value } }) => onUpdate(ScheduleKey.DESCRIPTION, value)}
        name={intl.get('description')}
        defaultValue={schedule.description}
        placeholder={intl.get('description')}
        css=""
      />
    </Field>
    <Field className={flexWrapping} label="Report interval" description="The interval frequency for emails to be sent">
      <Select
        value={intervals.filter((interval: any) => interval.value === schedule?.interval)}
        options={intervals}
        onChange={(selected: any) => {}}
      />
    </Field>
    <Field className={flexWrapping} label="Report interval" description="The interval frequency for emails to be sent">
      <Select
        value={intervals.filter((interval: any) => interval.value === schedule?.interval)}
        options={intervals}
        onChange={(selected: any) => {}}
      />
    </Field>
  </div>
);
