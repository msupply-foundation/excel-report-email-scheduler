import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { css } from 'emotion';
import { Field, Input } from '@grafana/ui';
import { ReportGroup } from './ReportSchedulesTab';
import { ReportGroupKey } from 'common/enums';

const container = css`
  display: flex;
  flex-direction: column;
  flex: 1;
  padding-right: 30px;
`;

type Props = {
  onUpdate: (key: ReportGroupKey, value: string) => void;
  reportGroup: ReportGroup | null;
};

export const EditReportGroupForm: FC<Props> = ({ onUpdate, reportGroup }) => (
  <div className={container}>
    <Field label={intl.get('name')} description={intl.get('group_name')}>
      <Input
        onChange={({ currentTarget: { value } }) => onUpdate(ReportGroupKey.NAME, value)}
        name={intl.get('name')}
        defaultValue={reportGroup.name}
        placeholder={intl.get('name')}
        css=""
      />
    </Field>

    <Field label={intl.get('description')} description={intl.get('group_description')}>
      <Input
        onChange={({ currentTarget: { value } }) => onUpdate(ReportGroupKey.DESCRIPTION, value)}
        name={intl.get('description')}
        defaultValue={reportGroup.description}
        placeholder={intl.get('description')}
        css=""
      />
    </Field>
  </div>
);
